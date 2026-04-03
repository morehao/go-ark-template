package testsetup

import (
	"fmt"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/morehao/goark/apps/iam/iamdao"
	"github.com/morehao/goark/pkg/contextkeys"
	"github.com/morehao/goark/pkg/testkit"
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

func NewContext(opts ...testkit.Option) *gin.Context {
	return testkit.NewContext(opts...)
}

func WithAuthByUserID(userID uint) testkit.Option {
	return func(ctx *gin.Context) {
		userEntity, err := iamdao.NewUserDao().GetByID(ctx, userID)
		if err != nil {
			panic(fmt.Sprintf("WithAuthByUserID: get user failed, userID=%d, err=%v", userID, err))
		}
		if userEntity == nil || userEntity.ID == 0 {
			panic(fmt.Sprintf("WithAuthByUserID: user not found, userID=%d", userID))
		}

		ctx.Set(string(contextkeys.KeyUserID), userID)
		ctx.Set(string(contextkeys.KeyTenantID), userEntity.TenantID)
		ctx.Set(string(contextkeys.KeyDeptID), userEntity.DeptID)
		ctx.Set(string(contextkeys.KeyPersonID), userEntity.PersonID)

		if userEntity.TenantID > 0 {
			tenantEntity, err := iamdao.NewTenantDao().GetByID(ctx, userEntity.TenantID)
			if err != nil {
				panic(fmt.Sprintf("WithAuthByUserID: get tenant failed, tenantID=%d, err=%v", userEntity.TenantID, err))
			}
			if tenantEntity != nil && tenantEntity.OrganizationID > 0 {
				ctx.Set(string(contextkeys.KeyOrgID), tenantEntity.OrganizationID)
			}
		}
	}
}
