package dockerDigest

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/GoogleContainerTools/skaffold/pkg/skaffold/build"
	"github.com/GoogleContainerTools/skaffold/pkg/skaffold/build/misc"
	"github.com/GoogleContainerTools/skaffold/pkg/skaffold/schema/latest"
	"github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"strings"
)

type customCommandTagger struct {
	tagCommand string
	artifacts  build.ArtifactGraph
}

type CommandContext struct {
	dependencies []*latest.Artifact
	artifact     *latest.Artifact
}

func NewCustomScriptTagger(tagCommand string, artifacts build.ArtifactGraph) *customCommandTagger {
	return &customCommandTagger{
		tagCommand: tagCommand,
		artifacts:  artifacts,
	}
}

func (c *customCommandTagger) GenerateTag(workingDir, imageName string) (string, error) {
	artifact := c.artifacts[imageName]

	dependenciesJson, err := json.Marshal(artifact)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return "", fmt.Errorf("cant marst the artifact dependencies information")
	}

	return runCommand(c.tagCommand, workingDir, string(dependenciesJson))
}

func runCommand(tagCommand string, workingDir, imageJson string) (string, error) {
	ctx := context.Background()
	cmd := exec.CommandContext(ctx, "sh", "-c", tagCommand)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "CONTEXT="+imageJson)
	cmd.Stderr = os.Stderr

	logrus.Infof("Running command: %s", cmd.Args)
	if out, err := cmd.Output(); err != nil {
		return "", fmt.Errorf("starting cmd: %w", err)
	} else {
		misc.HandleGracefulTermination(ctx, cmd)
		fixedString := strings.TrimSuffix(string(out), "\n")
		return fixedString, nil
	}

}
