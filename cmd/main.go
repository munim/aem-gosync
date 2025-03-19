package main

import (
	"flag"
	"fmt"
	"log"
	"aem-deployer/internal/monitor"
	"aem-deployer/internal/package"
	"aem-deployer/internal/deploy"
)

func main() {
	// Parse command-line arguments for directory paths
	directories := flag.String("dirs", "", "Comma-separated list of directories to monitor")
	flag.Parse()

	if *directories == "" {
		log.Fatal("No directories provided to monitor")
	}

	// Initialize the directory monitor
	dirMonitor := monitor.NewDirectoryMonitor(*directories)
	go dirMonitor.StartMonitoring()

	// Create a content package
	contentPackage := package.NewContentPackage()
	contentPackage.AddFiles(dirMonitor.GetChangedFiles())

	// Create the AEM content package
	if err := contentPackage.CreatePackage(); err != nil {
		log.Fatalf("Failed to create content package: %v", err)
	}

	// Deploy the package to AEM
	aemDeployer := deploy.NewAEMDeployer()
	if err := aemDeployer.Authenticate(); err != nil {
		log.Fatalf("Authentication failed: %v", err)
	}

	if err := aemDeployer.DeployPackage(contentPackage); err != nil {
		log.Fatalf("Failed to deploy package: %v", err)
	}

	fmt.Println("Content package deployed successfully")
}