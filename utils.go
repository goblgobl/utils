package utils

import (
	"encoding/binary"
	"reflect"
	"unsafe"
)

const (
	reqIdEncoding = "ABCDEFGHIJKLMNOPQRSTUVWXYZ234567"
)

var (
	NoopReleasable = noopReleasable{}
)

type Releasable interface {
	Release()
}

type noopReleasable struct {
}

func (_ noopReleasable) Release() {

}

func EncodeRequestId(requestId uint32, instanceId uint8) string {
	var data [4]byte
	id := data[:]

	binary.BigEndian.PutUint32(id, requestId)

	var encoded [8]byte
	encoded[7] = reqIdEncoding[instanceId&0x1F]
	encoded[6] = reqIdEncoding[(instanceId>>5|(id[3]<<3))&0x1F]
	encoded[5] = reqIdEncoding[(id[3]>>2)&0x1F]
	encoded[4] = reqIdEncoding[(id[3]>>7|(id[2]<<1))&0x1F]
	encoded[3] = reqIdEncoding[((id[2]>>4)|(id[1]<<4))&0x1F]
	encoded[2] = reqIdEncoding[(id[1]>>1)&0x1F]
	encoded[1] = reqIdEncoding[((id[1]>>6)|(id[0]<<2))&0x1F]
	encoded[0] = reqIdEncoding[id[0]>>3]
	return string(encoded[:])
}

func S2B(s string) (b []byte) {
	/* #nosec G103 */
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	/* #nosec G103 */
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh.Data = sh.Data
	bh.Cap = sh.Len
	bh.Len = sh.Len
	return b
}

func B2S(b []byte) string {
	/* #nosec G103 */
	return *(*string)(unsafe.Pointer(&b))
}

func Paging(perpage int, page int, defaultPerPage int) (int, int) {
	if perpage < 1 {
		perpage = defaultPerPage
	}

	if page < 1 {
		return perpage, 0
	}

	return perpage, (page - 1) * perpage
}
