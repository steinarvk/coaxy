package prereader

import (
	"bytes"
	"io/ioutil"
	"testing"
)

func TestPrereaderBasic(t *testing.T) {
	header, r, err := Preread(bytes.NewBufferString("helloworld"), 4)
	if err != nil {
		t.Fatal(err)
	}

	allbytes, err := ioutil.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}

	if got, want := string(header), "hell"; got != want {
		t.Errorf("Preread() returned header %q want %q", got, want)
	}

	if got, want := string(allbytes), "helloworld"; got != want {
		t.Errorf("Preread() left %q want %q", got, want)
	}
}

func TestPrereaderShort(t *testing.T) {
	header, r, err := Preread(bytes.NewBufferString("helloworld"), 32)
	if err != nil {
		t.Fatal(err)
	}

	allbytes, err := ioutil.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}

	if got, want := string(header), "helloworld"; got != want {
		t.Errorf("Preread() returned header %q want %q", got, want)
	}

	if got, want := string(allbytes), "helloworld"; got != want {
		t.Errorf("Preread() left %q want %q", got, want)
	}
}
