package snowflake

import (
	"context"
	_ "embed"
	"fmt"
	"strings"

	"github.com/centralmind/gateway/castx"
	"github.com/centralmind/gateway/connectors"
	"github.com/centralmind/gateway/model"
	"github.com/jmoiron/sqlx"
	_ "github.com/snowflakedb/gosnowflake"
	"golang.org/x/xerrors"
)

//go:embed readme.md
var docString string

func init() {
	connectors.Register(func(cfg Config) (connectors.Connector, error) {
		dsn, err := cfg.MakeDSN()
		if err != nil {
			return nil, xerrors.Errorf("unable to prepare Snowflake config: %w", err)
		}
		db, err := sqlx.Open("snowflake", dsn)
		if err != nil {
			return nil, xerrors.Errorf("unable to open Snowflake db: %w", err)
		}
		return &Connector{
			config: cfg,
			db:     db,
			base:   &connectors.BaseConnector{DB: db},
		}, nil
	})
}

type Config struct {
	Account   string
	Database  string
	User      string
	Password  string
	Warehouse string
	Schema    string
	Role      string
}

func (c Config) MakeDSN() (string, error) {

	dsn := fmt.Sprintf("%s:%s@%s/%s/%s?warehouse=%s&role=%s", c.User, c.Password, c.Account, c.Database, c.Schema, c.Warehouse, c.Role)

	return dsn, nil
}

func (c Config) Type() string {
	return "snowflake"
}

func (c Config) Doc() string {
	return docString
}

type Connector struct {
	config Config
	db     *sqlx.DB
	base   *connectors.BaseConnector
}

func (c Connector) Sample(ctx context.Context, table model.Table) ([]map[string]any, error) {
	rows, err := c.db.QueryxContext(ctx, fmt.Sprintf("SELECT * FROM %s LIMIT 5", table.Name))
	if err != nil {
		return nil, xerrors.Errorf("unable to query db: %w", err)
	}
	defer rows.Close()

	res := make([]map[string]any, 0, 5)
	for rows.Next() {
		row := make(map[string]any)
		if err := rows.MapScan(row); err != nil {
			return nil, xerrors.Errorf("unable to scan row: %w", err)
		}
		res = append(res, row)
	}
	return res, nil
}

func (c Connector) Discovery(ctx context.Context) ([]model.Table, error) {
	rows, err := c.db.Query("SELECT table_name FROM information_schema.tables WHERE table_schema = CURRENT_SCHEMA()")

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables []model.Table
	for rows.Next() {
		var tableName string

		fmt.Println("TableName")
		fmt.Println(tableName)

		if err := rows.Scan(&tableName); err != nil {
			return nil, err
		}
		columns, err := c.LoadsColumns(ctx, tableName)
		if err != nil {
			return nil, err
		}
		tables = append(tables, model.Table{Name: tableName, Columns: columns})
	}
	return tables, nil
}

func (c Connector) Ping(ctx context.Context) error {
	return c.db.PingContext(ctx)
}

func (c Connector) Query(ctx context.Context, endpoint model.Endpoint, params map[string]any) ([]map[string]any, error) {
	processed, err := castx.ParamsE(endpoint, params)
	if err != nil {
		return nil, xerrors.Errorf("unable to process params: %w", err)
	}

	rows, err := c.db.NamedQueryContext(ctx, endpoint.Query, processed)
	if err != nil {
		return nil, xerrors.Errorf("unable to query db: %w", err)
	}
	defer rows.Close()

	res := make([]map[string]any, 0)
	for rows.Next() {
		row := map[string]any{}
		if err := rows.MapScan(row); err != nil {
			return nil, xerrors.Errorf("unable to scan row: %w", err)
		}
		res = append(res, row)
	}
	return res, nil
}

func (c Connector) LoadsColumns(ctx context.Context, tableName string) ([]model.ColumnSchema, error) {
	rows, err := c.db.QueryContext(
		ctx,
		"SELECT COLUMN_NAME, DATA_TYPE FROM information_schema.columns WHERE table_name = ?",
		tableName,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var columns []model.ColumnSchema
	for rows.Next() {
		var name, dataType string
		if err := rows.Scan(&name, &dataType); err != nil {
			return nil, err
		}
		columns = append(columns, model.ColumnSchema{
			Name: name,
			Type: c.GuessColumnType(dataType),
		})
	}
	return columns, nil
}

// GuessColumnType implements TypeGuesser interface for Snowflake
func (c *Connector) GuessColumnType(sqlType string) model.ColumnType {
	upperType := strings.ToUpper(sqlType)

	// Array type
	if strings.Contains(upperType, "ARRAY") {
		return model.TypeArray
	}

	// Object types
	switch upperType {
	case "OBJECT", "VARIANT":
		return model.TypeObject
	}

	// String types
	switch upperType {
	case "STRING", "TEXT", "VARCHAR", "CHAR", "BINARY", "VARBINARY":
		return model.TypeString
	}

	// Numeric types
	switch upperType {
	case "NUMBER", "DECIMAL", "NUMERIC", "FLOAT", "DOUBLE":
		return model.TypeNumber
	}

	// Integer types
	switch upperType {
	case "INT", "INTEGER", "BIGINT", "SMALLINT", "TINYINT":
		return model.TypeInteger
	}

	// Boolean type
	if upperType == "BOOLEAN" {
		return model.TypeBoolean
	}

	// Date/Time types
	switch upperType {
	case "DATE", "TIME", "TIMESTAMP", "TIMESTAMP_LTZ", "TIMESTAMP_NTZ", "TIMESTAMP_TZ":
		return model.TypeDatetime
	}

	// Default to string for unknown types
	return model.TypeString
}

// InferResultColumns returns column information for the given query
func (c *Connector) InferResultColumns(ctx context.Context, query string) ([]model.ColumnSchema, error) {
	return c.base.InferResultColumns(ctx, query, c)
}

// InferQuery implements the Connector interface
func (c *Connector) InferQuery(ctx context.Context, query string) ([]model.ColumnSchema, error) {
	return c.base.InferResultColumns(ctx, query, c)
}
