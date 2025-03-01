package repository

import (
	"github.com/zhangyuanCloud/common"
	"github.com/zhangyuanCloud/common/logger"
	"github.com/zhangyuanCloud/common/database"
)

var {{.VarFieldName}}Repo *{{.ModelName}}Repo

type {{.ModelName}}Repo struct {
	common.BaseRepo
}

func New{{.ModelName}}Repo() *{{.ModelName}}Repo {
	if {{.VarFieldName}}Repo == nil {
		{{.VarFieldName}}Repo = &{{.ModelName}}Repo{
			BaseRepo: common.BaseRepo{
				TableName: database.TableName("{{.TableName}}"),
				Log:       logger.LOG.WithField("module", "{{.ModelName}}Repo"),
			},
		}
	}
	return {{.VarFieldName}}Repo
}