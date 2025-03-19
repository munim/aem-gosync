package monitor

import (
	"fmt"
	"log"
	"path/filepath"

	"aem-gosync/internal/pkg"

	"github.com/fsnotify/fsnotify"
)

type DirectoryMonitor struct {
	directories   []string
	changedFiles  []string
	changesSignal chan struct{}
	Done          chan struct{} // Added to signal when monitoring is complete
}

func NewDirectoryMonitor(directories []string) *DirectoryMonitor {
	return &DirectoryMonitor{
		directories:   directories,
		changedFiles:  []string{},
		changesSignal: make(chan struct{}, 1),
		Done:          make(chan struct{}), // Initialize the Done channel
	}
}

func (dm *DirectoryMonitor) StartMonitoring() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	for _, dir := range dm.directories {
		err = watcher.Add(dir)
		if err != nil {
			log.Fatal(err)
		}
	}

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				dm.OnChange(event)
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
		<-dm.Done // Wait for the Done signal
	}()

	fmt.Println("Monitoring directories:", dm.directories)
	<-done
}

func (dm *DirectoryMonitor) OnChange(event fsnotify.Event) {
	if event.Op&fsnotify.Write == fsnotify.Write {
		fmt.Println("Modified file:", event.Name)

		relativePath, err := filepath.Rel("jcr_root", event.Name)
		if err != nil {
			log.Printf("Failed to determine relative path for %s: %v", event.Name, err)
			return
		}

		err = DeployToAEM(relativePath)
		if err != nil {
			log.Printf("Failed to deploy %s to AEM: %v", relativePath, err)
		} else {
			fmt.Printf("Successfully deployed %s to AEM\n", relativePath)
		}

		dm.changedFiles = append(dm.changedFiles, relativePath)
		select {
		case dm.changesSignal <- struct{}{}:
		default:
		}
	}
}

func (dm *DirectoryMonitor) WaitForChanges() []string {
	select {
	case <-dm.changesSignal:
		changed := dm.changedFiles
		dm.changedFiles = []string{}
		return changed
	case <-dm.Done: // Exit if monitoring is complete
		return nil
	}
}

func DeployToAEM(relativePath string) error {
	fmt.Printf("Creating package for %s...\n", relativePath)

	// Create a ContentPackage instance
	contentPackage := &pkg.AEMContentPackage{
		Name:    "example-package",
		Version: "1.0",
	}
	contentPackage.AddFiles([]string{relativePath})

	// Create the package
	packagePath := "package.zip"
	err := contentPackage.CreatePackage(packagePath)
	if err != nil {
		return fmt.Errorf("failed to create package: %w", err)
	}

	fmt.Printf("Deploying package %s to AEM...\n", packagePath)
	// Example: Use HTTP client to send the package to AEM's package manager endpoint
	// ...implementation details...

	return nil
}
