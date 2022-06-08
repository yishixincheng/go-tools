package sql2code

import (
	"errors"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type DBService struct {
	DB       *sqlx.DB
	DBConfig *DBConfig
}

type DBConfig struct {
	Type     string
	Host     string
	User     string
	Password string
	Charset  string
}


type TableColumn struct {
	ColumnName    string
	DataType      string
	IsNullable    string
	ColumnKey     string
	ColumnType    string
	ColumnDefault *string
	ColumnComment string
}

var DBTypeToStructType = map[string]string{
	"int":        "int32",
	"tinyint":    "int8",
	"smallint":   "int",
	"mediumint":  "int64",
	"bigint":     "int64",
	"bit":        "int",
	"bool":       "bool",
	"enum":       "string",
	"set":        "string",
	"varchar":    "string",
	"char":       "string",
	"tinytext":   "string",
	"mediumtext": "string",
	"text":       "string",
	"longtext":   "string",
	"blob":       "string",
	"tinyblob":   "string",
	"mediumblob": "string",
	"longblob":   "string",
	"date":       "time.Time",
	"datetime":   "time.Time",
	"timestamp":  "time.Time",
	"time":       "time.Time",
	"float":      "float64",
	"double":     "float64",
}

func NewDBService(conf *DBConfig) *DBService {
	return &DBService{DBConfig: conf}
}

func (m *DBService) Connect() (err error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s)/information_schema?charset=%s&parseTime=True&loc=Local",
		m.DBConfig.User,
		m.DBConfig.Password,
		m.DBConfig.Host,
		m.DBConfig.Charset,
	)
	m.DB, err = sqlx.Open(m.DBConfig.Type, dsn)
	return
}

// ParseTableName 查询所有表
func (m *DBService) ParseTableName(dbName, tableName string) ([]string, error) {
	if strings.Contains(tableName, ",") {
		return strings.Split(tableName, ","), nil
	}
	if tableName != "*" {
		return []string{tableName}, nil
	}
	// * 遍历数据库下所有表
	query := "SELECT TABLE_NAME FROM TABLES WHERE TABLE_SCHEMA = ? "
	rows, err := m.DB.Query(query, dbName)
	if err != nil {
		return nil, err
	}
	if rows == nil {
		return nil, errors.New("未查询到表")
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var table string
		if err := rows.Scan(&table); err != nil {
			return  nil, err
		}
		tables = append(tables, table)
	}
	return tables, nil
}

func (m *DBService) GetColumns(dbName, tableName string) ([]*TableColumn, error) {
	query := "SELECT COLUMN_NAME, DATA_TYPE, COLUMN_KEY, " +
		"IS_NULLABLE, COLUMN_TYPE, COLUMN_DEFAULT, COLUMN_COMMENT " +
		"FROM COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? "
	rows, err := m.DB.Query(query, dbName, tableName)
	if err != nil {
		return nil, err
	}
	if rows == nil {
		return nil, errors.New("没有数据")
	}
	defer rows.Close()

	var columns []*TableColumn
	for rows.Next() {
		var column TableColumn
		err := rows.Scan(&column.ColumnName, &column.DataType, &column.ColumnKey, &column.IsNullable, &column.ColumnType, &column.ColumnDefault, &column.ColumnComment)
		if err != nil {
			return nil, err
		}

		columns = append(columns, &column)
	}

	return columns, nil
}
