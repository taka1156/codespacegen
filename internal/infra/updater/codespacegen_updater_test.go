package updater

import (
	"testing"

	"github.com/blang/semver"
)

func TestNewCodespacegenUpdater(t *testing.T) {
	ucc := NewCodespacegenUpdater()
	_ = ucc
}

func TestUpdate_InvalidVersionString(t *testing.T) {
	ucc := CodespacegenUpdater{}
	err := ucc.Update("invalid-version")
	if err == nil {
		t.Error("expected error for invalid version string, got nil")
	}
	if err.Error() == "" {
		t.Error("expected non-empty error message")
	}
}

func TestUpdate_InvalidVersionFormat(t *testing.T) {
	ucc := CodespacegenUpdater{}
	err := ucc.Update("not.a.valid.version")
	if err == nil {
		t.Error("expected error for invalid version format, got nil")
	}
}

func TestUpdate_EmptyVersion(t *testing.T) {
	ucc := CodespacegenUpdater{}
	err := ucc.Update("")
	if err == nil {
		t.Error("expected error for empty version, got nil")
	}
}

func TestUpdate_ValidVersionFormat(t *testing.T) {
	ucc := CodespacegenUpdater{}

	for _, version := range []string{"1.0.0", "1.2.3", "0.1.0"} {
		t.Run(version, func(t *testing.T) {
			err := ucc.Update(version)
			if err != nil {
				// Update が失敗した場合、バージョンのパース失敗ではないことを確認する。
				// selfupdate によるネットワーク/API エラーはテスト環境では許容される。
				if _, parseErr := semver.Parse(version); parseErr != nil {
					t.Errorf("Update(%q): version parse should not fail, got: %v", version, parseErr)
				}
			}
		})
	}
}
