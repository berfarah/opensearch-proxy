package proxy

import (
	"bytes"
	"html/template"
)

const definitionTemplate = `<OpenSearchDescription xmlns="http://a9.com/-/spec/opensearch/1.1/">
<ShortName>{{.Name}}</ShortName>
<Description>{{.Description}}</Description>
<InputEncoding>UTF-8</InputEncoding>
<Image width="16" height="16" type="image/png">{{.Favicon}}</Image>
<Url rel="results" type="text/html" method="get" template="{{.Host}}/search">
  <Param name="q" value="{searchTerms}"/>
</Url>
<Url type="application/x-suggestions+json" template="{{.Host}}/suggest?q={searchTerms}"/>
</OpenSearchDescription>`

func generateDefinition(c Configuration) ([]byte, error) {
	tmpl := template.Must(template.New("defintion").Parse(definitionTemplate))
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, c); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
