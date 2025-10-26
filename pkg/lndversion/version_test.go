package lndversion

import (
	"testing"
)

func TestParseLndVersion(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    *LndVersion
		expectError bool
	}{
		{
			name:  "basic version with beta",
			input: "0.19.3-beta commit=v0.19.3-beta",
			expected: &LndVersion{
				Major:      0,
				Minor:      19,
				Patch:      3,
				PreRelease: "beta",
				Commit:     "v0.19.3-beta",
				Raw:        "0.19.3-beta",
			},
		},
		{
			name:  "version with v prefix",
			input: "v0.18.0-beta",
			expected: &LndVersion{
				Major:      0,
				Minor:      18,
				Patch:      0,
				PreRelease: "beta",
				Commit:     "",
				Raw:        "v0.18.0-beta",
			},
		},
		{
			name:  "stable version without pre-release",
			input: "1.0.0",
			expected: &LndVersion{
				Major:      1,
				Minor:      0,
				Patch:      0,
				PreRelease: "",
				Commit:     "",
				Raw:        "1.0.0",
			},
		},
		{
			name:  "version with rc pre-release",
			input: "0.20.0-rc1 commit=v0.20.0-rc1",
			expected: &LndVersion{
				Major:      0,
				Minor:      20,
				Patch:      0,
				PreRelease: "rc1",
				Commit:     "v0.20.0-rc1",
				Raw:        "0.20.0-rc1",
			},
		},
		{
			name:        "empty string",
			input:       "",
			expectError: true,
		},
		{
			name:        "invalid format",
			input:       "invalid-version",
			expectError: true,
		},
		{
			name:        "missing patch version",
			input:       "0.19",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseLndVersion(tt.input)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if result.Major != tt.expected.Major {
				t.Errorf("Major: expected %d, got %d", tt.expected.Major, result.Major)
			}
			if result.Minor != tt.expected.Minor {
				t.Errorf("Minor: expected %d, got %d", tt.expected.Minor, result.Minor)
			}
			if result.Patch != tt.expected.Patch {
				t.Errorf("Patch: expected %d, got %d", tt.expected.Patch, result.Patch)
			}
			if result.PreRelease != tt.expected.PreRelease {
				t.Errorf("PreRelease: expected %s, got %s", tt.expected.PreRelease, result.PreRelease)
			}
			if result.Commit != tt.expected.Commit {
				t.Errorf("Commit: expected %s, got %s", tt.expected.Commit, result.Commit)
			}
			if result.Raw != tt.expected.Raw {
				t.Errorf("Raw: expected %s, got %s", tt.expected.Raw, result.Raw)
			}
		})
	}
}

