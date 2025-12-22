package util

import (
	configurables "hanamark/internal"
	"os"
	"path/filepath"
	"testing"
)

func TestFindTemplateUpward(t *testing.T) {
	// Create a temporary directory structure for testing
	tmpDir := t.TempDir()
	
	// Create directory structure:
	// tmpDir/
	//   templates/
	//     single.html
	//     home/
	//       blog/
	
	templatesDir := filepath.Join(tmpDir, "templates")
	homeDir := filepath.Join(templatesDir, "home")
	blogDir := filepath.Join(homeDir, "blog")
	
	// Create directories
	if err := os.MkdirAll(blogDir, 0755); err != nil {
		t.Fatalf("Failed to create test directories: %v", err)
	}
	
	// Create single.html in templates root
	singleTemplatePath := filepath.Join(templatesDir, "single.html")
	if err := os.WriteFile(singleTemplatePath, []byte("<html>test</html>"), 0644); err != nil {
		t.Fatalf("Failed to create single.html: %v", err)
	}
	
	// Test 1: Find template in current directory
	t.Run("FindInCurrentDir", func(t *testing.T) {
		result, err := FindTemplateUpward(templatesDir, templatesDir, "single.html")
		if err != nil {
			t.Errorf("Expected to find single.html, got error: %v", err)
		}
		if result != singleTemplatePath {
			t.Errorf("Expected path %s, got %s", singleTemplatePath, result)
		}
	})
	
	// Test 2: Find template in parent directory
	t.Run("FindInParentDir", func(t *testing.T) {
		result, err := FindTemplateUpward(homeDir, templatesDir, "single.html")
		if err != nil {
			t.Errorf("Expected to find single.html in parent, got error: %v", err)
		}
		if result != singleTemplatePath {
			t.Errorf("Expected path %s, got %s", singleTemplatePath, result)
		}
	})
	
	// Test 3: Find template in grandparent directory
	t.Run("FindInGrandparentDir", func(t *testing.T) {
		result, err := FindTemplateUpward(blogDir, templatesDir, "single.html")
		if err != nil {
			t.Errorf("Expected to find single.html in grandparent, got error: %v", err)
		}
		if result != singleTemplatePath {
			t.Errorf("Expected path %s, got %s", singleTemplatePath, result)
		}
	})
	
	// Test 4: File not found
	t.Run("FileNotFound", func(t *testing.T) {
		_, err := FindTemplateUpward(blogDir, templatesDir, "nonexistent.html")
		if err == nil {
			t.Error("Expected error for nonexistent file, got nil")
		}
	})
	
	// Test 5: Create list.html in home directory
	t.Run("FindListInHomeDir", func(t *testing.T) {
		listTemplatePath := filepath.Join(homeDir, "list.html")
		if err := os.WriteFile(listTemplatePath, []byte("<html>list</html>"), 0644); err != nil {
			t.Fatalf("Failed to create list.html: %v", err)
		}
		
		// Should find list.html in home directory when searching from blog
		result, err := FindTemplateUpward(blogDir, templatesDir, "list.html")
		if err != nil {
			t.Errorf("Expected to find list.html, got error: %v", err)
		}
		if result != listTemplatePath {
			t.Errorf("Expected path %s, got %s", listTemplatePath, result)
		}
	})
}
func TestCopyEmbedDir(t *testing.T) {
	tmpDir := t.TempDir()

	err := CopyEmbedDir(configurables.FS, "configurables", tmpDir)
	if err != nil {
		t.Fatalf("Failed to copy embed dir: %v", err)
	}

	// Check if _base.html exists in the destination
	basePath := filepath.Join(tmpDir, "templates", "_base.html")
	if _, err := os.Stat(basePath); os.IsNotExist(err) {
		t.Errorf("Expected _base.html to exist at %s, but it was not found", basePath)
	}
}
