package entity

import "github.com/taka1156/codespacegen/internal/utils"

const DefaultTemplateJsonPath = "codespacegen.json"

var DefaultTemplateJson = JsonConfig{
	Schema: "https://raw.githubusercontent.com/taka1156/codespacegen/master/codespacegen.schema.json",
	Common: new(CommonEntry{
		Timezone: utils.Ptr("Asia/Tokyo"),
		Locale:   utils.Ptr(DefaultLocale),
		VSCodeExtensions: utils.Ptr([]string{
			"MS-CEINTL.vscode-language-pack-ja",
			"streetsidesoftware.code-spell-checker",
			"username.errorlens",
		}),
	}),
	Langs: []*LangEntry{
		new(LangEntry{
			ProfileName: "go",
			Image:       "golang:1.24-alpine",
			VSCodeExtensions: utils.Ptr([]string{
				"golang.GO",
			}),
		}),
		new(LangEntry{
			ProfileName: "python",
			Image:       "python:3.12-alpine",
			VSCodeExtensions: utils.Ptr([]string{
				"ms-python.python",
			}),
		}),
		new(LangEntry{
			ProfileName: "node:biome",
			Image:       "node:24-alpine",
			VSCodeExtensions: utils.Ptr([]string{
				"biomejs.biome",
			}),
		}),
		new(LangEntry{
			ProfileName: "node:eslint",
			Image:       "node:24-alpine",
			VSCodeExtensions: utils.Ptr([]string{
				"dbaeumer.vscode-eslint",
				"esbenp.prettier-vscode",
			}),
		}),
		new(LangEntry{
			ProfileName: "node:react",
			Image:       "node:24-alpine",
			VSCodeExtensions: utils.Ptr([]string{
				"jawandarajbir.react-vscode-extension-pack",
				"dbaeumer.vscode-eslint",
				"stylelint.vscode-stylelint",
				"esbenp.prettier-vscode",
			}),
		}),
		new(LangEntry{
			ProfileName: "rust",
			Image:       "rust:1.72-alpine",
			VSCodeExtensions: utils.Ptr([]string{
				"Zerotaskx.rust-extension-pack",
			}),
		}),
		new(LangEntry{
			ProfileName: "moonbit",
			Image:       "ubuntu:24.04",
			RunCommand:  utils.Ptr("curl -fsSL https://cli.moonbitlang.com/install/unix.sh | bash"),
			VSCodeExtensions: utils.Ptr([]string{
				"moonbit.moonbit-lang",
			}),
		}),
		new(LangEntry{
			ProfileName: "gcc",
			Image:       "ubuntu:24.04",
			LinuxPackages: utils.Ptr([]string{
				"gcc",
				"make",
				"binutils",
				"libc6-dev",
			}),
			VSCodeExtensions: utils.Ptr([]string{
				"ms-vscode.cpptools",
			}),
		}),
	},
}
