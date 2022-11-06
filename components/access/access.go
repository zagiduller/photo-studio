package access

import (
	"fmt"
	"github.com/zagiduller/photo-studio/components"
	"github.com/zagiduller/photo-studio/components/users"
)

// @project photo-studio
// @created 10.08.2022

var (
	ErrTokenIsEmpty   = fmt.Errorf("token is empty")
	ErrLoginIDIsEmpty = fmt.Errorf("login id is empty")
)

type Access struct {
	components.Model

	User      users.User
	Token     string
	UserID    uint
	SessionID string
}

func NewAccess(userID uint, token string) *Access {
	return &Access{
		Token:  token,
		UserID: userID,
	}
}

func (a *Access) Validate() error {
	if a.GetDB() == nil {
		return fmt.Errorf("access.Validate: [%w] ", components.ErrorCodeDbIsNil)
	}
	if a.Token == "" {
		return fmt.Errorf("access.Validate: [%w] ", ErrTokenIsEmpty)
	}
	if a.UserID == 0 {
		return fmt.Errorf("access.Validate: [%w] ", ErrLoginIDIsEmpty)
	}

	return nil
}

func (a *Access) Save() error {
	if err := a.Validate(); err != nil {
		return fmt.Errorf("access.Save: [%w] ", err)
	}
	if err := a.GetDB().Create(a).Error; err != nil {
		return fmt.Errorf("access.Save: [%w] ", err)
	}
	return nil
}
