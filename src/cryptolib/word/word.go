package word

func U16toByte_16(a uint16) (b []byte) {
	b = make([]byte, 2)
	b[0] = byte(a >> 8)
	b[1] = byte(a)
	return
}
func U32toByte_32(a uint32) (b []byte) {
	b = make([]byte, 4)
	b[0] = byte(a >> 24)
	b[1] = byte(a >> 16)
	b[2] = byte(a >> 8)
	b[3] = byte(a)
	return
}
func BytetoU32_32(a []byte) (b uint32) {
	b = uint32(a[0])<<24 | uint32(a[1])<<16 | uint32(a[2])<<8 | uint32(a[3])<<0
	return
}

func U64toByte_64_asm(a uint64, b []byte) //by .asm
func U64toByte_64(a uint64) (b []byte) {
	b = make([]byte, 8)
	U64toByte_64_asm(a, b)
	return
}
func BytetoU64_64(a []byte) (b uint64) //by .asm

func U64toByte_256_asm(a [4]uint64, b []byte) //by .asm
func U64toByte_256(a [4]uint64) (b []byte) {
	b = make([]byte, 32)
	U64toByte_256_asm(a, b)
	return
}
func BytetoU64_256(x []byte) (y [4]uint64) //by .asm

// b is explored with big-endian, a is little-endian for u32
func U32toByte_256_asm(a [8]uint32, b []byte) //by .asm
func U32toByte_256(a [8]uint32) (b []byte) {
	b = make([]byte, 32)
	U32toByte_256_asm(a, b)
	return
}
func BytetoU32_256(a []byte) (b [8]uint32) //by .asm

func U64toU32_256(a [4]uint64) (b [8]uint32) //by .asm
func U32toU64_256(a [8]uint32) (b [4]uint64) //by .asm

func U32toByte_128_asm(a [4]uint32, b []byte) //by .asm
func U32toByte_128(a [4]uint32) (b []byte) {
	b = make([]byte, 16)
	U32toByte_128_asm(a, b)
	return
}
func BytetoU32_128(a []byte) (b [4]uint32) //by .asm

func BytetoU32_512(x []byte) (y [16]uint32) //by .asm
