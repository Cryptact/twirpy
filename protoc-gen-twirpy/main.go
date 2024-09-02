package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/verloop/twirpy/protoc-gen-twirpy/generator"
	"google.golang.org/protobuf/proto"
	plugin "google.golang.org/protobuf/types/pluginpb"
)

func buildCodeGeneratorRequest(r io.Reader) (*plugin.CodeGeneratorRequest, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("could not read from stdin: %w", err)
	}
	var req = &plugin.CodeGeneratorRequest{}
	if err = proto.Unmarshal(data, req); err != nil {
		return nil, fmt.Errorf("could not unmarshal proto: %w", err)
	}
	if len(req.GetFileToGenerate()) == 0 {
		return nil, fmt.Errorf("no files to generate")
	}
	return req, nil
}

func writeCodeGeneratorResponse(w io.Writer, resp *plugin.CodeGeneratorResponse) error {
	data, err := proto.Marshal(resp)
	if err != nil {
		return fmt.Errorf("could not unmarshal response proto: %w", err)
	}
	_, err = w.Write(data)
	return err
}

func main() {
	req, err := buildCodeGeneratorRequest(os.Stdin)
	if err != nil {
		log.Fatalln(err)
	}

	resp := generator.Generate(req)
	if resp == nil {
		resp = &plugin.CodeGeneratorResponse{}
	}

	err = writeCodeGeneratorResponse(os.Stdout, resp)
	if err != nil {
		log.Fatalln(err)
	}
}
