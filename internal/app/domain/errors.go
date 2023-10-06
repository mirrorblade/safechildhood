package domain

import "errors"

var (
	ErrUndefinedFeature = errors.New("cannot find feature with this coordinates")
	ErrIncorrectColor   = errors.New("incorrect color")
	ErrUserNotFound     = errors.New("user was not found")
	ErrTooManyPhotos    = errors.New("too many photos")
)
