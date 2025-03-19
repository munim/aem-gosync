package deploy

import (
	"fmt"
	"net/http"
	"os"
	"strings"
)

type AEMDeployer struct {
	AEMURL      string
	Username    string
	Password    string
	ContentPath string
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

func (d *AEMDeployer) DeployPackage(packagePath string) error {
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
	req.Header.Set("Content-Type", "application/octet-stream")

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