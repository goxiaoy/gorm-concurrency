package concurrency

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

type Version sql.NullString

func NewVersion() Version {
	return Version{
		Valid:  true,
		String: uuid.New().String(),
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
		// create new value if empty
		if cv, zero := v.Field.ValueOf(statement.ReflectValue); !zero {
			if cvv, ok := cv.(Version); ok {
				if cvv.Valid {
					return
				}
			}
		}
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
func (v VersionUpdateClause) ModifyStatement(statement *gorm.Statement) {
	if statement.SQL.Len() == 0 {
		if _, ok := statement.Clauses["concurrency_query"]; !ok {
			//build query
			if c, ok := statement.Clauses["WHERE"]; ok {
				if where, ok := c.Expression.(clause.Where); ok && len(where.Exprs) > 1 {
					for _, expr := range where.Exprs {
						if orCond, ok := expr.(clause.OrConditions); ok && len(orCond.Exprs) == 1 {
							where.Exprs = []clause.Expression{clause.And(where.Exprs...)}
							c.Expression = where
							statement.Clauses["WHERE"] = c
							break
						}
					}
				}
			}
			if cv, zero := v.Field.ValueOf(statement.ReflectValue); !zero {
				if cvv, ok := cv.(Version); ok {
					if cvv.Valid {
						statement.AddClause(clause.Where{Exprs: []clause.Expression{
							clause.Eq{Column: clause.Column{Table: clause.CurrentTable, Name: v.Field.DBName}, Value: cvv.String},
						}})
					}
				}
			}
			statement.Clauses["concurrency_query"] = clause.Clause{}
		}

		//set new value

		nv := NewVersion()
		if c, ok := statement.Clauses["SET"]; !ok {
			statement.AddClause(clause.Set{{Column: clause.Column{Name: v.Field.DBName}, Value: nv.String}})
		} else {
			if set, ok := c.Expression.(clause.Set); ok {
				a:= make([]clause.Assignment,len(set)+1)
				copy(a,set)
				a[len(a)-1] = clause.Assignment{Column: clause.Column{Name: v.Field.DBName}, Value: nv.String}
				c.Expression = clause.Set(a)
				statement.Clauses["SET"] = c
			}
		}
		v.Field.Set(statement.ReflectValue,nv.String)
	}

}
