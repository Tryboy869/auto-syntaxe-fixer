package fixer

import (
	"os"
	"path/filepath"
	"strings"
)

// DetectLanguages scans the repository and returns supported languages found
func DetectLanguages(repoPath string) []string {
	languageMap := make(map[string]bool)
	
	// Walk through all files in the repository
	filepath.Walk(repoPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Skip errors and continue
		}
		
		// Skip hidden directories and files
		if strings.HasPrefix(info.Name(), ".") {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		
		// Skip common non-source directories
		if info.IsDir() {
			skipDirs := []string{"node_modules", "vendor", "__pycache__", "target", "build", "dist"}
			for _, skipDir := range skipDirs {
				if info.Name() == skipDir {
					return filepath.SkipDir
				}
			}
			return nil
		}
		
		// Detect language based on file extension
		ext := strings.ToLower(filepath.Ext(path))
		
		switch ext {
		case ".go":
			languageMap["Go"] = true
		case ".py":
			languageMap["Python"] = true
		case ".js", ".jsx":
			languageMap["JavaScript"] = true
		case ".ts", ".tsx":
			languageMap["TypeScript"] = true
		case ".java":
			languageMap["Java"] = true
		case ".cpp", ".cc", ".cxx", ".c++":
			languageMap["C++"] = true
		case ".c":
			languageMap["C"] = true
		}
		
		return nil
	})
	
	// Convert map to slice
	var languages []string
	for lang := range languageMap {
		languages = append(languages, lang)
	}
	
	return languages
}

// GetFilesByExtension returns all files with specified extensions in the repository
func GetFilesByExtension(repoPath string, extensions []string) []string {
	var files []string
	
	// Convert extensions to lowercase for comparison
	extMap := make(map[string]bool)
	for _, ext := range extensions {
		extMap[strings.ToLower(ext)] = true
	}
	
	filepath.Walk(repoPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		
		// Skip hidden directories and files
		if strings.HasPrefix(info.Name(), ".") {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		
		// Skip common non-source directories
		if info.IsDir() {
			skipDirs := []string{"node_modules", "vendor", "__pycache__", "target", "build", "dist"}
			for _, skipDir := range skipDirs {
				if info.Name() == skipDir {
					return filepath.SkipDir
				}
			}
			return nil
		}
		
		// Check if file has one of the target extensions
		ext := strings.ToLower(filepath.Ext(path))
		if extMap[ext] {
			files = append(files, path)
		}
		
		return nil
	})
	
	return files
}

// IsSourceFile checks if a file is a source file based on its extension
func IsSourceFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	sourceExtensions := []string{
		".go", ".py", ".js", ".jsx", ".ts", ".tsx",
		".java", ".cpp", ".cc", ".cxx", ".c++", ".c",
		".h", ".hpp", ".rb", ".php", ".cs", ".swift",
		".kt", ".scala", ".rs", ".dart",
	}
	
	for _, sourceExt := range sourceExtensions {
		if ext == sourceExt {
			return true
		}
	}
	
	return false
}

// GetProjectStructure returns basic information about the project structure
func GetProjectStructure(repoPath string) map[string]int {
	structure := make(map[string]int)
	
	filepath.Walk(repoPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		
		if info.IsDir() {
			return nil
		}
		
		ext := strings.ToLower(filepath.Ext(path))
		if ext != "" {
			structure[ext]++
		}
		
		return nil
	})
	
	return structure
}
