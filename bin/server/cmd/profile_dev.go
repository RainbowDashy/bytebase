//go:build !release
// +build !release

package cmd

import (
	"fmt"
	"time"
)

func activeProfile(dataDir string, port int, isDemo bool) Profile {
	dsn := fmt.Sprintf("file:%s/bytebase_dev.db", dataDir)
	if isDemo {
		dsn = fmt.Sprintf("file:%s/bytebase_demo.db", dataDir)
	}
	return Profile{
		mode:                 "dev",
		port:                 port,
		dsn:                  dsn,
		seedDir:              "seed/test",
		forceResetSeed:       true,
		backupRunnerInterval: 10 * time.Second,
	}
}

// GetTestProfile will return a profile for testing.
func GetTestProfile(dataDir string) Profile {
	return Profile{
		mode:                 "dev",
		port:                 1234,
		dsn:                  fmt.Sprintf("file:%s/bytebase_test.db", dataDir),
		seedDir:              "seed/test",
		forceResetSeed:       true,
		backupRunnerInterval: 10 * time.Second,
	}
}
