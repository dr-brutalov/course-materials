// Optional Todo

package hscan

import (
	"testing"
)

var file = "wordlist.txt"

func TestGuessSingle(t *testing.T) {
	got := GuessSingle("77f62e3524cd583d698d51fa24fdff4f", file) // Currently function returns only number of open ports
	want := "foo"
	if got != want {
		t.Errorf("got %d, wanted %d", got, want)
	}

}
