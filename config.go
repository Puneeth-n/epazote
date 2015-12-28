package epazote

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type Epazote struct {
	Config   Config
	Services map[string]Service
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

type Scan struct {
	Paths                   []string
	Seconds, Minutes, Hours int
}

type Service struct {
	URL                     string
	Timeout                 int
	Seconds, Minutes, Hours int
	Log                     string
	Expect                  Expect
	IfStatus                map[string]Action `yaml:"if_status`
	IfHeader                map[string]Action `yaml:"if_header"`
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
		log.Println("No services found.")
		return nil
	}

	// add services to supervisor
	for k, v := range s {
		if !IsURL(v.URL) {
			log.Printf("[%s] %s - Verify URL: %q", Red(file), k, v.URL)
			continue
		}

		// how often to check for the service
		every := 60
		if v.Seconds > 0 {
			every = v.Seconds
		} else if v.Minutes > 0 {
			every = 60 * v.Minutes
		} else if v.Hours > 0 {
			every = 3600 * v.Hours
		}
		log.Println(every)
	}
	sk := GetScheduler()
	for k, v := range sk.Schedulers {
		log.Println(Red(k), v)
	}

	return nil
}
