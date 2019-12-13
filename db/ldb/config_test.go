package ldb

import (
	"testing"
)

func TestClone(t *testing.T) {
	cfg1 := &Config{Path: "/testpath", ReadOnly: true, Compression: true, FileSize: 16}
	cfg2 := cfg1.Clone()

	if cfg1 == cfg2 {
		t.Fatal("clone return old config")
	}

	if cfg1.Path != cfg2.Path {
		t.Fatal("Path corrupted")
	}

	if cfg1.Compression != cfg2.Compression {
		t.Fatal("Compression corrupted")
	}

	if cfg1.ReadOnly != cfg2.ReadOnly {
		t.Fatal("ReadOnly corrupted")
	}

	if cfg1.FileSize != cfg2.FileSize {
		t.Fatal("FileSize corrupted")
	}
}
