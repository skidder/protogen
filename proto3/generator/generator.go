package generator

import (
	"errors"

	"github.com/muxinc/protogen/proto3"
)

func ToProtobufSpec(spec *proto3.Spec) (string, error) {
	if spec == nil {
		return "", errors.New("Spec cannot be nil")
	}
	return spec.Write(0)
}
