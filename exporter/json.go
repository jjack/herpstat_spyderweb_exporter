package exporter

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-kit/log/level"
)

const (
	safetyRelayOff = "OFF (NORMAL OPERATION)"
	rampingOff     = "Not In Session"
)

type info struct {
	system  *system
	outputs *[]output
}

// information about the herpstat system itself
type system struct {
	Name        string          `json:"nickname"`
	IP          string          `json:"ip"`
	Mac         string          `json:"mac"`
	Firmware    json.RawMessage `json:"firmware"`
	SafetyRelay string          `json:"safetyrelay"`
	OutputCount float64         `json:"numberofoutputs"`
	PowerResets float64         `json:"powerresets"`
	Temp        float64         `json:"internaltemp"`
}

// information about an individual herpstat output
type output struct {
	ID            string  `json:"-"`
	Name          string  `json:"outputnickname"`
	Mode          string  `json:"outputmode"`
	Ramping       string  `json:"ramping,omitempty"`
	ErrorDesc     string  `json:"errorcodedescription,omitempty"`
	Power         float64 `json:"poweroutput,omitempty"`
	PowerLimit    float64 `json:"poweroutputLIMIT,omitempty"`
	ProbeTemp     float64 `json:"probereadingTEMP,omitempty"`
	ProbeHumidity float64 `json:"probereadingRH,omitempty"`
	AlarmEnabled  float64 `json:"enablehighlowalarm,omitempty"`
	AlarmHigh     float64 `json:"highalarm,omitempty"`
	AlarmLow      float64 `json:"lowalarm,omitempty"`
	RampEnd       float64 `json:"endoframpsetting,omitempty"`
	ErrorCode     float64 `json:"errorcode,omitempty"`
}

// UnmarshalJSON implements a custom JSON unmarshaler for our /RAWSTATUS data, which comes back in a format that's
// difficult to work with without making an arbitrary number of additional numbered [output] structs. Everything
// in "system" goes into [herpstat.exporter.info.system] and all of the numbered outputs (output1, output2, ...)
// are added to an array in [herpstat.exporter.info.outputs] in their numbered order.
//
//	{
//	  "system":{},
//	  "output1": {},
//	  "output2": {},
//	   ... etc ...
//	 }
func (info *info) UnmarshalJSON(data []byte) error {
	var mapped map[string]json.RawMessage

	level.Debug(logger).Log("msg", "unmarshaling raw data into map")

	if err := json.Unmarshal(data, &mapped); err != nil {
		level.Error(logger).Log("msg", "unable to unmarshal raw data into map")
		return err
	}

	level.Debug(logger).Log("msg", "unmarshaling system data")

	if err := json.Unmarshal(mapped["system"], &info.system); err != nil {
		level.Error(logger).Log("msg", "unable to unmarshal system data")
		return err
	}

	level.Debug(logger).Log(info.system)

	*info.outputs = make([]output, int(info.system.OutputCount))

	for key, value := range mapped {
		if key == "system" {
			continue
		}

		level.Debug(logger).Log("msg", "unmarshaling output data")

		id, err := strconv.Atoi(strings.TrimPrefix(key, "output"))
		if err != nil {
			return fmt.Errorf("%s doesn't look like 'output#' where # is a number: %s", key, err.Error())
		}

		if id > int(info.system.OutputCount) {
			return fmt.Errorf("output id %d is > than # of available outputs (%d)", id, int(info.system.OutputCount))
		}

		err = json.Unmarshal(value, &(*info.outputs)[id-1])
		if err != nil {
			level.Error(logger).Log("msg", fmt.Sprintf("unable to unmarshal output data for %s", key))
			return err
		}

		(*info.outputs)[id-1].ID = fmt.Sprintf("%d", id)

		level.Debug(logger).Log(&(*info.outputs)[id-1])
	}

	return nil
}

func (o *output) ramping() float64 {
	if o.Ramping == rampingOff {
		return 0
	}

	return 1
}

func (s *system) safetyrelay() float64 {
	if s.SafetyRelay == safetyRelayOff {
		return 0
	}

	return 1
}

func (s *system) infoLabelValues() []string {
	return []string{s.Name, s.IP, s.Mac, string(s.Firmware), fmt.Sprintf("%.0f", s.OutputCount)}
}

func (s *system) safetyRelayLabelValues() []string {
	return []string{s.Name, s.SafetyRelay}
}

func (s *system) labelValues() []string {
	return []string{s.Name}
}

func (o *output) infoLabelValues(system *string) []string {
	return []string{*system, o.ID, o.Name, o.Mode}
}

func (o *output) labelValues(system *string) []string {
	return []string{*system, o.ID}
}

func (o *output) errorLabelValues(system *string) []string {
	return []string{*system, o.ID, o.ErrorDesc}
}
