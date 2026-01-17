package pool

import (
	"context"
	"log"
	"time"

	"github.com/anurag-327/neuron-core/conn"
	"github.com/anurag-327/neuron-core/pkg/logger"
	"github.com/anurag-327/neuron-core/runtime"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
)

// InitDockerPool registers and warms up all supported language pools.
//
// This function should be invoked once during application startup.
func InitDockerPool(ctx context.Context, client *conn.DockerClient) error {
	appLogger := logger.GetGlobalLogger()
	appLogger.Info(time.Now(), "Initializing sandbox container pools...", nil)

	// Clean up any orphaned containers from previous runs
	if err := cleanupOrphanedContainers(ctx, client); err != nil {
		appLogger.Error(time.Now(), "Failed to cleanup orphaned containers", map[string]interface{}{"error": err})
	}

	for _, cfg := range runtime.LanguageRegistry {
		Manager.Register(cfg.Language, PoolConfig{
			Image:          cfg.Image,
			InitSize:       cfg.InitSize,
			MaxSize:        cfg.MaxSize,
			HealthCmd:      cfg.HealthCmd,
			HealthInterval: cfg.HealthInterval,
			ResourceLimits: ResourceLimits{
				MemoryKB:    cfg.ResourceLimits.MemoryKB,
				TimeMs:      cfg.ResourceLimits.TimeMs,
				NanoCPUs:    cfg.ResourceLimits.NanoCPUs,
				Pids:        cfg.ResourceLimits.Pids,
				ULimits:     cfg.ResourceLimits.ULimits,
				NetworkMode: cfg.ResourceLimits.NetworkMode,
			},
		})
	}

	if err := Manager.InitAll(ctx); err != nil {
		return err
	}

	log.Println("Container pools warmed and ready!")
	return nil
}

// cleanupOrphanedContainers removes all containers created by Neuron
// that are still running from previous worker instances
func cleanupOrphanedContainers(ctx context.Context, client *conn.DockerClient) error {
	// Get all language images we use
	images := make(map[string]bool)
	for _, cfg := range runtime.LanguageRegistry {
		images[cfg.Image] = true
	}

	// List all containers (running and stopped) that match our images
	filterArgs := filters.NewArgs()
	for image := range images {
		filterArgs.Add("ancestor", image)
	}

	containers, err := client.ContainerList(ctx, container.ListOptions{
		All:     true, // Include stopped containers
		Filters: filterArgs,
	})
	if err != nil {
		return err
	}

	if len(containers) == 0 {
		log.Println("✅ No orphaned containers found")
		return nil
	}

	log.Printf("🧹 Found %d orphaned containers, cleaning up...", len(containers))

	// Remove each container
	removed := 0
	for _, c := range containers {
		err := client.ContainerRemove(ctx, c.ID, container.RemoveOptions{
			Force: true, // Force remove even if running
		})
		if err != nil {
			log.Printf("Failed to remove container %s: %v", c.ID[:12], err)
		} else {
			removed++
			log.Printf("Removed orphaned container: %s (image: %s)", c.ID[:12], c.Image)
		}
	}

	log.Printf("Cleaned up %d/%d orphaned containers", removed, len(containers))
	return nil
}
