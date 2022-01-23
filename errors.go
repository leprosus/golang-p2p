package p2p

import "errors"

var (
	UnsupportedPackage    = errors.New("unsupported package type")
	UnexpectedPackage     = errors.New("unexpected package type")
	UnsupportedTopic      = errors.New("unsupported topic")
	ConnectionError       = errors.New("connection error")
	PresetConnectionError = errors.New("preset connection error")
	CipherKeyError        = errors.New("cipher key error")
)
