package hosts

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"time"

	"gopkg.in/yaml.v2"

	"github.com/Xuanwo/tiresias/config"
	"github.com/Xuanwo/tiresias/constants"
	"github.com/Xuanwo/tiresias/model"
	"github.com/Xuanwo/tiresias/utils"
)

const hostTemplate = `{{ .Address }} {{ .Name }}
`

// Hosts is used to update Hosts.
type Hosts struct {
	Path string `yaml:"path"`

	tmpl *template.Template
}

// Init will initiate Hosts.
func (h *Hosts) Init(c config.Endpoint) (err error) {
	// Load options
	content, err := yaml.Marshal(c.Options)
	if err != nil {
		return
	}
	err = yaml.Unmarshal(content, h)
	if err != nil {
		return
	}

	// Init template.
	h.tmpl, err = template.New("hosts").Parse(hostTemplate)
	if err != nil {
		return
	}

	// Init hosts file.
	hf, err := os.OpenFile(h.Path, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		log.Fatal(err)
	}
	defer hf.Close()
	// Seek to the start point of last update.
	cur, err := utils.Seek(hf)
	if err != nil {
		log.Fatal(err)
	}
	err = hf.Truncate(cur)
	if err != nil {
		log.Fatal(err)
	}
	_, err = hf.WriteString(fmt.Sprintf("%sGenerated by %s at %s%s\n",
		constants.CommentPrefix, constants.Name, time.Now(), constants.CommentSuffix))
	if err != nil {
		log.Fatal(err)
	}

	return
}

// Write will write servers into hosts.
func (h *Hosts) Write(s ...model.Server) (n int, err error) {
	hf, err := os.OpenFile(h.Path, os.O_RDWR|os.O_APPEND, 0600)
	if err != nil {
		log.Fatal(err)
	}
	defer hf.Close()

	for _, v := range s {
		err = h.tmpl.Execute(hf, v)
		if err != nil {
			log.Fatalf("Template generate failed for %v", err)
		}
	}

	return len(s), nil
}
