package policy

import (
	"fmt"
	"strings"
)

type CILArtifact struct {
	ModuleName string `json:"moduleName"`
	TypeName   string `json:"typeName"`
	Digest     string `json:"digest"`
	Content    string `json:"content"`
}

// GenerateCIL creates a candidate module scaffold. It deliberately emits only
// declarations in this slice: findings require an explicit reviewed rule
// mapping before they may become SELinux enforcement.
func GenerateCIL(ir PolicyIR) (CILArtifact, error) {
	if ir.ModuleName == "" || ir.TypeName == "" {
		return CILArtifact{}, fmt.Errorf("policy IR module and type names are required")
	}
	if !cilIdentifierPattern.MatchString(ir.ModuleName) || !cilIdentifierPattern.MatchString(ir.TypeName) {
		return CILArtifact{}, fmt.Errorf("policy IR contains an invalid CIL identifier")
	}
	if len(ir.Rules) > 0 {
		return CILArtifact{}, fmt.Errorf("policy enforcement rules are not enabled in this planning slice")
	}

	content := fmt.Sprintf("(block %s\n    (type %s)\n)\n", ir.ModuleName, ir.TypeName)
	digest, err := digestBytes([]byte(content))
	if err != nil {
		return CILArtifact{}, fmt.Errorf("digest generated CIL: %w", err)
	}
	return CILArtifact{
		ModuleName: ir.ModuleName,
		TypeName:   ir.TypeName,
		Digest:     digest,
		Content:    content,
	}, nil
}

// ValidateCIL performs a conservative, dependency-free structural check. It
// is not a replacement for a distribution's SELinux compiler.
func ValidateCIL(content string) error {
	if strings.TrimSpace(content) == "" {
		return fmt.Errorf("CIL content is empty")
	}
	if !strings.HasSuffix(content, "\n") {
		return fmt.Errorf("CIL content must end with a newline")
	}
	depth := 0
	for _, char := range content {
		switch char {
		case '(':
			depth++
		case ')':
			depth--
			if depth < 0 {
				return fmt.Errorf("CIL contains an unmatched closing parenthesis")
			}
		case '\x00':
			return fmt.Errorf("CIL contains a NUL byte")
		}
	}
	if depth != 0 {
		return fmt.Errorf("CIL parentheses are unbalanced")
	}
	if !strings.HasPrefix(content, "(block ") {
		return fmt.Errorf("CIL must start with a block declaration")
	}
	return nil
}
