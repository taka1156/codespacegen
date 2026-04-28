package entity

import "codespacegen/internal/utils"

const DefaultTemplateJsonPath = "codespacegen.json"

type TemplateJson struct {
	Schema     string      `json:"$schema,omitempty"`
	Common     CommonEntry `json:"common,omitempty"`
	Go         LangEntry   `json:"go,omitempty"`
	Python     LangEntry   `json:"python,omitempty"`
	NodeBiome  LangEntry   `json:"node:biome,omitempty"`
	NodeEslint LangEntry   `json:"node:eslint,omitempty"`
	NodeReact  LangEntry   `json:"node:react,omitempty"`
	Rust       LangEntry   `json:"rust,omitempty"`
	Moonbit    LangEntry   `json:"moonbit,omitempty"`
	Gcc        LangEntry   `json:"gcc,omitempty"`
}

var DefaultTemplateJson = TemplateJson{
	Schema: "https://raw.githubusercontent.com/taka1156/codespacegen/master/codespacegen.schema.json",
	Common: CommonEntry{
		Timezone: utils.Ptr("Asia/Tokyo"),
		Locale:   utils.Ptr(DefaultLocale),
		VSCodeExtensions: utils.Ptr([]string{
			"MS-CEINTL.vscode-language-pack-ja",
			"GitHub.copilot",
			"GitHub.copilot-chat",
			"streetsidesoftware.code-spell-checker",
			"username.errorlens",
		}),
	},
	Go: LangEntry{
		Image: "golang:1.24-alpine",
		VSCodeExtensions: utils.Ptr([]string{
			"golang.GO",
		}),
	},
	Python: LangEntry{
		Image: "python:3.12-alpine",
		VSCodeExtensions: utils.Ptr([]string{
			"ms-python.python",
		}),
	},
	NodeBiome: LangEntry{
		Image: "node:24-alpine",
		VSCodeExtensions: utils.Ptr([]string{
			"biomejs.biome",
		}),
	},
	NodeEslint: LangEntry{
		Image: "node:24-alpine",
		VSCodeExtensions: utils.Ptr([]string{
			"dbaeumer.vscode-eslint",
			"esbenp.prettier-vscode",
		}),
	},
	NodeReact: LangEntry{
		Image: "node:24-alpine",
		VSCodeExtensions: utils.Ptr([]string{
			"jawandarajbir.react-vscode-extension-pack",
			"dbaeumer.vscode-eslint",
			"stylelint.vscode-stylelint",
			"esbenp.prettier-vscode",
		}),
	},
	Rust: LangEntry{
		Image: "rust:1.72-alpine",
		VSCodeExtensions: utils.Ptr([]string{
			"Zerotaskx.rust-extension-pack",
		}),
	},
	Moonbit: LangEntry{
		Image:      "ubuntu:24.04",
		RunCommand: utils.Ptr("curl -fsSL https://cli.moonbitlang.com/install/unix.sh | bash"),
		VSCodeExtensions: utils.Ptr([]string{
			"moonbit.moonbit-lang",
		}),
	},
	Gcc: LangEntry{
		Image: "ubuntu:24.04",
		LinuxPackages: utils.Ptr([]string{
			"gcc",
			"make",
			"binutils",
			"libc6-dev",
		}),
		VSCodeExtensions: utils.Ptr([]string{
			"ms-vscode.cpptools",
		}),
	},
}
