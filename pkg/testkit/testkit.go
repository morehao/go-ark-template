package testkit

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"

	"github.com/gin-gonic/gin"
)

// Initializer defines the interface for test environment initialization.
type Initializer interface {
	Initialize() error
	Close() error
}

// InitializerFunc is a factory function that creates an Initializer.
type InitializerFunc func() (Initializer, error)

// Option configures a test gin.Context.
type Option func(*gin.Context)

// BaseInitializer provides a minimal base for app-specific initializers.
type BaseInitializer struct {
	AppName string
}

// NewBaseInitializer creates a BaseInitializer for the given app.
func NewBaseInitializer(appName string) (*BaseInitializer, error) {
	return &BaseInitializer{AppName: appName}, nil
}

var (
	mu           sync.RWMutex
	creators     = make(map[string]InitializerFunc)
	initializers = make(map[string]Initializer)
)

// RegisterInitializer registers a factory function for the named app.
func RegisterInitializer(name string, fn InitializerFunc) {
	mu.Lock()
	defer mu.Unlock()
	creators[name] = fn
}

// Initialize creates and initializes the named app's test environment.
func Initialize(name string) {
	mu.RLock()
	if _, ok := initializers[name]; ok {
		mu.RUnlock()
		return
	}
	mu.RUnlock()

	mu.Lock()
	defer mu.Unlock()

	if _, ok := initializers[name]; ok {
		return
	}

	fn, ok := creators[name]
	if !ok {
		panic(fmt.Sprintf("testkit: no initializer registered for %q", name))
	}

	init, err := fn()
	if err != nil {
		panic(fmt.Sprintf("testkit: create initializer for %q failed: %v", name, err))
	}

	if err := init.Initialize(); err != nil {
		panic(fmt.Sprintf("testkit: initialize %q failed: %v", name, err))
	}

	initializers[name] = init
}

// Close shuts down the named app's test environment.
func Close(name string) {
	mu.Lock()
	defer mu.Unlock()

	if init, ok := initializers[name]; ok {
		_ = init.Close()
		delete(initializers, name)
	}
}

// NewContext creates a test gin.Context with the given options applied.
func NewContext(opts ...Option) *gin.Context {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	for _, opt := range opts {
		opt(c)
	}
	return c
}
