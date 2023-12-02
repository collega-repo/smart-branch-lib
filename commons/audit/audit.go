package audit

import (
	"context"
	"errors"
	"fmt"
	"github.com/collega-repo/smart-branch-lib/commons"
	"github.com/collega-repo/smart-branch-lib/commons/info"
	"github.com/collega-repo/smart-branch-lib/configs"
	"github.com/collega-repo/smart-branch-lib/models"
	"gorm.io/gorm/schema"
	"reflect"
	"strings"
	"time"
)

type execType string

const (
	InsertType execType = "I"
	UpdateType execType = "U"
	DeleteType execType = "D"
)

type DataAudit struct {
	NewData  any
	ExecType execType
}

func GetDataLogAudits(ctx context.Context, listData []DataAudit) (logAudits []models.LogAudit, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf(`%v`, r)
			configs.Loggers.Err(err).Msg(`error get log audit data`)
		}
	}()
	var logAudit models.LogAudit
	for _, data := range listData {
		execType := data.ExecType
		logAudit, err = getDataLogAudit(ctx, data.NewData, execType)
		if err != nil {
			return nil, err
		}
		if len(logAudit.LogAuditDtls) > 0 {
			logAudits = append(logAudits, logAudit)
		}
	}
	return
}

func getDataLogAudit(ctx context.Context, newData any, execType execType) (logAudit models.LogAudit, err error) {
	var vOldData reflect.Value
	var vNewData reflect.Value
	var oldData any
	var tableName string
	var tableKeyString string
	var tableKey = make(map[string]any, 0)

	vNewData = reflect.ValueOf(newData)
	if vNewData.Kind() != reflect.Pointer {
		return logAudit, errors.New(`new data should be pointer`)
	}
	vNewData = vNewData.Elem()

	if schemaTable, ok := newData.(schema.Tabler); ok {
		tableName = schemaTable.TableName()
	}

	for i := 0; i < vNewData.NumField(); i++ {
		vNew := fmt.Sprintf(`%v`, vNewData.Field(i).Interface())
		keys := vNewData.Type().Field(i).Tag.Get(`gorm`)
		if strings.Contains(keys, `index:`) && strings.Contains(keys, `unique`) {
			for _, key := range strings.Split(keys, `;`) {
				if strings.HasPrefix(key, `column:`) {
					column := strings.ReplaceAll(key, `column:`, ``)
					tableKeyString = fmt.Sprintf(`,%s=%s%s`, column, vNew, tableKeyString)
					tableKey[column] = vNewData.Field(i).Interface()
				}
			}
		}
	}

	switch execType {
	case UpdateType, DeleteType:
		//copy pointer
		oldData = reflect.New(vNewData.Type()).Interface()
		if err = configs.DB.Where(tableKey).First(oldData).Error; err != nil {
			return logAudit, err
		}
		vOldData = reflect.ValueOf(oldData)
		vOldData = vOldData.Elem()
	case InsertType:
		vOldData = vNewData
	}

	reqInfo := info.GetRequestInfo(ctx)
	logAudit = models.LogAudit{
		ApplId:   "",
		BranchId: reqInfo.BranchId,
		DbName:   commons.Configs.Datasource.DB.Database,
		TableNm:  tableName,
		TableKey: strings.Replace(tableKeyString, `,`, ``, 1),
		ExecType: string(execType),
		ExecUser: reqInfo.UserId,
		IpAddr:   reqInfo.IpAddr,
		SpvId:    "",
	}

	extract(vNewData, vOldData, &logAudit, tableKey, execType)
	return logAudit, err
}

func extract(vNewData, vOldData reflect.Value, logAudit *models.LogAudit, tableKey map[string]any, execType execType) {
	switch execType {
	case InsertType, UpdateType:
		for i := 0; i < vNewData.NumField(); i++ {
			fieldNew := vNewData.Field(i)
			vNew := fmt.Sprintf(`%v`, fieldNew.Interface())
			switch execType {
			case UpdateType:
				vOld := fmt.Sprintf(`%v`, vOldData.Field(i).Interface())
				if vNew != vOld {
					for _, key := range strings.Split(vNewData.Type().Field(i).Tag.Get(`gorm`), `;`) {
						if strings.HasPrefix(key, `column:`) &&
							!strings.HasPrefix(key, `column:created_at`) &&
							!strings.HasPrefix(key, `column:created_by`) {
							logAuditDtl := models.LogAuditDtl{
								FieldNm: strings.ReplaceAll(key, `column:`, ``),
								NewData: vNew,
								OldData: vOld,
							}
							logAudit.LogAuditDtls = append(logAudit.LogAuditDtls, logAuditDtl)
							continue
						}
					}
				}
			case InsertType:
				for _, key := range strings.Split(vNewData.Type().Field(i).Tag.Get(`gorm`), `;`) {
					if strings.HasPrefix(key, `column:`) &&
						!strings.HasPrefix(key, `column:updated_at`) &&
						!strings.HasPrefix(key, `column:updated_by`) {
						logAuditDtl := models.LogAuditDtl{
							FieldNm: strings.ReplaceAll(key, `column:`, ``),
							NewData: vNew,
						}
						logAudit.LogAuditDtls = append(logAudit.LogAuditDtls, logAuditDtl)
						continue
					}
				}
			}
			if len(logAudit.LogAuditDtls) > 0 {
				structName := vNewData.Field(i).Type().Name()
				if structName == "MstModel" {
					extract(vNewData.Field(i), vOldData.Field(i), logAudit, tableKey, execType)
				}
			}
		}
	case DeleteType:
		for k, v := range tableKey {
			logAuditDtl := models.LogAuditDtl{
				FieldNm: k,
				NewData: fmt.Sprintf(`%v`, v),
			}
			logAudit.LogAuditDtls = append(logAudit.LogAuditDtls, logAuditDtl)
		}
	}
}

func CreateLogAudits(logAudits []models.LogAudit) {
	if len(logAudits) == 0 {
		return
	}

	go func() {
		var err error
		tx := configs.DB.Begin()
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf(`%v`, r)
			}
			if err != nil {
				configs.Loggers.Err(err).Msg(`error create log audit`)
				if err := tx.Rollback().Error; err != nil {
					configs.Loggers.Err(err).Msg(`error rollback transaction log audit`)
				}
				return
			}
			if err = tx.Commit().Error; err != nil {
				configs.Loggers.Err(err).Msg(`error commit transaction log audit`)
			} else {
				configs.Loggers.Info().Msg(`success create log audit`)
			}
		}()

		for _, logAudit := range logAudits {
			logAudit.ExecDt = time.Now()
			logAudit.TimeStamp = time.Now()
			if err = tx.Omit(`LogAuditDtls`).Create(&logAudit).Error; err != nil {
				return
			}

			for _, logAuditDtl := range logAudit.LogAuditDtls {
				logAuditDtl.LogId = logAudit.Id
				if err = tx.Create(&logAuditDtl).Error; err != nil {
					return
				}
			}
		}
	}()
}
