package main

import (
	"fmt"
	"github.com/Jumpscale/go-raml/raml"
	"log"
	"os"
)

func walkRamlResource(api *raml.Resource, url string) ([]string) {
	result := make([]string, 0)
	uri := url + api.URI
	if api.Get != nil {
		result = append(result, uri)
	}
	for _, v := range api.Nested {
		res := walkRamlResource(v, url + api.URI)
		result = append(result, res...)
	}
	return result
}

func WalkRaml(api *raml.APIDefinition) ([]string) {
	result := make([]string, 0)
	baseUri := api.BaseURI
	for _, value := range api.Resources {
		res := walkRamlResource(&value, baseUri)
		result = append(result, res...)
	}
	return result
}

func CreateCurlCommands(urls []string) ([]string) {
	result := make([]string, 0)
	for _, url := range urls {
		result = append(result, "curl " + url)
	}
	return result
}

func main() {
	args := os.Args[1:]
	if len(args) != 1 {
		log.Fatal("Wrong number of arguments!")
	}
	ramlFilePath := args[0]
	apiDef := new(raml.APIDefinition)
	err := raml.ParseFile(ramlFilePath, apiDef)
	if err != nil {
		log.Fatal(err)
	}
	urls := WalkRaml(apiDef)
	for _, v := range CreateCurlCommands(urls) {
		fmt.Println(v)
	}
}
