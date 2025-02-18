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
}

type GoTpl struct {
	Package string            `json:"package"`
	Tpl     template.Template `json:"tpl"`
}

func (t *Table) getModelName() string {
	if t.ModelName != "" {
		return t.ModelName
	}
	index := strings.Index(t.TableName, "_")
	ms := t.TableName[index+1:]
	t.ModelName = camelString(ms)
	return t.ModelName
}

func (t *Table) getModuleName() string {
	if t.ModuleName != "" {
		return t.ModuleName
	}
	index := strings.Index(t.TableName, "_")
	t.ModuleName = t.TableName[:index]
	return t.ModuleName
}

func (t *Table) BuildModelFields(projectName string) *TemplateModel {
	tmpl := &TemplateModel{
		Project:      projectName,
		ModelName:    t.getModelName(),
		VarFieldName: firstLowerCase(t.ModelName),
		ModuleName:   t.getModuleName(),
		TableName:    t.TableName,
		TableSchema:  t.TableSchema,
		PkColumn:     t.ColumnName,
		PkField:      camelString(t.ColumnName),
		Fields:       t.buildField(),
	}
	return tmpl
}

func (t *Table) buildField() []*ModelField {

	fields := make([]*ModelField, 0)
	for _, column := range t.Columns {
		isPk := column.ColumnName != t.ColumnName
		fields = append(fields, &ModelField{
			Name:       camelString(column.ColumnName),
			Type:       camelType(column.DataType),
			ColumnName: column.ColumnName,
			IsPk:       isPk,
			Tags:       column.buildFiledTags(isPk),
		})
	}
	return fields
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
