package plugindemo

import (
	"context"
	"fmt"
	"net/http"
	"text/template"

	"github.com/ua-parser/uap-go/uaparser"
)

// Config the plugin configuration.
type Config struct {
	Headers map[string]string `json:"headers,omitempty"`
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{
		Headers: make(map[string]string),
	}
}

// Demo a Demo plugin.
type Demo struct {
	next     http.Handler
	headers  map[string]string
	name     string
	template *template.Template
}

// New created a new Demo plugin.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	if len(config.Headers) == 0 {
		return nil, fmt.Errorf("headers cannot be empty")
	}

	return &Demo{
		headers:  config.Headers,
		next:     next,
		name:     name,
		template: template.New("demo").Delims("[[", "]]"),
	}, nil
}

var MOBILE []string = []string{
	"iPhone",
	"iPod",
	"Generic Smartphone",
	"Generic Feature Phone",
	"PlayStation Vita",
	"iOS-Device",
	"Windows Phone",
	"Windows Phone OS",
	"Symbian OS",
	"Bada",
	"Windows CE",
	"Windows Mobile",
	"Maemo",
	"IE Mobile",
	"Opera Mobile",
	"Opera Mini",
	"Chrome Mobile",
	"Chrome Mobile WebView",
	"Chrome Mobile iOS",
}

var PC []string = []string{
	"Windows 95",
	"Windows 98",
	"Solaris",
	"Chrome OS",
}

var TABLET []string = []string{
	"iPad",
	"BlackBerry Playbook",
	"Blackberry Playbook",
	"Kindle",
	"Kindle Fire",
	"Kindle Fire HD",
	"Galaxy Tab",
	"Xoom",
	"Dell Streak",
	"Generic_Android_Tablet",
}

var EMAIL_PROGRAM_FAMILIES []string = []string{
	"Outlook",
	"Windows Live Mail",
	"AirMail",
	"Apple Mail",
	"Outlook",
	"Thunderbird",
	"Lightning",
	"ThunderBrowse",
	"Windows Live Mail",
	"The Bat!",
	"Lotus Notes",
	"IBM Notes",
	"Barca",
	"MailBar",
	"kmail2",
	"YahooMobileMail",
}

func getVersionRecursive(parts ...string) string {
	if len(parts) == 0 {
		return ""
	}
	var version string
	if parts[0] == "" {
		version = "0"
	} else {
		version = parts[0]
	}
	if len(parts) > 1 && parts[1] != "" {
		version += "." + getVersionRecursive(parts[1:]...)
	}
	return version
}

func (a *Demo) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	parser, _ := uaparser.NewFromBytes(uaparser.DefinitionYaml)

	ua := req.Header.Get("User-Agent")
	client := parser.Parse(ua)

	var UAVersion string
	if client.UserAgent.Major == "" {
		UAVersion = client.UserAgent.Major
	} else {
		UAVersion = getVersionRecursive(client.UserAgent.Major, client.UserAgent.Minor, client.UserAgent.Patch)
	}
	var OSVersion string
	if client.Os.Major == "" {
		OSVersion = client.Os.Major
	} else {
		OSVersion = getVersionRecursive(client.Os.Major, client.Os.Minor, client.Os.Patch, client.Os.PatchMinor)
	}

	rw.Write([]byte(fmt.Sprintf("%s|%s|%s|%s|%s|%s|%s",
		client.UserAgent.Family,
		UAVersion,
		client.Os.Family,
		OSVersion,
		client.Device.Family,
		client.Device.Brand,
		client.Device.Model,
	)))
}

