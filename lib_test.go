package ekwsum

import (
	"reflect"
	"testing"
)

func TestNewEkwNumber(t *testing.T) {
	type ekwResult struct {
		ekw EkwNumber
		err error
	}
	testCases := []struct {
		testName       string
		rawInput       string
		expectedResult ekwResult
	}{
		{
			testName: "valid complete ekw  - parse should pass",
			rawInput: "PR1J/00104856/8",
			expectedResult: ekwResult{
				ekw: EkwNumber{
					geoID: "PR1J",
					lamID: "00104856",
					sumControl: sumControl{
						value: "8",
					},
				},
				err: nil,
			},
		},
		{
			testName: "valid ekw without sum control - parse should pass",
			rawInput: "PR1L/00022370",
			expectedResult: ekwResult{
				ekw: EkwNumber{
					geoID:      "PR1L",
					lamID:      "00022370",
					sumControl: sumControl{},
				},
				err: nil,
			},
		},
		{
			testName: "first part and sum control - parse should pass",
			rawInput: "PR1L/8",
			expectedResult: ekwResult{
				ekw: EkwNumber{
					geoID:      "PR1L",
					lamID:      "8",
					sumControl: sumControl{},
				},
				err: nil,
			},
		},
		{
			testName: "empty ekw - no way to parse",
			rawInput: "",
			expectedResult: ekwResult{
				ekw: EkwNumber{},
				err: ErrIncomplete,
			},
		},
		{
			testName: "only first part - missing second part required by parser",
			rawInput: "PR1L",
			expectedResult: ekwResult{
				ekw: EkwNumber{},
				err: ErrIncomplete,
			},
		},
		{
			testName: "unsupported input format - should end up with err",
			rawInput: "PR1J/00104856/8/00022370",
			expectedResult: ekwResult{
				ekw: EkwNumber{},
				err: ErrUnsupported,
			},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.testName, func(t *testing.T) {
			ekw, err := NewEkwNumber(tC.rawInput)
			if !reflect.DeepEqual(tC.expectedResult.ekw, ekw) {
				t.Fatalf("test %s | failed with different ekw: %v != %v", tC.testName, tC.expectedResult.ekw, ekw)
			}
			if tC.expectedResult.err != err {
				t.Fatalf("test: %s | failed with different err state: %v != %v", tC.testName, tC.expectedResult.err, err)
			}
		})
	}
}

func TestValidate(t *testing.T) {
	testCases := []struct {
		testName                 string
		givenEkw                 EkwNumber
		expectedValidationResult error
	}{
		{
			testName: "valid complete ekw  - validation should pass",
			givenEkw: EkwNumber{
				geoID: "PR1J",
				lamID: "00104856",
				sumControl: sumControl{
					value: "8",
				},
			},
			expectedValidationResult: nil,
		},
		{
			testName: "valid ekw without sum control - validation should pass",
			givenEkw: EkwNumber{
				geoID:      "PR1J",
				lamID:      "00104856",
				sumControl: sumControl{},
			},
			expectedValidationResult: nil,
		},
		{
			testName: "first part and sum control (given as raw) - validation should fail",
			givenEkw: EkwNumber{
				geoID: "PR1J",
				lamID: "8",
			},
			expectedValidationResult: ErrValidationSecondPartUnknownFormat,
		},
		{
			testName: "invalid geo area - validation should fail",
			givenEkw: EkwNumber{
				geoID: "FR1434252X",
				lamID: "00104856",
			},
			expectedValidationResult: ErrValidationFirstPartUnknownFormat,
		},
		{
			testName: "invalid sum control - validation should fail",
			givenEkw: EkwNumber{
				geoID: "PR1J",
				lamID: "00104856",
				sumControl: sumControl{
					value: "1",
				},
			},
			expectedValidationResult: ErrValidationSumControl,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.testName, func(t *testing.T) {
			validationResult := tC.givenEkw.Validate()

			if tC.expectedValidationResult != validationResult {
				t.Fatalf("test: %s | failed with different err state: %v != %v", tC.testName, tC.expectedValidationResult, validationResult)
			}
		})
	}
}

func TestSumControl(t *testing.T) {
	testCases := []struct {
		testName    string
		givenEkw    EkwNumber
		expectedSum string
	}{
		{
			testName: "valid complete ekw ",
			givenEkw: EkwNumber{
				geoID: "PR1J",
				lamID: "00104856",
				sumControl: sumControl{
					value: "8",
				},
			},
			expectedSum: "8",
		},
		{
			testName: "ekw with invalid sum control ",
			givenEkw: EkwNumber{
				geoID: "PR1L",
				lamID: "00022370",
				sumControl: sumControl{
					value: "8",
				},
			},
			expectedSum: "0",
		},
		{
			testName: "should never happen, it could be caught by previous validation",
			givenEkw: EkwNumber{
				geoID: "PR1J",
				lamID: "8",
			},
			expectedSum: "",
		},
		{
			testName: "unexpected input - encode failed",
			givenEkw: EkwNumber{
				geoID: "FR1434252X?",
				lamID: "001!04856",
			},
			expectedSum: "",
		},
		{
			testName: "trusted sum control",
			givenEkw: EkwNumber{
				geoID: "PR1J",
				lamID: "00104856",
				sumControl: sumControl{
					value: "300",
					valid: true,
				},
			},
			expectedSum: "300",
		},
	}

	for _, tC := range testCases {
		t.Run(tC.testName, func(t *testing.T) {
			sumControl := tC.givenEkw.SumControl()

			if tC.expectedSum != sumControl {
				t.Fatalf("test: %s | failed with different sum control: %s != %s", tC.testName, tC.expectedSum, sumControl)
			}
		})
	}
}
