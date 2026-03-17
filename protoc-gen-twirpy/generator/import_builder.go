package generator

import (
	"fmt"
	"strings"
)

type importBuilder struct {
	seenAliases    map[string]struct{}
	imports        map[string]*TwirpImport
	messagesToFile map[string]string
	currentFileDir string
}

func newImportBuilder(messagesToFile map[string]string, currentFileName string) *importBuilder {
	// Extract directory from proto filename: "crypto/v1.proto" -> "crypto"
	dir := ""
	parts := strings.Split(currentFileName, "/")
	if len(parts) > 1 {
		dirParts := parts[:len(parts)-1]
		for i, p := range dirParts {
			dirParts[i] = strings.ReplaceAll(p, "-", "_")
		}
		dir = strings.Join(dirParts, ".")
	}
	return &importBuilder{
		messagesToFile: messagesToFile,
		seenAliases:    make(map[string]struct{}),
		imports:        make(map[string]*TwirpImport),
		currentFileDir: dir,
	}
}

// Generate import info and return the qualified type
// https://github.com/protocolbuffers/protobuf/blob/b8764f0941a6a5d500c48671716f0de81eb1dcaf/src/google/protobuf/compiler/python/pyi_generator.cc#L132
func (ib *importBuilder) addImportAndQualify(typeToImport string) (string, error) {
	message := getSymbol(typeToImport)

	messageSlice := strings.Split(message, ".")
	messageName := messageSlice[len(messageSlice)-1]

	if file, ok := ib.messagesToFile[message]; ok {
		moduleName := getModuleName(file)

		if _, present := ib.imports[moduleName]; !present {
			moduleNameSlice := strings.Split(moduleName, ".")
			strippedModuleName := moduleNameSlice[len(moduleNameSlice)-1]
			modulePath := strings.Join(moduleNameSlice[:len(moduleNameSlice)-1], ".")
			// Use relative import when the module is in the same package as the file being generated
			if modulePath == ib.currentFileDir {
				modulePath = "."
			}
			alias := ib.generateAlias(strippedModuleName)
			ib.imports[moduleName] = &TwirpImport{
				From:   modulePath,
				Import: strippedModuleName,
				Alias:  alias,
			}
		}
		return ib.imports[moduleName].Alias + "." + messageName, nil
	}
	return "", fmt.Errorf("cannot map message %s to a file", message)
}

// Returns the Python module name expected for a given .proto filename.
// https://github.com/protocolbuffers/protobuf/blob/b8764f0941a6a5d500c48671716f0de81eb1dcaf/src/google/protobuf/compiler/python/helpers.cc#L31-L35
func getModuleName(filename string) string {
	// std::string basename = StripProto(filename);
	basename := strings.TrimSuffix(filename, ".proto")
	//  absl::StrReplaceAll({{"-", "_"}, {"/", "."}}, &basename);
	basename = strings.ReplaceAll(basename, "-", "_")
	basename = strings.ReplaceAll(basename, "/", ".")
	//  return absl::StrCat(basename, "_pb2");
	return basename + "_pb2"
}

// https://github.com/protocolbuffers/protobuf/blob/b8764f0941a6a5d500c48671716f0de81eb1dcaf/src/google/protobuf/compiler/python/pyi_generator.cc#L139-L143
func (ib *importBuilder) generateAlias(module string) string {
	alias := "_" + module

	// Generate a unique alias by adding _1 suffixes until we get an unused alias.
	for {
		if _, present := ib.seenAliases[alias]; !present {
			break
		}
		alias = alias + "_1"
	}

	ib.seenAliases[alias] = struct{}{}
	return alias
}
