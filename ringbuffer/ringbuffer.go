package ringbuffer

type RingBuffer struct {
	data    []byte
	written int
	cursor  int
	size    int
}

func New(size int) *RingBuffer {
	if size < 0 {
		panic("RingBuffer size should be > 0 ")
	}
	return &RingBuffer{size: size, data: make([]byte, size)}
}

// Write implements io.Writer interface
func (rb *RingBuffer) Write(buf []byte) (n int, err error) {
	rb.written += len(buf)
	if len(buf) > rb.size {
		buf = buf[len(buf)-rb.size:]
	}
	copy(rb.data[rb.cursor:], buf)
	remain := rb.size - rb.cursor
	if len(buf) > remain {
		copy(rb.data, buf[remain:])
	}
	rb.cursor = (rb.cursor + len(buf)) % rb.size
	return len(buf), nil
}

// Stirng() implements Stringer interface
func (rb *RingBuffer) String() string {
	return string(rb.Bytes())
}

func (rb *RingBuffer) Size() int {
	return rb.size
}
func (rb *RingBuffer) Reset() {
	rb.cursor = 0
	rb.written = 0
}

func (rb *RingBuffer) Written() int {
	return rb.written
}

func (rb *RingBuffer) Bytes() []byte {
	// written >= rb.size .the equality should be here.otherwise the when full-write happens,
	// the cursor should be 0,return the default rb.data[:rb.cursor] would be wrong
	if rb.written >= rb.size && rb.cursor == 0 {
		return rb.data
	}
	if rb.written > rb.size {
		result := make([]byte, rb.size)
		copy(result, rb.data[rb.cursor:])
		copy(result[rb.size-rb.cursor:], rb.data[:rb.cursor])
		return result
	}
	return rb.data[:rb.cursor]
}
