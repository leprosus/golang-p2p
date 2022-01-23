package p2p

type Binary interface {
	SetBytes(bs []byte)
	GetBytes() (bs []byte)
	SetGob(val interface{}) (err error)
	GetGob(val interface{}) (err error)
	SetJson(val interface{}) (err error)
	GetJson(val interface{}) (err error)
	String() (str string)
}
