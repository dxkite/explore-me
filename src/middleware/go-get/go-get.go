package goget

import (
	_ "embed"
	"io"
	"net/http"
	"strings"
	"text/template"

	"dxkite.cn/log"
)

//go:embed template.html
var tmpHtml string

type Package struct {
	Path   string `yaml:"path"`
	Import string `yaml:"import"`
	Repo   string `yaml:"repo"`
	Doc    string `yaml:"doc"`
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

		path := strings.Trim(r.URL.Path, "/")

		for _, pkg := range cfg.Spec {
			if strings.HasPrefix(pkg.Path, path) {
				if err := render(w, pkg); err != nil {
					log.Println(err)
				}
				return
			}
		}

		dft := cfg.Default
		dft.Path = strings.ReplaceAll(dft.Path, "{path}", path)
		dft.Import = strings.ReplaceAll(dft.Import, "{path}", path)
		dft.Repo = strings.ReplaceAll(dft.Repo, "{path}", path)
		dft.Doc = strings.ReplaceAll(dft.Doc, "{path}", path)

		if err := render(w, dft); err != nil {
			log.Println(err)
		}
	})
}
