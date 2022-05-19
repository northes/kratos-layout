package validator

import (
	"errors"

	"github.com/nyaruka/phonenumbers"
)

var (
	ErrAssertTypeToStringFailed = errors.New("assert type to string failed")
	ErrInvalidMobilePhoneNumber = errors.New("手机号码格式无效")
)

func IsMobilePhone(value interface{}) error {
	phone, ok := value.(string)
	if !ok {
		return ErrAssertTypeToStringFailed
	}

	phoneNumber, err := phonenumbers.Parse(phone, "CN")
	if err != nil {
		return err
	}

	if !phonenumbers.IsValidNumber(phoneNumber) {
		return ErrInvalidMobilePhoneNumber
	}

	return nil
}
