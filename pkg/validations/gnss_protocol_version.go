// SPDX-License-Identifier: GPL-2.0-or-later

package validations

import (
	"github.com/nocturnalastro/collector-framework/pkg/collectors/devices"
)

const (
	gnssProtID           = TGMEnvVerPath + "/gnss-protocol/"
	gnssProtIDescription = "GNSS protocol version is valid"
	MinProtoVersion      = "29.20"
)

func NewGNSSProtocol(gnss *devices.GPSVersions) *VersionCheck {
	return &VersionCheck{
		id:           gnssProtID,
		Version:      gnss.ProtoVersion,
		checkVersion: gnss.ProtoVersion,
		MinVersion:   MinProtoVersion,
		description:  gnssProtIDescription,
		order:        gnssProtOrdering,
	}
}
