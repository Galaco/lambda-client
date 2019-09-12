package filesystem

import "testing"

func TestGetFile(t *testing.T) {
	t.Skip()
}

func TestRegisterPakfile(t *testing.T) {
	t.Skip()
}

func TestRegisterVpk(t *testing.T) {
	t.Skip()
}

func TestUnregisterVpk(t *testing.T) {
	t.Skip()
}

func TestRegisterLocalDirectory(t *testing.T) {
	fs := NewFileSystem()
	dir := "foo/bar/baz"
	fs.RegisterLocalDirectory(dir)
	found := false
	for _, path := range fs.localDirectories {
		if path == dir {
			found = true
			break
		}
	}
	if found == false {
		t.Error("local filepath was not found in registered paths")
	}
}

func TestUnregisterLocalDirectory(t *testing.T) {
	fs := NewFileSystem()
	dir := "foo/bar/baz"
	fs.RegisterLocalDirectory(dir)
	fs.UnregisterLocalDirectory(dir)
	found := false
	for _, path := range fs.localDirectories {
		if path == dir {
			found = true
			break
		}
	}
	if found == true {
		t.Error("local filepath was not found in registered paths")
	}
}

func TestUnregisterPakfile(t *testing.T) {
	t.Skip()
}
