// +build ignore

// This is run during go generate and embeds ssh-agent-pipe into the ssh-agent-inject binary.
// * Compiles the ssh-agent-pipe binary for all target platforms (GOOS=linux, different GOARCHs).
// * Creates a .tar.gz with that binary having executable permissions set.
// * Stores that .tar.gz files in the assets/ folder as Go source files targeting the respective GOARCH.

package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"text/template"
)

var assetsTemplate = template.Must(template.New("").Funcs(template.FuncMap{
	"quote": strconv.Quote,
}).Parse(`package assets

// The .tar.gz archive containing the ssh-agent-pipe binary
const AgentArchive = {{.agent|quote}}
`))

func main() {
	os.RemoveAll("assets")
	os.MkdirAll("assets", 0700)
	archs := []string{
		"amd64",
		"arm",
		"arm64",
	}
	for _, arch := range archs {
		archive, err := buildAgentArchive(arch)
		if err != nil {
			log.Fatalln("Failed compiling ssh-agent-pipe", err)
		}
		genPlatformAssets(arch, archive)
	}
}

func buildAgentArchive(arch string) ([]byte, error) {
	os.Chdir("ssh-agent-pipe")
	defer os.Chdir("..")
	os.MkdirAll("dist", 0700)
	defer os.RemoveAll("dist")
	binaryPath := "dist/ssh-agent-pipe"
	cmd := exec.Command("go", "build", "-ldflags", "-s -w", "-o", binaryPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = append(os.Environ(), "CGO_ENABLED=0", "GOOS=linux", "GOARCH="+arch)
	if arch == "arm" {
		cmd.Env = append(cmd.Env, "GOARM=7")
	}
	if err := cmd.Run(); err != nil {
		return nil, err
	}
	contents, err := ioutil.ReadFile(binaryPath)
	if err != nil {
		return nil, err
	}
	buffer := &bytes.Buffer{}
	gw := gzip.NewWriter(buffer)
	tw := tar.NewWriter(gw)
	tw.WriteHeader(&tar.Header{
		Name: "ssh-agent-pipe",
		Size: int64(len(contents)),
		Mode: 0755,
	})
	tw.Write(contents)
	tw.Close()
	gw.Close()
	return buffer.Bytes(), nil
}

func genPlatformAssets(arch string, archive []byte) {
	file, err := os.Create(fmt.Sprintf("assets/assets_%s.go", arch))
	defer file.Close()
	if err != nil {
		log.Fatalln("Failed writing assets", err)
	}
	config := map[string]string{
		"agent": string(archive),
	}
	assetsTemplate.Execute(file, config)
}
