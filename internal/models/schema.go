package models

var Schema = []string{
`CREATE TABLE IF NOT EXISTS users (
	id INT AUTO_INCREMENT PRIMARY KEY,
	email VARCHAR(100) UNIQUE NOT NULL,
	password VARCHAR(100) NOT NULL,
	registration_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);`,
}
