package svcuser

import (
    "fmt"
    "strings"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/morehao/goark/apps/iam/dao"
    "github.com/morehao/goark/apps/iam/internal/dto/dtouser"
    "github.com/morehao/goark/apps/iam/internal/service/svcuser/strategy"
    "github.com/morehao/goark/apps/iam/model"
    "github.com/morehao/goark/pkg/code"
    "github.com/morehao/goark/pkg/dbclient"
    "github.com/morehao/golib/glog"
    "gorm.io/gorm"
)

type registerSvc struct {
}

func NewRegisterSvc() *registerSvc {
    return &registerSvc{}
}

func (svc *registerSvc) Register(ctx *gin.Context, req *dtouser.RegisterReq) (*dtouser.RegisterResp, error) {
    registerReq := &strategy.RegisterRequest{
        Username: req.Username,
        Password: req.Password,
        Mobile:   req.Mobile,
        Email:    req.Email,
        RealName: req.RealName,
    }

    selector := strategy.NewStrategySelector()
    regStrategy, err := selector.SelectStrategy(ctx, registerReq)
    if err != nil {
        return nil, err
    }

    preResult, err := regStrategy.PreRegister(ctx, registerReq)
    if err != nil {
        return nil, err
    }

    var userID uint
    txErr := dbclient.IamDB(ctx).Transaction(func(tx *gorm.DB) error {
        personID, err := svc.createOrGetPerson(ctx, tx, preResult, registerReq)
        if err != nil {
            return err
        }
        preResult.PersonID = personID

        userID, err = svc.createUser(ctx, tx, preResult, registerReq)
        if err != nil {
            return err
        }
        return nil
    })
    if txErr != nil {
        glog.Errorf(ctx, "[registerSvc.Register] Transaction fail, err:%v", txErr)
        return nil, code.GetError(code.UserCreateError)
    }

    if err := regStrategy.PostRegister(ctx, registerReq, userID, preResult); err != nil {
        glog.Errorf(ctx, "[registerSvc.Register] PostRegister fail, userID:%d, err:%v", userID, err)
    }

    return &dtouser.RegisterResp{
        UserID:       userID,
        PersonID:     preResult.PersonID,
        Status:       string(preResult.Status),
        PersonExists: preResult.PersonExists,
        Message:      preResult.Message,
    }, nil
}

func (svc *registerSvc) createOrGetPerson(ctx *gin.Context, tx *gorm.DB, result *strategy.RegisterResult, req *strategy.RegisterRequest) (uint, error) {
    email := ""
    if req.Email != "" {
        email = strings.TrimSpace(req.Email)
    }

    personEntity, _ := dao.NewPersonDao().WithTx(tx).GetByCond(ctx, &dao.PersonCond{Email: email})
    if personEntity != nil && personEntity.ID > 0 {
        result.PersonExists = true
        return personEntity.ID, nil
    }

    newPerson := &model.PersonEntity{
        Mobile:       strings.TrimSpace(req.Mobile),
        Email:        email,
        RealName:     req.RealName,
        PasswordHash: result.PasswordHash,
        CreatedBy:    0,
        UpdatedBy:    0,
    }
    if err := dao.NewPersonDao().WithTx(tx).Insert(ctx, newPerson); err != nil {
        glog.Errorf(ctx, "[createOrGetPerson] Insert fail, err:%v", err)
        return 0, code.GetError(code.UserCreateError)
    }
    return newPerson.ID, nil
}

func (svc *registerSvc) createUser(ctx *gin.Context, tx *gorm.DB, result *strategy.RegisterResult, req *strategy.RegisterRequest) (uint, error) {
    tenant, err := dao.NewTenantDao().WithTx(tx).GetByID(ctx, result.TenantID)
    if err != nil {
        glog.Errorf(ctx, "[createUser] GetByID tenant fail, err:%v", err)
        return 0, code.GetError(code.AuthRegisterError)
    }
    if tenant == nil {
        return 0, code.GetError(code.TenantNotExistError)
    }

    employeeNo, err := svc.generateEmployeeNo(ctx, tenant.TenantCode)
    if err != nil {
        return 0, err
    }

    userEntity := &model.UserEntity{
        TenantID:   result.TenantID,
        PersonID:    result.PersonID,
        EmployeeNo:  employeeNo,
        Username:    req.Username,
        Status:      result.Status,
        UserType:    model.UserTypeNormal,
        CreatedBy:   0,
        UpdatedBy:   0,
    }
    if err := dao.NewUserDao().WithTx(tx).Insert(ctx, userEntity); err != nil {
        glog.Errorf(ctx, "[createUser] Insert fail, err:%v", err)
        return 0, code.GetError(code.UserCreateError)
    }
    return userEntity.ID, nil
}

func (svc *registerSvc) generateEmployeeNo(ctx *gin.Context, tenantCode string) (string, error) {
    if len(tenantCode) < 2 {
        tenantCode = fmt.Sprintf("%-2s", tenantCode)
    }
    today := time.Now().Format("20060102")
    key := fmt.Sprintf("employee_no:%s:%s", tenantCode, today)

    seq, err := dbclient.RedisCli.Incr(ctx, key).Result()
    if err != nil {
        glog.Errorf(ctx, "[generateEmployeeNo] Redis Incr fail, err:%v", err)
        return "", code.GetError(code.AuthRegisterError)
    }

    expiry := time.Now().AddDate(0, 0, 1)
    if _, err := dbclient.RedisCli.ExpireAt(ctx, key, expiry).Result(); err != nil {
        glog.Errorf(ctx, "[generateEmployeeNo] Redis ExpireAt fail, err:%v", err)
    }

    return fmt.Sprintf("%s%s%04d", tenantCode[:2], today, seq), nil
}
