package ekwsum

import (
	"fmt"
	"regexp"
	"strings"
)

const (
	ekwPartsSeparator = "/"

	geoIDRawRegex = `^[a-zA-Z]{2}\d[a-zA-Z]$`
	lamIDRawRegex = `^\d{8}$`
)

var converterCharMapRules map[string]uint16 = map[string]uint16{
	"0": 0,
	"1": 1,
	"2": 2,
	"3": 3,
	"4": 4,
	"5": 5,
	"6": 6,
	"7": 7,
	"8": 8,
	"9": 9,
	"X": 10,
	"A": 11,
	"B": 12,
	"C": 13,
	"D": 14,
	"E": 15,
	"F": 16,
	"G": 17,
	"H": 18,
	"I": 19,
	"J": 20,
	"K": 21,
	"L": 22,
	"M": 23,
	"N": 24,
	"O": 25,
	"P": 26,
	"R": 27,
	"S": 28,
	"T": 29,
	"U": 30,
	"W": 31,
	"Y": 32,
	"Z": 33,
}

var calcWeights = []uint16{1, 3, 7, 1, 3, 7, 1, 3, 7, 1, 3, 7}

var (
	geoIDRegex = regexp.MustCompile(geoIDRawRegex)
	lamIDRegex = regexp.MustCompile(lamIDRawRegex)
)

var (
	ErrIncomplete  = fmt.Errorf("incomplete ekw number")
	ErrUnsupported = fmt.Errorf("unsupported ekw format")

	// validation err

	ErrValidationFirstPartUnknownFormat  = fmt.Errorf("validation: first part unknown format")
	ErrValidationSecondPartUnknownFormat = fmt.Errorf("validation: second part unknown format")
	ErrValidationSumControl              = fmt.Errorf("validation: sum control verification failed")
)

type EkwNumber struct {
	geoID      string
	lamID      string
	sumControl sumControl

	validationSettings ValidationSettings
	valid              bool
}

type ValidationSettings struct {
	validateControlSum bool
}

type ValidationOpts func(e *EkwNumber)

func WithValidationSum() ValidationOpts {
	return func(e *EkwNumber) {
		e.validationSettings.validateControlSum = true
	}
}

type sumControl struct {
	value string
	valid bool
}

// NewEkwNumber parses, simple validation included
// to parse correctly with required fields.
func NewEkwNumber(raw string, opts ...ValidationOpts) (EkwNumber, error) {
	parts := strings.Split(raw, ekwPartsSeparator)
	var ekw EkwNumber
	switch len(parts) {
	case 1:
		return EkwNumber{}, ErrIncomplete
	case 2:
		// assume that is ekw number without sum control
		ekw = EkwNumber{
			geoID: parts[0],
			lamID: parts[1],
		}
	case 3:
		// seems to be a complete ekw number
		ekw = EkwNumber{
			geoID: parts[0],
			lamID: parts[1],
			sumControl: sumControl{
				value: parts[2],
			},
		}
	default:
		return EkwNumber{}, ErrUnsupported
	}
	for _, o := range opts {
		o(&ekw)
	}
	return ekw, nil
}

// Validate checks whether first and second part are present and valid,
// also validates sum control if given.
//
// Note: Sum control validation fail is not a critical error so it can be
// skipped for calculation sum.
func (e *EkwNumber) Validate() error {
	if !geoIDRegex.MatchString(e.geoID) {
		return ErrValidationFirstPartUnknownFormat
	}
	if !lamIDRegex.MatchString(e.lamID) {
		return ErrValidationSecondPartUnknownFormat
	}
	e.valid = true

	// optional validation
	if e.validationSettings.validateControlSum {
		prevSumControl := e.sumControl.value
		if prevSumControl != e.SumControl() {
			return ErrValidationSumControl
		}
	}

	return nil
}

// SumControl estimates the sum control for given
// ekw id. It returns empty strings if calculation
// is not possible due to invalid ekw. To avoid this,
// we strongly recommend to run Validate before.
func (e *EkwNumber) SumControl() string {
	// in order to ensure that was validated
	if !e.valid {
		return ""
	}
	if e.sumControl.valid {
		return e.sumControl.value
	}
	base := strings.ToUpper(e.geoID + e.lamID)
	var encodedStr []uint16
	for _, c := range base {
		en, ok := converterCharMapRules[fmt.Sprintf("%c", c)]
		if !ok {
			// if provided char has not their equivalent, then
			// further encoding makes no sense so reject it now.
			return ""
		}
		encodedStr = append(encodedStr, en)
	}
	// unexpected length of weight pattern and ekw id
	// should never happen, but it does that is a bug
	if len(calcWeights) != len(encodedStr) {
		return ""
	}
	var overallSum uint16
	for idx := range calcWeights {
		encodedStr[idx] = encodedStr[idx] * calcWeights[idx]
		overallSum += encodedStr[idx]
	}
	e.sumControl.value = fmt.Sprintf("%d", overallSum%10)
	e.sumControl.valid = true
	return e.sumControl.value
}

// LamID returns the second part of ekw number,
// which is lam ID. It returns empty string if
// ekw is not valid or lam ID is not set.
// It is obligatory to run Validate before
// calling this method to ensure that ekw is valid.
func (e *EkwNumber) LamID() string {
	if !e.valid {
		return ""
	}
	return e.lamID
}
