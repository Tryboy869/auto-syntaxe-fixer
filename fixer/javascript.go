package fixer

import (
	"bytes"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

// FixJavaScriptFiles processes all JavaScript files in the repository and applies fixes
func FixJavaScriptFiles(repoPath string) (int, []string) {
	jsFiles := GetFilesByExtension(repoPath, []string{".js", ".jsx"})
	if len(jsFiles) == 0 {
		return 0, nil
	}

	var fixedFiles []string
	totalChanges := 0

	// Process each JavaScript file
	for _, file := range jsFiles {
		changes := 0
		
		// Read original file
		original, err := os.ReadFile(file)
		if err != nil {
			continue
		}

		// Apply JavaScript formatting
		formatted := formatJavaScriptFile(original)

		// Apply syntax fixes
		fixed, syntaxChanges := fixJavaScriptSyntax(formatted)
		changes += syntaxChanges

		// Apply import fixes
		finalCode, importChanges := fixJavaScriptImports(fixed)
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

// FixTypeScriptFiles processes all TypeScript files in the repository and applies fixes
func FixTypeScriptFiles(repoPath string) (int, []string) {
	tsFiles := GetFilesByExtension(repoPath, []string{".ts", ".tsx"})
	if len(tsFiles) == 0 {
		return 0, nil
	}

	var fixedFiles []string
	totalChanges := 0

	// Process each TypeScript file
	for _, file := range tsFiles {
		changes := 0
		
		// Read original file
		original, err := os.ReadFile(file)
		if err != nil {
			continue
		}

		// Apply TypeScript formatting (similar to JavaScript)
		formatted := formatJavaScriptFile(original)

		// Apply syntax fixes
		fixed, syntaxChanges := fixTypeScriptSyntax(formatted)
		changes += syntaxChanges

		// Apply import fixes
		finalCode, importChanges := fixJavaScriptImports(fixed)
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

// formatJavaScriptFile applies standard JavaScript/TypeScript formatting
func formatJavaScriptFile(source []byte) []byte {
	// Try to run prettier if available
	cmd := exec.Command("prettier", "--stdin-filepath", "file.js")
	cmd.Stdin = bytes.NewReader(source)
	
	var out bytes.Buffer
	cmd.Stdout = &out
	
	if err := cmd.Run(); err == nil {
		return out.Bytes()
	}
	
	// Fallback to manual formatting
	return formatJavaScriptManually(source)
}

// formatJavaScriptManually applies basic JavaScript formatting rules
func formatJavaScriptManually(source []byte) []byte {
	code := string(source)
	lines := strings.Split(code, "\n")
	
	// Fix indentation (convert tabs to 2 spaces for JS)
	for i, line := range lines {
		lines[i] = strings.ReplaceAll(line, "\t", "  ")
	}
	
	// Fix spacing around operators
	for i, line := range lines {
		// Add spaces around = operator
		re := regexp.MustCompile(`(\w)=(\w)`)
		line = re.ReplaceAllString(line, "$1 = $2")
		
		// Add spaces around comparison operators
		re = regexp.MustCompile(`(\w)==(\w)`)
		line = re.ReplaceAllString(line, "$1 == $2")
		
		re = regexp.MustCompile(`(\w)===(\w)`)
		line = re.ReplaceAllString(line, "$1 === $2")
		
		// Fix spacing around braces
		re = regexp.MustCompile(`(\w){`)
		line = re.ReplaceAllString(line, "$1 {")
		
		lines[i] = line
	}
	
	return []byte(strings.Join(lines, "\n"))
}

// fixJavaScriptSyntax applies basic syntax fixes to JavaScript code
func fixJavaScriptSyntax(source []byte) ([]byte, int) {
	code := string(source)
	changes := 0
	lines := strings.Split(code, "\n")
	
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		
		// Fix missing semicolons
		if len(trimmed) > 0 && !strings.HasSuffix(trimmed, ";") && 
		   !strings.HasSuffix(trimmed, "{") && !strings.HasSuffix(trimmed, "}") &&
		   !strings.HasPrefix(trimmed, "//") && !strings.HasPrefix(trimmed, "/*") &&
		   !strings.HasPrefix(trimmed, "if") && !strings.HasPrefix(trimmed, "for") &&
		   !strings.HasPrefix(trimmed, "while") && !strings.HasPrefix(trimmed, "function") &&
		   !strings.HasPrefix(trimmed, "class") && !strings.HasPrefix(trimmed, "else") {
			
			// Check if it's a statement that needs semicolon
			if strings.Contains(trimmed, "=") || strings.Contains(trimmed, "return") ||
			   strings.Contains(trimmed, "break") || strings.Contains(trimmed, "continue") ||
			   strings.Contains(trimmed, "throw") || strings.Contains(trimmed, "var") ||
			   strings.Contains(trimmed, "let") || strings.Contains(trimmed, "const") {
				lines[i] = line + ";"
				changes++
			}
		}
		
		// Fix missing parentheses in function calls
		if strings.Contains(trimmed, "console.log") && !strings.Contains(trimmed, "console.log(") {
			re := regexp.MustCompile(`console\.log\s+(.+)`)
			if re.MatchString(trimmed) {
				lines[i] = re.ReplaceAllString(line, "console.log($1)")
				changes++
			}
		}
		
		// Fix missing quotes in strings
		if strings.Contains(trimmed, "=") && !strings.Contains(trimmed, "\"") && !strings.Contains(trimmed, "'") {
			// Basic string detection
			re := regexp.MustCompile(`=\s*([a-zA-Z_][a-zA-Z0-9_]*)\s*;?$`)
			if re.MatchString(trimmed) {
				lines[i] = re.ReplaceAllString(line, `= "$1";`)
				changes++
			}
		}
		
		// Fix missing braces in single-line if statements
		if strings.HasPrefix(trimmed, "if") && !strings.Contains(trimmed, "{") {
			if i < len(lines)-1 {
				nextLine := strings.TrimSpace(lines[i+1])
				if nextLine != "" && !strings.HasPrefix(nextLine, "else") {
					lines[i] = line + " {"
					lines[i+1] = "  " + lines[i+1]
					if i+2 < len(lines) {
						lines = append(lines[:i+2], append([]string{"}"}, lines[i+2:]...)...)
					} else {
						lines = append(lines, "}")
					}
					changes++
				}
			}
		}
	}
	
	return []byte(strings.Join(lines, "\n")), changes
}

// fixTypeScriptSyntax applies TypeScript-specific syntax fixes
func fixTypeScriptSyntax(source []byte) ([]byte, int) {
	// Apply JavaScript fixes first
	code, changes := fixJavaScriptSyntax(source)
	
	// Add TypeScript-specific fixes
	codeStr := string(code)
	lines := strings.Split(codeStr, "\n")
	
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		
		// Add type annotations to function parameters (basic)
		if strings.Contains(trimmed, "function") && strings.Contains(trimmed, "(") {
			re := regexp.MustCompile(`function\s+(\w+)\s*\(([^)]*)\)`)
			if re.MatchString(trimmed) {
				// This is a simplified type annotation addition
				lines[i] = re.ReplaceAllString(line, "function $1($2): void")
				changes++
			}
		}
		
		// Add missing type declarations
		if strings.Contains(trimmed, "let ") || strings.Contains(trimmed, "const ") {
			re := regexp.MustCompile(`(let|const)\s+(\w+)\s*=\s*(\d+)`)
			if re.MatchString(trimmed) {
				lines[i] = re.ReplaceAllString(line, "$1 $2: number = $3")
				changes++
			}
		}
	}
	
	return []byte(strings.Join(lines, "\n")), changes
}

// fixJavaScriptImports fixes import statements and adds missing imports
func fixJavaScriptImports(source []byte) ([]byte, int) {
	code := string(source)
	lines := strings.Split(code, "\n")
	changes := 0
	
	// Sort imports (basic implementation)
	var importLines []string
	var requireLines []string
	var otherLines []string
	var inImportSection bool
	
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		
		if strings.HasPrefix(trimmed, "import ") {
			importLines = append(importLines, line)
			inImportSection = true
		} else if strings.HasPrefix(trimmed, "const ") && strings.Contains(trimmed, "require(") {
			requireLines = append(requireLines, line)
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
	missingImports := findMissingJavaScriptImports(codeStr)
	
	if len(missingImports) > 0 {
		importLines = append(importLines, missingImports...)
		changes += len(missingImports)
	}
	
	// Reconstruct file with sorted imports
	var newLines []string
	
	// Add ES6 imports first
	for _, imp := range importLines {
		newLines = append(newLines, imp)
	}
	
	// Add require statements
	for _, req := range requireLines {
		newLines = append(newLines, req)
	}
	
	// Add empty line after imports if there are any
	if len(importLines) > 0 || len(requireLines) > 0 {
		newLines = append(newLines, "")
	}
	
	// Add the rest of the code
	newLines = append(newLines, otherLines...)
	
	return []byte(strings.Join(newLines, "\n")), changes
}

// findMissingJavaScriptImports detects commonly used modules that aren't imported
func findMissingJavaScriptImports(code string) []string {
	var missing []string
	
	// Common patterns that need imports
	patterns := map[string]string{
		"React.":     "import React from 'react';",
		"useState":   "import { useState } from 'react';",
		"useEffect": "import { useEffect } from 'react';",
		"axios.":     "import axios from 'axios';",
		"lodash.":    "import _ from 'lodash';",
		"moment":     "import moment from 'moment';",
	}
	
	for pattern, importStmt := range patterns {
		if strings.Contains(code, pattern) && !strings.Contains(code, importStmt) {
			missing = append(missing, importStmt)
		}
	}
	
	return missing
}
