package git

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// CloneRepo clones a GitHub repository to the specified directory
func CloneRepo(repoURL, targetDir, token, branch string) (string, error) {
	// Create target directory if it doesn't exist
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create target directory: %v", err)
	}
	
	// Extract repo name from URL
	repoName := extractRepoName(repoURL)
	repoPath := filepath.Join(targetDir, repoName)
	
	// Remove existing directory if it exists
	os.RemoveAll(repoPath)
	
	// Prepare clone command
	var cmd *exec.Cmd
	
	if token != "" {
		// Use token for authentication
		authenticatedURL := addTokenToURL(repoURL, token)
		cmd = exec.Command("git", "clone", "--depth", "1", "-b", branch, authenticatedURL, repoPath)
	} else {
		// Clone without authentication (public repo)
		cmd = exec.Command("git", "clone", "--depth", "1", "-b", branch, repoURL, repoPath)
	}
	
	// Execute clone command
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to clone repository: %v, stderr: %s", err, stderr.String())
	}
	
	return repoPath, nil
}

// CreateBranchAndCommit creates a new branch and commits the changes
func CreateBranchAndCommit(repoPath, branchName string, modifiedFiles []string, token string) error {
	// Change to repository directory
	originalDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %v", err)
	}
	defer os.Chdir(originalDir)
	
	if err := os.Chdir(repoPath); err != nil {
		return fmt.Errorf("failed to change to repo directory: %v", err)
	}
	
	// Configure git user (required for commits)
	if err := configureGitUser(); err != nil {
		return fmt.Errorf("failed to configure git user: %v", err)
	}
	
	// Create and checkout new branch
	if err := createBranch(branchName); err != nil {
		return fmt.Errorf("failed to create branch: %v", err)
	}
	
	// Add modified files to staging
	if err := addFiles(modifiedFiles); err != nil {
		return fmt.Errorf("failed to add files: %v", err)
	}
	
	// Commit changes
	commitMessage := fmt.Sprintf("Auto-fix: Applied syntax fixes to %d files", len(modifiedFiles))
	if err := commitChanges(commitMessage); err != nil {
		return fmt.Errorf("failed to commit changes: %v", err)
	}
	
	// Push branch to remote (if token is provided)
	if token != "" {
		if err := pushBranch(branchName, token); err != nil {
			return fmt.Errorf("failed to push branch: %v", err)
		}
	}
	
	return nil
}

// extractRepoName extracts the repository name from a GitHub URL
func extractRepoName(repoURL string) string {
	// Remove .git suffix if present
	repoURL = strings.TrimSuffix(repoURL, ".git")
	
	// Split by / and get the last part
	parts := strings.Split(repoURL, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	
	return "repo"
}

// addTokenToURL adds authentication token to GitHub URL
func addTokenToURL(repoURL, token string) string {
	if strings.HasPrefix(repoURL, "https://github.com/") {
		return strings.Replace(repoURL, "https://github.com/", fmt.Sprintf("https://%s@github.com/", token), 1)
	}
	return repoURL
}

// configureGitUser sets up git user configuration for commits
func configureGitUser() error {
	// Set git user name
	cmd := exec.Command("git", "config", "user.name", "Auto-Syntax-Fixer")
	if err := cmd.Run(); err != nil {
		return err
	}
	
	// Set git user email
	cmd = exec.Command("git", "config", "user.email", "auto-syntax-fixer@example.com")
	if err := cmd.Run(); err != nil {
		return err
	}
	
	return nil
}

// createBranch creates and checks out a new branch
func createBranch(branchName string) error {
	// Create new branch
	cmd := exec.Command("git", "checkout", "-b", branchName)
	
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git checkout failed: %v, stderr: %s", err, stderr.String())
	}
	
	return nil
}

// addFiles adds specified files to git staging area
func addFiles(files []string) error {
	if len(files) == 0 {
		return nil
	}
	
	// Add all modified files
	args := append([]string{"add"}, files...)
	cmd := exec.Command("git", args...)
	
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git add failed: %v, stderr: %s", err, stderr.String())
	}
	
	return nil
}

// commitChanges commits the staged changes
func commitChanges(message string) error {
	cmd := exec.Command("git", "commit", "-m", message)
	
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git commit failed: %v, stderr: %s", err, stderr.String())
	}
	
	return nil
}

// pushBranch pushes the branch to remote repository
func pushBranch(branchName, token string) error {
	// Get remote origin URL
	cmd := exec.Command("git", "remote", "get-url", "origin")
	var out bytes.Buffer
	cmd.Stdout = &out
	
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to get remote URL: %v", err)
	}
	
	remoteURL := strings.TrimSpace(out.String())
	
	// Add token to remote URL if not already present
	if token != "" && !strings.Contains(remoteURL, "@") {
		remoteURL = addTokenToURL(remoteURL, token)
	}
	
	// Push branch
	cmd = exec.Command("git", "push", "-u", remoteURL, branchName)
	
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git push failed: %v, stderr: %s", err, stderr.String())
	}
	
	return nil
}

// GetRepoStatus returns the current status of the repository
func GetRepoStatus(repoPath string) (string, error) {
	originalDir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	defer os.Chdir(originalDir)
	
	if err := os.Chdir(repoPath); err != nil {
		return "", err
	}
	
	cmd := exec.Command("git", "status", "--porcelain")
	var out bytes.Buffer
	cmd.Stdout = &out
	
	if err := cmd.Run(); err != nil {
		return "", err
	}
	
	return out.String(), nil
}

// HasUncommittedChanges checks if there are any uncommitted changes
func HasUncommittedChanges(repoPath string) (bool, error) {
	status, err := GetRepoStatus(repoPath)
	if err != nil {
		return false, err
	}
	
	return strings.TrimSpace(status) != "", nil
}

// GetCurrentBranch returns the current branch name
func GetCurrentBranch(repoPath string) (string, error) {
	originalDir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	defer os.Chdir(originalDir)
	
	if err := os.Chdir(repoPath); err != nil {
		return "", err
	}
	
	cmd := exec.Command("git", "branch", "--show-current")
	var out bytes.Buffer
	cmd.Stdout = &out
	
	if err := cmd.Run(); err != nil {
		return "", err
	}
	
	return strings.TrimSpace(out.String()), nil
}
