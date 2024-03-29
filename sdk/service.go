// Copyright 2023 ztlcloud.com
// leovs @2023.12.12

package sdk

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	_const "github.com/leovs/yc-go-sdk/const"
	"github.com/leovs/yc-go-sdk/errors"
	"github.com/leovs/yc-go-sdk/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"gorm.io/plugin/dbresolver"
	"reflect"
	"strings"
)

// TPage 分页查询结构体
type TPage[T any] struct {
	Page     int   `json:"page"`
	PageSize int   `json:"pageSize"`
	Total    int64 `json:"total"`
	Lists    []T   `json:"lists"`
}

// CreateTPage 创建分页体
func CreateTPage[T any](lists []T) *TPage[T] {
	return &TPage[T]{
		Page:     0,
		PageSize: 0,
		Total:    0,
		Lists:    lists,
	}
}

// Service 服务
type Service[T any, R any] struct {
	Data T
	Ctx  *fiber.Ctx
	db   *gorm.DB
}

// WithContext 绑定上下文
func (s *Service[T, R]) WithContext(c *fiber.Ctx) *errors.Message {
	// 上下文取数据库
	s.Ctx = c
	idb := c.Locals("db")
	if idb == nil {
		return _const.DbConnectError
	}
	s.db = idb.(*gorm.DB)
	return nil
}

func (s *Service[T, R]) DB() *gorm.DB {
	if s.db == nil {
		s.db = Runtime.GetDb()
	}
	return s.db
}

// GetDb 获取指定DB
func (s *Service[T, R]) GetDb(dbName string) *gorm.DB {
	return s.DB().Clauses(dbresolver.Use(dbName)).Model(new(T))
}

// GetWriteDb 获取DB
func (s *Service[T, R]) GetWriteDb() *gorm.DB {
	return s.DB().Clauses(dbresolver.Write).Model(new(T))
}

// GetReadDb 获取DB
func (s *Service[T, R]) GetReadDb() *gorm.DB {
	return s.DB().Clauses(dbresolver.Read).Model(new(T))
}

func (s *Service[T, R]) GetCacheDb(ttl int64) *gorm.DB {
	return s.Cache(ttl).DB().Clauses(dbresolver.Read).Model(new(T))
}

// FindById 获取信息
func (s *Service[T, R]) FindById(id uint) (result *R, err *errors.Message) {
	return s.FindByField("id", id)
}

// FindByField 获取信息
func (s *Service[T, R]) FindByField(field string, value any) (result *R, err *errors.Message) {
	return s.FindByFields(map[string]interface{}{field: value})
}

// FindByFields 获取信息
func (s *Service[T, R]) FindByFields(fields any) (result *R, err *errors.Message) {
	if err := s.GetReadDb().Where(fields).First(&result).Error; err != nil {
		return nil, _const.NoDataReturn
	}
	return
}

// Update 更新信息
// omits 忽略字段
func (s *Service[T, R]) Update(id uint, data *T, omits ...string) *errors.Message {
	if err := GormErrorAs(s.GetWriteDb().Omit(omits...).Where("id = ?", id).
		Updates(&data)); err != nil {
		return _const.Failure.SetData(err.Error())
	}
	return nil
}

// Create 创建
// omits 忽略字段
func (s *Service[T, R]) Create(data *T, omits ...string) *errors.Message {
	if err := GormErrorAs(s.GetWriteDb().Omit(omits...).Create(&data)); err != nil {
		return _const.Failure.SetData(err.Error())
	}
	return nil
}

// Delete 删除
// omits 忽略字段
func (s *Service[T, R]) Delete(id uint, omits ...string) *errors.Message {
	if err := GormErrorAs(s.GetWriteDb().Where("id = ?", id).
		Omit(omits...).Delete(&s.Data)); err != nil {
		return _const.Failure.SetData(err.Error())
	}
	return nil
}

// Destroy 彻底删除
// omits 忽略字段
func (s *Service[T, R]) Destroy(id uint, omits ...string) *errors.Message {
	if err := GormErrorAs(s.GetWriteDb().Unscoped().Where("id = ?", id).
		Omit(omits...).Delete(&s.Data)); err != nil {
		return _const.Failure.SetData(err.Error())
	}
	return nil
}

// Search 分页查询
func (s *Service[T, R]) Search(page int, pageSize int, params any) (*TPage[R], *errors.Message) {
	result := CreateTPage([]R{})
	search := s.GetReadDb()

	result.Page = page
	result.PageSize = pageSize

	MakeCondition(params, search)

	// 获取总条数
	search.Count(&result.Total)

	// 获取分页数据
	if err := search.Offset((page - 1) * pageSize).Limit(pageSize).
		Find(&result.Lists).Error; err != nil {
		return nil, _const.Failure.SetMsg(err.Error())
	}
	return result, nil
}

// Cache 是否开启缓存
func (s *Service[T, R]) Cache(ttl int64) *Service[T, R] {
	s.DB().Statement.Settings.Store(_const.GormCacheEnablePrefix, true)
	s.DB().Statement.Settings.Store(_const.GormCacheTTLPrefix, ttl)
	return s
}

// MakeCondition 生成查询器
func MakeCondition(params any, gorm *gorm.DB) *gorm.DB {
	fieldStruct := reflect.TypeOf(params)
	valueStruct := reflect.ValueOf(params)
	for i := 0; i < fieldStruct.NumField(); i++ {
		field := fieldStruct.Field(i)
		value := valueStruct.FieldByName(field.Name)
		kind := field.Type.Kind()
		tag := field.Tag.Get("search")
		tagSetting := schema.ParseTagSetting(tag, ";")

		if value.IsZero() || tag == "-" {
			continue
		}

		if !value.CanSet() && kind == reflect.Struct {
			MakeCondition(valueStruct.Field(i).Interface(), gorm)
			continue
		}

		column := tagSetting["COLUMN"]
		switch tagSetting["CLAUSE"] {
		case "exact":
			gorm.Where(fmt.Sprintf("%s = ?", column), value.Interface())
		case "contains":
			gorm.Where(fmt.Sprintf("%s like ?", column), value.Interface())
		case "gt":
			gorm.Where(fmt.Sprintf("%s > ?", column), value.Interface())
		case "gte":
			gorm.Where(fmt.Sprintf("%s >= ?", column), value.Interface())
		case "lt":
			gorm.Where(fmt.Sprintf("%s < ?", column), value.Interface())
		case "lte":
			gorm.Where(fmt.Sprintf("%s <= ?", column), value.Interface())
		case "in":
			gorm.Where(fmt.Sprintf("%s in ?", column), value.Interface())
		case "order":
			gorm.Order(fmt.Sprintf("%s %s", column, utils.Ternary(value.Int() == 1, "asc", "desc")))
		}
	}
	return gorm
}

// GormErrorAs 检查Gorm执行错误
func GormErrorAs(gorm *gorm.DB) error {
	if gorm.Error != nil {
		e := gorm.Error.Error()
		if strings.Contains(e, "a foreign key constraint fails") {
			return _const.DataReferenced
		}
		if strings.Contains(e, "Duplicate entry") {
			return _const.DataExisted
		}
		return gorm.Error
	}
	if gorm.RowsAffected <= 0 {
		return _const.NoDataReturn
	}
	return nil
}
