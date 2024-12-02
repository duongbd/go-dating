package constant

import (
	"time"
)

const (
	JWT_CONFIG_SECRET_KEY      = "JWT_SECRET_KEY" // is the key to get the secret key from the environment
	JWT_DEFAULT_SECRET_VALUE   = "secret_key"     // is the default value of the secret key that uses for test and local
	JWT_CONFIG_DURATION_KEY    = "JWT_DURATION"   // is the key to get the duration from the environment
	JWT_DEFAULT_DURATION_VALUE = 72 * time.Hour   // is the default value of the duration that uses for test and local
)
