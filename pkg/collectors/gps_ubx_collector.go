// SPDX-License-Identifier: GPL-2.0-or-later

package collectors //nolint:dupl // new collector

import (
	"errors"
	"fmt"

	"github.com/nocturnalastro/collector-framework/pkg/collectors/contexts"
	"github.com/nocturnalastro/collector-framework/pkg/collectors/devices"
	"github.com/nocturnalastro/collector-framework/pkg/utils"
)

var (
	GPSCollectorName = "GNSS"
	gpsNavKey        = "gpsNav"
)

type GPSCollector struct {
	*ExecCollector
	interfaceName string
}

func (gps *GPSCollector) poll() error {
	gpsNav, err := devices.GetGPSNav(gps.ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch  %s %w", gpsNavKey, err)
	}
	err = gps.callback.Call(&gpsNav, gpsNavKey)
	if err != nil {
		return fmt.Errorf("callback failed %w", err)
	}
	return nil
}

// Poll collects information from the cluster then
// calls the callback.Call to allow that to persist it
func (gps *GPSCollector) Poll(resultsChan chan PollResult, wg *utils.WaitGroupCount) {
	defer func() {
		wg.Done()
	}()

	errorsToReturn := make([]error, 0)
	err := gps.poll()
	if err != nil {
		errorsToReturn = append(errorsToReturn, err)
	}
	resultsChan <- PollResult{
		CollectorName: GPSCollectorName,
		Errors:        errorsToReturn,
	}
}

// Returns a new GPSCollector based on values in the CollectionConstructor
func NewGPSCollector(constructor *CollectionConstructor) (Collector, error) {
	ctx, err := contexts.GetPTPDaemonContext(constructor.Clientset)
	if err != nil {
		return &GPSCollector{}, fmt.Errorf("failed to create DPLLCollector: %w", err)
	}

	ptpArgs, ok := constructor.CollectorArgs["PTP"]
	if !ok {
		return &GPSCollector{}, errors.New("no PTP args in collector args")
	}
	ptpInterfaceRaw, ok := ptpArgs["PtpInterface"]
	if !ok {
		return &GPSCollector{}, errors.New("no PtpInterface in PTP collector args")
	}

	ptpInterface, ok := ptpInterfaceRaw.(string)
	if !ok {
		return &GPSCollector{}, errors.New("PTP interface is not a string")
	}

	collector := GPSCollector{
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

func init() {
	RegisterCollector(GPSCollectorName, NewGPSCollector, optional)
}
