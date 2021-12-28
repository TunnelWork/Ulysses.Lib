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
	prefixUpdatedQuery := strings.ReplaceAll(query, "dbprefix_", tblPrefix)

	return db.Prepare(prefixUpdatedQuery)
}

func initDatabaseTable() error {
	stmtCreateUserTableIfNotExists, err := sqlStatement(`CREATE TABLE IF NOT EXISTS dbprefix_auth_user (
        id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
        email VARCHAR(128) NOT NULL,
        publickey VARCHAR(64) NOT NULL,
        role INT UNSIGNED NOT NULL DEFAULT 0,
        affiliation BIGINT UNSIGNED NOT NULL DEFAULT 0,
        PRIMARY KEY (id),
        UNIQUE KEY (email)
    )`)
	if err != nil {
		panic(err.Error())
	}
	defer stmtCreateUserTableIfNotExists.Close()

	_, err = stmtCreateUserTableIfNotExists.Exec()
	if err != nil {
		panic(err.Error())
	}

	stmtCreateUserInfoTableIfNotExists, err := sqlStatement(`CREATE TABLE IF NOT EXISTS dbprefix_auth_user_info (
        id BIGINT UNSIGNED NOT NULL,
        first_name VARCHAR(64) NOT NULL,
        last_name VARCHAR(64) NOT NULL,
        street_address VARCHAR(128) NOT NULL,
        suite VARCHAR(64) NOT NULL,
        city VARCHAR(64) NOT NULL,
        state VARCHAR(64) NOT NULL,
        country_iso VARCHAR(8) NOT NULL,
        zip_code VARCHAR(16) NOT NULL,
        PRIMARY KEY (id),
        CONSTRAINT FOREIGN KEY (id) REFERENCES dbprefix_auth_user(id) ON DELETE CASCADE
    )`)
	if err != nil {
		panic(err.Error())
	}
	defer stmtCreateUserInfoTableIfNotExists.Close()

	_, err = stmtCreateUserInfoTableIfNotExists.Exec()
	if err != nil {
		panic(err.Error())
	}

	stmtCreateAffiliationTableIfNotExists, err := sqlStatement(`CREATE TABLE IF NOT EXISTS dbprefix_auth_affiliation (
        id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
        name VARCHAR(64) NOT NULL,
        parent_id BIGINT UNSIGNED NOT NULL DEFAULT 0,
        owner_user_id BIGINT UNSIGNED NOT NULL,
        shared_wallet_id BIGINT UNSIGNED NOT NULL,
        street_address VARCHAR(128) NOT NULL,
        suite VARCHAR(64) NOT NULL,
        city VARCHAR(64) NOT NULL,
        state VARCHAR(64) NOT NULL,
        country_iso VARCHAR(8) NOT NULL,
        zip_code VARCHAR(16) NOT NULL,
        contact_email VARCHAR(128) NOT NULL,
        PRIMARY KEY (id),
        UNIQUE KEY (name)
    )`)
	if err != nil {
		panic(err.Error())
	}
	defer stmtCreateAffiliationTableIfNotExists.Close()

	_, err = stmtCreateAffiliationTableIfNotExists.Exec()
	if err != nil {
		panic(err.Error())
	}

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
		panic(err.Error())
	}
	defer stmtCreateTmpTableIfNotExists.Close()

	_, err = stmtCreateTmpTableIfNotExists.Exec()
	if err != nil {
		panic(err.Error())
	}

	stmtCreateMfaTableIfNotExists, err := sqlStatement(`CREATE TABLE IF NOT EXISTS dbprefix_auth_mfa (
        id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
        userID BIGINT UNSIGNED NOT NULL,
        extentionType VARCHAR(32) NOT NULL,
        extentionData TEXT NOT NULL,
        enabled BOOLEAN NOT NULL DEFAULT FALSE,
        PRIMARY KEY (id),
        UNIQUE KEY (userID, extentionType)
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;`)
	if err != nil {
		panic(err.Error())
	}
	defer stmtCreateMfaTableIfNotExists.Close()

	_, err = stmtCreateMfaTableIfNotExists.Exec()
	if err != nil {
		panic(err.Error())
	}
	return nil
}

/************ User Database ************/

func newUser(user *User) error {
	stmtInsertUser, err := sqlStatement(`INSERT INTO dbprefix_auth_user (email, publickey, role, affiliation) VALUES (?, ?, ?, ?);`)
	if err != nil {
		return err
	}
	defer stmtInsertUser.Close()

	result, err := stmtInsertUser.Exec(user.Email, user.PublicKey, user.Role, user.AffiliationID)
	if err != nil {
		return err
	}
	userid, err := result.LastInsertId()
	user.id = uint64(userid)
	return err
}

