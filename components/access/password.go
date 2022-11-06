package access

import (
	"fmt"
	"github.com/zagiduller/photo-studio/components"
	"golang.org/x/crypto/bcrypt"
)

// @project photo-studio
// @created 06.11.2022

const (
	PasswordTypeLoginPassword = "login_password"
)

var (
	ErrPasswordTypeIsEmpty     = fmt.Errorf("password type is empty")
	ErrPasswordValueIsEmpty    = fmt.Errorf("password value is empty")
	ErrPasswordValueIsTooShort = fmt.Errorf("password value is too short")
	ErrPasswordTypeIsExist     = fmt.Errorf("password type is exist")

	passwordTypeValidator = map[string]func(t, v string) error{
		PasswordTypeLoginPassword: func(t, v string) error {
			if v == "" {
				return fmt.Errorf("%s type: [%w] ", t, ErrPasswordValueIsEmpty)
			}
			if len(v) < 6 {
				return fmt.Errorf("%s type: [%w] ", t, ErrPasswordValueIsTooShort)
			}
			return nil
		},
	}
)

type Password struct {
	components.Model

	LoginID uint
	Type    string `gorm:"type:varchar(32)"`
	Value   string `gorm:"type:varchar(255)"`
}

func NewPassword(loginID uint, passwordType, value string) *Password {
	return &Password{
		LoginID: loginID,
		Type:    passwordType,
		Value:   value,
	}
}

func FindPasswordByTypeAndLoginID(loginID uint, passwordType string) (*Password, error) {
	db := components.GetDB()
	if db == nil {
		return nil, fmt.Errorf("password.FindPasswordByTypeAndLoginID: [%w] ", components.ErrorCodeDbIsNil)
	}
	if loginID == 0 {
		return nil, fmt.Errorf("password.FindPasswordByTypeAndLoginID: [%w] ", ErrLoginIDIsEmpty)
	}
	if passwordType == "" {
		return nil, fmt.Errorf("password.FindPasswordByTypeAndLoginID: [%w] ", ErrPasswordTypeIsEmpty)
	}

	p := &Password{
		Type:    passwordType,
		LoginID: loginID,
	}
	if err := db.Model(p).First(p).Error; err != nil {
		return nil, fmt.Errorf("password.FindPasswordByTypeAndLoginID: [%w] ", err)
	}
	if p == nil {
		return nil, fmt.Errorf("password.FindPasswordByTypeAndLoginID: [%w] ", components.ErrModelNotFound)
	}
	return p, nil
}

func (p *Password) Validate() error {
	if p == nil {
		return fmt.Errorf("password.Validate: [%w] ", components.ErrModelIsNil)
	}
	if p.GetDB() == nil {
		return fmt.Errorf("password.Validate: [%w] ", components.ErrorCodeDbIsNil)
	}
	if p.LoginID == 0 {
		return fmt.Errorf("password.Validate: [%w] ", ErrLoginIDIsEmpty)
	}
	if p.Type == "" {
		return fmt.Errorf("password.Validate: [%w] ", ErrPasswordTypeIsEmpty)
	}
	if p.Value == "" {
		return fmt.Errorf("password.Validate: [%w] ", ErrPasswordValueIsEmpty)
	}
	if err := passwordTypeValidator[p.Type](p.Type, p.Value); err != nil {
		return fmt.Errorf("password.Validate: [%w] ", err)
	}
	if _, err := FindPasswordByTypeAndLoginID(p.LoginID, p.Type); err == nil {
		return fmt.Errorf("password.Validate: [%w] ", ErrPasswordTypeIsExist)
	}

	return nil
}

func (p *Password) Save() error {
	if err := p.Validate(); err != nil {
		return fmt.Errorf("password.Save: [%w] ", err)
	}
	if err := p.GetDB().Create(p).Error; err != nil {
		return fmt.Errorf("password.Save: [%w] ", err)
	}
	return nil
}

func (p *Password) Check(password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(p.Value), []byte(password)); err != nil {
		return fmt.Errorf("Check: [%w] ", err)
	}
	return nil
}
