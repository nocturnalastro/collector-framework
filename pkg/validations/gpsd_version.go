package validations

// SPDX-License-Identifier: GPL-2.0-or-later

import (
	"strings"

	"github.com/nocturnalastro/collector-framework/pkg/collectors/devices"
)

const (
	gpsdID          = TGMEnvVerPath + "/gpsd/"
	gpsdDescription = "GPSD Version is valid"
	MinGSPDVersion  = "3.25"
)

func NewGPSDVersion(gpsdVer *devices.GPSVersions) *VersionCheck {
	parts := strings.Split(gpsdVer.GPSDVersion, " ")
	return &VersionCheck{
		id:           gpsdID,
		Version:      gpsdVer.GPSDVersion,
		checkVersion: strings.ReplaceAll(parts[0], "~", "-"),
		MinVersion:   MinGSPDVersion,
		description:  gpsdDescription,
		order:        gpsdVersionOrdering,
	}
}
