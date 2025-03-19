package pkg

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type AEMContentPackage struct { // Renamed for clarity
	Name    string
	Version string
	Files   []string
}

func NewAEMContentPackage() *AEMContentPackage { // Added constructor function
	return &AEMContentPackage{
		Name:    "default-package",
		Version: "1.0",
		Files:   []string{},
	}
}

func (pkg *AEMContentPackage) CreatePackage(outputPath string) error { // Updated receiver name
	fmt.Printf("Creating package %s...\n", outputPath)

	// Create the ZIP file
	zipFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create ZIP file: %w", err)
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// Add files to the ZIP under the typical AEM package structure
	for _, relativePath := range pkg.Files {
		jcrRootPath := filepath.Join("jcr_root", relativePath)
		fileToAdd, err := os.Open(relativePath) // Use the correct relative path
		if err != nil {
			return fmt.Errorf("failed to open file %s: %w", relativePath, err)
		}

		zipEntry, err := zipWriter.Create(jcrRootPath)
		if err != nil {
			fileToAdd.Close() // Ensure the file is closed before returning
			return fmt.Errorf("failed to create ZIP entry for %s: %w", jcrRootPath, err)
		}

		_, err = io.Copy(zipEntry, fileToAdd)
		fileToAdd.Close() // Close the file after copying
		if err != nil {
			return fmt.Errorf("failed to write file %s to ZIP: %w", jcrRootPath, err)
		}
	}

	// Add META-INF/vault directory and files
	metaInfVaultPath := "META-INF/vault"
	metaInfFiles := map[string]string{
		"filter.xml":     pkg.generateFilterXML(),
		"properties.xml": `<properties><entry key="packageType">content</entry></properties>`,
	}

	for fileName, content := range metaInfFiles {
		filePath := filepath.Join(metaInfVaultPath, fileName)
		zipEntry, err := zipWriter.Create(filePath)
		if err != nil {
			return fmt.Errorf("failed to create ZIP entry for %s: %w", filePath, err)
		}

		_, err = zipEntry.Write([]byte(content))
		if err != nil {
			return fmt.Errorf("failed to write content to %s: %w", filePath, err)
		}
	}

	fmt.Printf("Successfully created package %s\n", outputPath)
	return nil
}

func (pkg *AEMContentPackage) AddFiles(files []string) { // Updated receiver name
	pkg.Files = append(pkg.Files, files...)
}

func (pkg *AEMContentPackage) generateFilterXML() string { // Updated receiver name
	filters := ""
	for _, file := range pkg.Files {
		filters += fmt.Sprintf(`<filter root="/%s"/>`, file)
	}
	return fmt.Sprintf(`<workspaceFilter version="1.0">%s</workspaceFilter>`, filters)
}
