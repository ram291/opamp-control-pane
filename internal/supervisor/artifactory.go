package supervisor

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/open-telemetry/opamp-go/protobufs"
)

// BinaryUpgrader handles downloading and replacing collector binaries from Artifactory.
type BinaryUpgrader struct {
	config     *ArtifactoryConfig
	httpClient *http.Client
}

// NewBinaryUpgrader creates a new BinaryUpgrader.
func NewBinaryUpgrader(cfg *ArtifactoryConfig) *BinaryUpgrader {
	return &BinaryUpgrader{
		config:     cfg,
		httpClient: &http.Client{},
	}
}

// AvailableVersion represents a collector version available in Artifactory.
type AvailableVersion struct {
	Version     string
	ReleaseDate string
	SHA256      string
	Size        int64
}

// AvailableVersions lists available collector versions from Artifactory.
func (u *BinaryUpgrader) AvailableVersions(ctx context.Context) ([]AvailableVersion, error) {
	// This would query an Artifactory AQL (Artifactory Query Language) endpoint
	// or use the Artifactory REST API to list available versions.
	// For now, this is a placeholder that demonstrates the pattern.
	artifactPath := fmt.Sprintf("%s/%s/", u.config.BaseURL, u.config.RepoKey)

	req, err := http.NewRequestWithContext(ctx, "GET", artifactPath, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create artifact request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+u.config.APIToken)

	resp, err := u.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to query Artifactory: %w", err)
	}
	defer resp.Body.Close()

	// Parse response to extract versions (implementation depends on Artifactory API)
	// Placeholder return
	return []AvailableVersion{}, nil
}

// DownloadBinary downloads a collector binary from Artifactory, verifies its checksum,
// and saves it atomically.
func (u *BinaryUpgrader) DownloadBinary(ctx context.Context, version, expectedHash string) (string, error) {
	binaryName := fmt.Sprintf("otelcol_linux_amd64")
	downloadURL := fmt.Sprintf("%s/%s/%s/%s",
		u.config.BaseURL, u.config.RepoKey, version, binaryName)

	// Download to a temporary file
	tmpFile, err := os.CreateTemp("", "otelcol-*")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	tmpPath := tmpFile.Name()
	defer os.Remove(tmpPath)

	req, err := http.NewRequestWithContext(ctx, "GET", downloadURL, nil)
	if err != nil {
		tmpFile.Close()
		return "", fmt.Errorf("failed to create download request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+u.config.APIToken)

	resp, err := u.httpClient.Do(req)
	if err != nil {
		tmpFile.Close()
		return "", fmt.Errorf("failed to download binary: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		tmpFile.Close()
		return "", fmt.Errorf("download failed with status %d", resp.StatusCode)
	}

	// Compute SHA256 while writing
	hash := sha256.New()
	writer := io.MultiWriter(tmpFile, hash)

	if _, err := io.Copy(writer, resp.Body); err != nil {
		tmpFile.Close()
		return "", fmt.Errorf("failed to write binary: %w", err)
	}
	tmpFile.Close()

	// Verify checksum
	actualHash := hex.EncodeToString(hash.Sum(nil))
	if expectedHash != "" && actualHash != expectedHash {
		return "", fmt.Errorf("checksum mismatch: expected %s, got %s", expectedHash, actualHash)
	}

	return tmpPath, nil
}

// ReplaceBinary atomically replaces the collector binary.
// It moves the downloaded binary to the target location and makes it executable.
func (u *BinaryUpgrader) ReplaceBinary(downloadedPath, targetPath string) error {
	// Ensure target directory exists
	targetDir := filepath.Dir(targetPath)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return fmt.Errorf("failed to create target directory: %w", err)
	}

	// Atomic rename
	if err := os.Rename(downloadedPath, targetPath); err != nil {
		return fmt.Errorf("failed to replace binary: %w", err)
	}

	// Make executable
	if err := os.Chmod(targetPath, 0755); err != nil {
		return fmt.Errorf("failed to set executable permission: %w", err)
	}

	return nil
}

// ProcessPackageOffer processes a PackagesAvailable message from OpAMP server
// and downloads the collector binary from Artifactory.
func (u *BinaryUpgrader) ProcessPackageOffer(ctx context.Context, pkg *protobufs.PackagesAvailable, targetBinaryPath string) error {
	for name, available := range pkg.Packages {
		if name != "otelcol" {
			continue
		}

		version := available.GetVersion()
		var contentHash string
		if file := available.GetFile(); file != nil {
			contentHash = string(file.GetContentHash())
		}

		// Download the binary
		tmpPath, err := u.DownloadBinary(ctx, version, contentHash)
		if err != nil {
			return fmt.Errorf("failed to download version %s: %w", version, err)
		}

		// Atomically replace
		if err := u.ReplaceBinary(tmpPath, targetBinaryPath); err != nil {
			return fmt.Errorf("failed to replace binary with version %s: %w", version, err)
		}
	}

	return nil
}