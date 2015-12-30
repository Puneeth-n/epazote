package epazote

import (
	"fmt"
	"github.com/nbari/epazote/scheduler"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const herb = "\U0001f33f"

type Epazote struct {
	Config   Config
	Services Services
}

type Config struct {
	SMTP Email `yaml:"smtp"`
	HTTP Http  `yaml:"http"`
	Scan Scan  `yaml:"scan"`
}

type Email struct {
	Username string
	Password string
	Host     string
	Port     int
	Tls      bool
	Headers  map[string]string
}

type Http struct {
	Host     string
	pathPort int
}

type Every struct {
	Seconds, Minutes, Hours int
}

type Scan struct {
	Paths []string
	Every
}

type Services map[string]Service

type Service struct {
	URL     string
	Timeout int
	Every
	Log      string
	Expect   Expect
	IfStatus map[string]Action `yaml:"if_status`
	IfHeader map[string]Action `yaml:"if_header"`
}

type Expect struct {
	Status int
	Header map[string]string
	Body   string
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
			if r, err := filepath.EvalSymlinks(d); err != nil {
				return err
			} else {
				if _, err := os.Stat(r); os.IsNotExist(err) {
					return fmt.Errorf("Verify that directory: %s, exists and is readable.", r)
				}
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
			return fmt.Errorf("%s - Verify URL: %q", Red(x.Service), x.Err)
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

// Start Add services to scheduler
func (self *Epazote) Start(sk *scheduler.Scheduler) string {
	for k, v := range self.Services {
		// schedule service
		sk.AddScheduler(k, GetInterval(60, v.Every), Supervice(v))
	}

	if len(self.Config.Scan.Paths) > 0 {
		s := new(Scandir)
		for _, v := range self.Config.Scan.Paths {
			sk.AddScheduler(v, GetInterval(300, self.Config.Scan.Every), s.Scan(v))
		}
	}

	return fmt.Sprintf("Epazote %s   on %d services, scan paths: %s", herb, len(self.Services), strings.Join(self.Config.Scan.Paths, ","))
}

// GetInterval return the check interval in seconds
func GetInterval(d int, s Every) int {
	// default to 60 seconds
	if d == 0 {
		d = 60
	}
	every := d

	if s.Seconds > 0 {
		every = s.Seconds
	} else if s.Minutes > 0 {
		every = 60 * s.Minutes
	} else if s.Hours > 0 {
		every = 3600 * s.Hours
	}

	return every
}

func ParseScan(file string) error {
	yml_file, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	var s Services

	if err := yaml.Unmarshal(yml_file, &s); err != nil {
		return err
	}

	if len(s) == 0 {
		log.Printf("[%s] No services found.", Red(file))
		return nil
	}

	// get a Scheduler
	sk := GetScheduler()

	// add/update services to supervisor
	for k, v := range s {
		if !IsURL(v.URL) {
			log.Printf("[%s] %s - Verify URL: %q", Red(file), k, v.URL)
			continue
		}

		// schedule service
		sk.AddScheduler(k, GetInterval(60, v.Every), Supervice(v))
	}

	return nil
}
