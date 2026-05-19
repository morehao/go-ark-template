package docparser

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	docclient "github.com/morehao/goark/ragforge/client/docreader"
	"github.com/morehao/goark/ragforge/client/proto"
	"github.com/morehao/golib/glog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"
)

func getMaxMessageSize() int {
	if sizeStr := os.Getenv("MAX_FILE_SIZE_MB"); sizeStr != "" {
		if size, err := strconv.Atoi(sizeStr); err == nil && size > 0 {
			return size * 1024 * 1024
		}
	}
	return 50 * 1024 * 1024
}

type GRPCDocumentReader struct {
	mu     sync.RWMutex
	conn   *grpc.ClientConn
	client proto.DocReaderClient
	addr   string
}

func NewGRPCDocumentReader(addr string) (*GRPCDocumentReader, error) {
	p := &GRPCDocumentReader{}
	if addr != "" {
		if err := p.connect(addr); err != nil {
			return nil, err
		}
	}
	return p, nil
}

func (p *GRPCDocumentReader) connect(addr string) error {
	authConfig := docclient.LoadAuthConfigFromEnv()
	opts, err := authConfig.BuildDialOptions(getMaxMessageSize())
	if err != nil {
		return fmt.Errorf("failed to build docreader dial options: %w", err)
	}
	if authConfig.TLSEnabled {
		glog.Infof(context.Background(), "TLS enabled for docreader gRPC client")
	}
	if authConfig.AuthToken != "" {
		glog.Infof(context.Background(),
			"Token authentication enabled for docreader gRPC client (TLS=%v)",
			authConfig.TLSEnabled,
		)
	}

	resolver.SetDefaultScheme("dns")

	start := time.Now()
	conn, err := grpc.Dial(addr, opts...)
	if err != nil {
		return fmt.Errorf("failed to connect to docreader: %w", err)
	}
	glog.Infof(context.Background(), "Connected to docreader in %v", time.Since(start))

	p.conn = conn
	p.client = proto.NewDocReaderClient(conn)
	p.addr = addr
	return nil
}

func (p *GRPCDocumentReader) Reconnect(addr string) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.conn != nil {
		_ = p.conn.Close()
		p.conn = nil
		p.client = nil
		p.addr = ""
	}
	return p.connect(addr)
}

func (p *GRPCDocumentReader) IsConnected() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.conn != nil
}

func (p *GRPCDocumentReader) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.conn != nil {
		return p.conn.Close()
	}
	return nil
}

var errNotConnected = fmt.Errorf("docreader service not connected")

func (p *GRPCDocumentReader) Read(ctx context.Context, req *ReadRequest) (*ReadResult, error) {
	p.mu.RLock()
	client := p.client
	p.mu.RUnlock()
	if client == nil {
		return nil, errNotConnected
	}

	protoReq := &proto.ReadRequest{
		FileContent: req.FileContent,
		FileName:    req.FileName,
		FileType:    req.FileType,
		Url:         req.URL,
		Title:       req.Title,
		RequestId:   req.RequestID,
		Config: &proto.ReadConfig{
			ParserEngine:          req.ParserEngine,
			ParserEngineOverrides: req.ParserEngineOverrides,
		},
	}

	resp, err := client.Read(ctx, protoReq)
	if err != nil {
		return nil, fmt.Errorf("gRPC Read failed: %w", err)
	}
	return fromProtoReadResponse(resp), nil
}

func (p *GRPCDocumentReader) ListEngines(ctx context.Context, overrides map[string]string) ([]ParserEngineInfo, error) {
	p.mu.RLock()
	client := p.client
	p.mu.RUnlock()
	if client == nil {
		return nil, errNotConnected
	}

	resp, err := client.ListEngines(ctx, &proto.ListEnginesRequest{ConfigOverrides: overrides})
	if err != nil {
		return nil, fmt.Errorf("gRPC ListEngines failed: %w", err)
	}

	result := make([]ParserEngineInfo, 0, len(resp.GetEngines()))
	for _, e := range resp.GetEngines() {
		result = append(result, ParserEngineInfo{
			Name:              e.GetName(),
			Description:       e.GetDescription(),
			FileTypes:         e.GetFileTypes(),
			Available:         e.GetAvailable(),
			UnavailableReason: e.GetUnavailableReason(),
		})
	}
	return result, nil
}

func fromProtoReadResponse(resp *proto.ReadResponse) *ReadResult {
	result := &ReadResult{
		MarkdownContent: resp.GetMarkdownContent(),
		ImageDirPath:    resp.GetImageDirPath(),
		Metadata:        resp.GetMetadata(),
		Error:           resp.GetError(),
	}

	for _, ref := range resp.GetImageRefs() {
		result.ImageRefs = append(result.ImageRefs, ImageRef{
			Filename:    ref.GetFilename(),
			OriginalRef: ref.GetOriginalRef(),
			MimeType:    ref.GetMimeType(),
			StorageKey:  ref.GetStorageKey(),
			ImageData:   ref.GetImageData(),
		})
	}

	return result
}
