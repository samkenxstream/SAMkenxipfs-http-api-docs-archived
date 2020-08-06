package docs

import (
	"bytes"
	"fmt"
	"html"
	"regexp"
	"strings"
)

// Removes the "Default:..." part in the descriptions.
var fixDesc, _ = regexp.Compile(" Default: [a-zA-z0-9-_]+ ?\\.")

func genArgument(arg *Argument, aliasToArg bool) string {
	// These get handled by GenerateBodyBlock
	if arg.Type == "file" {
		return "\n"
	}

	buf := new(bytes.Buffer)
	alias := arg.Name
	if aliasToArg {
		alias = "arg"
	}

	fixedDescription := string(fixDesc.ReplaceAll([]byte(arg.Description), []byte("")))
	fixedDescription = html.EscapeString(fixedDescription)

	fmt.Fprintf(buf, "- `%s` [%s]: %s", alias, arg.Type, fixedDescription)
	if len(arg.Default) > 0 {
		fmt.Fprintf(buf, " Default: `%s`.", arg.Default)
	}
	if arg.Required {
		fmt.Fprintf(buf, ` Required: **yes**.`)
	} else {
		fmt.Fprintf(buf, ` Required: no.`)
	}
	fmt.Fprintln(buf)
	return buf.String()
}

func genParameter(arg *Argument, aliasToArg bool) string {
	// These get handled by GenerateBodyBlock
	if arg.Type == "file" {
		return "\n"
	}

	buf := new(bytes.Buffer)
	alias := arg.Name
	if aliasToArg {
		alias = "arg"
	}

	fixedDescription := string(fixDesc.ReplaceAll([]byte(arg.Description), []byte("")))
	fixedDescription = html.EscapeString(fixedDescription)

	//  TODO - Is  there a better way to handle this?
	// Some descriptions are coming back with colons in it which breaks the yaml
	fixedDescription = strings.Replace(fixedDescription, ":", "", -1)

	// This translates the types returned from go-ipfs into yaml format
	yamlTypes := map[string]string{
		"bool": "boolean", "int": "integer", "int64": "integer", "string": "string",
	}

	fmt.Fprintf(buf, `
          - name: %s
            in: query
            description: %s
            schema:
              type: %s
    `, alias, fixedDescription, yamlTypes[arg.Type])

	if len(arg.Default) > 0 {
		fmt.Fprintf(buf, `
              default: %s
        `, arg.Default)
	}

	if arg.Required {
		fmt.Fprintf(buf, `
            required: true
         `)
	} else {
		fmt.Fprintf(buf, `
            required: false
        `)
	}

	return buf.String()
}
