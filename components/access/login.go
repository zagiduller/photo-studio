package access

import (
	"errors"
	"fmt"
	"github.com/zagiduller/photo-studio/components"
	"github.com/zagiduller/photo-studio/components/users"
	"golang.org/x/crypto/bcrypt"
	"net/mail"
)

// @project photo-studio
// @created 06.11.2022

const (
	LoginTypeEmail = "email"
	LoginTypePhone = "phone"
)

var (
	ErrLoginTypeIsEmpty  = errors.New("login type is empty")
	ErrLoginValueIsEmpty = errors.New("login value is empty")

	LoginTypeValidators = map[string]func(value string) error{
		LoginTypeEmail: func(value string) error {
			_, err := mail.ParseAddress(value)
			return err
		},
	}
)

type Login struct {
	components.Model
	User     *users.User `gorm:"foreignKey:UserID"`
	Type     string      `json:"type" gorm:"type:varchar(32);"`
	Value    string      `json:"value" gorm:"type:varchar(255);unique;"`
	UserID   uint        `json:"user_id"`
	Verified bool        `json:"verified"`
}

func NewLogin(id uint, loginType, value string) *Login {
	return &Login{
		UserID: id,
		Type:   loginType,
		Value:  value,
	}
}

// CreateLoginPassword create password hash nad save
func CreateLoginPassword(login *Login, password string) error {
	if &login == nil || login.ID == 0 {
		return fmt.Errorf("CreateLoginPassword: [%w] ", components.ErrModelIsNil)
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return fmt.Errorf("CreateLoginPassword: [%w] ", err)
	}
	pwd := NewPassword(login.ID, PasswordTypeLoginPassword, string(hash))
	pwd.SetDB(login.GetDB())
	if err := pwd.Save(); err != nil {
		return fmt.Errorf("CreateLoginPassword: [%w] ", err)
	}
	return nil
}

func FindLoginByID(id uint) (*Login, error) {
	if id == 0 {
		return nil, fmt.Errorf("FindLoginByID: [%w] ", components.ErrModelIsNil)
	}
	db := components.GetDB()
	login := &Login{}
	if err := db.Model(&Login{}).Where("id = ?", id).Preload("User").First(login).Error; err != nil {
		return nil, fmt.Errorf("FindLoginByID: [%w] ", err)
	}
	return login, nil
}

func FindLoginsByUserID(id uint) ([]*Login, error) {
	if id == 0 {
		return nil, fmt.Errorf("FindLoginsByUserID: [%w] ", components.ErrModelIsNil)
	}
	db := components.GetDB()
	logins := make([]*Login, 0)
	if err := db.Model(&Login{}).Where("user_id = ?", id).Preload("User").Find(&logins).Error; err != nil {
		return nil, fmt.Errorf("FindLoginsByUserID: [%w] ", err)
	}
	return logins, nil
}

func FindLoginByValue(value string) (*Login, error) {
	login := &Login{}
	if value == "" {
		return nil, errors.New("FindLoginByValue: empty value ")
	}
	db := components.GetDB()
	if err := db.Model(&Login{}).Where("value = ?", value).Preload("User").First(login).Error; err != nil {
		return nil, fmt.Errorf("FindLoginByValue: [%w] ", err)
	}
	return login, nil
}

func (l *Login) Validate() error {
	if l == nil {
		return fmt.Errorf("login.Validate: [%w] ", components.ErrModelIsNil)
	}
	if l.GetDB() == nil {
		return fmt.Errorf("login.Validate: [%w] ", components.ErrorCodeDbIsNil)
	}
	if l.Type == "" {
		return fmt.Errorf("login.Validate: [%w] ", ErrLoginTypeIsEmpty)
	}
	if l.Value == "" {
		return fmt.Errorf("login.Validate: [%w] ", ErrLoginValueIsEmpty)
	}
	validator, ok := LoginTypeValidators[l.Type]
	if !ok {
		return fmt.Errorf("login.Validate: [%w] ", ErrLoginTypeIsEmpty)
	}
	if err := validator(l.Value); err != nil {
		return fmt.Errorf("login.Validate: [%w] ", err)
	}
	return nil
}

func (l *Login) Save() error {
	if err := l.Validate(); err != nil {
		return fmt.Errorf("login.Save: [%w] ", err)
	}
	if err := l.GetDB().Create(l).Error; err != nil {
		return fmt.Errorf("login.Save: [%w] ", err)
	}
	return nil
}

func (l *Login) CheckPassword(value string) error {
	if value == "" {
		return fmt.Errorf("login.Check: [%w] ", ErrPasswordValueIsEmpty)
	}
	if err := l.Validate(); err != nil {
		return fmt.Errorf("login.Check: [%w] ", err)
	}
	password, err := FindPasswordByTypeAndLoginID(l.ID, l.Type)
	if err != nil {
		return fmt.Errorf("login.Check: [%w] ", err)
	}
	if err := password.Check(value); err != nil {
		return fmt.Errorf("login.Check: [%w] ", err)
	}
	return nil
}

func (l *Login) GetUser() *users.User {
	return l.User
}
