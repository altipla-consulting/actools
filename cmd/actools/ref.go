package main

var containers = []string{
	"dev-go",
	"gcloud",
	"glide",
	"go",
	"gulp",
	"mysql",
	"node",
	"phpmyadmin",
	"protoc",
	"rambler",
	"sass-lint",
}

type toolConfig struct {
	Container string
	Cnf       *containerConfig
}

// TODO(ernesto): Estas configuraciones son más por contenedor que por herramienta,
// podemos llevar la lista abajo a manualRun y compartirla con un mejor nombre de
// variable.
var tools = map[string]*toolConfig{
	"gcloud": &toolConfig{
		Container: "gcloud",
		Cnf: &containerConfig{
			ShareWorkspace:    true,
			LocalUser:         true,
			ShareSSHSocket:    true,
			ShareGcloudConfig: true,
		},
	},
	"glide": &toolConfig{
		Container: "glide",
		Cnf: &containerConfig{
			ShareWorkspace:  true,
			LocalUser:       true,
			ShareSSHSocket:  true,
			ConfigureGopath: true,
		},
	},
	"go": &toolConfig{
		Container: "go",
		Cnf: &containerConfig{
			ShareWorkspace:  true,
			LocalUser:       true,
			ConfigureGopath: true,
		},
	},
	"gulp": &toolConfig{
		Container: "gulp",
		Cnf: &containerConfig{
			ShareWorkspace: true,
			LocalUser:      true,
		},
	},
	"node": &toolConfig{
		Container: "node",
		Cnf: &containerConfig{
			ShareWorkspace: true,
			LocalUser:      true,
		},
	},
	"protoc": &toolConfig{
		Container: "protoc",
		Cnf: &containerConfig{
			ShareWorkspace:  true,
			LocalUser:       true,
			ConfigureGopath: true,
		},
	},
	"rambler": &toolConfig{
		Container: "rambler",
		Cnf: &containerConfig{
			ShareWorkspace: true,
			LocalUser:      true,
		},
	},
	"sass-lint": &toolConfig{
		Container: "sass-lint",
		Cnf: &containerConfig{
			ShareWorkspace: true,
			LocalUser:      true,
		},
	},
	"npm": &toolConfig{
		Container: "node",
		Cnf: &containerConfig{
			ShareWorkspace: true,
			LocalUser:      true,
			ShareSSHSocket: true,
		},
	},
	"gsutil": &toolConfig{
		Container: "gcloud",
		Cnf: &containerConfig{
			ShareWorkspace:    true,
			LocalUser:         true,
			ShareGcloudConfig: true,
		},
	},
	"kubectl": &toolConfig{
		Container: "gcloud",
		Cnf: &containerConfig{
			ShareWorkspace:    true,
			LocalUser:         true,
			ShareGcloudConfig: true,
		},
	},
	"gofmt": &toolConfig{
		Container: "go",
		Cnf: &containerConfig{
			ShareWorkspace:  true,
			LocalUser:       true,
			ConfigureGopath: true,
		},
	},
	"mysql": &toolConfig{
		Container: "mysql",
		Cnf: &containerConfig{
			ShareWorkspace: true,
			LocalUser:      true,
			NoTTY:          true,
		},
	},
}

var manualRun = map[string]*containerConfig{
	"go": &containerConfig{
		ShareWorkspace:    true,
		LocalUser:         true,
		ConfigureGopath:   true,
		ShareGcloudConfig: true,
	},
	"node": &containerConfig{
		ShareWorkspace: true,
		LocalUser:      true,
	},
}