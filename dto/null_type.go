package dto

import (
	"database/sql"
	"encoding/json"
)

type JSONNullUnix struct {
	sql.NullInt64
}

func (j *JSONNullUnix) MarshalJSON() ([]byte, error) {
	if j.Valid {
		return json.Marshal(j.Int64)
	} else {
		return json.Marshal(nil)
	}
}

func (j *JSONNullUnix) UnmarshalJSON(data []byte) error {
	var t *int64
	if err := json.Unmarshal(data, &t); err != nil {
		return err
	}

	if t != nil {
		j.Valid = true
		j.Int64 = *t
	} else {
		j.Valid = false
	}

	return nil
}
