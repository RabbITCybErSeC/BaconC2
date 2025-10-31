package filesystem

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	command_handler "github.com/RabbITCybErSeC/BaconC2/pkg/commands/handlers"
	"github.com/RabbITCybErSeC/BaconC2/pkg/models"
)

func LsHandler(ctx *command_handler.CommandContext) models.CommandResult {
	var targetPath string

	if len(ctx.Command.Args) == 0 {
		targetPath = ctx.State.GetWorkingDirectory()
	} else {
		targetPath = ctx.Command.Args[0]
	}

	if !filepath.IsAbs(targetPath) {
		targetPath = filepath.Join(ctx.State.GetWorkingDirectory(), targetPath)
	}

	entries, err := os.ReadDir(targetPath)
	if err != nil {
		return models.CommandResult{
			ID:     ctx.Command.ID,
			Status: models.CommandStatusFailed,
			Output: fmt.Sprintf("Failed to list directory '%s': %v", targetPath, err),
		}
	}

	sort.Slice(entries, func(i, j int) bool {
		if entries[i].IsDir() != entries[j].IsDir() {
			return entries[i].IsDir()
		}
		return entries[i].Name() < entries[j].Name()
	})

	var output strings.Builder
	output.WriteString(fmt.Sprintf("Directory: %s\n\n", targetPath))
	output.WriteString(fmt.Sprintf("%-10s %-20s %-15s %s\n", "Type", "Modified", "Size", "Name"))
	output.WriteString(strings.Repeat("-", 80) + "\n")

	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}

		entryType := "FILE"
		if entry.IsDir() {
			entryType = "DIR"
		}

		size := formatSize(info.Size())
		if entry.IsDir() {
			size = "-"
		}

		modTime := info.ModTime().Format("2006-01-02 15:04:05")

		output.WriteString(fmt.Sprintf("%-10s %-20s %-15s %s\n",
			entryType,
			modTime,
			size,
			entry.Name(),
		))
	}

	output.WriteString(fmt.Sprintf("\nTotal: %d items\n", len(entries)))

	return models.CommandResult{
		ID:     ctx.Command.ID,
		Status: models.CommandStatusCompleted,
		Output: output.String(),
	}
}

func formatSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "KMGTPE"[exp])
}

func LsDetailedHandler(ctx *command_handler.CommandContext) models.CommandResult {
	var targetPath string

	if len(ctx.Command.Args) == 0 {
		targetPath = ctx.State.GetWorkingDirectory()
	} else {
		targetPath = ctx.Command.Args[0]
	}

	if !filepath.IsAbs(targetPath) {
		targetPath = filepath.Join(ctx.State.GetWorkingDirectory(), targetPath)
	}

	var output strings.Builder
	output.WriteString(fmt.Sprintf("Directory: %s\n\n", targetPath))

	err := filepath.WalkDir(targetPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if path == targetPath {
			return nil
		}

		relPath, _ := filepath.Rel(targetPath, path)
		if strings.Contains(relPath, string(filepath.Separator)) {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		info, err := d.Info()
		if err != nil {
			return nil
		}

		permissions := info.Mode().String()
		size := formatSize(info.Size())
		if d.IsDir() {
			size = "-"
		}
		modTime := info.ModTime().Format(time.RFC3339)

		output.WriteString(fmt.Sprintf("%s %15s %s %s\n",
			permissions,
			size,
			modTime,
			d.Name(),
		))

		return nil
	})

	if err != nil {
		return models.CommandResult{
			ID:     ctx.Command.ID,
			Status: models.CommandStatusFailed,
			Output: fmt.Sprintf("Failed to list directory: %v", err),
		}
	}

	return models.CommandResult{
		ID:     ctx.Command.ID,
		Status: models.CommandStatusCompleted,
		Output: output.String(),
	}
}

func NewLsHandler() command_handler.StatefulCommandHandler {
	return command_handler.StatefulCommandHandler{
		Name:    "ls",
		Handler: LsHandler,
	}
}
