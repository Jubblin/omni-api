package client

import (
	"fmt"
	"log"
	"os"

	"github.com/siderolabs/omni/client/pkg/client"
)

// NewOmniClient creates a new Omni client using environment variables.
func NewOmniClient() (*client.Client, error) {
	endpoint := os.Getenv("OMNI_ENDPOINT")
	if endpoint == "" {
		return nil, fmt.Errorf("OMNI_ENDPOINT environment variable is not set")
	}

	serviceAccount := os.Getenv("OMNI_SERVICE_ACCOUNT")
	if serviceAccount == "" {
		serviceAccount = os.Getenv("OMNI_SERVICE_ACCOUNT_KEY")
	}
	contextName := os.Getenv("OMNI_CONTEXT")
	identity := os.Getenv("OMNI_IDENTITY")
	keysDir := os.Getenv("OMNI_KEYS_DIR")
	
	log.Printf("Initializing Omni client with endpoint: %s\n", endpoint)
	
	var opts []client.Option

	if serviceAccount != "" {
		log.Println("Using Service Account authentication")
		opts = append(opts, client.WithServiceAccount(serviceAccount))
	} else if contextName != "" && identity != "" {
		log.Printf("Using PGP authentication (Context: %s, Identity: %s)\n", contextName, identity)
		opts = append(opts, client.WithUserAccount(contextName, identity))
		if keysDir != "" {
			log.Printf("Using custom keys directory: %s\n", keysDir)
			opts = append(opts, client.WithCustomKeysDir(keysDir))
		}
	} else {
		log.Println("Warning: No authentication method provided (Service Account or PGP)")
	}

	// You can add more options here, like insecure skip verify if needed
	if os.Getenv("OMNI_INSECURE") == "true" {
		opts = append(opts, client.WithInsecureSkipTLSVerify(true))
	}

	return client.New(endpoint, opts...)
}

