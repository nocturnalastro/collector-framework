// SPDX-License-Identifier: GPL-2.0-or-later

package collectors

import (
	"errors"
	"fmt"

	"github.com/nocturnalastro/collection-framework/pkg/collectors/contexts"
	"github.com/nocturnalastro/collection-framework/pkg/collectors/devices"
	"github.com/nocturnalastro/collection-framework/pkg/utils"
)

type DPLLFilesystemCollector struct {
	*ExecCollector
	interfaceName string
}

const (
	DPLLFilesystemCollectorName = "DPLL-Filesystem"
	DPLLInfo                    = "dpll-info-fs"
)

// Start sets up the collector so it is ready to be polled
func (dpll *DPLLFilesystemCollector) Start() error {
	dpll.running = true
	return nil
}

// polls for the dpll info then passes it to the callback
func (dpll *DPLLFilesystemCollector) poll() error {
	dpllInfo, err := devices.GetDevDPLLFilesystemInfo(dpll.ctx, dpll.interfaceName)

	if err != nil {
		return fmt.Errorf("failed to fetch %s %w", DPLLInfo, err)
	}
	err = dpll.callback.Call(&dpllInfo, DPLLInfo)
	if err != nil {
		return fmt.Errorf("callback failed %w", err)
	}
	return nil
}

// Poll collects information from the cluster then
// calls the callback.Call to allow that to persist it
func (dpll *DPLLFilesystemCollector) Poll(resultsChan chan PollResult, wg *utils.WaitGroupCount) {
	defer func() {
		wg.Done()
	}()
	errorsToReturn := make([]error, 0)
	err := dpll.poll()
	if err != nil {
		errorsToReturn = append(errorsToReturn, err)
	}
	resultsChan <- PollResult{
		CollectorName: DPLLFilesystemCollectorName,
		Errors:        errorsToReturn,
	}
}

// CleanUp stops a running collector
func (dpll *DPLLFilesystemCollector) CleanUp() error {
	dpll.running = false
	return nil
}

// Returns a new DPLLFilesystemCollector from the CollectionConstuctor Factory
func NewDPLLFilesystemCollector(constructor *CollectionConstructor) (Collector, error) {
	ctx, err := contexts.GetPTPDaemonContext(constructor.Clientset)
	if err != nil {
		return &DPLLFilesystemCollector{}, fmt.Errorf("failed to create DPLLFilesystemCollector: %w", err)
	}

	ptpArgs, ok := constructor.CollectorArgs["PTP"]
	if !ok {
		return &DPLLFilesystemCollector{}, errors.New("no PTP args in collector args")
	}
	ptpInterfaceRaw, ok := ptpArgs["PtpInterface"]
	if !ok {
		return &DPLLFilesystemCollector{}, errors.New("no PtpInterface in PTP collector args")
	}

	ptpInterface, ok := ptpInterfaceRaw.(string)
	if !ok {
		return &DPLLFilesystemCollector{}, errors.New("PTP interface is not a string")
	}

	err = devices.BuildFilesystemDPLLInfoFetcher(ptpInterface)
	if err != nil {
		return &DPLLFilesystemCollector{}, fmt.Errorf("failed to build fetcher for DPLLInfo %w", err)
	}

	collector := DPLLFilesystemCollector{
		ExecCollector: NewExecCollector(
			constructor.PollInterval,
			false,
			constructor.Callback,
			ctx,
		),
		interfaceName: ptpInterface,
	}
	return &collector, nil
}
