package deploy

import (
	"fmt"
	"net/http"
	"os"

	"aem-gosync/internal/pkg" // Import the pkg package
)

type AEMDeployer struct {
	AEMURL      string
	Username    string
	Password    string
	ContentPath string
}

// NewAEMDeployer initializes and returns a new AEMDeployer instance.
func NewAEMDeployer() (*AEMDeployer, error) {
	// Initialize the deployer with necessary configurations.
	// Replace the following with actual initialization logic.
	return &AEMDeployer{}, nil
}

func (d *AEMDeployer) Authenticate() error {
	req, err := http.NewRequest("GET", d.AEMURL, nil)
	if err != nil {
		return err
	}
	req.SetBasicAuth(d.Username, d.Password)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("authentication failed: %s", resp.Status)
	}
	return nil
}

func (d *AEMDeployer) DeployPackage(contentPackage *pkg.AEMContentPackage) error { // Updated to accept AEMContentPackage
	packagePath := contentPackage.Name + ".zip" // Generate the package path
	file, err := os.Open(packagePath)
	if err != nil {
		return err
	}
	defer file.Close()

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/bin/deploy", d.AEMURL), file)
	if err != nil {
		return err
	}
	req.SetBasicAuth(d.Username, d.Password)
	req.Header.Set("Content-Type", "application/zip") // Correct content type for ZIP files

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("deployment failed: %s", resp.Status)
	}
	return nil
}
