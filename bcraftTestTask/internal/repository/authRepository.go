package repository

import (
	"bcraftTestTask/internal/logging"
	"bcraftTestTask/internal/models"
	"bcraftTestTask/internal/properties"
	"database/sql"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"strings"
)

type AuthRepository struct {
	db *sql.DB
}

func NewAuthRepository(db_ *sql.DB) *AuthRepository {
	return &AuthRepository{
		db: db_,
	}
}

func (r *AuthRepository) Validate(account models.Account) error {
	if !strings.Contains(account.Email, "@") {
		return fmt.Errorf("Email address is required")
	}

	if len(account.Password) < 6 {
		return fmt.Errorf("Password is required")
	}

	row, err := r.db.Query("select id from accounts where email = $1;", account.Email)

	if err != nil {
		return err
	}

	if row != nil {
		return fmt.Errorf("Email address already in use by another user.")
	}

	return nil
}

func (r *AuthRepository) CreateAccount(account models.Account) error {
	config, err := properties.GetConfig()
	if err != nil {
		logger := logging.GetLogger()
		logger.Fatal(err)
	}

	err = r.Validate(account)
	if err != nil {
		return err
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(account.Password), bcrypt.DefaultCost)

	tk := &models.Token{UserId: account.Id}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(config.ProgramSettings.TokenPassword))

	_, err = r.db.Exec("insert into account (email, password, token) "+
		"values ($1, $2, $3);", account.Email, hashedPassword, tokenString)
	if err != nil {
		return err
	}

	return nil
}

func (r *AuthRepository) Login(email, password string) (*models.Account, error) {
	config, err := properties.GetConfig()
	if err != nil {
		logger := logging.GetLogger()
		logger.Fatal(err)
	}

	account := models.Account{}
	err = r.db.QueryRow("select id, email, password from accounts where email = $1;",
		email).Scan(&account.Id, &account.Email, &account.Password)

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("Email address not found")
		}
		return nil, fmt.Errorf("Connection error. Please retry")
	}

	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return nil, fmt.Errorf("Invalid login credentials. Please try again")
	}
	account.Password = ""

	tk := &models.Token{UserId: account.Id}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(config.ProgramSettings.TokenPassword))
	account.Token = tokenString

	return &account, nil
}
