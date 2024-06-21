package db

import (
	"database/sql"
	"errors"
	"fmt"
)

// HandleErro function handles database errors
func HandleError(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("object not found")
	}
	return err
}
