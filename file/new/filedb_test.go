package new

import (
	"bytes"
	"testing"
)

func TestHelloWorld(t *testing.T) {
	db := FileDB{
		p: NewMemoryPersistance(),
	}

	rc := nopReadCloser{bytes.NewBufferString("hello world")}

	err := db.Write("hash", rc)
	if err != nil {
		t.Error("got error on write", err.Error())
	}

	var b bytes.Buffer
	wc := nopWriteCloser{&b}

	err = db.Read("hash", wc)

	if err != nil {
		t.Error("got error on read", err.Error())
	}

	if b.String() != "hello world" {
		t.Error("got:", b.String(), "expected:", "hello world")
	}
}
