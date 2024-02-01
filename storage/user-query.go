package storage

import (
	"log"
	"time"
)

// UserAccount представляет модель данных для таблицы UserAccount.
type UserAccount struct {
	ID        int
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// UserAccessData представляет модель данных для таблицы UserAccessData.
type UserAccessData struct {
	ID             int
	IsEmailVerified bool
	RefreshToken   string
	UpdatedAt      time.Time
	UserID         int
}

// CreateUserAccount создает новую запись в таблице UserAccount.
func CreateUserAccount(db *Database, email, password string) (int, error) {
	var userID int

	err := db.QueryRow(`
		INSERT INTO UserAccount (Email, Password) 
		VALUES ($1, $2) 
		RETURNING ID
	`, email, password).Scan(&userID)

	if err != nil {
		log.Print(err)
		return 0, err
	}

	return userID, nil
}

// CreateUserAccessData создает новую запись в таблице UserAccessData.
func CreateUserAccessData(db *Database, isEmailVerified bool, verificationCode string , refreshToken string, userID int) error {
	_, err := db.Exec(`
		INSERT INTO UserAccessData (IsEmailVerified, VerificationCode, RefreshToken, UserID) 
		VALUES ($1, $2, $3, $4)
	`, isEmailVerified,verificationCode, refreshToken, userID)

	if err != nil {
		log.Print(err)
		return err
	}

	return nil
}

// GetUserByEmail получает данные пользователя по его электронной почте
func GetUserByEmail(db *Database, email string) (int, string, string, bool, string, string, error) {
	var ID int
	var password string
	var isEmailVerified bool
	var verificationCode string
	var refreshToken string

	// Выполняем SQL-запрос с использованием INNER JOIN для объединения данных из обеих таблиц
	query := `
		SELECT ua.ID, ua.Email, ua.Password, uad.IsEmailVerified, uad.VerificationCode, uad.RefreshToken
		FROM UserAccount ua
		INNER JOIN UserAccessData uad ON ua.ID = uad.UserID
		WHERE ua.Email = $1
	`

	err := db.QueryRow(query, email).Scan(&ID, &email, &password, &isEmailVerified, &verificationCode, &refreshToken)
	if err != nil {
		log.Fatal(err)
		return 0, "", "", false, "", "", err
	}

	return ID, email, password, isEmailVerified,verificationCode, refreshToken, nil
}


// UpdateAndGetUserAccessData обновляет поля IsEmailVerified и RefreshToken в таблице UserAccessData по UserID и возвращает обновленные данные
func UpdateRefresh(db *Database, userID int, refreshToken string) (UserAccessData, error) {
	// Выполняем SQL-запрос UPDATE
	updateQuery := `
		UPDATE UserAccessData
		SET RefreshToken = $1
		WHERE UserID = $2
		RETURNING ID, IsEmailVerified, RefreshToken, UpdatedAt, UserID
	`

	var updatedData UserAccessData
	err := db.QueryRow(updateQuery, refreshToken, userID).Scan(&updatedData.ID, &updatedData.IsEmailVerified, &updatedData.RefreshToken, &updatedData.UpdatedAt, &updatedData.UserID)
	if err != nil {
		log.Fatal(err)
		return UserAccessData{}, err
	}

	return updatedData, nil
}

func UpdateIsEmailVerified(db *Database, verificationCode string,  IsEmailVerified bool) (UserAccessData, error) {
	// Выполняем SQL-запрос UPDATE
	updateQuery := `
		UPDATE UserAccessData
		SET IsEmailVerified = $1
		WHERE VerificationCode = $2
		RETURNING ID, IsEmailVerified, RefreshToken, UpdatedAt, UserID
	`

	var updatedData UserAccessData
	err := db.QueryRow(updateQuery, IsEmailVerified, verificationCode).Scan(&updatedData.ID, &updatedData.IsEmailVerified, &updatedData.RefreshToken, &updatedData.UpdatedAt, &updatedData.UserID)
	if err != nil {
		log.Fatal(err)
		return UserAccessData{}, err
	}

	return updatedData, nil
}

func GetVerificationCode (db *Database, email string) (string, error) {
	var verificationCode string
	err := db.QueryRow("SELECT VerificationCode FROM UserAccessData WHERE UserID = (SELECT ID FROM UserAccount WHERE Email = $1)", email).Scan(&verificationCode)
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	return verificationCode, nil
}

func GetIDByEmail(db *Database, email string) (int, error) {
	var id int
	err := db.QueryRow("SELECT ID FROM UserAccount WHERE Email = $1", email).Scan(&id)
	if err != nil {
		log.Fatal(err)
		return 0, err
	}
	return id, nil
}

func CreateUserInfo (db *Database, ID int, FullName string, AvatarUrl string) (error) {

	err := db.QueryRow("INSERT INTO user_info (fullname, avatar_url, user_id) VALUES ($1, $2, $3) RETURNING ID", FullName, AvatarUrl, ID)
	if err != nil {
		return err.Err()
	}
	return nil
}


