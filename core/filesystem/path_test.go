package filesystem

import "testing"

func TestNormalisePath(t *testing.T) {
	path := "foo\\bar\\baz"
	expected := "foo/bar/baz"
	actual := NormalisePath(path)
	if expected != actual {
		t.Errorf("incorrect path normalised. Expected %s, but received: %s", expected, actual)
	}
}
