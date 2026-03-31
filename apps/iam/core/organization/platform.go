package organization

import (
	"github.com/gin-gonic/gin"
	"github.com/morehao/goark/apps/iam/iamdao"
	"github.com/morehao/goark/apps/iam/iammodel"
	"github.com/morehao/golib/glog"
)

const (
	PlatformTenantCode = "PLATFORM"
	PlatformDeptCode   = "PLATFORM"
)

func GetPlatformTenant(ctx *gin.Context) (*iammodel.TenantEntity, error) {
	tenantEntity, err := iamdao.NewTenantDao().GetByCond(ctx, &iamdao.TenantCond{
		TenantCode: PlatformTenantCode,
	})
	if err != nil {
		glog.Errorf(ctx, "[organization.GetPlatformTenant] GetByCond fail, err:%v", err)
		return nil, err
	}
	return tenantEntity, nil
}

func GetPlatformDept(ctx *gin.Context, tenantID uint) (*iammodel.DepartmentEntity, error) {
	deptEntity, err := iamdao.NewDepartmentDao().GetByCond(ctx, &iamdao.DepartmentCond{
		TenantID: tenantID,
		DeptCode: PlatformDeptCode,
	})
	if err != nil {
		glog.Errorf(ctx, "[organization.GetPlatformDept] GetByCond fail, err:%v, tenantID:%d", err, tenantID)
		return nil, err
	}
	return deptEntity, nil
}
