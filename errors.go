package p2p

import "errors"

var (
	UnsupportedPackage    = errors.New("unsupported package type")
	UnsupportedTopic      = errors.New("unsupported topic")
	ConnectionError       = errors.New("connection error")
	PresetConnectionError = errors.New("preset connection error")
)
