package concurrency

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"github.com/segmentio/ksuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

type Version sql.NullString

func NewVersion() Version {
	return Version{
		Valid:  true,
		String: ksuid.New().String(),
	}
}

// Scan implements the Scanner interface.
func (v *Version) Scan(value interface{}) error {
	return (*sql.NullString)(v).Scan(value)
}

// Value implements the driver Valuer interface.
func (v Version) Value() (driver.Value, error) {
	if !v.Valid {
		return nil, nil
	}
	return v.String, nil
}

func (v *Version) UnmarshalJSON(bytes []byte) error {
	if string(bytes) == "null" {
		v.Valid = false
		return nil
	}
	err := json.Unmarshal(bytes, &v.String)
	if err == nil {
		v.Valid = true
	}
	return err
}

func (v *Version) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.String)
	}
	return json.Marshal(nil)
}

func (v *Version) CreateClauses(field *schema.Field) []clause.Interface {
	return []clause.Interface{VersionCreateClause{Field: field}}
}

type VersionCreateClause struct {
	Field *schema.Field
}

func (v VersionCreateClause) Name() string {
	return ""
}

func (v VersionCreateClause) Build(builder clause.Builder) {

}

func (v VersionCreateClause) MergeClause(c *clause.Clause) {

}
func (v VersionCreateClause) ModifyStatement(statement *gorm.Statement) {
	if statement.SQL.Len() == 0 {
		nv := NewVersion()
		statement.AddClause(clause.Set{{Column: clause.Column{Name: v.Field.DBName}, Value: nv.String}})
		statement.SetColumn(v.Field.DBName, nv.String, true)
	}
}

func (v *Version) UpdateClauses(field *schema.Field) []clause.Interface {
	return []clause.Interface{VersionUpdateClause{Field: field}}
}

type VersionUpdateClause struct {
	Field *schema.Field
}

func (v VersionUpdateClause) Name() string {
	return ""
}

func (v VersionUpdateClause) Build(builder clause.Builder) {

}

func (v VersionUpdateClause) MergeClause(c *clause.Clause) {

}
func (v VersionUpdateClause) ModifyStatement(stmt *gorm.Statement) {
	if stmt.SQL.Len() == 0 {

		if _, ok := stmt.Clauses["concurrency_query"]; !ok {
			//build query
			if c, ok := stmt.Clauses["WHERE"]; ok {
				if where, ok := c.Expression.(clause.Where); ok && len(where.Exprs) > 1 {
					for _, expr := range where.Exprs {
						if orCond, ok := expr.(clause.OrConditions); ok && len(orCond.Exprs) == 1 {
							where.Exprs = []clause.Expression{clause.And(where.Exprs...)}
							c.Expression = where
							stmt.Clauses["WHERE"] = c
							break
						}
					}
				}
			}
			if cv, zero := v.Field.ValueOf(stmt.Context, stmt.ReflectValue); !zero {
				if cvv, ok := cv.(Version); ok {
					if cvv.Valid {
						stmt.AddClause(clause.Where{Exprs: []clause.Expression{
							clause.Eq{Column: clause.Column{Table: clause.CurrentTable, Name: v.Field.DBName}, Value: cvv.String},
						}})
					}
				}
			}
			stmt.Clauses["concurrency_query"] = clause.Clause{}
		}

		//set new value
		nv := NewVersion()
		stmt.SetColumn(v.Field.DBName, nv.String, true)

	}

}

// HasVersion embed struct
type HasVersion struct {
	Version Version
}

type VersionInterface interface {
	GetConcurrencyVersion() *string
}

func (v HasVersion) GetConcurrencyVersion() *string {
	if !v.Version.Valid {
		return nil
	}
	r := v.Version.String
	return &r
}
