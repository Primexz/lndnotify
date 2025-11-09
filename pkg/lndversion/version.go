package lndversion

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// LndVersion represents a parsed LND version
type LndVersion struct {
	Major      int    `json:"major"`
	Minor      int    `json:"minor"`
	Patch      int    `json:"patch"`
	PreRelease string `json:"preRelease"`
	Commit     string `json:"commit"`
	Raw        string `json:"raw"`
}

// String returns the string representation of the version
func (v LndVersion) String() string {
	if v.PreRelease != "" {
		return fmt.Sprintf("%d.%d.%d-%s", v.Major, v.Minor, v.Patch, v.PreRelease)
	}
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
}

// GitHubRelease represents a GitHub release response
type GitHubRelease struct {
	TagName string `json:"tag_name"`
	Name    string `json:"name"`
}

// ParseLndVersion parses an LND version string like "0.19.3-beta commit=v0.19.3-beta"
func ParseLndVersion(versionStr string) (*LndVersion, error) {
	if versionStr == "" {
		return nil, fmt.Errorf("empty version string")
	}

	// Clean up the input string
	versionStr = strings.TrimSpace(versionStr)

	// Extract commit information if present
	var commit string
	commitPattern := regexp.MustCompile(`commit=([^\s]+)`)
	if matches := commitPattern.FindStringSubmatch(versionStr); len(matches) > 1 {
		commit = matches[1]
		// Remove commit part from version string
		versionStr = commitPattern.ReplaceAllString(versionStr, "")
		versionStr = strings.TrimSpace(versionStr)
	}

	// Parse the main version part (e.g., "0.19.3-beta")
	versionPattern := regexp.MustCompile(`^v?(\d+)\.(\d+)\.(\d+)(?:-(.+))?`)
	matches := versionPattern.FindStringSubmatch(versionStr)

	if len(matches) < 4 {
		return nil, fmt.Errorf("invalid version format: %s", versionStr)
	}

	major, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil, fmt.Errorf("invalid major version: %s", matches[1])
	}

	minor, err := strconv.Atoi(matches[2])
	if err != nil {
		return nil, fmt.Errorf("invalid minor version: %s", matches[2])
	}

	patch, err := strconv.Atoi(matches[3])
	if err != nil {
		return nil, fmt.Errorf("invalid patch version: %s", matches[3])
	}

	preRelease := ""
	if len(matches) > 4 && matches[4] != "" {
		preRelease = matches[4]
	}

	return &LndVersion{
		Major:      major,
		Minor:      minor,
		Patch:      patch,
		PreRelease: preRelease,
		Commit:     commit,
		Raw:        versionStr,
	}, nil
}

// GetLatestVersionFromGitHub fetches the latest LND version from GitHub releases
func GetLatestVersionFromGitHub() (*LndVersion, error) {
	url := "https://api.github.com/repos/lightningnetwork/lnd/releases/latest"

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch latest version: %w", err)
	}
	defer resp.Body.Close() // nolint:errcheck

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned status code: %d", resp.StatusCode)
	}

	var release GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, fmt.Errorf("failed to decode GitHub response: %w", err)
	}

	return ParseLndVersion(release.TagName)
}

// CompareVersions compares two versions and returns:
// -1 if v1 < v2
//
//	0 if v1 == v2
//	1 if v1 > v2
func CompareVersions(v1, v2 *LndVersion) int {
	// Compare major version
	if v1.Major != v2.Major {
		if v1.Major < v2.Major {
			return -1
		}
		return 1
	}

	// Compare minor version
	if v1.Minor != v2.Minor {
		if v1.Minor < v2.Minor {
			return -1
		}
		return 1
	}

	// Compare patch version
	if v1.Patch != v2.Patch {
		if v1.Patch < v2.Patch {
			return -1
		}
		return 1
	}

	// Compare pre-release versions
	// No pre-release (stable) is greater than pre-release
	if v1.PreRelease == "" && v2.PreRelease != "" {
		return 1
	}
	if v1.PreRelease != "" && v2.PreRelease == "" {
		return -1
	}

	// Both have pre-release, compare lexicographically
	if v1.PreRelease != v2.PreRelease {
		if v1.PreRelease < v2.PreRelease {
			return -1
		}
		return 1
	}

	return 0
}

// IsOutdated checks if the local version is outdated compared to the latest GitHub version
func IsOutdated(localVersion *LndVersion) (bool, *LndVersion, error) {
	latestVersion, err := GetLatestVersionFromGitHub()
	if err != nil {
		return false, nil, fmt.Errorf("failed to get latest version: %w", err)
	}

	comparison := CompareVersions(localVersion, latestVersion)
	return comparison < 0, latestVersion, nil
}

// CheckVersion is a convenience function that parses a version string and checks if it's outdated
func CheckVersion(versionStr string) (bool, *LndVersion, *LndVersion, error) {
	localVersion, err := ParseLndVersion(versionStr)
	if err != nil {
		return false, nil, nil, fmt.Errorf("failed to parse local version: %w", err)
	}

	isOutdated, latestVersion, err := IsOutdated(localVersion)
	if err != nil {
		return false, localVersion, nil, err
	}

	return isOutdated, localVersion, latestVersion, nil
}
