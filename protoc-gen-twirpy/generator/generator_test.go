package generator

import (
	"fmt"

	"testing"

	"github.com/stretchr/testify/assert"
)

func TestImportBuilder(t *testing.T) {
	// The current file being generated is in twirp/twitch/example/
	testImportBuilder := newImportBuilder(map[string]string{
		"twirp.twitch.example.Hat":   "twirp/twitch/example/haberdasher.proto",
		"twirp.twitch.example.Price": "twirp/twitch/example/haberdasher.proto",
		"twirp.twitch.example.Color": "twirp/twitch/example_folder/haberdasher.proto",
		"twirp.twitch.example.Size":  "twirp/twitch/example/haberdasher_extension.proto",
		"google.protobuf.Empty":      "google/protobuf/empty.proto",
	}, "twirp/twitch/example/service.proto")

	testCases := []struct {
		typeToImport    string
		importKey       string
		qualifiedImport string
		expectedImport  *TwirpImport
	}{
		{
			// Same package → relative import
			typeToImport:    "twirp.twitch.example.Hat",
			importKey:       "twirp.twitch.example.haberdasher_pb2",
			qualifiedImport: "_haberdasher_pb2.Hat",
			expectedImport: &TwirpImport{
				From:   ".",
				Import: "haberdasher_pb2",
				Alias:  "_haberdasher_pb2",
			},
		},
		{
			// Same package, same file → reuses existing import
			typeToImport:    "twirp.twitch.example.Price",
			importKey:       "twirp.twitch.example.haberdasher_pb2",
			qualifiedImport: "_haberdasher_pb2.Price",
			expectedImport: &TwirpImport{
				From:   ".",
				Import: "haberdasher_pb2",
				Alias:  "_haberdasher_pb2",
			},
		},
		{
			// Different package → absolute import
			typeToImport:    "twirp.twitch.example.Color",
			importKey:       "twirp.twitch.example_folder.haberdasher_pb2",
			qualifiedImport: "_haberdasher_pb2_1.Color",
			expectedImport: &TwirpImport{
				From:   "twirp.twitch.example_folder",
				Import: "haberdasher_pb2",
				Alias:  "_haberdasher_pb2_1",
			},
		},
		{
			// Same package, different file → relative import
			typeToImport:    "twirp.twitch.example.Size",
			importKey:       "twirp.twitch.example.haberdasher_extension_pb2",
			qualifiedImport: "_haberdasher_extension_pb2.Size",
			expectedImport: &TwirpImport{
				From:   ".",
				Import: "haberdasher_extension_pb2",
				Alias:  "_haberdasher_extension_pb2",
			},
		},
		{
			// External package → absolute import
			typeToImport:    "google.protobuf.Empty",
			importKey:       "google.protobuf.empty_pb2",
			qualifiedImport: "_empty_pb2.Empty",
			expectedImport: &TwirpImport{
				From:   "google.protobuf",
				Import: "empty_pb2",
				Alias:  "_empty_pb2",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Test type %s", tc.typeToImport), func(t *testing.T) {
			qualified, err := testImportBuilder.addImportAndQualify(tc.typeToImport)
			assert.NoError(t, err)
			assert.Equal(t, tc.qualifiedImport, qualified)
			assert.Equal(t, tc.expectedImport, testImportBuilder.imports[tc.importKey])
		})
	}
}

func TestImportBuilderRootLevel(t *testing.T) {
	// currentFileName has no directory component: currentFileDir must be ""
	// and root-level dependencies must produce absolute imports, not "."
	ib := newImportBuilder(map[string]string{
		"mypackage.Foo": "types.proto",
	}, "service.proto")

	qualified, err := ib.addImportAndQualify("mypackage.Foo")
	assert.NoError(t, err)
	assert.Equal(t, "_types_pb2.Foo", qualified)
	// From must be "" (absolute import), not "." (relative)
	assert.Equal(t, &TwirpImport{
		From:   "",
		Import: "types_pb2",
		Alias:  "_types_pb2",
	}, ib.imports["types_pb2"])
}

func TestImportBuilderHyphenatedDirectory(t *testing.T) {
	// Hyphens in directory names are normalised to underscores so that the
	// directory string matches the Python module path produced by getModuleName.
	ib := newImportBuilder(map[string]string{
		"my_org.my_service.Request": "my-org/my-service/types.proto",
	}, "my-org/my-service/service.proto")

	qualified, err := ib.addImportAndQualify("my_org.my_service.Request")
	assert.NoError(t, err)
	assert.Equal(t, "_types_pb2.Request", qualified)
	// Same package after hyphen normalisation → relative import
	assert.Equal(t, &TwirpImport{
		From:   ".",
		Import: "types_pb2",
		Alias:  "_types_pb2",
	}, ib.imports["my_org.my_service.types_pb2"])
}
