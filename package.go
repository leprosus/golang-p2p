package p2p

type Package struct {
	Type PackageType
	Data
}

type PackageType uint8

const (
	Handshake PackageType = iota
	Resume
	Exchange
)

type ResumeStatus bool

const (
	ResumePossible   ResumeStatus = true
	ResumeImpossible ResumeStatus = false
)
