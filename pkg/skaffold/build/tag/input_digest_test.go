/*
Copyright 2020 The Skaffold Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package tag

import (
	"context"
	"io/ioutil"
	"testing"

	"github.com/GoogleContainerTools/skaffold/pkg/skaffold/docker"
	"github.com/GoogleContainerTools/skaffold/pkg/skaffold/schema/latest"
	"github.com/GoogleContainerTools/skaffold/testutil"
)

func TestInputDigest_GenerateTagWhenFileDoesntExist(t *testing.T) {
	testutil.Run(t, "", func(t *testutil.T) {
		mockDependenciesForArtifact := func(ctx context.Context, a *latest.Artifact, cfg docker.Config, r docker.ArtifactResolver) ([]string, error) {
			c := []string{"imput_digest.go"}
			return c, nil
		}
		getDependenciesForArtifacet = mockDependenciesForArtifact

		tagger, _ := NewInputDigestTagger(nil, nil)

		artifact := &latest.Artifact{
			ImageName: "image_name",
		}

		tag, _ := tagger.GenerateTag("", artifact)

		t.CheckDeepEqual("38e0b9de817f645c4bec37c0d4a3e58baecccb040f5718dc069a72c7385a0bed", tag)
	})
}

func TestInputDigest_GenerateCorrectChecksumForSingleFile(t *testing.T) {
	testutil.Run(t, "", func(t *testutil.T) {
		dir := t.TempDir()
		d1 := []byte("hello\ngo\n")
		ioutil.WriteFile(dir+"/temp.file", d1, 0644)

		hash, _ := fileHasher(dir + "/temp.file")
		t.CheckDeepEqual("a5565729485faee8479f7c0760817ea7", hash)
	})
}
