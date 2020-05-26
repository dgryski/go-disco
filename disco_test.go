package disco

import (
	"io/ioutil"
	"testing"
)

func TestHash(t *testing.T) {

	const filename = `testdata/0xa2a647993898a3df.txt`

	b, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Fatal(err)
	}

	got := BEBB4185_64(b, 0)

	want := uint64(0xa2a647993898a3df)

	if got != want {
		t.Errorf("Got %v, want %v", got, want)
	}
}
