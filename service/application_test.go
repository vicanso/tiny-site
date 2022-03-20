package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestApplicationStatus(t *testing.T) {
	assert := assert.New(t)
	originalApplicationStatus := applicationStatus.Load()
	defer SetApplicationStatus(originalApplicationStatus)

	assert.Equal(originalApplicationStatus, GetApplicationStatus())

	SetApplicationStatus(ApplicationStatusStopped)
	assert.False(ApplicationIsRunning())
	SetApplicationStatus(ApplicationStatusStopping)
	assert.False(ApplicationIsRunning())
	SetApplicationStatus(ApplicationStatusRunning)
	assert.True(ApplicationIsRunning())
}

func TestApplicationVersion(t *testing.T) {
	assert := assert.New(t)
	originalVersion := applicationVersion
	defer SetApplicationVersion(originalVersion)

	v := "2020"
	SetApplicationVersion(v)
	assert.Equal(v, GetApplicationVersion())
}

func TestApplicationBuildedAt(t *testing.T) {
	assert := assert.New(t)
	originalBuildedAt := applicationBuildedAt
	defer SetApplicationBuildedAt(originalBuildedAt)

	buildedAt := "2020"
	SetApplicationBuildedAt(buildedAt)
	assert.Equal(buildedAt, GetApplicationBuildedAt())
}
