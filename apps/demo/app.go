package demo

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/morehao/goark/demo/config"
	_ "github.com/morehao/goark/demo/docs"
	"github.com/morehao/goark/demo/internal/middleware"
	"github.com/morehao/goark/demo/internal/router"
	"github.com/morehao/goark/pkg/dbclient"
	"github.com/morehao/golib/biz/gserver/gindocs"
	"github.com/morehao/golib/biz/gserver/ginserver"
	"github.com/morehao/golib/filestore"
	"github.com/morehao/golib/glog"
	"github.com/morehao/golib/storage"
	_ "github.com/morehao/golib/storage/driver/local"
)

const AppName = "demo"

func Routers(engine *gin.Engine) {
	routerGroups := ginserver.NewRouterGroups(engine, AppName, ginserver.VersionGroup{
		Version: ginserver.ApiVersionV1,
		Middlewares: []gin.HandlerFunc{
			middleware.Example(),
		},
	})

	if config.Conf.Server.Env == "dev" {
		gindocs.Register(engine.Group("/"+AppName), AppName)
	}

	fileStore, err := initFileStore()
	if err != nil {
		panic(fmt.Errorf("demo.Routers: init file store failed: %w", err))
	}

	router.RegisterRouter(routerGroups, AppName, fileStore)
}

// initFileStore 根据配置初始化文件上传存储（storage 驱动 + filestore）。
func initFileStore() (*filestore.FileStore, error) {
	cfg := config.Conf.FileStorage

	st, err := storage.New(cfg.Driver, cfg.Storage)
	if err != nil {
		return nil, fmt.Errorf("storage.New: %w", err)
	}

	var storeOpts []filestore.StoreOption
	if cfg.Storage.SignSecret != "" {
		storeOpts = append(storeOpts, filestore.WithSignSecret(cfg.Storage.SignSecret))
	}

	fs, err := filestore.New(dbclient.DemoDB(context.TODO()), st, cfg.Bucket, storeOpts...)
	if err != nil {
		return nil, fmt.Errorf("filestore.New: %w", err)
	}
	glog.Infof(context.TODO(), "[demo.initFileStore] filestore init done, driver=%s, bucket=%s, local=%v", cfg.Driver, cfg.Bucket, fs.IsLocal())
	return fs, nil
}
