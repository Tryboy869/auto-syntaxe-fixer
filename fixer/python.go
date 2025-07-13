package fixer

import (
	"bytes"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

// FixPythonFiles processes all Python files in the repository and applies fixes
func FixPythonFiles(repoPath string) (int, []string) {
	pythonFiles := GetFilesByExtension(repoPath, []string{".py"})
	if len(pythonFiles) == 0 {
		return 0, nil
	}

	var fixedFiles []string
	totalChanges := 0

	// Process each Python file
	for _, file := range pythonFiles {
		changes := 0
		
		// Read original file
		original, err := os.ReadFile(file)
		if err != nil {
			continue
		}

		// Apply Python formatting
		formatted := formatPythonFile(original)

		// Apply syntax fixes
		fixed, syntaxChanges := fixPythonSyntax(formatted)
		changes += syntaxChanges

		// Apply import fixes
		finalCode, importChanges := fixPythonImports(fixed)
		changes += importChanges

		// Write back if changes were made
		if changes > 0 && !bytes.Equal(original, finalCode) {
			err = os.WriteFile(file, finalCode, 0644)
			if err == nil {
				fixedFiles = append(fixedFiles, file)
				totalChanges += changes
			}
		}
	}

	return totalChanges, fixedFiles
}

// formatPythonFile applies standard Python formatting
func formatPythonFile(source []byte) []byte {
	// Try to run autopep8 if available
	cmd := exec.Command("autopep8", "--aggressive", "--aggressive", "-")
	cmd.Stdin = bytes.NewReader(source)
	
	var out bytes.Buffer
	cmd.Stdout = &out
	
	if err := cmd.Run(); err == nil {
		return out.Bytes()
	}
	
	// Fallback to manual formatting
	return formatPythonManually(source)
}

// formatPythonManually applies basic Python formatting rules
func formatPythonManually(source []byte) []byte {
	code := string(source)
	lines := strings.Split(code, "\n")
	
	// Fix indentation (convert tabs to 4 spaces)
	for i, line := range lines {
		lines[i] = strings.ReplaceAll(line, "\t", "    ")
	}
	
	// Fix spacing around operators
	for i, line := range lines {
		// Add spaces around = operator
		re := regexp.MustCompile(`(\w)=(\w)`)
		line = re.ReplaceAllString(line, "$1 = $2")
		
		// Add spaces around comparison operators
		re = regexp.MustCompile(`(\w)==(\w)`)
		line = re.ReplaceAllString(line, "$1 == $2")
		
		re = regexp.MustCompile(`(\w)!=(\w)`)
		line = re.ReplaceAllString(line, "$1 != $2")
		
		lines[i] = line
	}
	
	return []byte(strings.Join(lines, "\n"))
}

// fixPythonSyntax applies basic syntax fixes to Python code
func fixPythonSyntax(source []byte) ([]byte, int) {
	code := string(source)
	changes := 0
	lines := strings.Split(code, "\n")
	
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		
		// Fix missing colons in control structures
		if strings.HasPrefix(trimmed, "if ") || strings.HasPrefix(trimmed, "elif ") ||
		   strings.HasPrefix(trimmed, "else") || strings.HasPrefix(trimmed, "for ") ||
		   strings.HasPrefix(trimmed, "while ") || strings.HasPrefix(trimmed, "def ") ||
		   strings.HasPrefix(trimmed, "class ") || strings.HasPrefix(trimmed, "try") ||
		   strings.HasPrefix(trimmed, "except") || strings.HasPrefix(trimmed, "finally") {
			
			if !strings.HasSuffix(trimmed, ":") && !strings.Contains(trimmed, "#") {
				lines[i] = line + ":"
				changes++
			}
		}
		
		// Fix missing parentheses in print statements (Python 2 to 3)
		if strings.Contains(trimmed, "print ") && !strings.Contains(trimmed, "print(") {
			re := regexp.MustCompile(`print\s+(.+)`)
			if re.MatchString(trimmed) {
				lines[i] = re.ReplaceAllString(line, "print($1)")
				changes++
			}
		}
		
		// Fix missing quotes in strings
		if strings.Contains(trimmed, "=") && !strings.Contains(trimmed, "\"") && !strings.Contains(trimmed, "'") {
			// Basic string detection (this is simplified)
			re := regexp.MustCompile(`=\s*([a-zA-Z_][a-zA-Z0-9_]*)\s*$`)
			if re.MatchString(trimmed) {
				lines[i] = re.ReplaceAllString(line, `= "$1"`)
				changes++
			}
		}
	}
	
	return []byte(strings.Join(lines, "\n")), changes
}

