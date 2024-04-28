package presentation

import (
	"Omnichannel-CRM/package/enum"
)

type UserInfo struct {
	Username  string `json:"username"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Role      int    `json:"role"`
	ID        string `json:"user_id"`
	City      string `json:"city"`
	Province  string `json:"province"`
}

type LoginViewModel struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (lvm *LoginViewModel) ValidateLogin() map[string]string {
	var errorMessages = make(map[string]string)

	if lvm.Username == "" {
		errorMessages["errorStatus"] = enum.USERNAME_REQUIRED_STATUS
		errorMessages["errorMessage"] = enum.USERNAME_REQUIRED_MESSAGE
		return errorMessages
	}

	if lvm.Password == "" {
		errorMessages["errorStatus"] = enum.PASSWORD_REQUIRED_STATUS
		errorMessages["errorMessage"] = enum.PASSWORD_REQUIRED_MESSAGE
		return errorMessages
	}

	return errorMessages
}
