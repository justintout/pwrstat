package pwrstat

import (
	"testing"
	"time"
)

const docStatus = `Listing current UPS propertyies and status as following.

Properties:
        Model Name.................. UPS CP585
        Firmware Number............. BFH8102-6O1.5
        Rating Voltage.............. 120 V
        Rating Power................ 515 VA (335 Watt)

Current UPS status:
        State....................... Normal
        Power Supply by............. Utility Power
        Utility Voltage............. 111 V
        Output Voltage.............. 110 V
        Battery Capacity............ 100 %
        Remaining Runtime........... 60 min.
        Load........................ 0 Watt(0 %)
        Test Result................. Passed at 2011/01/27 13:17:15
        Last Power Event............ Blackout at 2011/01/27 13:21:15`

var docStatusExpected = StatusResult{
	Properties: properties{
		ModelName:      "UPS CP585",
		FirmwareNumber: "BFH8102-6O1.5",
		RatingVoltage:  120,
		RatingPower: ratingPower{
			Watts:   335,
			VoltAmp: 515,
		},
	},
	CurrentStatus: currentStatus{
		State:            "Normal",
		PowerSupplyBy:    "Utility Power",
		UtilityVoltage:   111,
		OutputVoltage:    110,
		BatteryCapacity:  1,
		RemainingRuntime: 60,
		Load: load{
			Watts:   0,
			Percent: 0,
		},
		TestResult: testResult{
			State:     "Passed",
			Timestamp: time.Date(2011, time.January, 27, 13, 17, 15, 0, time.Local),
		},
		LastPowerEvent: lastPowerEvent{
			State:     "Blackout",
			Timestamp: time.Date(2011, time.January, 27, 13, 21, 15, 0, time.Local),
		},
	},
}

const realStatus = `
The UPS information shows as following:

        Properties:
            Model Name................... CST1500SUC
            Firmware Number.............. CR02201A9713
            Rating Voltage............... 120 V
            Rating Power................. 900 Watt(1500 VA)

        Current UPS status:
            State........................ Normal
            Power Supply by.............. Utility Power
            Utility Voltage.............. 118 V
            Output Voltage............... 118 V
            Battery Capacity............. 100 %
            Remaining Runtime............ 134 min.
            Load......................... 27 Watt(3 %)
            Line Interaction............. None
            Test Result.................. Passed at 2024/07/16 01:01:08
            Last Power Event............. None
`

var realStatusExpected = StatusResult{
	Properties: properties{
		ModelName:      "CST1500SUC",
		FirmwareNumber: "CR02201A9713",
		RatingVoltage:  120,
		RatingPower: ratingPower{
			Watts:   900,
			VoltAmp: 1500,
		},
	},
	CurrentStatus: currentStatus{
		State:            "Normal",
		PowerSupplyBy:    "Utility Power",
		UtilityVoltage:   118,
		OutputVoltage:    118,
		BatteryCapacity:  1,
		RemainingRuntime: 134,
		Load: load{
			Watts:   27,
			Percent: 0.03,
		},
		LineInteraction: "None",
		TestResult: testResult{
			State:     "Passed",
			Timestamp: time.Date(2024, time.July, 16, 1, 1, 8, 0, time.Local),
		},
		LastPowerEvent: lastPowerEvent{
			State: "None",
		},
	},
}

