package reply

// ping pong
type PongReply struct {
}

var pongBytes = []byte("+PONG\r\n")

func (p *PongReply) ToBytes() []byte {
	return pongBytes
}

func MakePongReply() *PongReply {
	return &PongReply{}
}

// ok
type OkReply struct {
}

var okBytes = []byte("+OK\r\n")

func (o *OkReply) ToBytes() []byte {
	return okBytes
}

func MakeOkReply() *OkReply {
	return &OkReply{}
}

type EmptyMultiBulkReply struct {
}

var emptyMultiBulkBytes = []byte("*0\r\n")

func (e *EmptyMultiBulkReply) ToBytes() []byte {
	return emptyMultiBulkBytes
}

func MakeEmptyMultiBulkReply() *EmptyMultiBulkReply {
	return &EmptyMultiBulkReply{}
}

type NullBulkReply struct {
}

var nullBulkBytes = []byte("$-1\r\n")

func (e *NullBulkReply) ToBytes() []byte {
	return nullBulkBytes
}

func MakeNullBulkReply() *NullBulkReply {
	return &NullBulkReply{}
}

type NoReply struct {
}

var noBytes = []byte("")

func (e *NoReply) ToBytes() []byte {
	return noBytes
}

func MakeNoReply() *NoReply {
	return &NoReply{}
}
