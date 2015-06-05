// Package localization abstracts i18n implementation and provides it as a
// infrastructure.
package localization

import (
	"errors"

	"github.com/nicksnyder/go-i18n/i18n"
)

var (
	// ErrNoTranslationFiles denotes no translation files were requested: i.e. []string{}.
	ErrNoTranslationFiles = errors.New("No translation files provided")
	// T is the default translator function
	T i18n.TranslateFunc
)

// *****************************************************************************
// Localized Error Messages
// *****************************************************************************

// Error represents a localized error. This abstraction is useful for offering
// a type check to determine if an error should be shared with a user--or not.
// For example, an error like "You gave us the wrong password" should be shared
// with a user; however, an error like "Unable to connect to database" is more
// of a systems issue and we should only tell the user we had an internal error.
type Error struct {
	Err error
}

func (e Error) Error() string {
	return e.Err.Error()
}

// NewError returns a new localized error.
func NewError(message string) *Error {
	return &Error{Err: errors.New(message)}
}

// *****************************************************************************

// Init configures the localizaton page. Pass the valid path(s) to i18n
// JSON files as defined by the `github.com/nicksnyder/go-i18n` package.
func Init(files []string, defaultLocale string) (err error) {
	if len(files) == 0 {
		return ErrNoTranslationFiles
	}

	for _, path := range files {
		i18n.MustLoadTranslationFile(path)
	}

	T = NewTranslationFunc(defaultLocale)

	return nil
}

// NewTranslationFunc returns a new translation function
func NewTranslationFunc(locale string) i18n.TranslateFunc {
	return i18n.MustTfunc(locale)
}
