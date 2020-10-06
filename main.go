package main

// #include <string.h>
// #include <stdbool.h>
// #include <mysql.h>
// #cgo CFLAGS: -O3 -I/usr/include/mysql -fno-omit-frame-pointer
import "C"
import (
	"encoding/hex"
	"hash"
	"log"
	"os"
	"unsafe"

	"golang.org/x/crypto/sha3"
)

// main function is needed even for generating shared object files
func main() {}

var l = log.New(os.Stderr, "sha3: ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Llongfile)

func msg(message *C.char, s string) {
	m := C.CString(s)
	defer C.free(unsafe.Pointer(m))

	C.strcpy(message, m)
}

//export Sha3_init
func Sha3_init(initid *C.UDF_INIT, args *C.UDF_ARGS, message *C.char) C.bool {
	if args.arg_count != 2 {
		msg(message, "`sha3`() requires 2 parameters: the message part, and the number of bits")
		return C.bool(true)
	}

	argsTypes := (*[2]uint32)(unsafe.Pointer(args.arg_type))

	argsTypes[0] = C.STRING_RESULT
	argsTypes[1] = C.INT_RESULT

	initid.maybe_null = C.bool(true)

	return C.bool(false)
}

//export Sha3
func Sha3(initid *C.UDF_INIT, args *C.UDF_ARGS, result *C.char, length *uint64, isNull *C.char, message *C.char) *C.char {
	*isNull = 1

	argsArgs := (*[1 << 30]*C.char)(unsafe.Pointer(args.args))[:2:2]

	bits := (*[1 << 30]*uint64)(unsafe.Pointer(args.args))[1]

	if argsArgs[0] == nil || argsArgs[1] == nil {
		return nil
	}

	var h hash.Hash
	switch *bits {
	case 224:
		h = sha3.New224()
	case 256:
		h = sha3.New256()
	case 384:
		h = sha3.New384()
	case 512:
		h = sha3.New512()
	default:
		return nil
	}

	argsLengths := (*[1 << 30]uint64)(unsafe.Pointer(args.lengths))[:1:1]
	s := C.GoStringN(argsArgs[0], C.int(argsLengths[0]))

	b := make([]byte, *bits/8)
	h.Write([]byte(s))
	h.Sum(b[:0])
	s = hex.EncodeToString(b)

	*length = uint64(len(s))
	*isNull = 0
	return C.CString(s)
}

//export unhex_sha3_init
func unhex_sha3_init(initid *C.UDF_INIT, args *C.UDF_ARGS, message *C.char) C.bool {
	if args.arg_count != 2 {
		msg(message, "`unhex_sha3`() requires 2 parameters: the message part, and the number of bits")
		return C.bool(true)
	}

	argsTypes := (*[2]uint32)(unsafe.Pointer(args.arg_type))

	argsTypes[0] = C.STRING_RESULT
	argsTypes[1] = C.INT_RESULT

	initid.maybe_null = C.bool(true)

	return C.bool(false)
}

//export unhex_sha3
func unhex_sha3(initid *C.UDF_INIT, args *C.UDF_ARGS, result *C.char, length *uint64, isNull *C.char, message *C.char) *C.char {
	*isNull = 1

	argsArgs := (*[1 << 30]*C.char)(unsafe.Pointer(args.args))[:2:2]

	bits := (*[1 << 30]*uint64)(unsafe.Pointer(args.args))[1]

	if argsArgs[0] == nil || argsArgs[1] == nil {
		return nil
	}

	var h hash.Hash
	switch *bits {
	case 224:
		h = sha3.New224()
	case 256:
		h = sha3.New256()
	case 384:
		h = sha3.New384()
	case 512:
		h = sha3.New512()
	default:
		return nil
	}

	argsLengths := (*[1 << 30]uint64)(unsafe.Pointer(args.lengths))[:1:1]
	s := C.GoStringN(argsArgs[0], C.int(argsLengths[0]))

	b := make([]byte, *bits/8)
	h.Write([]byte(s))
	h.Sum(b[:0])

	*length = uint64(len(b))
	*isNull = 0
	return C.CString(string(b))
}
