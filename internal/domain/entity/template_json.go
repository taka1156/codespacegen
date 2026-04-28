package entity

import "codespacegen/internal/utils"

const DefaultTemplateJsonPath = "codespacegen.json"

type TemplateJson struct {
	Schema     string    `json:"$schema,omitempty"`
	Common     JsonEntry `json:"common,omitempty"`
	Go         JsonEntry `json:"go,omitempty"`
	Python     JsonEntry `json:"python,omitempty"`
	NodeBiome  JsonEntry `json:"node:biome,omitempty"`
	NodeEslint JsonEntry `json:"node:eslint,omitempty"`
	NodeReact  JsonEntry `json:"node:react,omitempty"`
	Rust       JsonEntry `json:"rust,omitempty"`
	Moonbit    JsonEntry `json:"moonbit,omitempty"`
	Gcc        JsonEntry `json:"gcc,omitempty"`
}

var DefaultTemplateJson = TemplateJson{
	Schema: "https://raw.githubusercontent.com/taka1156/codespacegen/master/codespacegen.schema.json",
	Common: JsonEntry{
		Timezone: "Asia/Tokyo",
		Locale:   utils.Ptr(DefaultLocale),
		VSCodeExtensions: []string{
			"MS-CEINTL.vscode-language-pack-ja",
			"GitHub.copilot",
			"GitHub.copilot-chat",
			"streetsidesoftware.code-spell-checker",
			"username.errorlens",
		},
	},
	Go: JsonEntry{
		Image: "golang:1.24-alpine",
		VSCodeExtensions: []string{
			"golang.GO",
		},
	},
	Python: JsonEntry{
		Image: "python:3.12-alpine",
		VSCodeExtensions: []string{
			"ms-python.python",
		},
	},
	NodeBiome: JsonEntry{
		Image: "node:24-alpine",
		VSCodeExtensions: []string{
			"biomejs.biome",
		},
	},
	NodeEslint: JsonEntry{
		Image: "node:24-alpine",
		VSCodeExtensions: []string{
			"dbaeumer.vscode-eslint",
			"esbenp.prettier-vscode",
		},
	},
	NodeReact: JsonEntry{
		Image: "node:24-alpine",
		VSCodeExtensions: []string{
			"jawandarajbir.react-vscode-extension-pack",
			"dbaeumer.vscode-eslint",
			"stylelint.vscode-stylelint",
			"esbenp.prettier-vscode",
		},
	},
	Rust: JsonEntry{
		Image: "rust:1.72-alpine",
		VSCodeExtensions: []string{
			"Zerotaskx.rust-extension-pack",
		},
	},
	Moonbit: JsonEntry{
		Image:      "ubuntu:24.04",
		RunCommand: "curl -fsSL https://cli.moonbitlang.com/install/unix.sh | bash",
		VSCodeExtensions: []string{
			"moonbit.moonbit-lang",
		},
	},
	Gcc: JsonEntry{
		Image:      "ubuntu:24.04",
		RunCommand: "apt install -y gcc make git binutils libc6-dev",
		VSCodeExtensions: []string{
			"ms-vscode.cpptools",
		},
	},
}
