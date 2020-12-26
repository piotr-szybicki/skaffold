package tag

import (
	"context"
	"fmt"
	"github.com/GoogleContainerTools/skaffold/pkg/skaffold/build/misc"
	"github.com/GoogleContainerTools/skaffold/pkg/skaffold/schema/latest"
	"github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"strings"
)

type customCommandTagger struct{}

func NewCustomCommandTagger() *customCommandTagger {
	return &customCommandTagger{}
}

func (t *customCommandTagger) GenerateTag(workspace string, image *latest.Artifact) (string, error) {
	command, err := runCommand(image.CustomTagger)
	if err != nil {
		return "", err
	}

	return command, nil
}

func runCommand(customTagger *latest.CustomTagger) (string, error) {
	ctx := context.Background()
	cmd := exec.CommandContext(ctx, "sh", "-c", customTagger.Command)
	cmd.Env = os.Environ()

	for _, envVar := range customTagger.EnvVar {
		cmd.Env = append(cmd.Env, envVar)
	}
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
