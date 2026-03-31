// Package version holds build-time metadata (overridden via -ldflags).
package version

// Version is the semantic version (e.g. from git tag).
var Version = "0.0.0-dev"

// Commit is the VCS revision.
var Commit = "unknown"

// BuildDate is the build timestamp (RFC3339 recommended).
var BuildDate = "unknown"
