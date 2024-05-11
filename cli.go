package onepassword

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"github.com/pkg/errors"
)

type Kind string

const (
	KindPassword    Kind = "password"
	KindCredential  Kind = "credential"
	KindUnspecified Kind = "unspecified"
)

func NewKindFromString(s string) (Kind, error) {
	switch s {
	case "password":
		return KindPassword, nil
	case "credential":
		return KindCredential, nil
	}
	return KindUnspecified, errors.Errorf("unexpected kind: `%s`", s)
}

type OnePassword interface {
	GetByLabel(ctx context.Context, kind Kind, label string) (string, error)
	GetPasswordByLabel(ctx context.Context, label string) (string, error)
	GetCredentialByLabel(ctx context.Context, label string) (string, error)
}

type onePassword struct {
	opBinary string
}

func New() OnePassword {
	return &onePassword{
		opBinary: "op",
	}
}

func (op *onePassword) GetByLabel(ctx context.Context, kind Kind, label string) (string, error) {
	type responseModel struct {
		ID              string `json:"id"`
		Type            string `json:"type"`
		Purpose         string `json:"purpose"`
		Label           string `json:"label"`
		Value           string `json:"value"`
		Reference       string `json:"reference"`
		PasswordDetails struct {
			Strength string `json:"strength"`
		} `json:"password_details"`
	}

	cmd := exec.CommandContext(ctx, op.opBinary, "item", "get", label, "--format=json", fmt.Sprintf("--fields=%s", string(kind)))
	cmd.Stderr = os.Stderr

	data, err := cmd.Output()
	if err != nil {
		return "", errors.Wrap(err, "error running 1password CLI")
	}

	var result responseModel
	return result.Value, errors.Wrap(json.Unmarshal(data, &result), "error decoding op result")
}

func (op *onePassword) GetPasswordByLabel(ctx context.Context, label string) (string, error) {
	return op.GetByLabel(ctx, KindPassword, label)
}

func (op *onePassword) GetCredentialByLabel(ctx context.Context, label string) (string, error) {
	return op.GetByLabel(ctx, KindCredential, label)
}
