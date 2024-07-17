package pwrstat

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	modelNameLabel        = "Model Name"
	firmwareNumberLabel   = "Firmware Number"
	ratingVoltageLabel    = "Rating Voltage"
	ratingPowerLabel      = "Rating Power"
	stateLabel            = "State"
	powerSupplyByLabel    = "Power Supply by"
	utilityVoltageLabel   = "Utility Voltage"
	outputVoltageLabel    = "Output Voltage"
	batteryCapacityLabel  = "Battery Capacity"
	remainingRuntimeLabel = "Remaining Runtime"
	loadLabel             = "Load"
	lineInteractionLabel  = "Line Interaction"
	testResultLabel       = "Test Result"
	lastPowerEventLabel   = "Last Power Event"
)

var (
	lineRegexp        = regexp.MustCompile(`([a-zA-Z ]+)\.\.+ ([a-zA-Z0-9 \(\)%/:\.\-]+)`)
	wattRegexp        = regexp.MustCompile(`([0-9]+) Watt`)
	vaRegexp          = regexp.MustCompile(`([0-9]+) VA`)
	loadRegexp        = regexp.MustCompile(`([0-9]+) Watt\(([0-9]+) %\)`)
	eventRegexp       = regexp.MustCompile(`([a-zA-Z]+) at (\d\d\d\d/\d\d/\d\d \d\d:\d\d:\d\d)`)
	eventTimestampFmt = "2006/01/02 15:04:05"
)

type StatusResult struct {
	Properties    properties    `json:"properties"`
	CurrentStatus currentStatus `json:"current_status"`
}

type properties struct {
	ModelName      string      `json:"model_name"`
	FirmwareNumber string      `json:"firmware_number"`
	RatingVoltage  int         `json:"rating_voltage"`
	RatingPower    ratingPower `json:"rating_power"`
}

type ratingPower struct {
	Watts   int `json:"watts"`
	VoltAmp int `json:"volt_amps"`
}

type currentStatus struct {
	State            string         `json:"state"`
	PowerSupplyBy    string         `json:"power_supply_by"`
	UtilityVoltage   int            `json:"utility_voltage"`
	OutputVoltage    int            `json:"output_voltage"`
	BatteryCapacity  float64        `json:"battery_capacity_pct"`
	RemainingRuntime int            `json:"remaining_runtime_minutes"`
	Load             load           `json:"load"`
	LineInteraction  string         `json:"line_interaction"`
	TestResult       testResult     `json:"test_result"`
	LastPowerEvent   lastPowerEvent `json:"last_power_event"`
}

type load struct {
	Watts   int     `json:"watts"`
	Percent float64 `json:"percent"`
}

type testResult = event
type lastPowerEvent = event

type event struct {
	State string `json:"state"`
	// Timestamp of the last test, if available.
	// This will not be available if a test is in progress or has not been run.
	// It will be the zero value if timestamp is unavailable.
	Timestamp time.Time `json:"timestamp"`
}

