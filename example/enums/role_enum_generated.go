// Code generated by ggenums; DO NOT EDIT.
package enums

import (
	"encoding/json"
	"fmt"
)

type RoleEnum string

const (
	RoleAdmin RoleEnum = "admin"
	RoleUser  RoleEnum = "user"
	RoleGuest RoleEnum = "guest"
)

var AllRoles = []RoleEnum{
	RoleAdmin,
	RoleUser,
	RoleGuest,
}

func (e RoleEnum) String() string {
	return string(e)
}

func (e RoleEnum) Validate() error {
	switch e {
	case RoleAdmin, RoleUser, RoleGuest:
		return nil
	default:
		return fmt.Errorf("invalid Role: %s", e)
	}
}

func ParseRole(s string) (RoleEnum, error) {
	e := RoleEnum(s)
	if err := e.Validate(); err != nil {
		return "", err
	}
	return e, nil
}

func (e RoleEnum) MarshalJSON() ([]byte, error) {
	if err := e.Validate(); err != nil {
		return []byte("null"), nil
	}
	return json.Marshal(string(e))
}

func (e *RoleEnum) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	parsed, err := ParseRole(s)
	if err != nil {
		return err
	}

	*e = parsed
	return nil
}
