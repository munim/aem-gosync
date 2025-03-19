package main

import (
	"aem-gosync/internal/deploy"  // Ensure this package exists and is correctly implemented
	"aem-gosync/internal/monitor" // Ensure this package exists and is correctly implemented
	"aem-gosync/internal/pkg"     // Ensure this package exists and is correctly implemented
	"flag"
	"fmt"
	"log"
	"strings"
)

func main() {
	// Parse command-line arguments for directory paths
	directories := flag.String("dirs", "", "Comma-separated list of directories to monitor")
	flag.Parse()

	if *directories == "" {
		log.Fatal("No directories provided to monitor")
	}

	// Initialize the directory monitor
	dirMonitor := monitor.NewDirectoryMonitor(strings.Split(*directories, ","))
	dirMonitor.Done = make(chan struct{}) // Ensure the Done channel is initialized
	go func() {
		dirMonitor.StartMonitoring()
		close(dirMonitor.Done) // Ensure the done channel is closed
	}()

	// Wait for changes and retrieve changed files
	fmt.Println("Monitoring started. Waiting for changes...")
	changedFiles := dirMonitor.WaitForChanges() // Updated to use WaitForChanges

	// Create a content package
	contentPackage := pkg.NewAEMContentPackage() // Fixed to use 'NewAEMContentPackage'
	contentPackage.AddFiles(changedFiles)

	// Create the AEM content package
	outputPath := "output.zip" // Define the output path for the package
	if err := contentPackage.CreatePackage(outputPath); err != nil {
		log.Fatalf("Failed to create content package: %v", err)
	}

	// Deploy the package to AEM
	aemDeployer, err := deploy.NewAEMDeployer() // Ensure NewAEMDeployer is defined in the 'deploy' package
	if err != nil {
		log.Fatalf("Failed to initialize AEM deployer: %v", err)
	}

	if err := aemDeployer.Authenticate(); err != nil {
		log.Fatalf("Authentication failed: %v", err)
	}

	if err := aemDeployer.DeployPackage(contentPackage); err != nil {
		log.Fatalf("Failed to deploy package: %v", err)
	}

	fmt.Println("Content package deployed successfully")
}
