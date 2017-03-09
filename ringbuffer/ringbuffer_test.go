package ringbuffer

import (
	"testing"
	"log"
)

func TestBuffer_Short(t *testing.T) {
	rb := New(1024)
	source := "hello,world"
	rb.Write([]byte(source))
	if rb.String() != source {
		log.Fatal("fuck hell wrong")
	}
}

func TestBuffer_Short2(t *testing.T) {
	rb := New(11)
	source := "hello,go"
	rb.Write([]byte(source))
	if rb.String() != source {
		log.Fatal("fuck hell wrong")
	}
}

func TestBuffer_Full(t *testing.T) {
	rb := New(5)
	source := "hello"
	rb.Write([]byte(source))
	if string(rb.Bytes()) != source {
		log.Fatal("full write fail")
	}
}

func TestBuffer_Overflow(t *testing.T) {
	rb := New(5)
	source := "hello world"
	rb.Write([]byte(source))
	if string(rb.Bytes()) != "world" {
		log.Fatal("overflow write fail")
	}
}

func TestBuffer_Overflow2(t *testing.T) {
	rb := New(11)
	source := "hello world,hello,golang"
	rb.Write([]byte(source))
	if string(rb.Bytes()) != "ello,golang" {
		log.Fatal("overflow write fail")
	}
}

func TestBuffer_Written(t *testing.T) {
	rb := New(11)
	source := "hello world,hello,golang"
	rb.Write([]byte(source))
	if rb.Written()!= len([]byte(source)) {
		log.Fatal("Writeen() fail")
	}
}
