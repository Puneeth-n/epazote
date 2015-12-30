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

func TestConfigGetIntervalDefault(t *testing.T) {
	e := Every{}
	i := GetInterval(0, e)
	if i != 60 {
		t.Error("Expected 60")
	}
}

func TestConfigGetIntervalSeconds(t *testing.T) {
	e := Every{1, 0, 0}
	i := GetInterval(30, e)
	if i != 1 {
		t.Error("Expected 1")
	}
}

func TestConfigGetIntervalMinutes(t *testing.T) {
	e := Every{0, 1, 0}
	i := GetInterval(30, e)
	if i != 60 {
		t.Error("Expected 60")
	}
}

func TestConfigGetIntervalHours(t *testing.T) {
	e := Every{0, 0, 1}
	i := GetInterval(30, e)
	if i != 3600 {
		t.Error("Expected 3600")
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
