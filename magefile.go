//go:build mage
// +build mage

package main

import (
	"fmt"
	"go/build"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/openimsdk/gomake/mageutil"
)

var Default = Build

func Build() {
	platforms := os.Getenv("PLATFORMS")
	if platforms == "" {
		platforms = mageutil.DetectPlatform()
	}

	for _, platform := range strings.Split(platforms, " ") {
		mageutil.CompileForPlatform(platform)
	}

	mageutil.PrintGreen("All binaries under cmd and tools were successfully compiled.")
}

func Start() {
	mageutil.InitForSSC()
	err := setMaxOpenFiles()
	if err != nil {
		mageutil.PrintRed("setMaxOpenFiles failed " + err.Error())
		os.Exit(1)
	}
	mageutil.StartToolsAndServices()
}

func Stop() {
	mageutil.StopAndCheckBinaries()
}

func Check() {
	mageutil.CheckAndReportBinariesStatus()
}

func getWorkDirToolPath(name string) string {
	toolPath := ""
	workDir, err := os.Getwd()
	if err != nil {
		log.Println("Error", err)
		return toolPath
	}
	toolsPath := filepath.Join(workDir, "tools")
	filepath.Walk(toolsPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if strings.TrimSuffix(info.Name(), filepath.Ext(info.Name())) == name {
			toolPath = path
		}
		return nil
	})

	return toolPath
}

// ------------------

var protoModules = []string{
	"admin",
	"chat",
	"common",
	"wrapperspb",
	"sdkws",
}

// install proto plugin
func InstallDepend() error {
	log.SetOutput(os.Stdout)
	// log.SetFlags(log.Lshortfile)

	// log.Println("installing protoc-gen-go and protoc-gen-go-grpc")
	log.Println("installing protobuf dependencies in Go.")

	cmds := [][]string{
		{"install", "google.golang.org/protobuf/cmd/protoc-gen-go@latest"},
		{"install", "google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest"},
	}

	for _, cmdArgs := range cmds {
		cmd := exec.Command("go", cmdArgs...)

		// log.Println("running command:", "go", cmdArgs)
		connectStd(cmd)

		if err := cmd.Run(); err != nil {
			log.Printf("command %v error: %v", cmdArgs, err)
			return err
		}
	}

	return nil
}

// Generate Go code from protobuf files.
func GenGo() error {
	log.SetOutput(os.Stdout)
	// log.SetFlags(log.Lshortfile)
	log.Println("Generating Go code from proto files")

	// goOutDir := filepath.Join(protoDir, GO)

	protoc, err := getToolPath("protoc")
	if err != nil {
		return err
	}

	for _, module := range protoModules {
		args := []string{
			"--proto_path=./pkg/protocol",
			"--proto_path=./",
			"--go_out=" + filepath.Join("./pkg/protocol", module),
			"--go-grpc_out=" + filepath.Join("./pkg/protocol", module),
			"--go_opt=module=github.com/openimsdk/chat/pkg/protocol/" + strings.Join([]string{module}, "/"),
			"--go-grpc_opt=module=github.com/openimsdk/chat/pkg/protocol/" + strings.Join([]string{module}, "/"),
			filepath.Join("./pkg/protocol", module, module) + ".proto",
		}
		// log.Println("protoc args", args)

		cmd := exec.Command(protoc, args...)
		connectStd(cmd)

		if err := cmd.Run(); err != nil {
			log.Printf("Error generating Go code for module %s: %v\n", module, err)
			continue
		}
	}

	if err := removeOmitemptyTags(); err != nil {
		log.Println("Remove Omitempty is Error", err)
		return err
	} else {
		log.Println("Remove Omitempty is Success")
	}

	return nil
}

func connectStd(cmd *exec.Cmd) {
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
}

func getToolPath(name string) (string, error) {
	// Get in work dir.
	toolPath := getWorkDirToolPath(name)
	if toolPath != "" {
		return toolPath, nil
	}

	// Get in env path.
	if p, err := exec.LookPath(name); err == nil {
		return p, nil
	}

	// check under gopath
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		gopath = build.Default.GOPATH
	}
	p := filepath.Join(gopath, "bin", name)

	if _, err := os.Stat(p); err != nil {
		return "", err
	}

	return p, nil
}

func removeOmitemptyTags() error {
	// protoGoDir := filepath.Join(protoDir, GO) // "./proto/go"

	re := regexp.MustCompile(`,\s*omitempty`)

	return filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println("access path error:", err)
			return err
		}
		if !info.IsDir() && strings.HasSuffix(path, ".pb.go") {
			input, err := os.ReadFile(path)
			if err != nil {
				fmt.Println("ReadFile error. Path: %s, Error %v", path, err)
				return err
			}

			output := re.ReplaceAllString(string(input), "")

			// check replace is happened
			if string(input) != output {
				err = os.WriteFile(path, []byte(output), info.Mode())
				if err != nil {
					fmt.Printf("Error writing file: %s, error: %v\n", path, err)
					return err
				}
				// fmt.Println("Modified file:", path)
			}
		}

		return nil
	})
}
