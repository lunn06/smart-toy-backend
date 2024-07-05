package sql

import (
	"errors"
	"log"
	"time"

	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lunn06/smart-toy-backend/internal/config"
	"github.com/lunn06/smart-toy-backend/internal/models"
)

var (
	DB *sqlx.DB

	insertUserRequest = `
		INSERT INTO users (email, password) VALUES($1, $2) RETURNING id
	`

	insertTokenRequest = `
		INSERT INTO jwt_tokens (uuid, token) VALUES ($1, $2) RETURNING uuid
	`

	insertUserTokenRequest = `
		INSERT INTO users_tokens (user_id, token_uuid) VALUES ($1, $2)
	`

	updateTokenRequest = `
		UPDATE jwt_tokens SET creation_time=$1 WHERE uuid=$2
	`

	getUserByRefreshTokenRequest = `
		SELECT * FROM jwt_tokens WHERE uuid=(
		    SELECT token_uuid FROM users_tokens WHERE user_id=$1
		)
	`

	getTokenByUserRequest = `
		SELECT * FROM users WHERE id=(
		    SELECT user_id FROM users_tokens WHERE token_uuid=$1
		)
	`

	getTokenRequest = `
		SELECT * FROM jwt_tokens WHERE uuid=$1
	`

	deleteTokenRequest = `
		DELETE FROM jwt_tokens WHERE uuid=$1
	`

	getUserRequest = `
		SELECT * FROM users WHERE email=$1
	`
)

func Init() {
	DB = MustCreate(config.CFG)

	// databaseDefaults := config.MustLoadDatabaseDefaults("configs/database_defaults.yaml")

	// tx := DB.MustBegin()
	// for _, role := range databaseDefaults.Roles {
	// 	tx.MustExec(
	// 		`INSERT INTO roles VALUES (
	// 		$1, $2, $3, $4
	// 		) ON CONFLICT (id) DO NOTHING `,
	// 		role.Id, role.Name, role.CanRemoveUsers, role.CanRemoveOthersVideos,
	// 	)
	// }
	// if err := tx.Commit(); err != nil {
	// 	log.Fatalf("Init() error = %v", err)
	// }
}

func MustCreate(cfg config.Config) *sqlx.DB {
	if DB != nil {
		return DB
	}

	dbConf := mysql.Config{
		User: cfg.Database.User,
		Passwd: cfg.Database.Password,
		Addr: cfg.Database.Address,
		DBName: cfg.Database.Name,
	}

	db, err := sqlx.Connect(cfg.Database.Driver, dbConf.FormatDSN())
	if err != nil {
		log.Fatalf("MustCreate() Error: %v", err)
	}

	for _, stmt := range models.Schema {
		db.MustExec(stmt)
	}

	return db
}

func checkDBConnection() error {
	if DB == nil {
		return errors.New("no DB connection")
	}
	return nil
}

func InsertUser(user models.School) (int, error) {
	if err := checkDBConnection(); err != nil {
		return -1, err
	}

	tx := DB.MustBegin()

	var lastInsertIndex int
	err := tx.QueryRow(insertUserRequest, user.Email, user.Password).Scan(&lastInsertIndex)
	if err != nil {
		return -1, err
	}

	if err := tx.Commit(); err != nil {
		return -1, err
	}

	return lastInsertIndex, err
}

func InsertToken(userId int, jwtToken string) (string, error) {
	if err := checkDBConnection(); err != nil {
		return "", err
	}
	tx := DB.MustBegin()

	var tempUUID string
	tokenUUID, err := uuid.NewV7()

	if err != nil {
		return "", err
	}

	err = tx.QueryRow(insertTokenRequest, tokenUUID, jwtToken).Scan(&tempUUID)

	if err != nil {
		return "", err
	}

	_, err = tx.Exec(insertUserTokenRequest, userId, tempUUID)
	if err != nil {
		return "", err
	}

	if err := tx.Commit(); err != nil {
		return "", err
	}

	return tempUUID, err
}

func UpdateTokenTime(token models.JwtToken) error {
	if err := checkDBConnection(); err != nil {
		return err
	}

	_, err := DB.Exec(updateTokenRequest, time.Now(), token.Uuid)
	if err != nil {
		return err
	}

	return nil
}

func GetTokenByUser(user models.School) (models.JwtToken, error) {
	var token models.JwtToken

	if err := checkDBConnection(); err != nil {
		return token, err
	}

	err := DB.Get(&token, getTokenByUserRequest, user.Id)
	if err != nil {
		return token, err
	}

	return token, nil
}

func GetUserByRefreshToken(tokenUUID string) (models.School, error) {
	var user models.School

	if err := checkDBConnection(); err != nil {
		return user, err
	}

	err := DB.Get(&user, getUserByRefreshTokenRequest, tokenUUID)
	if err != nil {
		return user, err
	}

	return user, nil
}

func GetToken(tokenUUID string) (models.JwtToken, error) {
	var token models.JwtToken

	if err := checkDBConnection(); err != nil {
		return token, err
	}

	err := DB.Get(&token, getTokenRequest, tokenUUID)
	if err != nil {
		return token, err
	}

	return token, nil
}

func PopToken(tokenUUID string) (models.JwtToken, error) {
	var token models.JwtToken

	if err := checkDBConnection(); err != nil {
		return token, err
	}

	err := DB.Get(&token, getTokenRequest, tokenUUID)
	if err != nil {
		return token, err
	}

	_, err = DB.Exec(deleteTokenRequest, tokenUUID)
	if err != nil {
		return token, err
	}

	return token, nil
}

func GetUser(email string) (models.School, error) {
	var user models.School

	if err := checkDBConnection(); err != nil {
		return user, err
	}

	err := DB.Get(&user, getUserRequest, email)

	if err != nil {
		return user, err
	}

	return user, nil
}
