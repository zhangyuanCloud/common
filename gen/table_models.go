package gen

import (
	"fmt"
	"strings"
	"text/template"
)

type Table struct {
	TableSchema  string `json:"TABLE_SCHEMA" orm:"column(TABLE_SCHEMA)"`
	TableName    string `json:"TABLE_NAME" orm:"column(TABLE_NAME)"`
	TableComment string `json:"TABLE_COMMENT" orm:"column(TABLE_COMMENT)"`
	ColumnName   string `json:"COLUMN_NAME" orm:"column(COLUMN_NAME)"`
	ModelName    string
	ModuleName   string
	Columns      []*TableColumn
}

type TableColumn struct {
	TableSchema   string `json:"TABLE_SCHEMA" orm:"column(TABLE_SCHEMA)"`
	TableName     string `json:"TABLE_NAME" orm:"column(TABLE_NAME)"`
	ColumnName    string `json:"COLUMN_NAME" orm:"column(COLUMN_NAME)"`
	DataType      string `json:"COLUMN_TYPE" orm:"column(DATA_TYPE)"`
	ColumnComment string `json:"COLUMN_COMMENT" orm:"column(COLUMN_COMMENT)"`
	IsNullable    string `json:"IS_NULLABLE" orm:"column(IS_NULLABLE)"`
}

type GoTpl struct {
	Package string            `json:"package"`
	Tpl     template.Template `json:"tpl"`
}

func (t *Table) parseTableName(prefix string, isModule bool) {
	tn := t.TableName
	if prefix != "" {
		tn = strings.Replace(t.TableName, prefix, "", 1)
		t.TableName = tn
	}
	split := strings.Split(tn, "_")
	if isModule {
		t.ModuleName = split[0]
		t.ModelName = camelStr(split[1:]...)
		return
	}
	t.ModelName = camelStr(split...)
}

func (t *Table) BuildModelFields(projectName string) *TemplateModel {
	tmpl := &TemplateModel{
		Project:      projectName,
		ModelName:    t.ModelName,
		VarFieldName: firstLowerCase(t.ModelName),
		ModuleName:   t.ModuleName,
		TableName:    t.TableName,
		TableSchema:  t.TableSchema,
		PkColumn:     t.ColumnName,
		PkField:      camelString(t.ColumnName),
	}

	tmpl.Fields, tmpl.IsTime = t.buildField()

	return tmpl
}

func (t *Table) buildField() ([]*ModelField, bool) {

	fields := make([]*ModelField, 0)
	isTime := false
	for _, column := range t.Columns {
		isPk := column.ColumnName != t.ColumnName
		fieldType := camelType(column.DataType)
		isTime = isTime || fieldType == "time.Time"
		fields = append(fields, &ModelField{
			Name:       camelStr(strings.Split(column.ColumnName, "_")...),
			Type:       fieldType,
			ColumnName: column.ColumnName,
			IsPk:       isPk,
			Tags:       column.buildFiledTags(isPk),
			FormTags:   column.buildFormTags(),
		})
	}
	return fields, isTime
}

func (c *TableColumn) buildFiledTags(isPk bool) []*Tag {

	ormTag := &Tag{Name: "orm"}
	if !isPk {
		ormTag.Value = fmt.Sprintf("pk,column(%s)", c.ColumnName)
	} else {
		ormTag.Value = fmt.Sprintf("column(%s)", c.ColumnName)
	}

	return []*Tag{ormTag,
		{
			Name:  "json",
			Value: camelJSONTag(c.ColumnName),
		},
		{
			Name:  "comment",
			Value: c.ColumnComment,
		},
	}
}

func (c *TableColumn) buildFormTags() []*Tag {
	tags := []*Tag{
		{
			Name:  "form",
			Value: camelJSONTag(c.ColumnName),
		},
		{
			Name:  "json",
			Value: camelJSONTag(c.ColumnName),
		},
		{
			Name:  "comment",
			Value: c.ColumnComment,
		},
	}
	if c.IsNullable == "NO" {
		tags = append(tags, &Tag{
			Name:  "binding",
			Value: "required",
		})
	}
	return tags
}
