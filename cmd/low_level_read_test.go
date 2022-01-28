package main

import (
	"fmt"
	"testing"
)

func TestValidBytes(t *testing.T) {
	var magicBytes = []byte{0x66, 0x67, 0x62, SupportedVersion, 0x66, 0x67, 0x62, 0x01}

	expectedVersion := fmt.Sprintf("%d.0.1", SupportedVersion)

	version, err := Version(magicBytes)

	if expectedVersion != version || err != nil {
		t.Fatalf(`Expected version was %s, result was %s and error is %v`, expectedVersion, version, err)
	}
}

func TestInvalidBytes(t *testing.T) {
	var magicBytes = []byte{0x99, 0x67, 0x62, SupportedVersion, 0x66, 0x67, 0x62, 0x01}

	version, err := Version(magicBytes)

	if version != "" || err != ErrInvalidFile {
		t.Fatalf(`Expected version was %q, result was %q and error is %v`, "", version, err)
	}
}

func TestInvalidVersion(t *testing.T) {
	var magicBytes = []byte{0x66, 0x67, 0x62, 2, 0x66, 0x67, 0x62, 0x01}

	version, err := Version(magicBytes)

	if version != "" || err != ErrUnsupportedVersion {
		t.Fatalf(`Expected version was %q, result was %q and error is %v`, "", version, err)
	}
}
