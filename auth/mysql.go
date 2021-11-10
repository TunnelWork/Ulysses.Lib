package auth

import (
	"database/sql"
	"strings"
)

var (
	tblPrefix = "ulysses_"
	db        *sql.DB
)

/************ Helper Functions ************/
func sqlStatement(query string) (*sql.Stmt, error) {
	prefixUpdatedQuery := strings.Replace(query, "dbprefix_", tblPrefix, -1)

	return db.Prepare(prefixUpdatedQuery)
}

func initDatabaseTable(db *sql.DB) error {
	stmtCreateTmpTableIfNotExists, err := sqlStatement(`CREATE TABLE IF NOT EXISTS dbprefix_tmp_auth (
        id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
        userID BIGINT UNSIGNED NOT NULL,
        extentionType VARCHAR(32) NOT NULL,
        indexKey VARCHAR(32) NOT NULL,
        storedValue TEXT NOT NULL,
        expiry DATETIME NOT NULL,
        PRIMARY KEY (id),
        UNIQUE KEY (userID, extentionType, indexKey)
    )`)
	if err != nil {
		return err
	}
	defer stmtCreateTmpTableIfNotExists.Close()

	_, err = stmtCreateTmpTableIfNotExists.Exec()
	if err != nil {
		return err
	}

	stmtClearTmpTable, err := sqlStatement(`DELETE FROM dbprefix_tmp_auth WHERE expiry < NOW();`)
	if err != nil {
		return err
	}
	defer stmtClearTmpTable.Close()

	_, err = stmtClearTmpTable.Exec()
	if err != nil {
		return err
	}

	stmtCreateTableIfNotExists, err := sqlStatement(`CREATE TABLE IF NOT EXISTS dbprefix_auth_mfa (
		id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
		userID BIGINT UNSIGNED NOT NULL,
		extentionType VARCHAR(32) NOT NULL,
		extentionData TEXT NOT NULL,
		enabled BOOLEAN NOT NULL DEFAULT FALSE,
		PRIMARY KEY (id),
		UNIQUE KEY (userID, extentionType)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;`)
	if err != nil {
		return err
	}
	defer stmtCreateTableIfNotExists.Close()

	_, err = stmtCreateTableIfNotExists.Exec()
	return err
}

// Create
func InitMFA(userID uint64, extentionType string, extentionData string) error {
	stmtInsertExtention, err := sqlStatement(`INSERT INTO dbprefix_auth_mfa (userID, extentionType, extentionData) VALUES (?, ?, ?);`)
	if err != nil {
		return err
	}
	defer stmtInsertExtention.Close()

	_, err = stmtInsertExtention.Exec(userID, extentionType, extentionData)
	return err
}

// Read
func CheckoutMFA(userID uint64, extentionType string) (string, error) {
	stmtCheckoutExtentionData, err := sqlStatement(`SELECT extentionData FROM dbprefix_auth_mfa WHERE userID = ? AND extentionType = ?;`)
	if err != nil {
		return "", err
	}
	defer stmtCheckoutExtentionData.Close()

	var extentionData string
	err = stmtCheckoutExtentionData.QueryRow(userID, extentionType).Scan(&extentionData)
	if err != nil {
		return "", err
	}

	return extentionData, nil
}

// Read
func MFAEnabled(userID uint64, extentionType string) (bool, error) {
	stmtCheckIfEnabled, err := sqlStatement(`SELECT enabled FROM dbprefix_auth_mfa WHERE userID = ? AND extentionType = ?;`)
	if err != nil {
		return false, err
	}
	defer stmtCheckIfEnabled.Close()

	var enabled bool
	err = stmtCheckIfEnabled.QueryRow(userID, extentionType).Scan(&enabled)
	if err != nil {
		return false, err
	}

	return enabled, nil
}

