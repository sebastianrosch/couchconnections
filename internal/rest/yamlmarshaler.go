package rest

import (
	"bytes"
	"io"

	yaml "github.com/ghodss/yaml"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
)

// YamlMarshaler is used to marshal and unmarshal yaml content.
type YamlMarshaler struct {
}

type yamlDecoder struct {
	reader io.Reader
}

func (y *yamlDecoder) Decode(v interface{}) error {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(y.reader)
	if err != nil {
		return err
	}

	json, err := yaml.YAMLToJSON(buf.Bytes())
	if err != nil {
		return err
	}

	n := runtime.JSONPb{}
	return n.Unmarshal(json, v)
}

type yamlEncoder struct {
	writer io.Writer
}

func (y *yamlEncoder) Encode(v interface{}) error {
	jsonpb := runtime.JSONPb{}
	json, err := jsonpb.Marshal(v)
	if err != nil {
		return err
	}
	yaml, err := yaml.JSONToYAML(json)
	if err != nil {
		return err
	}
	_, err = y.writer.Write(yaml)
	if err != nil {
		return err
	}
	return nil
}

// Marshal marshals "v" into byte sequence.
func (y *YamlMarshaler) Marshal(v interface{}) ([]byte, error) {
	jsonpb := runtime.JSONPb{
		EmitDefaults: true,
	}
	json, err := jsonpb.Marshal(v)
	if err != nil {
		return nil, err
	}
	return yaml.JSONToYAML(json)
}

// Unmarshal unmarshals "data" into "v".
// "v" must be a pointer value.
func (y *YamlMarshaler) Unmarshal(data []byte, v interface{}) error {
	json, _ := yaml.YAMLToJSON(data)
	n := runtime.JSONPb{}
	return n.Unmarshal(json, v)
}

// NewDecoder returns a Decoder which reads byte sequence from "r".
func (y *YamlMarshaler) NewDecoder(r io.Reader) runtime.Decoder {
	return &yamlDecoder{reader: r}
}

// NewEncoder returns an Encoder which writes bytes sequence into "w".
func (y *YamlMarshaler) NewEncoder(w io.Writer) runtime.Encoder {
	return &yamlEncoder{writer: w}
}

// ContentType returns the Content-Type which this marshaler is responsible for.
func (y *YamlMarshaler) ContentType() string {
	return "application/yaml"
}
