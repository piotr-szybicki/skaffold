package dockerDigest

import (
	"context"
	"github.com/GoogleContainerTools/skaffold/pkg/skaffold/build"
	"github.com/GoogleContainerTools/skaffold/pkg/skaffold/build/cache"
	cfg "github.com/GoogleContainerTools/skaffold/pkg/skaffold/config"
)

type DockerTagger struct {
	dpListener cache.DependencyLister
	artifacts  build.ArtifactGraph
	runMode    cfg.RunMode
}

func NewDockerDigestTagger(dpListener cache.DependencyLister, artifacts build.ArtifactGraph, runMode cfg.RunMode) *DockerTagger {
	return &DockerTagger{
		dpListener: dpListener,
		artifacts:  artifacts,
		runMode:    runMode,
	}
}

func (d *DockerTagger) GenerateTag(_, imageName string) (string, error) {
	ctx := context.Background()

	hasherFunc := cache.NewArtifactHasherFunc(d.artifacts, d.dpListener, d.runMode)

	return hasherFunc.Hash(ctx, d.artifacts[imageName])

}
