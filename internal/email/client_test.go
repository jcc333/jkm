package email

import (
	"testing"
	"os"
	
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	
	"github.com/jclemer/jkm/internal/configure"
)

func TestEmailClientWithTestServer(t *testing.T) {
	// This test will be skipped unless we explicitly run in integration test mode
	if testing.Short() {
		t.Skip("Skipping integration test")
	}
	
	// The test expects a running GreenMail server from docker-compose
	// with the configured JKM_ environment variables
	
	// Check if we're running in test environment
	if os.Getenv("JKM_IMAP_SERVER") != "imap" {
		t.Skip("Skipping test: not running in Docker test environment")
	}
	
	// Load configuration
	config, err := configure.Load()
	require.NoError(t, err)
	require.NotEmpty(t, config.IMAPServer)
	
	t.Run("Can connect to test server", func(t *testing.T) {
		// Create a new email client
		client, err := NewClient(config)
		require.NoError(t, err)
		require.NotNil(t, client)
		
		// Test basic connectivity
		err = client.Connect()
		assert.NoError(t, err)
		
		// Clean up
		defer client.Close()
	})
}