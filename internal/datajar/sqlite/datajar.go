package sqlite

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func FetchStore() (map[string]interface{}, error) {
	const filename = "/Users/josh/Library/Group Containers/group.dk.simonbs.DataJar/Store/DataJar.sqlite"
	_, err := os.Stat(filename)
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite3", fmt.Sprintf("file:%s?mode=ro", filename))
	if err != nil {
		return nil, err
	}
	defer db.Close()

	sqlRootQuery := `
	SELECT Z_PK, ZKEY
	FROM ZSTOREVALUE
	WHERE ZPARENT IS NULL AND ZORPHANDATE IS NULL
	ORDER BY ZORDER
	`
	rows, err := db.Query(sqlRootQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]interface{})
	for rows.Next() {
		var pk int
		var key string
		err = rows.Scan(&pk, &key)
		if err != nil {
			return nil, err
		}

		value, err := getValue(db, pk)
		if err != nil {
			return nil, err
		}
		result[key] = value
	}
	return result, nil
}

func getValue(db *sql.DB, pk int) (interface{}, error) {
	sqlValueTypeQuery := "SELECT ZRAWVALUETYPE FROM ZSTOREVALUE WHERE Z_PK IS ?"
	row := db.QueryRow(sqlValueTypeQuery, pk)

	var rawValueType string
	err := row.Scan(&rawValueType)
	if err != nil {
		return nil, err
	}

	switch rawValueType {
	case "text":
		sqlValueQuery := "SELECT ZSTRINGVALUE FROM ZSTOREVALUE WHERE Z_PK IS ?"
		row := db.QueryRow(sqlValueQuery, pk)
		var stringValue string
		err := row.Scan(&stringValue)
		return stringValue, err

	case "number":
		sqlValueQuery := "SELECT ZDOUBLEVALUE FROM ZSTOREVALUE WHERE Z_PK IS ?"
		row := db.QueryRow(sqlValueQuery, pk)
		var doubleValue float64
		err := row.Scan(&doubleValue)
		return doubleValue, err

	case "boolean":
		sqlValueQuery := "SELECT ZBOOLEANVALUE FROM ZSTOREVALUE WHERE Z_PK IS ?"
		row := db.QueryRow(sqlValueQuery, pk)
		var booleanValue bool
		err := row.Scan(&booleanValue)
		return booleanValue, err

	case "file":
		sqlValueQuery := "SELECT ZFILENAME FROM ZSTOREVALUE WHERE Z_PK IS ?"
		row := db.QueryRow(sqlValueQuery, pk)
		var filenameValue string
		err := row.Scan(&filenameValue)
		return filenameValue, err

	case "array":
		sqlChildrenQuery := "SELECT Z_PK FROM ZSTOREVALUE WHERE ZPARENT IS ? ORDER BY ZORDER"
		rows, err := db.Query(sqlChildrenQuery, pk)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		result := make([]interface{}, 0)
		for rows.Next() {
			var pk int
			err = rows.Scan(&pk)
			if err != nil {
				return nil, err
			}
			value, err := getValue(db, pk)
			if err != nil {
				return nil, err
			}
			result = append(result, value)
		}
		return result, nil

	case "dictionary":
		sqlChildrenQuery := "SELECT Z_PK, ZKEY FROM ZSTOREVALUE WHERE ZPARENT IS ? ORDER BY ZORDER"
		rows, err := db.Query(sqlChildrenQuery, pk)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		result := make(map[string]interface{})
		for rows.Next() {
			var pk int
			var key string
			err = rows.Scan(&pk, &key)
			if err != nil {
				return nil, err
			}
			value, err := getValue(db, pk)
			if err != nil {
				return nil, err
			}
			result[key] = value
		}
		return result, nil

	default:
		return nil, fmt.Errorf("unknown raw value type: %s", rawValueType)
	}
}
