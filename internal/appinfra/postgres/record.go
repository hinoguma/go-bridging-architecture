package postgres

import (
	"database/sql"
	"fmt"
)

type SQLRecordID struct {
	str *string
	int *int64
}

func (value SQLRecordID) Value() any {
	if value.str != nil {
		return *value.str
	}
	if value.int != nil {
		return *value.int
	}
	return nil
}

func (value SQLRecordID) IsString() bool {
	return value.str != nil
}

func (value SQLRecordID) IsInt() bool {
	return value.int != nil
}

func (value SQLRecordID) String() string {
	if value.str != nil {
		return *value.str
	}
	if value.int != nil {
		return fmt.Sprintf("%d", *value.int)
	}
	return ""
}

func (value SQLRecordID) Int() int64 {
	if value.int != nil {
		return *value.int
	}
	return 0
}

func SQLStrID(id string) SQLRecordID {
	return SQLRecordID{
		str: &id,
	}
}

func SQLIntID(id int64) SQLRecordID {
	return SQLRecordID{
		int: &id,
	}
}

type SQLDatabaseRecord interface {
	ToMap() map[string]interface{}
	SetBySQLRow(row *sql.Row) error
}
