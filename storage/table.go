package storage

import (
	"fmt"
	"log"
)


func createTables(db *Database){
	// Создание таблицы UserAccount
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS UserAccount (
			ID SERIAL PRIMARY KEY,
			Email VARCHAR(255) UNIQUE NOT NULL,
			Password VARCHAR(255) NOT NULL,
			CreatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			UpdatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		log.Fatal(err)
	}

	// Создание таблицы UserAccessData с внешним ключом, указывающим на UserAccount
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS UserAccessData (
			ID SERIAL PRIMARY KEY,
			IsEmailVerified BOOLEAN NOT NULL,
			VerificationCode VARCHAR(255) NOT NULL,
			RefreshToken VARCHAR(255) NOT NULL,
			UpdatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			UserID INT UNIQUE,
			FOREIGN KEY (UserID) REFERENCES UserAccount(ID)
		)
	`)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Таблицы и индексы успешно созданы.")
}