func TestLndVersionString(t *testing.T) {
	tests := []struct {
		name     string
		version  LndVersion
		expected string
	}{
		{
			name: "version with pre-release",
			version: LndVersion{
				Major:      0,
				Minor:      19,
				Patch:      3,
				PreRelease: "beta",
			},
			expected: "0.19.3-beta",
		},
		{
			name: "stable version",
			version: LndVersion{
				Major: 1,
				Minor: 0,
				Patch: 0,
			},
			expected: "1.0.0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.version.String()
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestCompareVersions(t *testing.T) {
	tests := []struct {
		name     string
		v1       *LndVersion
		v2       *LndVersion
		expected int
	}{
		{
			name:     "v1 < v2 (major)",
			v1:       &LndVersion{Major: 0, Minor: 19, Patch: 3},
			v2:       &LndVersion{Major: 1, Minor: 0, Patch: 0},
			expected: -1,
		},
		{
			name:     "v1 > v2 (major)",
			v1:       &LndVersion{Major: 1, Minor: 0, Patch: 0},
			v2:       &LndVersion{Major: 0, Minor: 19, Patch: 3},
			expected: 1,
		},
		{
			name:     "v1 < v2 (minor)",
			v1:       &LndVersion{Major: 0, Minor: 18, Patch: 0},
			v2:       &LndVersion{Major: 0, Minor: 19, Patch: 0},
			expected: -1,
		},
		{
			name:     "v1 > v2 (minor)",
			v1:       &LndVersion{Major: 0, Minor: 19, Patch: 0},
			v2:       &LndVersion{Major: 0, Minor: 18, Patch: 0},
			expected: 1,
		},
		{
			name:     "v1 < v2 (patch)",
			v1:       &LndVersion{Major: 0, Minor: 19, Patch: 2},
			v2:       &LndVersion{Major: 0, Minor: 19, Patch: 3},
			expected: -1,
		},
		{
			name:     "v1 > v2 (patch)",
			v1:       &LndVersion{Major: 0, Minor: 19, Patch: 3},
			v2:       &LndVersion{Major: 0, Minor: 19, Patch: 2},
			expected: 1,
		},
		{
			name:     "equal versions",
			v1:       &LndVersion{Major: 0, Minor: 19, Patch: 3},
			v2:       &LndVersion{Major: 0, Minor: 19, Patch: 3},
			expected: 0,
		},
		{
			name:     "stable > pre-release",
			v1:       &LndVersion{Major: 0, Minor: 19, Patch: 3, PreRelease: ""},
			v2:       &LndVersion{Major: 0, Minor: 19, Patch: 3, PreRelease: "beta"},
			expected: 1,
		},
		{
			name:     "pre-release < stable",
			v1:       &LndVersion{Major: 0, Minor: 19, Patch: 3, PreRelease: "beta"},
			v2:       &LndVersion{Major: 0, Minor: 19, Patch: 3, PreRelease: ""},
			expected: -1,
		},
		{
			name:     "alpha < beta",
			v1:       &LndVersion{Major: 0, Minor: 19, Patch: 3, PreRelease: "alpha"},
			v2:       &LndVersion{Major: 0, Minor: 19, Patch: 3, PreRelease: "beta"},
			expected: -1,
		},
		{
			name:     "beta > alpha",
			v1:       &LndVersion{Major: 0, Minor: 19, Patch: 3, PreRelease: "beta"},
			v2:       &LndVersion{Major: 0, Minor: 19, Patch: 3, PreRelease: "alpha"},
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CompareVersions(tt.v1, tt.v2)
			if result != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, result)
			}
		})
	}
}

// Example test for the GitHub integration (this would require network access)
func TestGetLatestVersionFromGitHub(t *testing.T) {
	// Skip this test in normal CI/CD environments
	if testing.Short() {
		t.Skip("skipping GitHub API test in short mode")
	}

	version, err := GetLatestVersionFromGitHub()
	if err != nil {
		t.Logf("Failed to get latest version from GitHub: %v", err)
		t.Skip("GitHub API not accessible")
		return
	}

	if version == nil {
		t.Error("expected version but got nil")
		return
	}

	// Basic validation that we got a reasonable version
	if version.Major < 0 || version.Minor < 0 || version.Patch < 0 {
		t.Errorf("invalid version numbers: %d.%d.%d", version.Major, version.Minor, version.Patch)
	}

	t.Logf("Latest LND version from GitHub: %s", version.String())
}

func TestCheckVersion(t *testing.T) {
	// Test with a very old version that should be outdated
	isOutdated, localVer, latestVer, err := CheckVersion("0.1.0-beta")

	if err != nil {
		t.Logf("Failed to check version: %v", err)
		t.Skip("GitHub API not accessible")
		return
	}

	if !isOutdated {
		t.Error("expected version 0.1.0-beta to be outdated")
	}

	if localVer == nil {
		t.Error("expected local version but got nil")
	}

	if latestVer == nil {
		t.Error("expected latest version but got nil")
	}

	t.Logf("Local: %s, Latest: %s, Outdated: %v", localVer.String(), latestVer.String(), isOutdated)
}
