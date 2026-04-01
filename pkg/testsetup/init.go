package testsetup

import (
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/morehao/golib/biz/testkit"
)

type Initializer = testkit.Initializer
type InitializerFunc = testkit.InitializerFunc

var (
	mu                  sync.RWMutex
	initializerCreators = map[string]InitializerFunc{
		AppNameDemo: newDemoappInitializer,
		AppNameIam:  newIamappInitializer,
	}
	registeredApps = make(map[string]bool)
)

func init() {
	gin.SetMode(gin.TestMode)
}

func RegisterApp(appName string, initFunc InitializerFunc) {
	mu.Lock()
	defer mu.Unlock()
	if registeredApps[appName] {
		return
	}
	testkit.RegisterInitializer(appName, initFunc)
	registeredApps[appName] = true
}

func Initialize(appName string) {
	mu.Lock()
	if !registeredApps[appName] {
		if creator, ok := initializerCreators[appName]; ok {
			testkit.RegisterInitializer(appName, creator)
			registeredApps[appName] = true
		}
	}
	mu.Unlock()
	testkit.Initialize(appName)
}

func Close(appName string) {
	testkit.Close(appName)
}

func Init(appName string) {
	Initialize(appName)
}

func Done(appName string) {
	Close(appName)
}

type Option = testkit.Option

var (
	WithUserID            = testkit.WithUserID
	WithCompanyID         = testkit.WithCompanyID
	WithRequestID         = testkit.WithRequestID
	WithKeyValue          = testkit.WithKeyValue
	WithHeader            = testkit.WithHeader
	WithHeaders           = testkit.WithHeaders
	WithMethod            = testkit.WithMethod
	WithURL               = testkit.WithURL
	WithQueryParam        = testkit.WithQueryParam
	WithQueryParams       = testkit.WithQueryParams
	WithContentType       = testkit.WithContentType
	WithJSON              = testkit.WithJSON
	WithFormData          = testkit.WithFormData
	WithMultipartFormData = testkit.WithMultipartFormData
	WithAuth              = testkit.WithAuth
	WithBearerToken       = testkit.WithBearerToken
	WithClientIP          = testkit.WithClientIP
	WithBody              = testkit.WithBody
)

func NewContext(opts ...Option) *gin.Context {
	return testkit.NewContext(opts...)
}
