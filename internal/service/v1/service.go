package service

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	buildinfo "github.com/sebastianrosch/livingroompresentations/pkg/build-info"
	v1 "github.com/sebastianrosch/livingroompresentations/rpc/livingroom-api/v1"
)

type LivingRoomService struct {
}

// ------------------
// Utility endpoints.
// ------------------

// GetVersion returns the service version.
func (s *LivingRoomService) GetVersion(ctx context.Context, req *empty.Empty) (*v1.Version, error) {
	buildInfo := buildinfo.NewDefaultBuildInfo()

	return &v1.Version{
		Version:  buildInfo.Version,
		Branch:   buildInfo.Branch,
		Revision: buildInfo.Revision,
	}, nil
}
