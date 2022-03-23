package version

import (
	"bytes"
	"github.com/Masterminds/sprig"
	"runtime"
	"strings"
	"text/template"
)

var (
	Version   string
	Revision  string
	BuildDate string
	GoVersion = runtime.Version()
)

var versionInfoTmpl = `
Build version: {{.version | default "N/A" }}
Build date: {{.buildDate | default "N/A" }}
Build commit: {{.revision | default "N/A" }}
`

func Print() string {
	m := map[string]string{
		"version":   Version,
		"revision":  Revision,
		"buildDate": BuildDate,
	}
	t := template.Must(template.New("version").Funcs(sprig.TxtFuncMap()).Parse(versionInfoTmpl))

	var buf bytes.Buffer
	if err := t.ExecuteTemplate(&buf, "version", m); err != nil {
		panic(err)
	}

	return strings.TrimSpace(buf.String())
}
