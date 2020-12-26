package main

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/creekorful/mvnparser"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type dependencies struct {
	remoteDependencies    []string
	localFileDEpendencies []string
}

func main() {
	pomStr := os.Getenv("POM_LOCATION")

	//slice holding all the projects that can be a part of a main project
	var allProjects []string

	projectFolders := strings.Split(os.Getenv("PROJECT_FOLDERS"), ":")
	for _, folder := range projectFolders {
		folders, _ := ioutil.ReadDir(folder)
		for _, folder_ := range folders {
			if folder_.IsDir() {
				allProjects = append(allProjects, folder+"/"+folder_.Name())
			}
		}
	}

	localGroupIds := strings.Split(os.Getenv("POM_LOCATAL_DEPENDENCIES"), ":")
	pom, err := CalculateHashesForProjectAndSubProjects(pomStr, allProjects, localGroupIds)

	if err != nil {
		log.Fatal("can't calculate the hash: %s", err)
	} else {
		artifactHash, err := encode(pom)
		if err != nil {
			log.Fatal("can't calculate the hash: %s", err)
		}
		fmt.Println(artifactHash)
	}

}

func encode(inputs []string) (string, error) {
	// get a key for the hashes
	hasher := sha256.New()
	enc := json.NewEncoder(hasher)
	if err := enc.Encode(inputs); err != nil {
		return "", err
	}
	return hex.EncodeToString(hasher.Sum(nil)), nil
}

func CalculateHashesForProjectAndSubProjects(pomLocation string, allProjects, localGroupIds []string) ([]string, error) {
	var hashes []string
	pom, err := ioutil.ReadFile(pomLocation + "/pom.xml")
	if err != nil {
		log.Fatalf("unable to read pom file. Reason: %s", err)
		return nil, err
	}

	//add hash of a pom file
	pomHashed, _ := fileHasher(pomLocation + "/pom.xml")
	hashes = append(hashes, pomHashed)

	// Load project from string
	var project mvnparser.MavenProject
	if err := xml.Unmarshal(pom, &project); err != nil {
		log.Fatalf("unable to unmarshal pom file. Reason: %s", err)
		return nil, err
	}

	// walk src folder and hash every file
	srcHashes, err := calculateHashForSrcFolder(pomLocation)
	if err != nil {
		return nil, err
	}
	hashes = append(hashes, srcHashes...)

	// iterate over dependencies
	for _, dep := range project.Dependencies {
		if contains(localGroupIds, dep.GroupId) {
			for _, x := range allProjects {
				if strings.HasSuffix(x, dep.ArtifactId) {
					extract, err := CalculateHashesForProjectAndSubProjects(x, allProjects, localGroupIds)
					if err != nil {
						return nil, err
					}
					hashes = append(hashes, extract...)
					break
				}
			}
		}
	}

	return hashes, nil
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func fileHasher(p string) (string, error) {
	h := md5.New()
	fi, err := os.Lstat(p)
	if err != nil {
		return "", err
	}
	h.Write([]byte(fi.Mode().String()))
	h.Write([]byte(fi.Name()))
	if fi.Mode().IsRegular() {
		f, err := os.Open(p)
		if err != nil {
			return "", err
		}
		defer f.Close()
		if _, err := io.Copy(h, f); err != nil {
			return "", err
		}
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func calculateHashForSrcFolder(pomLocation string) ([]string, error) {
	var srcHashes []string
	err := filepath.Walk(pomLocation+"/src", func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			fileHash, err := fileHasher(path)
			if err != nil {
				return err
			}
			srcHashes = append(srcHashes, fileHash)
		}
		return err
	})
	sort.Strings(srcHashes)
	return srcHashes, err
}
