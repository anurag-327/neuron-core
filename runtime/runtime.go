package runtime

import (
	"fmt"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-units"
)

type EntryFile struct {
	FileName  string
	Extension string
}

type FileNames struct {
	FileName string // "main" or "Main"
	FullName string // "main.cpp"
	PathBase string // "/host/job/main"
	PathFull string // "/host/job/main.cpp"
}

type ResourceLimits struct {
	MemoryKB    int64
	TimeMs      time.Duration
	NanoCPUs    int64
	Pids        int64
	ULimits     []*units.Ulimit
	NetworkMode container.NetworkMode
}

type RuntimeConfig struct {
	Language       string
	Image          string
	InitSize       int
	MaxSize        int
	HealthCmd      []string
	HealthInterval time.Duration
	ResourceLimits ResourceLimits
	Validator      func(code string) error

	EntryFile  EntryFile
	CompileCmd func(n FileNames) string
	RunCmd     func(n FileNames) string
}

var commonResouceLimits = ResourceLimits{
	MemoryKB: 256 * 1024,
	TimeMs:   5 * time.Second,
	NanoCPUs: 1000000000,
	Pids:     100,
	ULimits: []*units.Ulimit{
		{Name: "nofile", Soft: 64, Hard: 64},
	},
	NetworkMode: "none",
}

var LanguageRegistry = map[string]RuntimeConfig{
	"cpp": {
		Language:       "cpp",
		Image:          "neuron-cpp",
		InitSize:       2,
		MaxSize:        4,
		HealthCmd:      []string{"echo", "ok"},
		ResourceLimits: commonResouceLimits,
		Validator:      ValidateAndSanitizeCpp,
		EntryFile: EntryFile{
			FileName:  "main",
			Extension: "cpp",
		},

		CompileCmd: func(n FileNames) string {
			return fmt.Sprintf("g++ -O2 -pipe -std=gnu++20 %s -o %s", n.FullName, n.FileName)
		},
		RunCmd: func(n FileNames) string {
			return fmt.Sprintf("./%s < input.txt", n.FileName)
		},
	},
	"python": {
		Language:       "python",
		Image:          "neuron-python",
		InitSize:       2,
		MaxSize:        4,
		HealthCmd:      []string{"echo", "ok"},
		ResourceLimits: commonResouceLimits,
		//		Validator:      ValidateAndSanitizePython,
		EntryFile: EntryFile{
			FileName:  "main",
			Extension: "py",
		},

		CompileCmd: nil,
		RunCmd: func(n FileNames) string {
			return fmt.Sprintf("python3 %s < input.txt", n.FullName)
		},
	},
	"javascript": {
		Language:       "javascript",
		Image:          "neuron-node",
		InitSize:       2,
		MaxSize:        4,
		HealthCmd:      []string{"echo", "ok"},
		ResourceLimits: commonResouceLimits,
		Validator:      ValidateAndSanitizeJS,
		EntryFile: EntryFile{
			FileName:  "main",
			Extension: "js",
		},

		CompileCmd: nil,
		RunCmd: func(n FileNames) string {
			return fmt.Sprintf("node %s < input.txt", n.FullName)
		},
	},
	"java": {
		Language:       "java",
		Image:          "neuron-java",
		InitSize:       1,
		MaxSize:        4,
		HealthCmd:      []string{"echo", "ok"},
		ResourceLimits: commonResouceLimits,
		Validator:      ValidateAndSanitizeJava,
		EntryFile: EntryFile{
			FileName:  "Main",
			Extension: "java",
		},

		CompileCmd: func(n FileNames) string {
			return fmt.Sprintf("javac %s", n.FullName)
		},
		RunCmd: func(n FileNames) string {
			// e.g. java Main < input.txt
			// n.FileName is "Main" (if EntryFile is Main.java)
			return fmt.Sprintf("java %s < input.txt", n.FileName)
		},
	},
}
