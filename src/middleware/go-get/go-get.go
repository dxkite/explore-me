package goget

import (
	_ "embed"
	"io"
	"log"
	"net/http"
	"strings"
	"text/template"
)

//go:embed template.html
var tmpHtml string

type Package struct {
	Name string `yaml:"name"`
	Repo string `yaml:"repo"`
	Doc  string `yaml:"doc"`
}

type PackageConfig struct {
	Spec    []Package `yaml:"spec"`
	Default Package   `yaml:"default"`
}

var tmp *template.Template

func init() {
	if t, err := template.New("go-get").Parse(tmpHtml); err != nil {
		panic(err)
	} else {
		tmp = t
	}
}

func render(w io.Writer, data Package) error {
	return tmp.Execute(w, data)
}

func Middleware(fn func() *PackageConfig, handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("go-get", r.URL.String())
		cfg := fn()

		goGet := r.URL.Query().Get("go-get") == "1"
		if !goGet {
			handler.ServeHTTP(w, r)
			return
		}

		sp := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		name := sp[0]

		for _, pkg := range cfg.Spec {
			if pkg.Name == name {
				if err := render(w, pkg); err != nil {
					log.Println(err)
				}
				return
			}
		}

		dft := cfg.Default
		dft.Name = name
		dft.Repo = strings.ReplaceAll(dft.Repo, "{name}", name)
		dft.Doc = strings.ReplaceAll(dft.Doc, "{name}", name)

		if err := render(w, dft); err != nil {
			log.Println(err)
		}
	})
}
