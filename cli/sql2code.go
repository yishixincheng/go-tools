package cli

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/yishixincheng/go-tools/internal/sql2code"
	"log"
)

var (
	dbName    string
	tableName string
)

var sql2codeCmd = &cobra.Command{
	Use:   "sql2code",
	Short: "sql转code",
	Long:  "sql转code",
	Run:   handleSql2Code,
}

// init -
func init()  {
	flags := sql2codeCmd.Flags()
	flags.StringVarP(&dbName, "database","d", viper.GetString("mysql.database"), "数据库名称")
	flags.StringVarP(&tableName, "table", "t", "","表名")
	registerMap["sql2code"] = sql2codeCmd
}

// handleSql2Code -
func handleSql2Code(cmd *cobra.Command, args []string) {
	fmt.Println(tableName, dbName, cmd.CommandPath())
	dbType := viper.GetString("dbType")
	if dbType == "" {
		dbType = "mysql"
	}
	dbService := sql2code.NewDBService(&sql2code.DBConfig{
		Type: dbType,
		Host: viper.GetString(dbType + ".host"),
		User: viper.GetString(dbType + ".user"),
		Password: viper.GetString(dbType + ".password"),
		Charset: viper.GetString(dbType + ".charset"),
	})
	if err := dbService.Connect(); err != nil {
		log.Fatalf("数据库链接错误：%v", err)
	}
	tables, err := dbService.ParseTableName(dbName, tableName)
	if err != nil {
		log.Fatalf("查询表错误：%v", err)
	}
	if tables == nil || len(tables) == 0 {
		log.Fatalf("数据库%s下没有查询到表", dbName)
	}
	template := sql2code.NewStructTemplate()
	for i := range tables {
		table := tables[i]
		columns, err := dbService.GetColumns(dbName, table)
		if err != nil {
			log.Fatalf("查询表结构报错：%v", err)
		}
		wrapColumnData := template.AssemblyColumns(columns)
		err = template.SaveToModelFile(dbName, table, wrapColumnData)
		if err != nil {
			log.Fatalf("template.Generate err: %v", err)
		}
	}
}

