package main

import (
	"context"
	"fmt"
	"reflect"

	docopt "github.com/docopt/docopt-go"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	gogo_proto "github.com/gogo/protobuf/proto"
	_ "github.com/golang/protobuf/ptypes/wrappers"
	yaml "gopkg.in/yaml.v2"


	_ "github.com/envoyproxy/go-control-plane/envoy/admin/v2alpha"
	_ "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	_ "github.com/envoyproxy/go-control-plane/envoy/config/bootstrap/v2"

	"github.com/stripe/skycfg"
)

type protoRegistry struct{}

func (*protoRegistry) UnstableProtoMessageType(name string) (reflect.Type, error) {
	if t := proto.MessageType(name); t != nil {
		return t, nil
	}
	if t := gogo_proto.MessageType(name); t != nil {
		return t, nil
	}
	return nil, nil
}

func (*protoRegistry) UnstableEnumValueMap(name string) map[string]int32 {
	if ev := proto.EnumValueMap(name); ev != nil {
		return ev
	}
	if ev := gogo_proto.EnumValueMap(name); ev != nil {
		return ev
	}
	return nil
}

func main() {

	usage := `envoy-proto: creates an appoptics-style envoy configuration from a skycfg file

Usage:
  envoy-proto <input-file>
  envoy-proto -h | --help

Options:
  <input-file>  skycfg file to be read
  -h --help     Show this screen.`

	arguments, _ := docopt.ParseDoc(usage)

	ctx := context.Background()
	config, err := skycfg.Load(ctx, arguments["<input-file>"].(string), skycfg.WithProtoRegistry(&protoRegistry{}))
	if err != nil {
		panic(err)
	}
	messages, err := config.Main(ctx)
	if err != nil {
		panic(err)
	}
	for _, msg := range messages {
		var jsonMarshaler = &jsonpb.Marshaler{OrigName: true}

		marshaled, err := jsonMarshaler.MarshalToString(msg)
		sep := ""
		var yamlMap yaml.MapSlice
		if err := yaml.Unmarshal([]byte(marshaled), &yamlMap); err != nil {
			panic(fmt.Sprintf("yaml.Unmarshal: %v", err))
		}
		yamlMarshaled, err := yaml.Marshal(yamlMap)
		if err != nil {
			panic(fmt.Sprintf("yaml.Marshal: %v", err))
		}
		marshaled = string(yamlMarshaled)
		sep = "---\n"
		fmt.Printf("%s%s\n", sep, marshaled)
	}
}
