package nulls

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

type Bool struct {
	sql.NullBool
}

func NewBool(n interface{}) Bool {
	if n == nil {
		return Bool{
			sql.NullBool{
				Valid: false,
			},
		}
	}
	return Bool{
		sql.NullBool{
			Bool:  n.(bool),
			Valid: true,
		},
	}
}

func (n Bool) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return json.Marshal(n.Bool)
	}
	return json.Marshal(nil)
}

func (n *Bool) UnmarshalJSON(data []byte) error {
	var b *bool
	if err := json.Unmarshal(data, &b); err != nil {
		return err
	}
	if b != nil {
		n.Valid = true
		n.Bool = *b
	} else {
		n.Valid = false
	}
	return nil
}

type Int64 struct {
	sql.NullInt64
}

func NewInt64(n interface{}) Int64 {
	if n == nil {
		return Int64{
			sql.NullInt64{
				Valid: false,
			},
		}
	}
	return Int64{
		sql.NullInt64{
			Int64: n.(int64),
			Valid: true,
		},
	}
}

func (n *Int64) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return json.Marshal(n.Int64)
	}
	return json.Marshal(nil)
}

func (n *Int64) UnmarshalJSON(data []byte) error {
	var i *int64
	if err := json.Unmarshal(data, &i); err != nil {
		return err
	}
	if i != nil {
		n.Valid = true
		n.Int64 = *i
	} else {
		n.Valid = false
	}
	return nil
}

type Float64 struct {
	sql.NullFloat64
}

func NewFloat64(n interface{}) Float64 {
	if n == nil {
		return Float64{
			sql.NullFloat64{
				Valid: false,
			},
		}
	}
	return Float64{
		sql.NullFloat64{
			Float64: n.(float64),
			Valid:   true,
		},
	}
}

func (n *Float64) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return json.Marshal(n.Float64)
	}
	return json.Marshal(nil)
}

func (n *Float64) UnmarshalJSON(data []byte) error {
	var f *float64
	if err := json.Unmarshal(data, &f); err != nil {
		return err
	}
	if f != nil {
		n.Valid = true
		n.Float64 = *f
	} else {
		n.Valid = false
	}
	return nil
}

type String struct {
	sql.NullString
}

func NewString(n interface{}) String {
	if n == nil {
		return String{
			sql.NullString{
				Valid: false,
			},
		}
	}
	return String{
		sql.NullString{
			String: n.(string),
			Valid:  true,
		},
	}
}

func (n String) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return json.Marshal(n.String)
	}
	return json.Marshal(nil)
}

func (n *String) UnmarshalJSON(data []byte) error {
	var s *string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	if s != nil {
		n.Valid = true
		n.String = *s
	} else {
		n.Valid = false
	}
	return nil
}

type Time struct {
	Time  time.Time
	Valid bool
}

func (n *Time) Scan(value interface{}) error {
	if value == nil {
		n.Time = time.Time{}
		n.Valid = false
		return nil
	}
	timeValue, ok := value.(time.Time)
	if !ok {
		return fmt.Errorf("null: cannot scan type %T into null.Time: %v", value, value)
	}
	n.Time = timeValue
	n.Valid = true
	return nil
}

func (n Time) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Time, nil
}

func NewTime(n interface{}) Time {
	if n == nil {
		return Time{
			Valid: false,
		}
	}
	return Time{
		Time:  n.(time.Time),
		Valid: true,
	}
}

func (n Time) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return json.Marshal(n.Time)
	}
	return json.Marshal(nil)
}

func (n *Time) UnmarshalJSON(data []byte) error {
	var t *time.Time
	if err := json.Unmarshal(data, &t); err != nil {
		return err
	}
	if t != nil {
		n.Valid = true
		n.Time = *t
	} else {
		n.Valid = false
	}
	return nil
}
