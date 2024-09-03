package generator

import (
	"bytes"
	"fmt"
	"path"
	"strings"

	"google.golang.org/protobuf/proto"
	descriptor "google.golang.org/protobuf/types/descriptorpb"
	plugin "google.golang.org/protobuf/types/pluginpb"
)

func Generate(r *plugin.CodeGeneratorRequest) *plugin.CodeGeneratorResponse {
	resp := &plugin.CodeGeneratorResponse{}
	resp.SupportedFeatures = proto.Uint64(uint64(plugin.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL))

	responseFiles, err := generateFiles(r)
	if err != nil {
		resp.Error = proto.String(err.Error())
	}
	resp.File = responseFiles

	return resp
}

func generateFiles(r *plugin.CodeGeneratorRequest) ([]*plugin.CodeGeneratorResponse_File, error) {
	responseFiles := []*plugin.CodeGeneratorResponse_File{}

	protoFiles := r.GetProtoFile()

	// Build a map of the descriptors of all .proto files indexed by fully qualified names
	protoFileDescriptors := make(map[string]*descriptor.FileDescriptorProto)
	for _, fd := range protoFiles {
		protoFileDescriptors[fd.GetName()] = fd
	}

	messagesToFiles := buildMessagesToFiles(protoFiles)

	for _, fileName := range r.GetFileToGenerate() {
		fd, ok := protoFileDescriptors[fileName]
		if !ok {
			return nil, fmt.Errorf("File[%s][descriptor]: could not find descriptor", fileName)
		}

		// Skip the file if there is no service in it
		if len(fd.GetService()) == 0 {
			continue
		}

		templateVars, err := buildTwirpServiceDescription(messagesToFiles, fd)
		if err != nil {
			return nil, fmt.Errorf("File[%s][descriptor]: %w", fileName, err)
		}

		twirpFile, err := generateTwirpFile(templateVars)
		if err != nil {
			return nil, fmt.Errorf("File[%s][descriptor]: %w", fileName, err)
		}
		responseFiles = append(responseFiles, twirpFile)
	}

	return responseFiles, nil
}

func buildTwirpServiceDescription(messagesToFiles map[string]string, fd *descriptor.FileDescriptorProto) (*ProtoFileDescription, error) {
	name := fd.GetName()

	vars := &ProtoFileDescription{
		FileName: name,
	}

	imports := newImportBuilder(messagesToFiles)

	for _, service := range fd.GetService() {
		serviceURL := fmt.Sprintf("%s.%s", fd.GetPackage(), service.GetName())
		twirpService := &TwirpService{
			Name:       service.GetName(),
			ServiceURL: serviceURL,
		}

		for _, method := range service.GetMethod() {
			qualifiedInput, err := imports.addImportAndQualify(method.GetInputType())
			if err != nil {
				return nil, err
			}

			qualifiedOutput, err := imports.addImportAndQualify(method.GetOutputType())
			if err != nil {
				return nil, err
			}

			twirpMethod := &TwirpMethod{
				ServiceURL:      serviceURL,
				ServiceName:     twirpService.Name,
				Name:            method.GetName(),
				Input:           getSymbol(method.GetInputType()),
				Output:          getSymbol(method.GetOutputType()),
				QualifiedInput:  qualifiedInput,
				QualifiedOutput: qualifiedOutput,
			}

			twirpService.Methods = append(twirpService.Methods, twirpMethod)
		}
		vars.Services = append(vars.Services, twirpService)
	}

	for _, importStmt := range imports.imports {
		vars.Imports = append(vars.Imports, importStmt)
	}

	return vars, nil
}

func generateTwirpFile(vars *ProtoFileDescription) (*plugin.CodeGeneratorResponse_File, error) {
	var buf = &bytes.Buffer{}
	if err := TwirpTemplate.Execute(buf, vars); err != nil {
		return nil, err
	}

	return &plugin.CodeGeneratorResponse_File{
		Name:    proto.String(strings.TrimSuffix(vars.FileName, path.Ext(vars.FileName)) + "_twirp.py"),
		Content: proto.String(buf.String()),
	}, nil
}

func getSymbol(name string) string {
	return strings.TrimPrefix(name, ".")
}

func buildMessagesToFiles(fds []*descriptor.FileDescriptorProto) map[string]string {
	mapOut := make(map[string]string)
	for _, fd := range fds {
		for _, msg := range fd.GetMessageType() {
			mapOut[fd.GetPackage()+"."+msg.GetName()] = fd.GetName()
		}
	}
	return mapOut
}
