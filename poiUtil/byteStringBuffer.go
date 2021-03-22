package poiUtil

import (
	"github.com/pkg/errors"
	"unicode/utf8"
	"unsafe"
)

// 改造 strings.Builder
// 功能 1 : replace 能对 原字符串替换，而不用重新分配内存
// 功能 2 : Reset 修改 Reset 功能到达重复利用原 buf

// ByteStringBuffer 注意事项
// ByteStringBuffer 对象不要相互拷贝
// 对要使用的字符串操作，空间上的预分配可以提升程序的性能
// 配合 Reset 函数以达到重复利用已分配的内存，进一步提升性能

// A Builder is used to efficiently build a string using Write methods.
// It minimizes memory copying. The zero value is ready to use.
// Do not copy a non-zero Builder.
type ByteStringBuffer struct {
	addr *ByteStringBuffer // of receiver, to detect copies by value
	buf  []byte
}

// noescape hides a pointer from escape analysis.  noescape is
// the identity function but escape analysis doesn't think the
// output depends on the input. noescape is inlined and currently
// compiles down to zero instructions.
// USE CAREFULLY!
// This was copied from the runtime; see issues 23382 and 7921.
//go:nosplit
func noescape(p unsafe.Pointer) unsafe.Pointer {
	x := uintptr(p)
	return unsafe.Pointer(x ^ 0)
}

func (b *ByteStringBuffer) copyCheck() {
	if b.addr == nil {
		// This hack works around a failing of Go's escape analysis
		// that was causing b to escape and be heap allocated.
		// See issue 23382.
		// TODO: once issue 7921 is fixed, this should be reverted to
		// just "b.addr = b".
		b.addr = (*ByteStringBuffer)(noescape(unsafe.Pointer(b)))
	} else if b.addr != b {
		panic("strings: illegal use of non-zero Builder copied by value")
	}
}

// String returns the accumulated string.
// 注意返回的是指针的东西，如果赋值，bsb.WriteString("ABCdef") a := bsb.String() 只是指针的拷贝
// Reset后，写入数据，bsb.WriteString("123") 会体现在 a 变量里 a => "123def" 而不是 "ABCdef"
// 如果要拷贝数据 CopyString()函数
func (b *ByteStringBuffer) String() string {
	return *(*string)(unsafe.Pointer(&b.buf))
}

func (b *ByteStringBuffer) CopyString() string {
	return string(b.buf[0:len(b.buf)])
}

// Len returns the number of accumulated bytes; b.Len() == len(b.String()).
func (b *ByteStringBuffer) Len() int { return len(b.buf) }

// Reset resets the Builder to be empty.
func (b *ByteStringBuffer) Reset() {
	//b.addr = nil
	b.buf = b.buf[:0]
}

// grow copies the buffer to a new, larger buffer so that there are at least n
// bytes of capacity beyond len(b.buf).
func (b *ByteStringBuffer) grow(n int) {
	buf := make([]byte, len(b.buf), 2*cap(b.buf)+n)
	copy(buf, b.buf)
	b.buf = buf
}

// Grow grows b's capacity, if necessary, to guarantee space for
// another n bytes. After Grow(n), at least n bytes can be written to b
// without another allocation. If n is negative, Grow panics.
func (b *ByteStringBuffer) Grow(n int) {
	b.copyCheck()
	if n < 0 {
		panic("strings.Builder.Grow: negative count")
	}
	if cap(b.buf)-len(b.buf) < n {
		b.grow(n)
	}
}

// Write appends the contents of p to b's buffer.
// Write always returns len(p), nil.
func (b *ByteStringBuffer) Write(p []byte) (int, error) {
	b.copyCheck()
	b.buf = append(b.buf, p...)
	return len(p), nil
}

// WriteByte appends the byte c to b's buffer.
// The returned error is always nil.
func (b *ByteStringBuffer) WriteByte(c byte) error {
	b.copyCheck()
	b.buf = append(b.buf, c)
	return nil
}

// WriteRune appends the UTF-8 encoding of Unicode code point r to b's buffer.
// It returns the length of r and a nil error.
func (b *ByteStringBuffer) WriteRune(r rune) (int, error) {
	b.copyCheck()
	if r < utf8.RuneSelf {
		b.buf = append(b.buf, byte(r))
		return 1, nil
	}
	l := len(b.buf)
	if cap(b.buf)-l < utf8.UTFMax {
		b.grow(utf8.UTFMax)
	}
	n := utf8.EncodeRune(b.buf[l:l+utf8.UTFMax], r)
	b.buf = b.buf[:l+n]
	return n, nil
}

// WriteString appends the contents of s to b's buffer.
// It returns the length of s and a nil error.
func (b *ByteStringBuffer) WriteString(s string) (int, error) {
	b.copyCheck()
	b.buf = append(b.buf, s...)
	return len(s), nil
}

// Replace fill some charactor by index
// it returns false and error, if pos not in buf range
func (b *ByteStringBuffer) Replace(pos int, c byte, count int) (bool, error) {
	if pos < 0 || pos >= b.Len() || pos+count > b.Len() {
		return false, errors.New("pos or pos+count out of buf range")
	}

	for i := 0; i < count; i++ {
		b.buf[pos+i] = c
	}
	return true, nil
}

// slice buf
func (b *ByteStringBuffer) Slice(start, end int) (bool, error) {
	if start < 0 || start > b.Len() || end > b.Len() {
		return false, errors.New("start or end out of range")
	}

	b.buf = b.buf[start:end]
	return true, nil
}