func getUserByID(userID uint64) (*User, error) {
	stmtGetUserByID, err := sqlStatement(`SELECT id, email, publickey, role, affiliation FROM dbprefix_auth_user WHERE id = ?;`)
	if err != nil {
		return nil, err
	}
	defer stmtGetUserByID.Close()

	var user User
	err = stmtGetUserByID.QueryRow(userID).Scan(&user.id, &user.Email, &user.PublicKey, &user.Role, &user.AffiliationID)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func getUserByEmail(email string) (*User, error) {
	stmtGetUserByEmailPassword, err := sqlStatement(`SELECT id, email, publickey, role, affiliation FROM dbprefix_auth_user WHERE email = ?;`)
	if err != nil {
		return nil, err
	}
	defer stmtGetUserByEmailPassword.Close()

	var user User
	err = stmtGetUserByEmailPassword.QueryRow(email).Scan(&user.id, &user.Email, &user.PublicKey, &user.Role, &user.AffiliationID)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func getUsersByAffiliationID(affiliationID uint64) ([]*User, error) {
	stmtGetUsersByAffiliationID, err := sqlStatement(`SELECT id, email, publickey, role, affiliation FROM dbprefix_auth_user WHERE affiliation = ?;`)
	if err != nil {
		return nil, err
	}
	defer stmtGetUsersByAffiliationID.Close()

	rows, err := stmtGetUsersByAffiliationID.Query(affiliationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		var user User
		err = rows.Scan(&user.id, &user.Email, &user.PublicKey, &user.Role, &user.AffiliationID)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func listUserID() ([]uint64, error) {
	stmtListUserID, err := sqlStatement(`SELECT id FROM dbprefix_auth_user;`)
	if err != nil {
		return nil, err
	}
	defer stmtListUserID.Close()

	rows, err := stmtListUserID.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userIDs []uint64
	for rows.Next() {
		var userID uint64
		err = rows.Scan(&userID)
		if err != nil {
			return nil, err
		}
		userIDs = append(userIDs, userID)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return userIDs, nil
}

func listUserIDByAffiliationID(affiliationID uint64) ([]uint64, error) {
	stmtListUserIDByAffiliationID, err := sqlStatement(`SELECT id FROM dbprefix_auth_user WHERE affiliation = ?;`)
	if err != nil {
		return nil, err
	}
	defer stmtListUserIDByAffiliationID.Close()

	rows, err := stmtListUserIDByAffiliationID.Query(affiliationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userIDs []uint64
	for rows.Next() {
		var userID uint64
		err = rows.Scan(&userID)
		if err != nil {
			return nil, err
		}
		userIDs = append(userIDs, userID)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return userIDs, nil
}

func emailExists(email string) (bool, error) {
	stmtCheckEmailExists, err := sqlStatement(`SELECT id FROM dbprefix_auth_user WHERE email = ?;`)
	if err != nil {
		return false, err
	}
	defer stmtCheckEmailExists.Close()

	var id uint64
	err = stmtCheckEmailExists.QueryRow(email).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		} else {
			return false, err
		}
	}

	return true, nil
}

func updateUser(user *User) error {
	stmtUpdateUser, err := sqlStatement(`UPDATE dbprefix_auth_user SET email = ?, publickey = ?, role = ?, affiliation = ? WHERE id = ?;`)
	if err != nil {
		return err
	}
	defer stmtUpdateUser.Close()

	_, err = stmtUpdateUser.Exec(user.Email, user.PublicKey, user.Role, user.AffiliationID, user.id)
	return err
}

func wipeUserData(user *User) error {
	stmtWipeUserData, err := sqlStatement(`DELETE FROM dbprefix_auth_user WHERE id = ?;`)
	if err != nil {
		return err
	}
	defer stmtWipeUserData.Close()

	_, err = stmtWipeUserData.Exec(user.id)
	return err
}

/************ UserInfo Database ************/

func newUserInfo(user *User, info *UserInfo) error {
	stmtInsertUserInfo, err := sqlStatement(`INSERT INTO dbprefix_auth_user_info (id, first_name, last_name, street_address, suite, city, state, country_iso, zip_code) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);`)
	if err != nil {
		return err
	}
	defer stmtInsertUserInfo.Close()

	_, err = stmtInsertUserInfo.Exec(user.id, info.FirstName, info.LastName, info.StreetAddress, info.Suite, info.City, info.State, info.CountryISO, info.ZipCode)
	return err
}

func getUserInfo(userID uint64) (*UserInfo, error) {
	stmtGetUserInfo, err := sqlStatement(`SELECT first_name, last_name, street_address, suite, city, state, country_iso, zip_code FROM dbprefix_auth_user_info WHERE id = ?;`)
	if err != nil {
		return nil, err
	}
	defer stmtGetUserInfo.Close()

	var info UserInfo
	err = stmtGetUserInfo.QueryRow(userID).Scan(&info.FirstName, &info.LastName, &info.StreetAddress, &info.Suite, &info.City, &info.State, &info.CountryISO, &info.ZipCode)
	if err != nil {
		return nil, err
	}

	return &info, nil
}

func updateUserInfo(user *User, info *UserInfo) error {
	stmtUpdateUserInfo, err := sqlStatement(`UPDATE dbprefix_auth_user_info SET first_name = ?, last_name = ?, street_address = ?, suite = ?, city = ?, state = ?, country_iso = ?, zip_code = ? WHERE id = ?;`)
	if err != nil {
		return err
	}
	defer stmtUpdateUserInfo.Close()

	_, err = stmtUpdateUserInfo.Exec(info.FirstName, info.LastName, info.StreetAddress, info.Suite, info.City, info.State, info.CountryISO, info.ZipCode, user.id)
	return err
}

/************ Affiliation Database ************/

func getAffiliationByID(affiliationID uint64) (*Affiliation, error) {
	stmtGetAffiliationByID, err := sqlStatement(`SELECT 
        id, 
        name, 
        parent_id,
        owner_user_id,
        shared_wallet_id,
        street_address,
        suite,
        city,
        state,
        country_iso,
        zip_code,
        contact_email,
    FROM dbprefix_auth_affiliation WHERE id = ?;`)
	if err != nil {
		return nil, err
	}
	defer stmtGetAffiliationByID.Close()

	var affiliation Affiliation
	err = stmtGetAffiliationByID.QueryRow(affiliationID).Scan(&affiliation.id, &affiliation.Name, &affiliation.ParentID, &affiliation.OwnerUserID, &affiliation.SharedWalletID, &affiliation.StreetAddress, &affiliation.Suite, &affiliation.City, &affiliation.State, &affiliation.CountryISO, &affiliation.ZipCode, &affiliation.ContactEmail)
	if err != nil {
		return nil, err
	}

	return &affiliation, nil
}

func newAffiliation(affiliation *Affiliation) error {
	stmtInsertAffiliation, err := sqlStatement(`INSERT INTO dbprefix_auth_affiliation (
        name,
        parent_id,
        owner_user_id,
        shared_wallet_id,
        street_address,
        suite,
        city,
        state,
        country_iso,
        zip_code,
        contact_email,
    ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`)
	if err != nil {
		return err
	}
	defer stmtInsertAffiliation.Close()

	_, err = stmtInsertAffiliation.Exec(affiliation.Name, affiliation.OwnerUserID, affiliation.SharedWalletID, affiliation.StreetAddress, affiliation.Suite, affiliation.City, affiliation.State, affiliation.CountryISO, affiliation.ZipCode, affiliation.ContactEmail)
	return err
}

func updateAffiliation(affiliation *Affiliation) error {
	stmtUpdateAffiliation, err := sqlStatement(`UPDATE dbprefix_auth_affiliation SET
        name = ?,
        owner_user_id = ?,
        shared_wallet_id = ?,
        street_address = ?,
        suite = ?,
        city = ?,
        state = ?,
        country_iso = ?,
        zip_code = ?,
        contact_email = ?
    WHERE id = ?;`)
	if err != nil {
		return err
	}
	defer stmtUpdateAffiliation.Close()

	_, err = stmtUpdateAffiliation.Exec(affiliation.Name, affiliation.OwnerUserID, affiliation.SharedWalletID, affiliation.StreetAddress, affiliation.Suite, affiliation.City, affiliation.State, affiliation.CountryISO, affiliation.ZipCode, affiliation.ContactEmail, affiliation.id)
	return err
}

/************ MFA Database ************/

// Create
func InitMFA(userID uint64, extentionType, extentionData string) error {
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
func UpdateMFA(userID uint64, extentionType, extentionData string) error {
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

/************ Temporary Database ************/
// Create
func InsertTmpEntry(userID uint64, extentionType, indexKey, storedValue string) error {
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
func ReadTmpEntry(userID uint64, extentionType, indexKey string) (string, error) {
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
func UpdateTmpEntry(userID uint64, extentionType, indexKey, storedValue string) error {
	stmtUpdateTmpEntry, err := sqlStatement(`UPDATE dbprefix_tmp_auth SET storedValue = ? WHERE userID = ? AND extentionType = ? AND indexKey = ?;`)
	if err != nil {
		return err
	}
	defer stmtUpdateTmpEntry.Close()

	_, err = stmtUpdateTmpEntry.Exec(storedValue, userID, extentionType, indexKey)
	return err
}

// Delete
func DeleteTmpEntry(userID uint64, extentionType, indexKey string) error {
	stmtDeleteTmpEntry, err := sqlStatement(`DELETE FROM dbprefix_tmp_auth WHERE userID = ? AND extentionType = ? AND indexKey = ?;`)
	if err != nil {
		return err
	}
	defer stmtDeleteTmpEntry.Close()

	_, err = stmtDeleteTmpEntry.Exec(userID, extentionType, indexKey)
	return err
}

func PurgeExpiredTmpEntry() error {
	stmtClearTmpTable, err := sqlStatement(`DELETE FROM dbprefix_tmp_auth WHERE expiry < NOW();`)
	if err != nil {
		return err
	}
	defer stmtClearTmpTable.Close()

	_, err = stmtClearTmpTable.Exec()
	return err
}

/************ Internal ************/
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