// fixPythonImports fixes import statements and adds missing imports
func fixPythonImports(source []byte) ([]byte, int) {
	// Try to run isort if available
	cmd := exec.Command("isort", "--stdout", "-")
	cmd.Stdin = bytes.NewReader(source)
	
	var out bytes.Buffer
	cmd.Stdout = &out
	
	if err := cmd.Run(); err == nil {
		result := out.Bytes()
		if !bytes.Equal(source, result) {
			return result, 1
		}
		return source, 0
	}
	
	// Fallback to manual import fixing
	return fixPythonImportsManually(source)
}

// fixPythonImportsManually applies basic import fixes
func fixPythonImportsManually(source []byte) ([]byte, int) {
	code := string(source)
	lines := strings.Split(code, "\n")
	changes := 0
	
	// Sort imports (basic implementation)
	var importLines []string
	var fromImportLines []string
	var otherLines []string
	var inImportSection bool
	
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		
		if strings.HasPrefix(trimmed, "import ") {
			importLines = append(importLines, line)
			inImportSection = true
		} else if strings.HasPrefix(trimmed, "from ") {
			fromImportLines = append(fromImportLines, line)
			inImportSection = true
		} else if trimmed == "" && inImportSection {
			// Empty line in import section
			continue
		} else {
			if inImportSection {
				inImportSection = false
			}
			otherLines = append(otherLines, line)
		}
	}
	
	// Check if we need to add missing imports
	codeStr := strings.Join(otherLines, "\n")
	missingImports := findMissingPythonImports(codeStr)
	
	if len(missingImports) > 0 {
		importLines = append(importLines, missingImports...)
		changes += len(missingImports)
	}
	
	// Reconstruct file with sorted imports
	var newLines []string
	
	// Add regular imports first
	for _, imp := range importLines {
		newLines = append(newLines, imp)
	}
	
	// Add from imports
	for _, imp := range fromImportLines {
		newLines = append(newLines, imp)
	}
	
	// Add empty line after imports if there are any
	if len(importLines) > 0 || len(fromImportLines) > 0 {
		newLines = append(newLines, "")
	}
	
	// Add the rest of the code
	newLines = append(newLines, otherLines...)
	
	return []byte(strings.Join(newLines, "\n")), changes
}

// findMissingPythonImports detects commonly used modules that aren't imported
func findMissingPythonImports(code string) []string {
	var missing []string
	
	// Common patterns that need imports
	patterns := map[string]string{
		"os.":       "import os",
		"sys.":      "import sys",
		"json.":     "import json",
		"re.":       "import re",
		"time.":     "import time",
		"datetime.": "import datetime",
		"random.":   "import random",
		"math.":     "import math",
		"urllib.":   "import urllib",
		"requests.": "import requests",
	}
	
	for pattern, importStmt := range patterns {
		if strings.Contains(code, pattern) && !strings.Contains(code, importStmt) {
			missing = append(missing, importStmt)
		}
	}
	
	return missing
}

// fixPythonIndentation fixes indentation issues in Python code
func fixPythonIndentation(source []byte) ([]byte, int) {
	code := string(source)
	lines := strings.Split(code, "\n")
	changes := 0
	
	expectedIndent := 0
	
	for i, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		
		trimmed := strings.TrimSpace(line)
		
		// Calculate expected indentation
		if strings.HasSuffix(trimmed, ":") {
			expectedIndent += 4
		} else if strings.HasPrefix(trimmed, "return") || 
		         strings.HasPrefix(trimmed, "break") || 
		         strings.HasPrefix(trimmed, "continue") {
			if i < len(lines)-1 {
				nextLine := strings.TrimSpace(lines[i+1])
				if nextLine != "" && !strings.HasPrefix(nextLine, "def ") && 
				   !strings.HasPrefix(nextLine, "class ") {
					expectedIndent -= 4
				}
			}
		}
		
		// Fix indentation
		currentIndent := len(line) - len(strings.TrimLeft(line, " "))
		if currentIndent != expectedIndent && expectedIndent >= 0 {
			lines[i] = strings.Repeat(" ", expectedIndent) + trimmed
			changes++
		}
	}
	
	return []byte(strings.Join(lines, "\n")), changes
}
