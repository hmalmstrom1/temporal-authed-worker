package main_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.temporal.io/api/enums/v1"
	"go.temporal.io/api/workflowservice/v1"
	"go.temporal.io/sdk/client"
)

func TestSystem(t *testing.T) {
	// 1. Test User Creation (Kratos)
	t.Run("CreateTestUser", func(t *testing.T) {
		url := "http://localhost:4434/identities"
		email := fmt.Sprintf("testuser_%d@example.com", time.Now().Unix())

		body := map[string]interface{}{
			"schema_id": "default",
			"traits": map[string]interface{}{
				"email": email,
			},
		}
		jsonBody, _ := json.Marshal(body)

		resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonBody))
		require.NoError(t, err, "Failed to call Kratos Admin API")
		defer resp.Body.Close()

		assert.Equal(t, http.StatusCreated, resp.StatusCode, "Failed to create user in Kratos")
	})

	// 2. Test Frontend Reachability
	t.Run("FrontendReachability", func(t *testing.T) {
		url := "http://localhost:8080"
		resp, err := http.Get(url)
		require.NoError(t, err, "Failed to reach Temporal UI")
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode, "Temporal UI returned non-200 status")
	})

	// 3. Test Worker Connectivity
	t.Run("WorkerConnectivity", func(t *testing.T) {
		// Connect to Temporal Server
		c, err := client.Dial(client.Options{
			HostPort: "localhost:7233",
		})
		require.NoError(t, err, "Failed to create Temporal client")
		defer c.Close()

		// Verify Worker is polling 'my-task-queue'
		// We use DescribeTaskQueue to check if there are pollers
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		var resp *workflowservice.DescribeTaskQueueResponse
		for {
			resp, err = c.DescribeTaskQueue(ctx, "my-task-queue", enums.TASK_QUEUE_TYPE_ACTIVITY)
			require.NoError(t, err, "Failed to describe task queue")

			if len(resp.Pollers) > 0 {
				break
			}

			select {
			case <-ctx.Done():
				break
			case <-time.After(1 * time.Second):
				continue
			}
		}

		assert.NotEmpty(t, resp.Pollers, "No pollers found on 'my-task-queue'. Worker might not be connected.")
	})
	// 4. Test Namespace Registration
	t.Run("NamespaceRegistration", func(t *testing.T) {
		// Connect to Temporal Namespace Client
		nc, err := client.NewNamespaceClient(client.Options{
			HostPort: "localhost:7233",
		})
		require.NoError(t, err, "Failed to create Temporal namespace client")
		defer nc.Close()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		_, err = nc.Describe(ctx, "default")
		assert.NoError(t, err, "Default namespace should exist")
	})
}
