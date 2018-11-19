package proxy

import (
	"bytes"
	"html/template"
)

const rootTemplate = `<!DOCTYPE html><html>
		<head>
<link type="application/opensearchdescription+xml" rel="search" href="{{.Host}}/search.xml" title="{{.Name}}">
<link rel="icon" type="image/x-icon" href="{{.Favicon}}" />
	  </head>
	  <body>Add your the search engine by right clicking the URL bar and looking for the option above!</body></html>`

func generateRoot(c Configuration) ([]byte, error) {
	tmpl := template.Must(template.New("root").Parse(rootTemplate))
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, c); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
