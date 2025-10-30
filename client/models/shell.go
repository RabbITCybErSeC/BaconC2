package models

import "fmt"

type ShellType int

const (
	ShellTypeUnknown ShellType = iota
	ShellTypeCmd
	ShellTypePowerShell
	ShellTypeBash
	ShellTypeSh
	ShellTypeZsh
	ShellTypeFish
)

func (s ShellType) String() string {
	switch s {
	case ShellTypeCmd:
		return "cmd"
	case ShellTypePowerShell:
		return "powershell"
	case ShellTypeBash:
		return "bash"
	case ShellTypeSh:
		return "sh"
	case ShellTypeZsh:
		return "zsh"
	case ShellTypeFish:
		return "fish"
	default:
		return "unknown"
	}
}

func ParseShellType(s string) (ShellType, error) {
	switch s {
	case "cmd":
		return ShellTypeCmd, nil
	case "powershell":
		return ShellTypePowerShell, nil
	case "bash":
		return ShellTypeBash, nil
	case "sh":
		return ShellTypeSh, nil
	case "zsh":
		return ShellTypeZsh, nil
	case "fish":
		return ShellTypeFish, nil
	default:
		return ShellTypeUnknown, fmt.Errorf("invalid shell type: %s", s)
	}
}
