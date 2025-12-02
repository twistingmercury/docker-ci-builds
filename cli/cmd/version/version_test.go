/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package version_test

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/twistingmercury/cli/cmd/version"
)

func TestSimpleExample(t *testing.T) {
	cmd := version.VersionCmd()
	require.NotNil(t, cmd, "VersionCmd() should not return nil")
	assert.IsType(t, &cobra.Command{}, cmd, "VersionCmd() should return a *cobra.Command")
}
