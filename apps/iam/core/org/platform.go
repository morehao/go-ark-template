package org

import (
	"github.com/gin-gonic/gin"
	"github.com/morehao/goark/apps/iam/dao"
	"github.com/morehao/goark/apps/iam/model"
	"github.com/morehao/golib/glog"
)

const (
	PlatformTenantCode = "PLATFORM"
	PlatformDeptCode   = "PLATFORM"
)

func GetPlatformTenant(ctx *gin.Context) (*model.TenantEntity, error) {
	tenantEntity, err := dao.NewTenantDao().GetByCond(ctx, &dao.TenantCond{
		TenantCode: PlatformTenantCode,
	})
	if err != nil {
		glog.Errorf(ctx, "[org.GetPlatformTenant] GetByCond fail, err:%v", err)
		return nil, err
	}
	return tenantEntity, nil
}

func GetPlatformDept(ctx *gin.Context, tenantID uint) (*model.DepartmentEntity, error) {
	deptEntity, err := dao.NewDepartmentDao().GetByCond(ctx, &dao.DepartmentCond{
		TenantID: tenantID,
		DeptCode: PlatformDeptCode,
	})
	if err != nil {
		glog.Errorf(ctx, "[org.GetPlatformDept] GetByCond fail, err:%v, tenantID:%d", err, tenantID)
		return nil, err
	}
	return deptEntity, nil
}