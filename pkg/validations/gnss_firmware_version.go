// SPDX-License-Identifier: GPL-2.0-or-later

package validations

import (
	"strings"

	"github.com/nocturnalastro/collector-framework/pkg/collectors/devices"
)

const (
	gnssID          = TGMEnvVerPath + "/gnss-firmware/"
	gnssDescription = "GNSS Version is valid"
)

var (
	MinGNSSVersion = "2.20"
)

func NewGNSS(gnss *devices.GPSVersions) *VersionCheck {
	parts := strings.Split(gnss.FirmwareVersion, " ")
	return &VersionCheck{
		id:           gnssID,
		Version:      gnss.FirmwareVersion,
		checkVersion: parts[1],
		MinVersion:   MinGNSSVersion,
		description:  gnssDescription,
		order:        gnssVersionOrdering,
	}
}
