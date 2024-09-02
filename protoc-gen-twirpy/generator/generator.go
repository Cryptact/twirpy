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

	// Build a map of the descriptors of all .proto files indexed by fully qualified names
	protoFileDescriptors := make(map[string]*descriptor.FileDescriptorProto)
	for _, fd := range r.GetProtoFile() {
		protoFileDescriptors[fd.GetName()] = fd
	}

	for _, fileName := range r.GetFileToGenerate() {
		fd, ok := protoFileDescriptors[fileName]
		if !ok {
			resp.Error = proto.String("File[" + fileName + "][descriptor]: could not find descriptor")
			return resp
		}

		// Skip the file if there is no service in it
		if len(fd.GetService()) == 0 {
			continue
		}

		twirpFile, err := GenerateTwirpFile(fd)
		if err != nil {
			resp.Error = proto.String("File[" + fileName + "][generate]: " + err.Error())
			return resp
		}
		resp.File = append(resp.File, twirpFile)
	}
	return resp
}

func GenerateTwirpFile(fd *descriptor.FileDescriptorProto) (*plugin.CodeGeneratorResponse_File, error) {

	name := fd.GetName()

	vars := TwirpTemplateVariables{
		FileName: name,
	}

	for _, service := range fd.GetService() {
		serviceURL := fmt.Sprintf("%s.%s", fd.GetPackage(), service.GetName())
		twirpService := &TwirpService{
			Name:       service.GetName(),
			ServiceURL: serviceURL,
		}

		for _, method := range service.GetMethod() {
			twirpMethod := &TwirpMethod{
				ServiceURL:  serviceURL,
				ServiceName: twirpService.Name,
				Name:        method.GetName(),
				Input:       getSymbol(method.GetInputType()),
				Output:      getSymbol(method.GetOutputType()),
			}

			twirpService.Methods = append(twirpService.Methods, twirpMethod)
		}
		vars.Services = append(vars.Services, twirpService)
	}

	var buf = &bytes.Buffer{}
	if err := TwirpTemplate.Execute(buf, vars); err != nil {
		return nil, err
	}

	resp := &plugin.CodeGeneratorResponse_File{
		Name:    proto.String(strings.TrimSuffix(name, path.Ext(name)) + "_twirp.py"),
		Content: proto.String(buf.String()),
	}
	return resp, nil
}

func getSymbol(name string) string {
	return strings.TrimPrefix(name, ".")
}
