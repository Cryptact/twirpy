package generator

import (
	"fmt"
	"strings"
)

type importBuilder struct {
	seenAliases    map[string]struct{}
	aliasMappings  map[string]string
	imports        map[string]*TwirpImport
	messagesToFile map[string]string
}

func newImportBuilder(messagesToFile map[string]string) *importBuilder {
	return &importBuilder{
		messagesToFile: messagesToFile,
		seenAliases:    make(map[string]struct{}),
		aliasMappings:  make(map[string]string),
		imports:        make(map[string]*TwirpImport),
	}
}

func (ib *importBuilder) addImportAndQualify(typeToImport string) (string, error) {
	message := getSymbol(typeToImport)

	if file, ok := ib.messagesToFile[message]; ok {
		pyFile := asPythonPath(file)

		pathSlice := strings.Split(pyFile, ".")
		path := strings.Join(pathSlice[:len(pathSlice)-1], ".")
		module := pathSlice[len(pathSlice)-1]
		if _, present := ib.imports[pyFile]; !present {
			alias := ib.generateAlias(module)
			ib.aliasMappings[pyFile] = alias
			ib.imports[pyFile] = &TwirpImport{
				From:   path,
				Import: module,
				Alias:  alias,
			}
		}
		return ib.aliasMappings[pyFile] + "." + getMessageName(message), nil
	}
	return "", fmt.Errorf("cannot map message %s to a file", message)
}

func (ib *importBuilder) generateAlias(module string) string {
	alias := "_" + module

	for {
		if _, present := ib.seenAliases[alias]; !present {
			break
		}
		alias = "_" + alias
	}

	ib.seenAliases[alias] = struct{}{}
	return alias
}

func getMessageName(name string) string {
	parts := strings.Split(name, ".")
	return parts[len(parts)-1]
}

func asPythonPath(fileName string) string {
	asPath := strings.ReplaceAll(fileName, "/", ".")
	return strings.Replace(asPath, ".proto", "_pb2", 1)
}
