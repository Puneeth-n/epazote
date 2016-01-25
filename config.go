package epazote

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
)

const herb = "\U0001f33f"

type Epazote struct {
	Config   Config
	Services Services
	debug    bool
}

type Config struct {
	SMTP Email `yaml:"smtp"`
	Scan Scan  `yaml:"scan"`
}

type Email struct {
	Username string
	Password string
	Server   string
	Port     int
	Headers  map[string]string
}

type Every struct {
	Seconds, Minutes, Hours int
}

type Scan struct {
	Paths []string
	Every `yaml:",inline"`
}

type Services map[string]Service

type Test struct {
	Test  string `json:"test,omitempty"`
	IfNot Action `yaml:"if_not" json:"-"`
}

type Service struct {
	Name     string `json:"name" yaml:"-"`
	URL      string `json:"url,omitempty"`
	Test     `yaml:",inline" json:",omitempty"`
	Timeout  int `json:"-"`
	Every    `yaml:",inline" json:"-"`
	Log      string            `json:"-"`
	Expect   Expect            `json:"-"`
	IfStatus map[int]Action    `yaml:"if_status" json:"-"`
	IfHeader map[string]Action `yaml:"if_header" json:"-"`
}

type Expect struct {
	Status int
	Header map[string]string
	Body   string
	body   *regexp.Regexp
	IfNot  Action `yaml:"if_not"`
}

type Action struct {
	Cmd    string
	Notify string
	Msg    string
}

func New(file string) (*Epazote, error) {
	yml_file, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var ez Epazote

	if err := yaml.Unmarshal(yml_file, &ez); err != nil {
		return nil, err
	}

	return &ez, nil
}

// CheckPaths verify that directories exist and are readable
func (self *Epazote) CheckPaths() error {
	if len(self.Config.Scan.Paths) > 0 {
		for k, d := range self.Config.Scan.Paths {
			if _, err := os.Stat(d); os.IsNotExist(err) {
				return fmt.Errorf("Verify that directory: %s, exists and is readable.", d)
			}
			if r, err := filepath.EvalSymlinks(d); err != nil {
				return err
			} else {
				self.Config.Scan.Paths[k] = r
			}
		}
		return nil
	}
	return nil
}

// VerifyUrls, we can't supervice unreachable services
func (self *Epazote) VerifyUrls() error {
	ch := AsyncGet(self.Services)
	for i := 0; i < len(self.Services); i++ {
		x := <-ch
		if x.Err != nil {
			// if not a valid URL check if service contains a test & if_not
			if len(self.Services[x.Service].Test.Test) > 0 {
				if len(self.Services[x.Service].Test.IfNot.Cmd) == 0 {
					return fmt.Errorf("%s - Verify test, missing cmd", Red(x.Service))
				}
			} else {
				return fmt.Errorf("%s - Verify URL: %q", Red(x.Service), x.Err)
			}
		}
	}
	return nil
}

// PathOrServices check if at least one path or service is set
func (self *Epazote) PathsOrServices() error {
	if len(self.Config.Scan.Paths) == 0 && len(self.Services) == 0 {
		return fmt.Errorf("%s", Red("No services to supervices or paths to scan."))
	}
	return nil
}

// GetInterval return the check interval in seconds
func GetInterval(d int, e Every) int {
	// default to 60 seconds
	if d < 1 {
		d = 60
	}

	every := d

	if e.Seconds > 0 {
		return e.Seconds
	} else if e.Minutes > 0 {
		return 60 * e.Minutes
	} else if e.Hours > 0 {
		return 3600 * e.Hours
	}

	return every
}

func ParseScan(file string) (Services, error) {
	yml_file, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var s Services

	if err := yaml.Unmarshal(yml_file, &s); err != nil {
		return nil, err
	}

	if len(s) == 0 {
		return nil, fmt.Errorf("[%s] No services found.", Red(file))
	}

	return s, nil
}
