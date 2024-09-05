package generator

import (
	"fmt"

	"testing"

	"github.com/stretchr/testify/assert"
)

func TestImportBuilder(t *testing.T) {
	testImportBuilder := newImportBuilder(map[string]string{
		"twirp.twitch.example.Hat":   "twirp/twitch/example/haberdasher.proto",
		"twirp.twitch.example.Price": "twirp/twitch/example/haberdasher.proto",
		"twirp.twitch.example.Color": "twirp/twitch/example_folder/haberdasher.proto",
		"twirp.twitch.example.Size":  "twirp/twitch/example/haberdasher_extension.proto",
		"google.protobuf.Empty":      "google/protobuf/empty.proto",
	})

	testCases := []struct {
		typeToImport    string
		importKey       string
		qualifiedImport string
		expectedImport  *TwirpImport
	}{
		{
			typeToImport:    "twirp.twitch.example.Hat",
			importKey:       "twirp.twitch.example.haberdasher_pb2",
			qualifiedImport: "_haberdasher_pb2.Hat",
			expectedImport: &TwirpImport{
				From:   "twirp.twitch.example",
				Import: "haberdasher_pb2",
				Alias:  "_haberdasher_pb2",
			},
		},
		{
			typeToImport:    "twirp.twitch.example.Price",
			importKey:       "twirp.twitch.example.haberdasher_pb2",
			qualifiedImport: "_haberdasher_pb2.Price",
			expectedImport: &TwirpImport{
				From:   "twirp.twitch.example",
				Import: "haberdasher_pb2",
				Alias:  "_haberdasher_pb2",
			},
		},
		{
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
			typeToImport:    "twirp.twitch.example.Size",
			importKey:       "twirp.twitch.example.haberdasher_extension_pb2",
			qualifiedImport: "_haberdasher_extension_pb2.Size",
			expectedImport: &TwirpImport{
				From:   "twirp.twitch.example",
				Import: "haberdasher_extension_pb2",
				Alias:  "_haberdasher_extension_pb2",
			},
		},
		{
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
