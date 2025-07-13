package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"auto-syntax-fixer/fixer"
	"auto-syntax-fixer/git"
)

type Config struct {
	RepoURL    string
	Token      string
	Branch     string
	DryRun     bool
	OutputPath string
}

func main() {
	config := parseFlags()
	
	if config.RepoURL == "" {
		log.Fatal("Repository URL is required (--repo)")
	}

	// Create temporary directory for cloning
	tempDir := filepath.Join(os.TempDir(), "auto-syntax-fixer")
	os.RemoveAll(tempDir)
	
	fmt.Printf("🚀 Starting auto-syntax-fixer for: %s\n", config.RepoURL)
	
	// Clone repository
	fmt.Println("📥 Cloning repository...")
	repoPath, err := git.CloneRepo(config.RepoURL, tempDir, config.Token, config.Branch)
	if err != nil {
		log.Fatalf("Failed to clone repository: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Detect languages in the repository
	fmt.Println("🔍 Detecting languages...")
	languages := fixer.DetectLanguages(repoPath)
	if len(languages) == 0 {
		fmt.Println("No supported languages found")
		return
	}
	
	fmt.Printf("Found languages: %s\n", strings.Join(languages, ", "))

	// Apply fixes for each detected language
	totalChanges := 0
	var fixedFiles []string
	
	for _, lang := range languages {
		fmt.Printf("🔧 Processing %s files...\n", lang)
		
		var changes int
		var files []string
		
		switch lang {
		case "Go":
			changes, files = fixer.FixGoFiles(repoPath)
		case "Python":
			changes, files = fixer.FixPythonFiles(repoPath)
		case "JavaScript":
			changes, files = fixer.FixJavaScriptFiles(repoPath)
		case "TypeScript":
			changes, files = fixer.FixTypeScriptFiles(repoPath)
		}
		
		totalChanges += changes
		fixedFiles = append(fixedFiles, files...)
		
		if changes > 0 {
			fmt.Printf("  ✅ Fixed %d issues in %d files\n", changes, len(files))
		} else {
			fmt.Printf("  ✨ No issues found\n")
		}
	}

	// If no changes were made, exit
	if totalChanges == 0 {
		fmt.Println("🎉 No fixes needed! Repository is already clean.")
		return
	}

	// Create branch and commit changes (if not dry run)
	if !config.DryRun {
		fmt.Println("📝 Creating fix branch and committing changes...")
		
		branchName := "auto-syntax-fixes"
		err = git.CreateBranchAndCommit(repoPath, branchName, fixedFiles, config.Token)
		if err != nil {
			log.Fatalf("Failed to commit changes: %v", err)
		}
		
		fmt.Printf("✅ Changes committed to branch: %s\n", branchName)
		fmt.Printf("🚀 Push the branch manually or create a PR\n")
	} else {
		fmt.Println("🔍 Dry run mode - no changes committed")
	}

	// Generate summary report
	generateReport(languages, totalChanges, fixedFiles, config.OutputPath)
}

func parseFlags() Config {
	var config Config
	
	flag.StringVar(&config.RepoURL, "repo", "", "GitHub repository URL (required)")
	flag.StringVar(&config.Token, "token", "", "GitHub token for private repos")
	flag.StringVar(&config.Branch, "branch", "main", "Branch to work on")
	flag.BoolVar(&config.DryRun, "dry-run", false, "Run without making changes")
	flag.StringVar(&config.OutputPath, "output", "", "Output file for report (default: stdout)")
	
	flag.Parse()
	
	return config
}

func generateReport(languages []string, totalChanges int, fixedFiles []string, outputPath string) {
	report := fmt.Sprintf(`
🔧 Auto-Syntax-Fixer Report
==========================

Languages processed: %s
Total fixes applied: %d
Files modified: %d

Modified files:
%s

Summary: %s
`, 
		strings.Join(languages, ", "),
		totalChanges,
		len(fixedFiles),
		strings.Join(fixedFiles, "\n"),
		func() string {
			if totalChanges == 0 {
				return "No issues found - repository is clean!"
			}
			return fmt.Sprintf("Successfully fixed %d syntax issues", totalChanges)
		}(),
	)

	if outputPath != "" {
		err := os.WriteFile(outputPath, []byte(report), 0644)
		if err != nil {
			log.Printf("Failed to write report to file: %v", err)
		} else {
			fmt.Printf("📄 Report saved to: %s\n", outputPath)
		}
	} else {
		fmt.Print(report)
	}
}
