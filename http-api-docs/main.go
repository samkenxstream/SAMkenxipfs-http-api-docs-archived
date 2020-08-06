// This is an utility to generate documentation from go-ipfs commands
package main

import (
	"fmt"

	docs "github.com/ipfs/http-api-docs"
)

func main() {
	endpoints := docs.AllEndpoints()
	formatter := new(docs.OpenAPIFormatter)
	fmt.Println(docs.GenerateDocs(endpoints, formatter))
}
