package epazote

import (
	//	"fmt"
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
	_, err := ParseScan("test/no-exist.yml")
	if err == nil {
		t.Error(err)
	}
}

func TestParseScanBadYaml(t *testing.T) {
	_, err := ParseScan("test/bad.yml")
	if err == nil {
		t.Error(err)
	}
}

func TestParseScanEmpty(t *testing.T) {
	_, err := ParseScan("test/empty.yml")
	if err == nil {
		t.Error(err)
	}
}

func TestParseScanBadUrl(t *testing.T) {
	_, err := ParseScan("test/bad-url.yml")
	if err != nil {
		t.Error(err)
	}
}

func TestParseScanEvery(t *testing.T) {
	s, err := ParseScan("test/every.yml")

	if err != nil {
		t.Error(err)
	}

	switch {
	case s["service 1"].Every.Seconds != 30:
		t.Error("Expecting 60 got: ", s["service 1"].Every.Minutes)
	case s["service 2"].Every.Minutes != 1:
		t.Error("Expecting 1 got:", s["service 2"].Every.Minutes)
	case s["service 3"].Every.Hours != 2:
		t.Error("Expecting 2 got:", s["service 3"].Every.Hours)
	}
}

func TestCheckPathsNe(t *testing.T) {
	cfg, err := New("test/epazote-checkpaths-ne.yml")
	if err != nil {
		t.Error(err, cfg)
	}

	// scan check config and clean paths
	err = cfg.CheckPaths()
	if err == nil {
		t.Error("Expecting: Verify that directory: nonexist, exists and is readable.")
	}
}

func TestCheckPaths(t *testing.T) {
	cfg, err := New("test/epazote-checkpaths.yml")
	if err != nil {
		t.Error(err, cfg)
	}

	// scan check config and clean paths
	err = cfg.CheckPaths()
	if err != nil {
		t.Error(err)
	}
}

func TestCheckPathsEmpty(t *testing.T) {
	cfg, err := New("test/epazote-checkpaths-empty.yml")
	if err != nil {
		t.Error(err, cfg)
	}

	// scan check config and clean paths
	err = cfg.CheckPaths()
	if err != nil {
		t.Error(err)
	}
}

func TestCheckVerifyUrlsOk(t *testing.T) {
	cfg, err := New("test/every.yml")
	if err != nil {
		t.Error(err, cfg)
	}

	// scan check config and clean paths
	err = cfg.CheckPaths()
	if err != nil {
		t.Error(err)
	}

	err = cfg.VerifyUrls()
	if err != nil {
		t.Error(err)
	}
}

func TestCheckVerifyBadUrls(t *testing.T) {
	cfg, err := New("test/epazote.yml")
	if err != nil {
		t.Error(err, cfg)
	}

	// scan check config and clean paths
	err = cfg.CheckPaths()
	if err != nil {
		t.Error(err)
	}

	err = cfg.VerifyUrls()
	if err == nil {
		t.Error(err)
	}
}

func TestPathsOrServicesEmpty(t *testing.T) {
	cfg, err := New("test/empty.yml")
	if err != nil {
		t.Error(err)
	}
	err = cfg.PathsOrServices()
	if err == nil {
		t.Error(err)
	}
}

func TestPathsOrServices(t *testing.T) {
	cfg, err := New("test/epazote.yml")
	if err != nil {
		t.Error(err)
	}
	err = cfg.PathsOrServices()
	if err != nil {
		t.Error(err)
	}
}
