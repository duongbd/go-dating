package constant

import "errors"

type UserGender string

// UserGenderMale and UserGenderFemale are the possible values foo user's gender
const (
	UserGenderMale   UserGender = "MALE"
	UserGenderFemale UserGender = "FEMALE"

	DefaultDiscoveryDistance = 5000
)

var ErrEmailInUse = errors.New("email already in use") // ErrEmailInUse is the error message when the email is already in use
