package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
)

type cliOptions struct {
	outDir string
	stdout bool
}

var targetToFiles = map[string][]string{
	"bash": {
		".bash_profile",
		".bash_prompt",
		".aliases",
		".exports",
		".functions",
	},
	"vim": {
		".vim/vimrc",
		".vim/plugin.vim",
		".vim/gvimrc",
		".vim/.ycm_extra_conf.py",
		".vim/colors/peaksea.vim",
	},
	"git": {
		".gitconfig",
		".gitignore",
	},
	"tig": {
		".tigrc",
	},
	"tmux": {
		".tmux.conf",
		".tmux-theme.conf",
	},
	"kube": {
		"config/kube-ps1.sh",
	},
	"sshrc": {
		"config/sshrc.d/vimrc",
		"config/sshrc.d/.vimrc",
	},
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "config":
		if err := runConfig(os.Args[2:]); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	default:
		printUsage()
		os.Exit(1)
	}
}

func runConfig(args []string) error {
	if runtime.GOOS != "linux" {
		return fmt.Errorf("this program only supports linux, current os=%s", runtime.GOOS)
	}

	if len(args) == 0 {
		return errors.New("missing target, run: config list")
	}

	if args[0] == "list" {
		printTargets()
		return nil
	}

	target := strings.ToLower(args[0])

	fs := flag.NewFlagSet("config", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	opts := cliOptions{}
	fs.StringVar(&opts.outDir, "out", "generated", "output directory")
	fs.BoolVar(&opts.stdout, "stdout", false, "print generated content to stdout")

	if err := fs.Parse(args[1:]); err != nil {
		return err
	}

	files, err := resolveTarget(target)
	if err != nil {
		return err
	}

	generated := 0
	for _, rel := range files {
		content, ok := embeddedConfigs[rel]
		if !ok {
			return fmt.Errorf("embedded content not found for %s", rel)
		}
		body := []byte(content)

		if opts.stdout {
			printToStdout(rel, body)
		} else {
			outPath := filepath.Join(opts.outDir, filepath.FromSlash(rel))
			if err := writeFile(outPath, body); err != nil {
				return fmt.Errorf("write %s failed: %w", outPath, err)
			}
			fmt.Printf("generated: %s\n", outPath)
		}

		generated++
	}

	fmt.Printf("done: target=%s files=%d\n", target, generated)
	return nil
}

func resolveTarget(target string) ([]string, error) {
	if target == "all" {
		merged := make([]string, 0)
		seen := map[string]bool{}
		for _, files := range targetToFiles {
			for _, f := range files {
				if !seen[f] {
					seen[f] = true
					merged = append(merged, f)
				}
			}
		}
		sort.Strings(merged)
		return merged, nil
	}

	files, ok := targetToFiles[target]
	if !ok {
		return nil, fmt.Errorf("unknown target: %s", target)
	}

	cp := append([]string(nil), files...)
	return cp, nil
}

func writeFile(path string, body []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, body, 0o644)
}

func printToStdout(rel string, body []byte) {
	fmt.Printf("===== %s =====\n", rel)
	fmt.Print(string(body))
	if len(body) == 0 || body[len(body)-1] != '\n' {
		fmt.Println()
	}
}

func printTargets() {
	keys := make([]string, 0, len(targetToFiles)+1)
	for k := range targetToFiles {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	fmt.Println("available targets:")
	for _, k := range keys {
		fmt.Printf("- %s\n", k)
	}
	fmt.Println("- all")
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  config list")
	fmt.Println("  config <target> [--out generated] [--stdout]")
	fmt.Println("")
	fmt.Println("Note: linux only. Config content is embedded in this binary.")
	fmt.Println("")
	printTargets()
}