func parseStatus(status string) (s StatusResult, err error) {
	lines := strings.Split(status, "\n")
	for _, l := range lines {
		if !strings.Contains(l, "...") {
			continue
		}
		matches := lineRegexp.FindStringSubmatch(strings.TrimSpace(l))
		name, value := matches[1], matches[2]
		switch name {
		case modelNameLabel:
			s.Properties.ModelName = value
		case firmwareNumberLabel:
			s.Properties.FirmwareNumber = value
		case ratingVoltageLabel:
			v, err := parseVolts(value)
			if err != nil {
				return s, fmt.Errorf("failed to parse rating voltage %q: %w", value, err)
			}
			s.Properties.RatingVoltage = v
		case ratingPowerLabel:
			w, va, err := parseRatingPower(value)
			if err != nil {
				return s, fmt.Errorf("failed to parse rating power %q: %w", value, err)
			}
			s.Properties.RatingPower.Watts = w
			s.Properties.RatingPower.VoltAmp = va
		case stateLabel:
			s.CurrentStatus.State = value
		case powerSupplyByLabel:
			s.CurrentStatus.PowerSupplyBy = value
		case utilityVoltageLabel:
			v, err := parseVolts(value)
			if err != nil {
				return s, fmt.Errorf("failed to parse utility voltage %q: %w", value, err)
			}
			s.CurrentStatus.UtilityVoltage = v
		case outputVoltageLabel:
			v, err := parseVolts(value)
			if err != nil {
				return s, fmt.Errorf("failed to parse output voltage %q: %w", value, err)
			}
			s.CurrentStatus.OutputVoltage = v
		case batteryCapacityLabel:
			v, err := parsePercent(value)
			if err != nil {
				return s, fmt.Errorf("failed to parse battery capacity %q: %w", value, err)
			}
			s.CurrentStatus.BatteryCapacity = v
		case remainingRuntimeLabel:
			v, err := parseMinutes(value)
			if err != nil {
				return s, fmt.Errorf("failed to parse remaining runtime %q: %w", value, err)
			}
			s.CurrentStatus.RemainingRuntime = v
		case loadLabel:
			w, p, err := parseLoad(value)
			if err != nil {
				return s, fmt.Errorf("failed to parse load %q: %w", value, err)
			}
			s.CurrentStatus.Load.Watts = w
			s.CurrentStatus.Load.Percent = p
		case lineInteractionLabel:
			s.CurrentStatus.LineInteraction = value
		case testResultLabel:
			e, err := parseEvent(value)
			if err != nil {
				return s, fmt.Errorf("failed to parse test result %q: %w", value, err)
			}
			s.CurrentStatus.TestResult = e
		case lastPowerEventLabel:
			e, err := parseEvent(value)
			if err != nil {
				return s, fmt.Errorf("failed to parse last power event %q: %w", value, err)
			}
			s.CurrentStatus.LastPowerEvent = e
		}
	}
	return s, err
}

func parseVolts(value string) (int, error) {
	s := strings.Split(value, " ")
	if len(s) != 2 || s[1] != "V" {
		return 0, errors.New("value not in volts")
	}
	return strconv.Atoi(s[0])
}

func parseRatingPower(value string) (int, int, error) {
	m := wattRegexp.FindStringSubmatch(value)
	if len(m) != 2 {
		return 0, 0, errors.New("can't find rating power watts")
	}
	w, err := strconv.Atoi(m[1])
	if err != nil {
		return 0, 0, fmt.Errorf("failed to parse watt value: %w", err)
	}
	m = vaRegexp.FindStringSubmatch(value)
	if len(m) != 2 {
		return w, 0, errors.New("can't find rating power va")
	}
	va, err := strconv.Atoi(m[1])
	if err != nil {
		return w, 0, fmt.Errorf("failed to parse va value: %w", err)
	}
	return w, va, nil
}

func parsePercent(value string) (float64, error) {
	s := strings.Split(value, " ")
	if len(s) != 2 || s[1] != "%" {
		return 0, errors.New("value not in percentage")
	}
	f, err := strconv.ParseFloat(s[0], 64)
	if err != nil {
		return 0, err
	}
	return f / 100, err
}

func parseMinutes(value string) (int, error) {
	s := strings.Split(value, " ")
	if len(s) != 2 || !strings.HasPrefix(s[1], "min") {
		return 0, errors.New("value not in minutes")
	}
	return strconv.Atoi(s[0])
}

func parseLoad(value string) (int, float64, error) {
	m := loadRegexp.FindStringSubmatch(value)
	if len(m) != 3 {
		return 0, 0, errors.New("value doesn't match load format")
	}
	w, err := strconv.Atoi(m[1])
	if err != nil {
		return 0, 0, fmt.Errorf("failed to parse watt value: %w", err)
	}
	p, err := strconv.ParseFloat(m[2], 64)
	if err != nil {
		return w, 0, fmt.Errorf("failed to parse percent value: %w", err)
	}
	return w, p / 100, nil
}

func parseEvent(value string) (e event, err error) {
	if !strings.Contains(value, "at") {
		e.State = value
		return e, nil
	}
	m := eventRegexp.FindStringSubmatch(value)
	if len(m) != 3 {
		return e, errors.New("value doesn't match event format")
	}
	e.State = m[1]
	e.Timestamp, err = time.ParseInLocation(eventTimestampFmt, m[2], time.Local)
	return e, err
}
