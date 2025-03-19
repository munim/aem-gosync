package package

import (
    "archive/zip"
    "fmt"
    "io"
    "os"
    "path/filepath"
)

type ContentPackage struct {
    Name        string
    Version     string
    Files       []string
}

func (cp *ContentPackage) CreatePackage(outputPath string) error {
    fmt.Printf("Creating package %s...\n", outputPath)

    // Create the ZIP file
    zipFile, err := os.Create(outputPath)
    if (err != nil) {
        return fmt.Errorf("failed to create ZIP file: %w", err)
    }
    defer zipFile.Close()

    zipWriter := zip.NewWriter(zipFile)
    defer zipWriter.Close()

    // Add files to the ZIP under the typical AEM package structure
    for _, relativePath := range cp.Files {
        jcrRootPath := filepath.Join("jcr_root", relativePath)
        fileToAdd, err := os.Open(jcrRootPath)
        if err != nil {
            return fmt.Errorf("failed to open file %s: %w", jcrRootPath, err)
        }
        defer fileToAdd.Close()

        zipEntry, err := zipWriter.Create(jcrRootPath)
        if err != nil {
            return fmt.Errorf("failed to create ZIP entry for %s: %w", jcrRootPath, err)
        }

        _, err = io.Copy(zipEntry, fileToAdd)
        if err != nil {
            return fmt.Errorf("failed to write file %s to ZIP: %w", jcrRootPath, err)
        }
    }

    // Add META-INF/vault directory and files
    metaInfVaultPath := "META-INF/vault"
    metaInfFiles := map[string]string{
        "filter.xml": cp.generateFilterXML(),
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

func (cp *ContentPackage) AddFiles(files []string) {
    cp.Files = append(cp.Files, files...)
}

func (cp *ContentPackage) generateFilterXML() string {
    filters := ""
    for _, file := range cp.Files {
        filters += fmt.Sprintf(`<filter root="/%s"/>`, file)
    }
    return fmt.Sprintf(`<workspaceFilter version="1.0">%s</workspaceFilter>`, filters)
}