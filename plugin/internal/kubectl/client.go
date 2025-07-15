package kubectl

import (
	"fmt"
	"kubectl-login/internal/config"
	"os/exec"
)

type Client struct {
	cfg *config.Config
}

func NewClient(cfg *config.Config) *Client {
	return &Client{cfg: cfg}
}

func (c *Client) CreateConfig(token string) error {
	cmd := exec.Command("kubectl", "config", "set-credentials", c.cfg.Username,
		"--auth-provider=oidc",
		"--auth-provider-arg=idp-issuer-url="+c.cfg.Server,
		"--auth-provider-arg=client-id=kubernetes",
		"--auth-provider-arg=id-token="+token)

	if err := cmd.Run(); err != nil {
		fmt.Println("Error setting creds", token)
		return err
	}

	cmd = exec.Command("kubectl", "config", "set-context", "--current", "--user="+c.cfg.Username)
	if err := cmd.Run(); err != nil {
		fmt.Println("error setting context")
		return err
	}

	return nil
}
