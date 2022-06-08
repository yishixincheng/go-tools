package sql2code

import (
	"fmt"
	"github.com/yishixincheng/go-tools/internal/utils"
	"log"
	"os"
	"text/template"
)

const strcutTpl = `
package {{.Package}}

{{.ImportPackage}}

type {{.TableName | ToCamelCase}} struct {
{{ $gorm := len .Gorm }}{{ if gt $gorm 0 }}     gorm.Model{{ end }}
{{range .Columns}}	{{ $length := len .Comment}} {{ if gt $length 0 }}// {{.Comment}} {{else}}// {{.Name}} {{ end }}
	{{ $typeLen := len .Type }} {{ if gt $typeLen 0 }}{{.Name | ToCamelCase}}	{{.Type}}	{{.Tag}}{{ else }}{{.Name}}{{ end }}
{{end}}}

func (model {{.TableName | ToCamelCase}}) TableName() string {
	return "{{.TableName}}"
}

`

type StructTemplate struct {
	strcutTpl string
}

type ImportPackage string

type WrapColumnData struct {
	StructColumnList []*StructColumn
	ImportPackageList []ImportPackage
}

type StructColumn struct {
	Name    string
	Type    string
	Tag     string
	Comment string
}

type StructTemplateDB struct {
	Package   string
	TableName string
	Columns   []*StructColumn
	ImportPackage  string
	Gorm string
}

func NewStructTemplate() *StructTemplate {
	return &StructTemplate{strcutTpl: strcutTpl}
}

func (t *StructTemplate) AssemblyColumns(tbColumns []*TableColumn, gorm string) *WrapColumnData {
	tplColumns := make([]*StructColumn, 0, len(tbColumns))
	importPackageList := make([]ImportPackage, 0, 1)

	if gorm == "true" || gorm == "1" {
		importPackageList = append(importPackageList, "github.com/jinzhu/gorm")
	}

	for _, column := range tbColumns {
		gormTag := ""
		if gorm == "true" || gorm == "1" {
			gormTag += "gorm:\"type:"+column.ColumnType
			if column.IsNullable == "YES" {
				gormTag += "; null"
			} else {
				gormTag += "; not null"
			}
			if column.ColumnDefault != nil {
				gormTag += "; default:'" + *column.ColumnDefault + "'"
			}
			gormTag += "\" "
		}

		tag := fmt.Sprintf("`%s"+"json:"+"\"%s\""+"`", gormTag, column.ColumnName)

		tplColumns = append(tplColumns, &StructColumn{
			Name:    column.ColumnName,
			Type:    DBTypeToStructType[column.DataType],
			Tag:     tag,
			Comment: utils.StripLineBreakChar(column.ColumnComment),
		})
		if DBTypeToStructType[column.DataType] == "time.Time" {
			importPackageList = append(importPackageList, "time")
		}
	}

	return &WrapColumnData{
		StructColumnList: tplColumns,
		ImportPackageList: importPackageList,
	}
}

func (t *StructTemplate) SaveToModelFile(dbName string, tableName string, gorm string, wrapColumnData *WrapColumnData) error {
	tpl := template.Must(template.New("sql2struct").Funcs(template.FuncMap{
		"ToCamelCase": utils.UnderscoreToUpperCamelCase,
	}).Parse(t.strcutTpl))

	packageName := utils.ToLower(dbName)

	tplDB := StructTemplateDB{
		Package: packageName,
		TableName: tableName,
		Columns:   wrapColumnData.StructColumnList,
	}
	if gorm == "true" || gorm == "1" {
		tplDB.Gorm = "true"
	}
	if len(wrapColumnData.ImportPackageList) != 0 {
		var importPage = "import (\n"
		for i := range wrapColumnData.ImportPackageList {
			importPage += "\t\""+string(wrapColumnData.ImportPackageList[i]) + "\"\n"
		}
		importPage += ")\n"
		tplDB.ImportPackage = importPage
	}

	dir := "./dest/model/" + packageName
	if err := utils.CreateDir(dir); err != nil {
		fmt.Println("创建文件夹失败", err)
	}
	filePath := dir + "/" + utils.ToLower(utils.StripUnderscore(tableName)) + ".go"

	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("打开文件失败%v", err)
	}
	defer f.Close()

	err = tpl.Execute(f, tplDB)
	if err != nil {
		return err
	}

	fmt.Println("已完成" + tableName + "表转化，保存路径为" + filePath)
	return nil
}

