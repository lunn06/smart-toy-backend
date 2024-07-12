package sql

import (
	"errors"
	"log"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lunn06/smart-toy-backend/internal/config"
	"github.com/lunn06/smart-toy-backend/internal/models"
)

var (
	DB *sqlx.DB

	insertUserRequest = `
		INSERT INTO users (email, password) VALUES(?, ?)
	`

	insertTokenRequest = `
		INSERT INTO jwt_tokens (uuid, token) VALUES (?, ?)
	`

	insertUserTokenRequest = `
		INSERT INTO users_tokens (user_id, token_uuid) VALUES (?, ?)
	`

	updateTokenRequest = `
		UPDATE jwt_tokens SET creation_time=? WHERE uuid=?
	`

	getUserByRefreshTokenRequest = `
		SELECT * FROM jwt_tokens WHERE uuid=(
		    SELECT token_uuid FROM users_tokens WHERE user_id=?
		)
	`

	getTokenByUserRequest = `
		SELECT * FROM users WHERE id=(
		    SELECT user_id FROM users_tokens WHERE token_uuid=?
		)
	`

	getTokenRequest = `
		SELECT * FROM jwt_tokens WHERE uuid=?
	`

	deleteTokenRequest = `
		DELETE FROM jwt_tokens WHERE uuid=?
	`

	getUserByEmailRequest = `
		SELECT * FROM users WHERE email=?
	`

	getUserByIdRequest = `
		SELECT * FROM users WHERE id=?
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
		User:      cfg.Database.User,
		Passwd:    cfg.Database.Password,
		Net:       "tcp",
		Addr:      cfg.Database.Address,
		DBName:    cfg.Database.Name,
		ParseTime: true,
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

func InsertUser(user models.User) error {
	if err := checkDBConnection(); err != nil {
		return err
	}

	_, err := DB.Exec(insertUserRequest, user.Email, user.Password)
	if err != nil {
		return err
	}

	return nil
}

func InsertToken(userId int, jwtToken string) (string, error) {
	if err := checkDBConnection(); err != nil {
		return "", err
	}
	tx := DB.MustBegin()

	tokenUUIDStruct, err := uuid.NewV7()
	tokenUUID := tokenUUIDStruct.String()

	if err != nil {
		return "", err
	}

	_, err = tx.Exec(insertTokenRequest, tokenUUID, jwtToken)
	if err != nil {
		return "", err
	}

	_, err = tx.Exec(insertUserTokenRequest, userId, tokenUUID)
	if err != nil {
		return "", err
	}

	if err := tx.Commit(); err != nil {
		return "", err
	}

	return tokenUUID, err
}

func UpdateTokenTime(token models.JwtToken) error {
	if err := checkDBConnection(); err != nil {
		return err
	}

	_, err := DB.Exec(updateTokenRequest, time.Now(), token.UserId) // TODO fix back
	if err != nil {
		return err
	}

	return nil
}

func GetTokenByUser(user models.User) (models.JwtToken, error) {
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

func GetUserByRefreshToken(tokenUUID string) (models.User, error) {
	var user models.User

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

func GetUserByEmail(email string) (models.User, error) {
	var user models.User

	if err := checkDBConnection(); err != nil {
		return user, err
	}

	rows := DB.QueryRowx(getUserByEmailRequest, email)

	err := rows.Scan(&user.Id, &user.Email, &user.Password, &user.RegistrationTime)

	if err != nil {
		return user, err
	}

	return user, nil
}

func GetUserById(userId int) (models.User, error) {
	var user models.User

	if err := checkDBConnection(); err != nil {
		return user, err
	}

	// err := DB.Get(&user, getUserByIdRequest, userId)
	rows := DB.QueryRowx(getUserByIdRequest, userId)

	err := rows.Scan(&user.Id, &user.Email, &user.Password, &user.RegistrationTime)

	if err != nil {
		return user, err
	}

	return user, nil
}
