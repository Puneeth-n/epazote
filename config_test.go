package epazote

import (
	"testing"
)

func TestConfigNewBadFile(t *testing.T) {
	_, err := New("test/no-exist.yml")
	if err == nil {
		t.Error(err)
	}
}

func TestConfigNewBadYaml(t *testing.T) {
	_, err := New("test/bad.yml")
	if err == nil {
		t.Error(err)
	}
}

func TestConfigNew(t *testing.T) {
	_, err := New("test/epazote.yml")
	if err != nil {
		t.Error(err)
	}
}

func TestParseScanBadFile(t *testing.T) {
	err := ParseScan("test/no-exist.yml")
	if err == nil {
		t.Error(err)
	}
}

func TestParseScanBadYaml(t *testing.T) {
	err := ParseScan("test/bad.yml")
	if err == nil {
		t.Error(err)
	}
}

func TestParseScanEmpty(t *testing.T) {
	err := ParseScan("test/empty.yml")
	if err != nil {
		t.Error(err)
	}
}

func TestParseScanBadUrl(t *testing.T) {
	err := ParseScan("test/bad-url.yml")
	if err != nil {
		t.Error(err)
	}
}

func TestParseScanEvery(t *testing.T) {
	err := ParseScan("test/every.yml")
	if err != nil {
		t.Error(err)
	}
}
