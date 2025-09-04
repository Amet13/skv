package provider

import "errors"

// ErrNotFound indicates the requested secret does not exist in the provider.
var ErrNotFound = errors.New("secret not found")