// Update
func ConfirmMFA(userID uint64, extentionType string) error {
	stmtConfirmExtention, err := sqlStatement(`UPDATE dbprefix_auth_mfa SET enabled = TRUE WHERE userID = ? AND extentionType = ?;`)
	if err != nil {
		return err
	}
	defer stmtConfirmExtention.Close()

	_, err = stmtConfirmExtention.Exec(userID, extentionType)
	return err
}

// Update
func UpdateMFA(userID uint64, extentionType string, extentionData string) error {
	stmtUpdateExtention, err := sqlStatement(`UPDATE dbprefix_auth_mfa SET extentionData = ? WHERE userID = ? AND extentionType = ?;`)
	if err != nil {
		return err
	}
	defer stmtUpdateExtention.Close()

	_, err = stmtUpdateExtention.Exec(extentionData, userID, extentionType)
	return err
}

// Delete
func ClearMFA(userID uint64, extentionType string) error {
	stmtClearExtention, err := sqlStatement(`DELETE FROM dbprefix_auth_mfa WHERE userID = ? AND extentionType = ?;`)
	if err != nil {
		return err
	}
	defer stmtClearExtention.Close()

	_, err = stmtClearExtention.Exec(userID, extentionType)
	return err
}

/**** Temporary Data Storage ****/
// Create
func InsertTmpEntry(userID uint64, extentionType string, indexKey string, storedValue string) error {
	stmtInsertTmpEntry, err := sqlStatement(`INSERT INTO dbprefix_tmp_auth 
	(userID, extentionType, indexKey, storedValue, expiry) 
    VALUES (?, ?, ?, ?, NOW() + INTERVAL 1 DAY );`)
	if err != nil {
		return err
	}
	defer stmtInsertTmpEntry.Close()

	_, err = stmtInsertTmpEntry.Exec(userID, extentionType, indexKey, storedValue)
	return err
}

// Read
func ReadTmpEntry(userID uint64, extentionType string, indexKey string) (string, error) {
	stmtReadTmpEntry, err := sqlStatement(`SELECT storedValue FROM dbprefix_tmp_auth WHERE userID = ? AND extentionType = ? AND indexKey = ?;`)
	if err != nil {
		return "", err
	}
	defer stmtReadTmpEntry.Close()

	var storedValue string
	err = stmtReadTmpEntry.QueryRow(userID, extentionType, indexKey).Scan(&storedValue)
	if err != nil {
		return "", err
	}

	return storedValue, nil
}

// Update
func UpdateTmpEntry(userID uint64, extentionType string, indexKey string, storedValue string) error {
	stmtUpdateTmpEntry, err := sqlStatement(`UPDATE dbprefix_tmp_auth SET storedValue = ? WHERE userID = ? AND extentionType = ? AND indexKey = ?;`)
	if err != nil {
		return err
	}
	defer stmtUpdateTmpEntry.Close()

	_, err = stmtUpdateTmpEntry.Exec(storedValue, userID, extentionType, indexKey)
	return err
}

// Delete
func DeleteTmpEntry(userID uint64, extentionType string, indexKey string) error {
	stmtDeleteTmpEntry, err := sqlStatement(`DELETE FROM dbprefix_tmp_auth WHERE userID = ? AND extentionType = ? AND indexKey = ?;`)
	if err != nil {
		return err
	}
	defer stmtDeleteTmpEntry.Close()

	_, err = stmtDeleteTmpEntry.Exec(userID, extentionType, indexKey)
	return err
}

/*** Internal ***/
func checkEnabledMFA(userID uint64) ([]string, error) {
	stmtCheckEnabledMFA, err := sqlStatement(`SELECT extentionType FROM dbprefix_auth_mfa WHERE userID = ? AND enabled = TRUE;`)
	if err != nil {
		return nil, err
	}
	defer stmtCheckEnabledMFA.Close()

	rows, err := stmtCheckEnabledMFA.Query(userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var extentionTypes []string
	for rows.Next() {
		var extentionType string
		err = rows.Scan(&extentionType)
		if err != nil {
			return nil, err
		}
		extentionTypes = append(extentionTypes, extentionType)
	}

	return extentionTypes, nil
}