func TestParseStatus(t *testing.T) {
	tests := []struct {
		in       string
		expected StatusResult
	}{
		{docStatus, docStatusExpected},
		{realStatus, realStatusExpected},
	}
	for _, tt := range tests {
		got, err := parseStatus(tt.in)
		if err != nil {
			t.Error(err)
		}
		switch {
		case got.Properties.ModelName != tt.expected.Properties.ModelName:
			t.Errorf("Properties.ModelName; got: %s, want: %s", got.Properties.ModelName, tt.expected.Properties.ModelName)
		case got.Properties.FirmwareNumber != tt.expected.Properties.FirmwareNumber:
			t.Errorf("Properties.FirmwareNumber; got: %s, want: %s", got.Properties.FirmwareNumber, tt.expected.Properties.FirmwareNumber)
		case got.Properties.RatingVoltage != tt.expected.Properties.RatingVoltage:
			t.Errorf("Properties.RatingVoltage; got: %d, want: %d", got.Properties.RatingVoltage, tt.expected.Properties.RatingVoltage)
		case got.Properties.RatingPower.VoltAmp != tt.expected.Properties.RatingPower.VoltAmp:
			t.Errorf("Properties.RatingPower.VoltAmp; got: %d, want: %d", got.Properties.RatingPower.VoltAmp, tt.expected.Properties.RatingPower.VoltAmp)
		case got.Properties.RatingPower.Watts != tt.expected.Properties.RatingPower.Watts:
			t.Errorf("Properties.RatingPower.Watts; got: %d, want: %d", got.Properties.RatingPower.Watts, tt.expected.Properties.RatingPower.VoltAmp)
		case got.CurrentStatus.State != tt.expected.CurrentStatus.State:
			t.Errorf("CurrentStatus.State; got: %s, want: %s", got.CurrentStatus.State, tt.expected.CurrentStatus.State)
		case got.CurrentStatus.PowerSupplyBy != tt.expected.CurrentStatus.PowerSupplyBy:
			t.Errorf("CurrentStatus.PowerSupplyBy; got: %s, want: %s", got.CurrentStatus.PowerSupplyBy, tt.expected.CurrentStatus.PowerSupplyBy)
		case got.CurrentStatus.UtilityVoltage != tt.expected.CurrentStatus.UtilityVoltage:
			t.Errorf("CurrentStatus.UtilityVoltage; got: %d, want: %d", got.CurrentStatus.UtilityVoltage, tt.expected.CurrentStatus.UtilityVoltage)
		case got.CurrentStatus.OutputVoltage != tt.expected.CurrentStatus.OutputVoltage:
			t.Errorf("CurrentStatus.OutputVoltage; got: %d, want: %d", got.CurrentStatus.OutputVoltage, tt.expected.CurrentStatus.OutputVoltage)
		case got.CurrentStatus.BatteryCapacity != tt.expected.CurrentStatus.BatteryCapacity:
			t.Errorf("CurrentStatus.BatteryCapacity; got: %f, want: %f", got.CurrentStatus.BatteryCapacity, tt.expected.CurrentStatus.BatteryCapacity)
		case got.CurrentStatus.RemainingRuntime != tt.expected.CurrentStatus.RemainingRuntime:
			t.Errorf("CurrentStatus.RemainingRuntime; got: %d, want: %d", got.CurrentStatus.RemainingRuntime, tt.expected.CurrentStatus.RemainingRuntime)
		case got.CurrentStatus.Load.Watts != tt.expected.CurrentStatus.Load.Watts:
			t.Errorf("CurrentStatus.Load.Watts; got: %d, want: %d", got.CurrentStatus.Load.Watts, tt.expected.CurrentStatus.Load.Watts)
		case got.CurrentStatus.Load.Percent != tt.expected.CurrentStatus.Load.Percent:
			t.Errorf("CurrentStatus.Load.Percent; got: %f, want: %f", got.CurrentStatus.Load.Percent, tt.expected.CurrentStatus.Load.Percent)
		case got.CurrentStatus.LineInteraction != tt.expected.CurrentStatus.LineInteraction:
			t.Errorf("CurrentStatus.LineInteraction; got: %s, want: %s", got.CurrentStatus.LineInteraction, tt.expected.CurrentStatus.LineInteraction)
		case got.CurrentStatus.TestResult.State != tt.expected.CurrentStatus.TestResult.State:
			t.Errorf("CurrentStatus.TestResult.State; got: %s, want: %s", got.CurrentStatus.TestResult.State, tt.expected.CurrentStatus.TestResult.State)
		case got.CurrentStatus.TestResult.Timestamp != tt.expected.CurrentStatus.TestResult.Timestamp:
			t.Errorf("CurrentStatus.TestResult.Timestamp; got: %v, want: %v", got.CurrentStatus.TestResult.Timestamp, tt.expected.CurrentStatus.TestResult.Timestamp)
		case got.CurrentStatus.LastPowerEvent.State != tt.expected.CurrentStatus.LastPowerEvent.State:
			t.Errorf("CurrentStatus.LastPowerEvent.State; got: %s, want: %s", got.CurrentStatus.LastPowerEvent.State, tt.expected.CurrentStatus.LastPowerEvent.State)
		case got.CurrentStatus.LastPowerEvent.Timestamp != tt.expected.CurrentStatus.LastPowerEvent.Timestamp:
			t.Errorf("CurrentStatus.LastPowerEvent.Timestamp; got: %v, want: %v", got.CurrentStatus.LastPowerEvent.Timestamp, tt.expected.CurrentStatus.LastPowerEvent.Timestamp)

		}
	}
}
