// SPDX-License-Identifier: GPL-2.0-or-later

package collectors

import (
	"errors"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/nocturnalastro/collection-framework/pkg/collectors/contexts"
	"github.com/nocturnalastro/collection-framework/pkg/collectors/devices"
)

const (
	DPLLCollectorName = "DPLL"
)

// Returns a new DPLLCollector from the CollectionConstuctor Factory
func NewDPLLCollector(constructor *CollectionConstructor) (Collector, error) {
	ctx, err := contexts.GetPTPDaemonContext(constructor.Clientset)
	if err != nil {
		return &DPLLNetlinkCollector{}, fmt.Errorf("failed to create DPLLCollector: %w", err)
	}

	ptpArgs, ok := constructor.CollectorArgs["PTP"]
	if !ok {
		return &DPLLNetlinkCollector{}, errors.New("no PTP args in collector args")
	}
	ptpInterfaceRaw, ok := ptpArgs["PtpInterface"]
	if !ok {
		return &DPLLNetlinkCollector{}, errors.New("no PtpInterface in PTP collector args")
	}

	ptpInterface, ok := ptpInterfaceRaw.(string)
	if !ok {
		return &DPLLNetlinkCollector{}, errors.New("PTP interface is not a string")
	}

	dpllFSExists, err := devices.IsDPLLFileSystemPresent(ctx, ptpInterface)
	log.Debug("DPLL FS exists: ", dpllFSExists)
	if dpllFSExists && err == nil {
		return NewDPLLFilesystemCollector(constructor)
	} else {
		return NewDPLLNetlinkCollector(constructor)
	}
}

func init() {
	RegisterCollector(DPLLCollectorName, NewDPLLCollector, optional)
}
