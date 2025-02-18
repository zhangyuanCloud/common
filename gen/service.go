package gen

import (
	"fmt"
	"github.com/beego/beego/v2/client/orm"
	"github.com/zhangyuanCloud/common/logger"
	"os"
	"strings"
	"text/template"
)

const (
	PackageModels     = "models"
	PackageService    = "service"
	PackageRepo       = "repository"
	PackageController = "controller"
	PackageValidate   = "validate"
)

var tablePackages = []string{PackageModels, PackageService, PackageRepo, PackageController, PackageValidate}

var tableTplMap = map[string]template.Template{}

func init() {

	basePath := "gen/template/"

	for _, tablePackage := range tablePackages {
		tplPath := basePath + tablePackage + ".go.tpl"
		tmpl, err := template.ParseFiles(tplPath)
		if err != nil {
			fmt.Println("模板文件读取失败:"+tplPath, err)
			continue
		}
		tableTplMap[tablePackage] = *tmpl
	}
}

func ReadTableSchema(schema string) ([]*Table, error) {
	o := orm.NewOrm()
	var tables []*Table
	qs := o.Raw("SELECT a.TABLE_SCHEMA, a.TABLE_NAME,a.TABLE_COMMENT, b.COLUMN_NAME  " +
		"FROM information_schema.`TABLES` a " +
		"LEFT JOIN information_schema.KEY_COLUMN_USAGE b ON a.TABLE_SCHEMA=b.TABLE_SCHEMA and a.TABLE_NAME=b.TABLE_NAME" +
		" WHERE b.CONSTRAINT_NAME = 'PRIMARY'  AND a.TABLE_SCHEMA = '" + schema + "'")
	_, err := qs.QueryRows(&tables)
	return tables, err
}

func ReadTableColumns(schema, tableName string) ([]*TableColumn, error) {
	var columns []*TableColumn
	sql := fmt.Sprintf("SELECT * FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = '%s' AND TABLE_NAME = '%s'", schema, tableName)
	_, err := orm.NewOrm().Raw(sql).QueryRows(&columns)
	return columns, err
}

func BuildTableTplCode(tplm *TemplateModel) error {

	for _, tablePackage := range tablePackages {
		modelTpl := tableTplMap[tablePackage]
		modelFile := tplm.getFile(tablePackage)
		model, err := os.OpenFile(modelFile, os.O_CREATE|os.O_WRONLY, 0755)
		if err != nil {
			logger.LOG.Error("GO文件打开失败:"+modelFile, err)
			return err
		}

		if err := modelTpl.Execute(model, tplm); err != nil {
			logger.LOG.Error("GO文件写入失败:"+modelFile, err)
			return err
		}
		logger.LOG.Debugf("GO文件写入成功 Package：%s 文件路径：%s", tablePackage, modelFile)
	}
	return nil
}

// 驼峰命名转换
func camelString(ms string) string {
	if ms == "" {
		return ms
	}

	data := make([]byte, 0, len(ms))
	flag, num := true, len(ms)-1
	for i := 0; i <= num; i++ {
		d := ms[i]
		if d == '_' {
			flag = true
			continue
		} else if flag {
			if d >= 'a' && d <= 'z' {
				d = d - 32
			}
			flag = false
		}
		data = append(data, d)
	}
	return string(data[:])
}

func convertName(fieldName string) string {
	temp := strings.Split(fieldName, "_")
	var upperStr string
	for y := 0; y < len(temp); y++ {
		vv := []rune(temp[y])
		for i := 0; i < len(vv); i++ {
			if i == 0 {
				upperStr += strings.ToUpper(string(vv[i]))
			} else {
				upperStr += string(vv[i])
			}
		}

	}
	return upperStr
}

func firstLowerCase(str string) string {
	if len(str) <= 0 {
		return ""
	}
	return strings.ToLower(string(str[0])) + str[1:]
}
func camelJSONTag(fieldName string) string {
	temp := strings.Split(fieldName, "_")
	var upperStr string
	for y := 0; y < len(temp); y++ {
		vv := []rune(temp[y])
		if y != 0 {
			for i := 0; i < len(vv); i++ {
				if i == 0 {
					vv[i] -= 32
					upperStr += string(vv[i]) // + string(vv[i+1])
				} else {
					upperStr += string(vv[i])
				}
			}
		}
	}
	return temp[0] + upperStr
}

func camelType(fieldType string) string {
	var key string
	if fieldType == "" {
		return key
	}

	switch fieldType {
	case "varchar":
		key = "string"
	case "char":
		key = "string"
	case "text":
		key = "string"
	case "double":
		key = "float32"
	case "float":
		key = "float64"
	case "date":
		key = "time.Time"
	case "datetime":
		key = "time.Time"
	case "timestamp":
		key = "time.Time"
	case "bigint":
		key = "int64"
	case "int":
		key = "int"
	case "tinyint":
		key = "int32"
	default:
		key = fieldType
	}

	return key
}
