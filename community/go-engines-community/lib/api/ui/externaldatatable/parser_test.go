package externaldatatable

import (
	"slices"
	"strings"
	"testing"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/externaldata"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
)

func TestParseNumber(t *testing.T) {
	testCases := []struct {
		name                 string
		input                string
		thousandsSeparator   string
		decimalSeparator     string
		expectedResult       float64
		expectedErrorMessage string
	}{
		{
			name:               "given simple integer with no separators expect valid parsing",
			input:              "1000",
			thousandsSeparator: "",
			decimalSeparator:   "",
			expectedResult:     1000,
		},
		{
			name:               "given negative integer with no separators expect valid parsing",
			input:              "-1000",
			thousandsSeparator: "",
			decimalSeparator:   "",
			expectedResult:     -1000,
		},
		{
			name:               "given decimal number with no separators expect valid parsing",
			input:              "12.34",
			thousandsSeparator: "",
			decimalSeparator:   "",
			expectedResult:     12.34,
		},
		{
			name:               "given negative decimal number with no separators expect valid parsing",
			input:              "-12.34",
			thousandsSeparator: "",
			decimalSeparator:   "",
			expectedResult:     -12.34,
		},
		{
			name:               "given zero with no separators expect valid parsing",
			input:              "0",
			thousandsSeparator: "",
			decimalSeparator:   "",
			expectedResult:     0,
		},
		{
			name:                 "given negative zero with no separators expect error",
			input:                "-0",
			thousandsSeparator:   "",
			decimalSeparator:     "",
			expectedErrorMessage: errStringIsNotAValidNumber.Error(),
		},
		{
			name:               "given number with many decimal places expect valid parsing",
			input:              "123.123456789",
			thousandsSeparator: "",
			decimalSeparator:   "",
			expectedResult:     123.123456789,
		},
		{
			name:               "given negative number with many decimal places expect valid parsing",
			input:              "-123.123456789",
			thousandsSeparator: "",
			decimalSeparator:   "",
			expectedResult:     -123.123456789,
		},
		{
			name:               "given small decimal with no separators expect valid parsing",
			input:              "0.001",
			thousandsSeparator: "",
			decimalSeparator:   "",
			expectedResult:     0.001,
		},
		{
			name:               "given negative small decimal with no separators expect valid parsing",
			input:              "-0.001",
			thousandsSeparator: "",
			decimalSeparator:   "",
			expectedResult:     -0.001,
		},
		{
			name:                 "given number with spaces and no separators expect error",
			input:                "1 000",
			thousandsSeparator:   "",
			decimalSeparator:     "",
			expectedErrorMessage: errStringIsNotAValidNumber.Error(),
		},
		{
			name:                 "given number with comma and no separators expect error",
			input:                "12,34",
			thousandsSeparator:   "",
			decimalSeparator:     "",
			expectedErrorMessage: errStringIsNotAValidNumber.Error(),
		},
		{
			name:                 "given empty string expect error",
			input:                "",
			thousandsSeparator:   "",
			decimalSeparator:     "",
			expectedErrorMessage: errStringIsNotAValidNumber.Error(),
		},
		{
			name:                 "given only negative sign expect error",
			input:                "-",
			thousandsSeparator:   "",
			decimalSeparator:     "",
			expectedErrorMessage: errStringIsNotAValidNumber.Error(),
		},
		{
			name:                 "given double negative expect error",
			input:                "--123",
			thousandsSeparator:   "",
			decimalSeparator:     "",
			expectedErrorMessage: errStringIsNotAValidNumber.Error(),
		},
		{
			name:                 "given negative sign in middle expect error",
			input:                "12-34",
			thousandsSeparator:   "",
			decimalSeparator:     "",
			expectedErrorMessage: errStringIsNotAValidNumber.Error(),
		},
		{
			name:                 "given negative sign at end expect error",
			input:                "123-",
			thousandsSeparator:   "",
			decimalSeparator:     "",
			expectedErrorMessage: errStringIsNotAValidNumber.Error(),
		},
		{
			name:                 "given non-numeric characters expect error",
			input:                "abc",
			thousandsSeparator:   "",
			decimalSeparator:     "",
			expectedErrorMessage: errStringIsNotAValidNumber.Error(),
		},
		{
			name:                 "given number with leading zeros expect error",
			input:                "00000",
			thousandsSeparator:   "",
			decimalSeparator:     "",
			expectedErrorMessage: errStringIsNotAValidNumber.Error(),
		},
		{
			name:                 "given negative number with leading zeros expect error",
			input:                "-00000",
			thousandsSeparator:   "",
			decimalSeparator:     "",
			expectedErrorMessage: errStringIsNotAValidNumber.Error(),
		},
		{
			name:                 "given number with leading zero expect error",
			input:                "01234",
			thousandsSeparator:   "",
			decimalSeparator:     "",
			expectedErrorMessage: errStringIsNotAValidNumber.Error(),
		},
		{
			name:                 "given negative number with leading zero expect error",
			input:                "-01234",
			thousandsSeparator:   "",
			decimalSeparator:     "",
			expectedErrorMessage: errStringIsNotAValidNumber.Error(),
		},
		{
			name:                 "given number with multiple leading zeros expect error",
			input:                "00123",
			thousandsSeparator:   "",
			decimalSeparator:     "",
			expectedErrorMessage: errStringIsNotAValidNumber.Error(),
		},
		{
			name:                 "given negative number with multiple leading zeros expect error",
			input:                "-00123",
			thousandsSeparator:   "",
			decimalSeparator:     "",
			expectedErrorMessage: errStringIsNotAValidNumber.Error(),
		},
		{
			name:                 "given zero followed by digits expect error",
			input:                "0123456789",
			thousandsSeparator:   "",
			decimalSeparator:     "",
			expectedErrorMessage: errStringIsNotAValidNumber.Error(),
		},
		{
			name:                 "given negative zero followed by digits expect error",
			input:                "-0123456789",
			thousandsSeparator:   "",
			decimalSeparator:     "",
			expectedErrorMessage: errStringIsNotAValidNumber.Error(),
		},
		{
			name:                 "given very long leading zeros expect error",
			input:                "000000000000000000001",
			thousandsSeparator:   "",
			decimalSeparator:     "",
			expectedErrorMessage: errStringIsNotAValidNumber.Error(),
		},
		{
			name:                 "given negative very long leading zeros expect error",
			input:                "-000000000000000000001",
			thousandsSeparator:   "",
			decimalSeparator:     "",
			expectedErrorMessage: errStringIsNotAValidNumber.Error(),
		},
		{
			name:               "given integer with space thousands separator expect valid parsing",
			input:              "1 000",
			thousandsSeparator: "space",
			decimalSeparator:   "",
			expectedResult:     1000,
		},
		{
			name:               "given negative integer with space thousands separator expect valid parsing",
			input:              "-1 000",
			thousandsSeparator: "space",
			decimalSeparator:   "",
			expectedResult:     -1000,
		},
		{
			name:               "given large number with space thousands separator expect valid parsing",
			input:              "1 000 000",
			thousandsSeparator: "space",
			decimalSeparator:   "",
			expectedResult:     1000000,
		},
		{
			name:               "given negative large number with space thousands separator expect valid parsing",
			input:              "-1 000 000",
			thousandsSeparator: "space",
			decimalSeparator:   "",
			expectedResult:     -1000000,
		},
		{
			name:               "given integer with comma thousands separator expect valid parsing",
			input:              "1,000",
			thousandsSeparator: "comma",
			decimalSeparator:   "",
			expectedResult:     1000,
		},
		{
			name:               "given negative integer with comma thousands separator expect valid parsing",
			input:              "-1,000",
			thousandsSeparator: "comma",
			decimalSeparator:   "",
			expectedResult:     -1000,
		},
		{
			name:               "given large number with comma thousands separator expect valid parsing",
			input:              "1,000,000",
			thousandsSeparator: "comma",
			decimalSeparator:   "",
			expectedResult:     1000000,
		},
		{
			name:               "given negative large number with comma thousands separator expect valid parsing",
			input:              "-1,000,000",
			thousandsSeparator: "comma",
			decimalSeparator:   "",
			expectedResult:     -1000000,
		},
		{
			name:               "given integer with dot thousands separator expect valid parsing",
			input:              "1.000",
			thousandsSeparator: "dot",
			decimalSeparator:   "",
			expectedResult:     1000,
		},
		{
			name:               "given negative integer with dot thousands separator expect valid parsing",
			input:              "-1.000",
			thousandsSeparator: "dot",
			decimalSeparator:   "",
			expectedResult:     -1000,
		},
		{
			name:               "given large number with dot thousands separator expect valid parsing",
			input:              "1.000.000",
			thousandsSeparator: "dot",
			decimalSeparator:   "",
			expectedResult:     1000000,
		},
		{
			name:               "given negative large number with dot thousands separator expect valid parsing",
			input:              "-1.000.000",
			thousandsSeparator: "dot",
			decimalSeparator:   "",
			expectedResult:     -1000000,
		},
		{
			name:                 "given invalid two digit grouping expect error",
			input:                "1 00",
			thousandsSeparator:   "space",
			decimalSeparator:     "",
			expectedErrorMessage: errStringIsNotAValidNumber.Error(),
		},
		{
			name:                 "given negative with invalid two digit grouping expect error",
			input:                "-1 00",
			thousandsSeparator:   "space",
			decimalSeparator:     "",
			expectedErrorMessage: errStringIsNotAValidNumber.Error(),
		},
		{
			name:                 "given double separator expect error",
			input:                "1,,000",
			thousandsSeparator:   "comma",
			decimalSeparator:     "",
			expectedErrorMessage: errStringIsNotAValidNumber.Error(),
		},
		{
			name:                 "given integer with leading zero and thousands separator expect error",
			input:                "01,000",
			thousandsSeparator:   "comma",
			decimalSeparator:     "",
			expectedErrorMessage: errStringIsNotAValidNumber.Error(),
		},
		{
			name:                 "given negative integer with leading zero and thousands separator expect error",
			input:                "-01,000",
			thousandsSeparator:   "comma",
			decimalSeparator:     "",
			expectedErrorMessage: errStringIsNotAValidNumber.Error(),
		},
		{
			name:               "given integer with dot decimal separator expect valid parsing",
			input:              "1000",
			thousandsSeparator: "",
			decimalSeparator:   "dot",
			expectedResult:     1000,
		},
		{
			name:               "given negative integer with dot decimal separator expect valid parsing",
			input:              "-1000",
			thousandsSeparator: "",
			decimalSeparator:   "dot",
			expectedResult:     -1000,
		},
		{
			name:               "given decimal number with dot separator expect valid parsing",
			input:              "1234.56",
			thousandsSeparator: "",
			decimalSeparator:   "dot",
			expectedResult:     1234.56,
		},
		{
			name:               "given negative decimal number with dot separator expect valid parsing",
			input:              "-1234.56",
			thousandsSeparator: "",
			decimalSeparator:   "dot",
			expectedResult:     -1234.56,
		},
		{
			name:               "given integer with comma decimal separator expect valid parsing",
			input:              "1000",
			thousandsSeparator: "",
			decimalSeparator:   "comma",
			expectedResult:     1000,
		},
		{
			name:               "given negative integer with comma decimal separator expect valid parsing",
			input:              "-1000",
			thousandsSeparator: "",
			decimalSeparator:   "comma",
			expectedResult:     -1000,
		},
		{
			name:               "given decimal number with comma separator expect valid parsing",
			input:              "1234,56",
			thousandsSeparator: "",
			decimalSeparator:   "comma",
			expectedResult:     1234.56,
		},
		{
			name:               "given negative decimal number with comma separator expect valid parsing",
			input:              "-1234,56",
			thousandsSeparator: "",
			decimalSeparator:   "comma",
			expectedResult:     -1234.56,
		},
		{
			name:                 "given comma in number with dot decimal separator expect error",
			input:                "1000,50",
			thousandsSeparator:   "",
			decimalSeparator:     "dot",
			expectedErrorMessage: errStringIsNotAValidNumber.Error(),
		},
		{
			name:                 "given dot in number with comma decimal separator expect error",
			input:                "1000.50",
			thousandsSeparator:   "",
			decimalSeparator:     "comma",
			expectedErrorMessage: errStringIsNotAValidNumber.Error(),
		},
		{
			name:                 "given leading decimal point expect error",
			input:                ".123",
			thousandsSeparator:   "",
			decimalSeparator:     "dot",
			expectedErrorMessage: errStringIsNotAValidNumber.Error(),
		},
		{
			name:                 "given negative leading decimal point expect error",
			input:                "-.123",
			thousandsSeparator:   "",
			decimalSeparator:     "dot",
			expectedErrorMessage: errStringIsNotAValidNumber.Error(),
		},
		{
			name:                 "given trailing decimal point expect error",
			input:                "123.",
			thousandsSeparator:   "",
			decimalSeparator:     "dot",
			expectedErrorMessage: errStringIsNotAValidNumber.Error(),
		},
		{
			name:                 "given negative trailing decimal point expect error",
			input:                "-123.",
			thousandsSeparator:   "",
			decimalSeparator:     "dot",
			expectedErrorMessage: errStringIsNotAValidNumber.Error(),
		},
		{
			name:                 "given decimal with leading zero in integer part expect error",
			input:                "01.50",
			thousandsSeparator:   "",
			decimalSeparator:     "dot",
			expectedErrorMessage: errStringIsNotAValidNumber.Error(),
		},
		{
			name:                 "given negative decimal with leading zero in integer part expect error",
			input:                "-01.50",
			thousandsSeparator:   "",
			decimalSeparator:     "dot",
			expectedErrorMessage: errStringIsNotAValidNumber.Error(),
		},
		{
			name:                 "given decimal with leading zero and comma separator expect error",
			input:                "01,50",
			thousandsSeparator:   "",
			decimalSeparator:     "comma",
			expectedErrorMessage: errStringIsNotAValidNumber.Error(),
		},
		{
			name:                 "given negative decimal with leading zero and comma separator expect error",
			input:                "-01,50",
			thousandsSeparator:   "",
			decimalSeparator:     "comma",
			expectedErrorMessage: errStringIsNotAValidNumber.Error(),
		},
		{
			name:               "given valid zero with decimal expect valid parsing",
			input:              "0.0",
			thousandsSeparator: "",
			decimalSeparator:   "dot",
			expectedResult:     0,
		},
		{
			name:               "given valid zero with decimal and comma separator expect valid parsing",
			input:              "0,0",
			thousandsSeparator: "",
			decimalSeparator:   "comma",
			expectedResult:     0,
		},
		{
			name:               "given valid zero with many decimal places expect valid parsing",
			input:              "0.000000",
			thousandsSeparator: "",
			decimalSeparator:   "dot",
			expectedResult:     0,
		},
		{
			name:               "given valid decimal starting with zero expect valid parsing",
			input:              "0.123",
			thousandsSeparator: "",
			decimalSeparator:   "dot",
			expectedResult:     0.123,
		},
		{
			name:               "given negative valid decimal starting with zero expect valid parsing",
			input:              "-0.123",
			thousandsSeparator: "",
			decimalSeparator:   "dot",
			expectedResult:     -0.123,
		},
		{
			name:               "given number with space thousands and dot decimal expect valid parsing",
			input:              "1 234.56",
			thousandsSeparator: "space",
			decimalSeparator:   "dot",
			expectedResult:     1234.56,
		},
		{
			name:               "given negative number with space thousands and dot decimal expect valid parsing",
			input:              "-1 234.56",
			thousandsSeparator: "space",
			decimalSeparator:   "dot",
			expectedResult:     -1234.56,
		},
		{
			name:               "given number with space thousands and comma decimal expect valid parsing",
			input:              "1 234 567,89",
			thousandsSeparator: "space",
			decimalSeparator:   "comma",
			expectedResult:     1234567.89,
		},
		{
			name:               "given negative number with space thousands and comma decimal expect valid parsing",
			input:              "-1 234 567,89",
			thousandsSeparator: "space",
			decimalSeparator:   "comma",
			expectedResult:     -1234567.89,
		},
		{
			name:               "given number with comma thousands and dot decimal expect valid parsing",
			input:              "1,234.56",
			thousandsSeparator: "comma",
			decimalSeparator:   "dot",
			expectedResult:     1234.56,
		},
		{
			name:               "given negative number with comma thousands and dot decimal expect valid parsing",
			input:              "-1,234.56",
			thousandsSeparator: "comma",
			decimalSeparator:   "dot",
			expectedResult:     -1234.56,
		},
		{
			name:               "given number with dot thousands and comma decimal expect valid parsing",
			input:              "1.234,56",
			thousandsSeparator: "dot",
			decimalSeparator:   "comma",
			expectedResult:     1234.56,
		},
		{
			name:               "given negative number with dot thousands and comma decimal expect valid parsing",
			input:              "-1.234,56",
			thousandsSeparator: "dot",
			decimalSeparator:   "comma",
			expectedResult:     -1234.56,
		},
		{
			name:               "given very large number expect valid parsing",
			input:              "999 999 999 999.99",
			thousandsSeparator: "space",
			decimalSeparator:   "dot",
			expectedResult:     999999999999.99,
		},
		{
			name:               "given negative very large number expect valid parsing",
			input:              "-999 999 999 999.99",
			thousandsSeparator: "space",
			decimalSeparator:   "dot",
			expectedResult:     -999999999999.99,
		},
		{
			name:               "given inconsistent grouping start expect valid parsing",
			input:              "12 345 678",
			thousandsSeparator: "space",
			decimalSeparator:   "dot",
			expectedResult:     12345678,
		},
		{
			name:               "given negative inconsistent grouping start expect valid parsing",
			input:              "-12 345 678",
			thousandsSeparator: "space",
			decimalSeparator:   "dot",
			expectedResult:     -12345678,
		},
		{
			name:               "given same separator for thousands and decimal expect valid parsing",
			input:              "0.0",
			thousandsSeparator: "dot",
			decimalSeparator:   "dot",
			expectedResult:     0,
		},
		{
			name:               "given negative same separator for thousands and decimal expect valid parsing",
			input:              "-0.0",
			thousandsSeparator: "dot",
			decimalSeparator:   "dot",
			expectedResult:     0,
		},
		{
			name:                 "given invalid grouping inconsistent end expect error",
			input:                "1 234 56",
			thousandsSeparator:   "space",
			decimalSeparator:     "dot",
			expectedErrorMessage: errStringIsNotAValidNumber.Error(),
		},
		{
			name:                 "given multiple decimal points expect error",
			input:                "1.44.3",
			thousandsSeparator:   "space",
			decimalSeparator:     "dot",
			expectedErrorMessage: errStringIsNotAValidNumber.Error(),
		},
		{
			name:                 "given leading separator expect error",
			input:                " 1000",
			thousandsSeparator:   "space",
			decimalSeparator:     "dot",
			expectedErrorMessage: errStringIsNotAValidNumber.Error(),
		},
		{
			name:                 "given trailing separator expect error",
			input:                "1000 ",
			thousandsSeparator:   "space",
			decimalSeparator:     "dot",
			expectedErrorMessage: errStringIsNotAValidNumber.Error(),
		},
		{
			name:                 "given negative with space before minus expect error",
			input:                " -1000",
			thousandsSeparator:   "space",
			decimalSeparator:     "dot",
			expectedErrorMessage: errStringIsNotAValidNumber.Error(),
		},
		{
			name:                 "given invalid thousands separator expect error",
			input:                "42",
			thousandsSeparator:   "\xcb",
			decimalSeparator:     "",
			expectedErrorMessage: "thousands separator must be one of \",\", \".\", \" \" or empty",
		},
		{
			name:                 "given invalid decimal separator expect error",
			input:                "42",
			thousandsSeparator:   "",
			decimalSeparator:     "\xcb",
			expectedErrorMessage: "decimal separator must be one of \".\", \",\" or empty",
		},
		{
			name:                 "given invalid both separators expect error",
			input:                "42",
			thousandsSeparator:   "+",
			decimalSeparator:     "[",
			expectedErrorMessage: "thousands separator must be one of \",\", \".\", \" \" or empty",
		},
	}

	p := NewParser()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rawResult, err := p.Parse(ColumnConfig{
				BaseColumnConfig: BaseColumnConfig{
					Type:               externaldata.ColumnTypeNumber,
					ThousandsSeparator: tc.thousandsSeparator,
					DecimalSeparator:   tc.decimalSeparator,
				},
			}, tc.input)
			if err != nil {
				if tc.expectedErrorMessage == "" {
					t.Errorf("error is not expectedResult: %v", err)
				} else if err.Error() != tc.expectedErrorMessage {
					t.Errorf("expectedResult error containing %q but got %q", tc.expectedErrorMessage, err.Error())
				}
			} else if tc.expectedErrorMessage != "" {
				t.Error("expectedResult error but got none")
			}

			result, ok := rawResult.(float64)
			if !ok {
				t.Error("result is not a float64")
			}

			if result != tc.expectedResult {
				t.Errorf("for input %q: expectedResult %v, got %v", tc.input, tc.expectedResult, result)
			}
		})
	}
}

func FuzzParseNumber(f *testing.F) {
	f.Add("1000", "", "")
	f.Add("1,000.50", "comma", "dot")
	f.Add("1 000,50", "space", "comma")
	f.Add("1.000.000", "dot", "comma")
	f.Add("abc", "space", "dot")
	f.Add("123.45.67", "space", "dot")
	f.Add("", "", "")
	f.Add("1  000", "space", "dot")
	f.Add("1..000", "dot", "comma")
	f.Add(".123", "", "dot")
	f.Add("123.", "", "dot")
	f.Add("1"+strings.Repeat("0", 1000), "", "")
	f.Add("-1"+strings.Repeat("0", 1000), "", "")
	f.Add("-1000", "", "")
	f.Add("-1,000.50", "comma", "dot")
	f.Add("--123", "", "")
	f.Add("12-34", "", "")
	f.Add("00000", "", "")
	f.Add("-0", "", "")
	f.Add("01234", "", "")
	f.Add("-01234", "", "")
	f.Add("01,000", "comma", "")
	f.Add("-01,000", "comma", "")
	f.Add("01.50", "", "dot")
	f.Add("-01.50", "", "dot")

	p := NewParser()

	f.Fuzz(func(t *testing.T, input string, thousandsSeparator string, decimalSeparator string) {
		rawResult, err := p.Parse(ColumnConfig{
			BaseColumnConfig: BaseColumnConfig{
				Type:               externaldata.ColumnTypeNumber,
				ThousandsSeparator: thousandsSeparator,
				DecimalSeparator:   decimalSeparator,
			},
		}, input)

		result, ok := rawResult.(float64)
		if !ok {
			t.Error("result is not a float64")
		}

		if err != nil && result != 0 {
			t.Errorf("parse returned non-zero result %v with error %v for input %q, thousands %q, decimal %q", result, err, input, thousandsSeparator, decimalSeparator)
		}
	})
}

func TestParseBool(t *testing.T) {
	testCases := []struct {
		name                 string
		input                string
		expectedResult       bool
		expectedErrorMessage string
	}{
		{
			name:           "given 'yes' expect true",
			input:          "yes",
			expectedResult: true,
		},
		{
			name:           "given 'y' expect true",
			input:          "y",
			expectedResult: true,
		},
		{
			name:           "given 'oui' expect true",
			input:          "oui",
			expectedResult: true,
		},
		{
			name:           "given 'true' expect true",
			input:          "true",
			expectedResult: true,
		},
		{
			name:           "given '1' expect true",
			input:          "1",
			expectedResult: true,
		},
		{
			name:           "given 'YES' expect true",
			input:          "YES",
			expectedResult: true,
		},
		{
			name:           "given 'Y' expect true",
			input:          "Y",
			expectedResult: true,
		},
		{
			name:           "given 'OUI' expect true",
			input:          "OUI",
			expectedResult: true,
		},
		{
			name:           "given 'TRUE' expect true",
			input:          "TRUE",
			expectedResult: true,
		},
		{
			name:           "given 'True' expect true",
			input:          "True",
			expectedResult: true,
		},
		{
			name:           "given 'Yes' expect true",
			input:          "Yes",
			expectedResult: true,
		},
		{
			name:           "given 'Oui' expect true",
			input:          "Oui",
			expectedResult: true,
		},
		{
			name:           "given 'yEs' expect true",
			input:          "yEs",
			expectedResult: true,
		},
		{
			name:           "given 'no' expect false",
			input:          "no",
			expectedResult: false,
		},
		{
			name:           "given 'n' expect false",
			input:          "n",
			expectedResult: false,
		},
		{
			name:           "given 'non' expect false",
			input:          "non",
			expectedResult: false,
		},
		{
			name:           "given 'false' expect false",
			input:          "false",
			expectedResult: false,
		},
		{
			name:           "given '0' expect false",
			input:          "0",
			expectedResult: false,
		},
		{
			name:           "given 'NO' expect false",
			input:          "NO",
			expectedResult: false,
		},
		{
			name:           "given 'N' expect false",
			input:          "N",
			expectedResult: false,
		},
		{
			name:           "given 'NON' expect false",
			input:          "NON",
			expectedResult: false,
		},
		{
			name:           "given 'FALSE' expect false",
			input:          "FALSE",
			expectedResult: false,
		},
		{
			name:           "given 'False' expect false",
			input:          "False",
			expectedResult: false,
		},
		{
			name:           "given 'No' expect false",
			input:          "No",
			expectedResult: false,
		},
		{
			name:           "given 'Non' expect false",
			input:          "Non",
			expectedResult: false,
		},
		{
			name:           "given 'fAlSe' expect false",
			input:          "fAlSe",
			expectedResult: false,
		},
		{
			name:                 "given empty string expect error",
			input:                "",
			expectedResult:       false,
			expectedErrorMessage: errStringIsNotAValidBool.Error(),
		},
		{
			name:                 "given space expect error",
			input:                "space",
			expectedResult:       false,
			expectedErrorMessage: errStringIsNotAValidBool.Error(),
		},
		{
			name:                 "given tab expect error",
			input:                "\t",
			expectedResult:       false,
			expectedErrorMessage: errStringIsNotAValidBool.Error(),
		},
		{
			name:                 "given newline expect error",
			input:                "\n",
			expectedResult:       false,
			expectedErrorMessage: errStringIsNotAValidBool.Error(),
		},
		{
			name:                 "given multiple spaces expect error",
			input:                "   ",
			expectedResult:       false,
			expectedErrorMessage: errStringIsNotAValidBool.Error(),
		},
		{
			name:                 "given 'yes' with leading space expect error",
			input:                " yes",
			expectedResult:       false,
			expectedErrorMessage: errStringIsNotAValidBool.Error(),
		},
		{
			name:                 "given 'yes' with trailing space expect error",
			input:                "yes ",
			expectedResult:       false,
			expectedErrorMessage: errStringIsNotAValidBool.Error(),
		},
		{
			name:                 "given 'true' with spaces expect error",
			input:                " true ",
			expectedResult:       false,
			expectedErrorMessage: errStringIsNotAValidBool.Error(),
		},
		{
			name:                 "given '2' expect error",
			input:                "2",
			expectedResult:       false,
			expectedErrorMessage: errStringIsNotAValidBool.Error(),
		},
		{
			name:                 "given '01' expect error",
			input:                "01",
			expectedResult:       false,
			expectedErrorMessage: errStringIsNotAValidBool.Error(),
		},
		{
			name:                 "given '1.0' expect error",
			input:                "1.0",
			expectedResult:       false,
			expectedErrorMessage: errStringIsNotAValidBool.Error(),
		},
		{
			name:                 "given 'ye' expect error",
			input:                "ye",
			expectedResult:       false,
			expectedErrorMessage: errStringIsNotAValidBool.Error(),
		},
		{
			name:                 "given 'tru' expect error",
			input:                "tru",
			expectedResult:       false,
			expectedErrorMessage: errStringIsNotAValidBool.Error(),
		},
		{
			name:                 "given 'yep' expect error",
			input:                "yep",
			expectedResult:       false,
			expectedErrorMessage: errStringIsNotAValidBool.Error(),
		},
		{
			name:                 "given 'nope' expect error",
			input:                "nope",
			expectedResult:       false,
			expectedErrorMessage: errStringIsNotAValidBool.Error(),
		},
		{
			name:                 "given '@' expect error",
			input:                "@",
			expectedResult:       false,
			expectedErrorMessage: errStringIsNotAValidBool.Error(),
		},
		{
			name:                 "given 'ñ' expect error",
			input:                "ñ",
			expectedResult:       false,
			expectedErrorMessage: errStringIsNotAValidBool.Error(),
		},
		{
			name:                 "given emoji expect error",
			input:                "👍",
			expectedResult:       false,
			expectedErrorMessage: errStringIsNotAValidBool.Error(),
		},
		{
			name:                 "given very long string expect error",
			input:                strings.Repeat("a", 1000),
			expectedResult:       false,
			expectedErrorMessage: errStringIsNotAValidBool.Error(),
		},
	}

	p := NewParser()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rawResult, err := p.Parse(ColumnConfig{
				BaseColumnConfig: BaseColumnConfig{
					Type: externaldata.ColumnTypeBoolean,
				},
			}, tc.input)
			if err != nil {
				if tc.expectedErrorMessage == "" {
					t.Errorf("error is not expectedResult: %v", err)
				} else if err.Error() != tc.expectedErrorMessage {
					t.Errorf("expectedResult error containing %q but got %q", tc.expectedErrorMessage, err.Error())
				}
			} else if tc.expectedErrorMessage != "" {
				t.Error("expectedResult error but got none")
			}

			result, ok := rawResult.(bool)
			if !ok {
				t.Error("result is not a bool")
			}

			if result != tc.expectedResult {
				t.Errorf("for input %q: expectedResult %v, got %v", tc.input, tc.expectedResult, result)
			}
		})
	}
}

func FuzzParseBool(f *testing.F) {
	f.Add("yes")
	f.Add("y")
	f.Add("oui")
	f.Add("true")
	f.Add("no")
	f.Add("n")
	f.Add("non")
	f.Add("false")
	f.Add("0")
	f.Add("YES")
	f.Add("TRUE")
	f.Add("FALSE")
	f.Add("")
	f.Add("abc")
	f.Add("2")
	f.Add("-1")
	f.Add("random")
	f.Add("yep")
	f.Add("nope")
	f.Add(" yes")
	f.Add("yes ")
	f.Add("👍")
	f.Add(strings.Repeat("a", 100))

	p := NewParser()

	f.Fuzz(func(t *testing.T, input string) {
		rawResult, err := p.Parse(ColumnConfig{
			BaseColumnConfig: BaseColumnConfig{
				Type: externaldata.ColumnTypeBoolean,
			},
		}, input)

		result, ok := rawResult.(bool)
		if !ok {
			t.Error("result is not a bool")
		}

		if err != nil && result != false {
			t.Errorf("parseBool returned non-false result %v with error %v for input %q", result, err, input)
		}
	})
}

func TestParseStringArray(t *testing.T) {
	testCases := []struct {
		name                 string
		input                string
		arrayType            int
		separator            string
		expectedResult       []string
		expectedErrorMessage string
	}{
		{
			name:           "given valid json array with strings expect successful parsing",
			input:          `["a", "b", "c"]`,
			arrayType:      jsonArray,
			separator:      "",
			expectedResult: []string{"a", "b", "c"},
		},
		{
			name:           "given valid json array with single element expect successful parsing",
			input:          `["single"]`,
			arrayType:      jsonArray,
			separator:      "",
			expectedResult: []string{"single"},
		},
		{
			name:           "given valid empty json array expect successful parsing",
			input:          `[]`,
			arrayType:      jsonArray,
			separator:      "",
			expectedResult: []string{},
		},
		{
			name:           "given valid json array with special characters expect successful parsing",
			input:          `["hello world", "test,with,commas", "test\"quotes\""]`,
			arrayType:      jsonArray,
			separator:      "",
			expectedResult: []string{"hello world", "test,with,commas", `test"quotes"`},
		},
		{
			name:           "given valid json array with unicode expect successful parsing",
			input:          `["café", "🎉"]`,
			arrayType:      jsonArray,
			separator:      "",
			expectedResult: []string{"café", "🎉"},
		},
		{
			name:                 "given invalid json syntax expect error",
			input:                `["a", "b"`,
			arrayType:            jsonArray,
			separator:            "",
			expectedResult:       nil,
			expectedErrorMessage: "string is not a valid JSON: cannot parse JSON: cannot parse array: unexpected end of array; unparsed tail: \"\"",
		},
		{
			name:                 "given json object instead of array expect error",
			input:                `{"a": "b"}`,
			arrayType:            jsonArray,
			separator:            "",
			expectedResult:       nil,
			expectedErrorMessage: "string is not a JSON array",
		},
		{
			name:                 "given json array with non-string elements expect error",
			input:                `["a", 123, "c"]`,
			arrayType:            jsonArray,
			separator:            "",
			expectedResult:       nil,
			expectedErrorMessage: "string array element at index 1 is not a string",
		},
		{
			name:                 "given empty string for json array expect error",
			input:                "",
			arrayType:            jsonArray,
			separator:            "",
			expectedResult:       nil,
			expectedErrorMessage: "empty string is not a valid JSON array",
		},
		{
			name:           "given comma separated string expect successful separator array parsing",
			input:          "a,b,c",
			arrayType:      separatorArray,
			separator:      ",",
			expectedResult: []string{"a", "b", "c"},
		},
		{
			name:           "given semicolon separated string expect successful separator array parsing",
			input:          "a;b;c",
			arrayType:      separatorArray,
			separator:      ";",
			expectedResult: []string{"a", "b", "c"},
		},
		{
			name:           "given string with multi-character separator expect successful separator array parsing",
			input:          "a::b::c",
			arrayType:      separatorArray,
			separator:      "::",
			expectedResult: []string{"a", "b", "c"},
		},
		{
			name:           "given single element string expect successful separator array parsing",
			input:          "single",
			arrayType:      separatorArray,
			separator:      ",",
			expectedResult: []string{"single"},
		},
		{
			name:           "given string with empty elements expect successful separator array parsing",
			input:          "a,,c",
			arrayType:      separatorArray,
			separator:      ",",
			expectedResult: []string{"a", "", "c"},
		},
		{
			name:           "given string with leading and trailing separators expect successful separator array parsing",
			input:          ",a,b,c,",
			arrayType:      separatorArray,
			separator:      ",",
			expectedResult: []string{"", "a", "b", "c", ""},
		},
		{
			name:           "given empty string for separator array expect successful parsing",
			input:          "",
			arrayType:      separatorArray,
			separator:      ",",
			expectedResult: []string{},
		},
		{
			name:                 "given separator array with empty separator expect error",
			input:                "test",
			arrayType:            separatorArray,
			separator:            "",
			expectedResult:       nil,
			expectedErrorMessage: "string array separator is required",
		},
		{
			name:                 "given invalid array type expect error",
			input:                "test",
			arrayType:            999,
			separator:            ",",
			expectedResult:       nil,
			expectedErrorMessage: "invalid array type: 999",
		},
	}

	p := NewParser()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			raw, err := p.Parse(ColumnConfig{
				BaseColumnConfig: BaseColumnConfig{
					Type:                 externaldata.ColumnTypeStringArray,
					StringArrayType:      tc.arrayType,
					StringArraySeparator: tc.separator,
				},
			}, tc.input)
			if err != nil {
				if tc.expectedErrorMessage == "" {
					t.Errorf("error is not expectedResult: %v", err)
				} else if err.Error() != tc.expectedErrorMessage {
					t.Errorf("expectedResult error containing %q but got %q", tc.expectedErrorMessage, err.Error())
				}
			} else if tc.expectedErrorMessage != "" {
				t.Error("expectedResult error but got none")
			} else {
				result, ok := raw.([]string)
				if !ok {
					t.Errorf("result is not a string array")
				}

				if !slices.Equal(result, tc.expectedResult) {
					t.Errorf("expectedResult %v, got %v", tc.expectedResult, result)
				}
			}
		})
	}
}

func FuzzParseStringArray(f *testing.F) {
	f.Add(`["a", "b", "c"]`, jsonArray, "")
	f.Add("a,b,c", separatorArray, ",")
	f.Add(`[]`, jsonArray, "")
	f.Add("", separatorArray, ",")
	f.Add(`["test"]`, jsonArray, "")
	f.Add("single", separatorArray, ",")
	f.Add("a;b;c", separatorArray, ";")
	f.Add(`["hello world", "test"]`, jsonArray, "")
	f.Add("a::b::c", separatorArray, "::")
	f.Add(`["café", "🎉"]`, jsonArray, "")

	p := NewParser()

	f.Fuzz(func(t *testing.T, input string, arrayType int, separator string) {
		result, err := p.Parse(ColumnConfig{
			BaseColumnConfig: BaseColumnConfig{
				Type:                 externaldata.ColumnTypeStringArray,
				StringArrayType:      arrayType,
				StringArraySeparator: separator,
			},
		}, input)

		if err != nil {
			resSlice, ok := utils.IsStringSlice(result)

			if result != nil && !ok && len(resSlice) != 0 {
				t.Errorf("parseStringArray returned non-nil or not empty string slice result %v with error %v for input %q", result, err, input)
			}
		} else {
			_, ok := utils.IsStringSlice(result)
			if !ok {
				t.Errorf("result is not a string array")
			}
		}
	})
}

func TestParseDatetime(t *testing.T) {
	testCases := []struct {
		name                 string
		input                string
		expectedResult       int64
		expectedErrorMessage string
	}{
		{
			name:           "given valid datetime with milliseconds and Z expect valid parsing",
			input:          "1990-12-31T00:00:00.00Z",
			expectedResult: time.Date(1990, 12, 31, 0, 0, 0, 0, time.UTC).Unix(),
		},
		{
			name:           "given valid datetime with timezone offset expect valid parsing",
			input:          "1990-12-31T00:00:00-00:00",
			expectedResult: time.Date(1990, 12, 31, 0, 0, 0, 0, time.UTC).Unix(),
		},
		{
			name:           "given valid datetime with Z expect valid parsing",
			input:          "1990-12-31T00:00:00Z",
			expectedResult: time.Date(1990, 12, 31, 0, 0, 0, 0, time.UTC).Unix(),
		},
		{
			name:           "given valid datetime with milliseconds and positive timezone expect valid parsing",
			input:          "1990-12-31T00:00:00.00+00:00",
			expectedResult: time.Date(1990, 12, 31, 0, 0, 0, 0, time.UTC).Unix(),
		},
		{
			name:           "given valid datetime with different timezone expect valid parsing",
			input:          "2023-05-15T14:30:45+05:00",
			expectedResult: time.Date(2023, 5, 15, 9, 30, 45, 0, time.UTC).Unix(),
		},
		{
			name:           "given valid datetime with negative timezone expect valid parsing",
			input:          "2023-05-15T14:30:45-05:00",
			expectedResult: time.Date(2023, 5, 15, 19, 30, 45, 0, time.UTC).Unix(),
		},
		{
			name:           "given valid datetime with milliseconds expect valid parsing",
			input:          "2023-01-01T12:00:00.99Z",
			expectedResult: time.Date(2023, 1, 1, 12, 0, 0, 990000000, time.UTC).Unix(),
		},
		{
			name:                 "given empty string expect error",
			input:                "",
			expectedErrorMessage: errStringIsNotAValidDatetime.Error(),
		},
		{
			name:                 "given invalid format without timezone expect error",
			input:                "1990-12-31T00:00:00",
			expectedErrorMessage: errStringIsNotAValidDatetime.Error(),
		},
		{
			name:                 "given invalid format with only date expect error",
			input:                "1990-12-31",
			expectedErrorMessage: errStringIsNotAValidDatetime.Error(),
		},
		{
			name:                 "given invalid format with space instead of T expect error",
			input:                "1990-12-31 00:00:00Z",
			expectedErrorMessage: errStringIsNotAValidDatetime.Error(),
		},
		{
			name:                 "given invalid month expect error",
			input:                "1990-13-31T00:00:00Z",
			expectedErrorMessage: errStringIsNotAValidDatetime.Error(),
		},
		{
			name:                 "given invalid day expect error",
			input:                "1990-12-32T00:00:00Z",
			expectedErrorMessage: errStringIsNotAValidDatetime.Error(),
		},
		{
			name:                 "given invalid hour expect error",
			input:                "1990-12-31T25:00:00Z",
			expectedErrorMessage: errStringIsNotAValidDatetime.Error(),
		},
		{
			name:                 "given invalid minute expect error",
			input:                "1990-12-31T00:60:00Z",
			expectedErrorMessage: errStringIsNotAValidDatetime.Error(),
		},
		{
			name:                 "given invalid second expect error",
			input:                "1990-12-31T00:00:60Z",
			expectedErrorMessage: errStringIsNotAValidDatetime.Error(),
		},
		{
			name:                 "given non-numeric characters expect error",
			input:                "abcd-ef-ghTij:kl:mnZ",
			expectedErrorMessage: errStringIsNotAValidDatetime.Error(),
		},
		{
			name:                 "given random string expect error",
			input:                "not a datetime",
			expectedErrorMessage: errStringIsNotAValidDatetime.Error(),
		},
		{
			name:                 "given partial datetime expect error",
			input:                "1990-12-31T00:",
			expectedErrorMessage: errStringIsNotAValidDatetime.Error(),
		},
		{
			name:                 "given datetime with wrong timezone format expect error",
			input:                "1990-12-31T00:00:00+0000",
			expectedErrorMessage: errStringIsNotAValidDatetime.Error(),
		},
		{
			name:                 "given datetime with invalid timezone expect error",
			input:                "1990-12-31T00:00:00+25:00",
			expectedErrorMessage: errStringIsNotAValidDatetime.Error(),
		},
		{
			name:           "given datetime with single digit milliseconds expect valid parsing",
			input:          "1990-12-31T00:00:00.0Z",
			expectedResult: time.Date(1990, 12, 31, 0, 0, 0, 0, time.UTC).Unix(),
		},
		{
			name:           "given datetime with three digit milliseconds expect valid parsing",
			input:          "1990-12-31T00:00:00.000Z",
			expectedResult: time.Date(1990, 12, 31, 0, 0, 0, 0, time.UTC).Unix(),
		},
		{
			name:                 "given datetime with extra characters expect error",
			input:                "1990-12-31T00:00:00Z extra",
			expectedErrorMessage: errStringIsNotAValidDatetime.Error(),
		},
		{
			name:                 "given datetime with leading whitespace expect error",
			input:                " 1990-12-31T00:00:00Z",
			expectedErrorMessage: errStringIsNotAValidDatetime.Error(),
		},
		{
			name:                 "given datetime with trailing whitespace expect error",
			input:                "1990-12-31T00:00:00Z ",
			expectedErrorMessage: errStringIsNotAValidDatetime.Error(),
		},
	}

	p := NewParser()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rawResult, err := p.Parse(ColumnConfig{
				BaseColumnConfig: BaseColumnConfig{
					Type: externaldata.ColumnTypeDateTime,
				},
			}, tc.input)
			if err != nil {
				if tc.expectedErrorMessage == "" {
					t.Errorf("error is not expectedResult: %v", err)
				} else if err.Error() != tc.expectedErrorMessage {
					t.Errorf("expectedResult error containing %q but got %q", tc.expectedErrorMessage, err.Error())
				}
			} else if tc.expectedErrorMessage != "" {
				t.Error("expectedResult error but got none")
			}

			result, ok := rawResult.(int64)
			if !ok {
				t.Errorf("result is not an int64")
			}

			if result != tc.expectedResult {
				t.Errorf("for input %q: expectedResult %v, got %v", tc.input, tc.expectedResult, result)
			}
		})
	}
}

func FuzzParseDatetime(f *testing.F) {
	f.Add("1990-12-31T00:00:00.00Z")
	f.Add("1990-12-31T00:00:00-00:00")
	f.Add("1990-12-31T00:00:00Z")
	f.Add("1990-12-31T00:00:00+05:00")
	f.Add("2023-05-15T14:30:45.99+00:00")
	f.Add("")
	f.Add("1990-12-31T00:00:00")
	f.Add("1990-12-31")
	f.Add("1990-12-31 00:00:00Z")
	f.Add("not a datetime")
	f.Add("1990-13-31T00:00:00Z")
	f.Add("1990-12-32T00:00:00Z")
	f.Add("1990-12-31T25:00:00Z")
	f.Add("1990-12-31T00:60:00Z")
	f.Add("1990-12-31T00:00:60Z")
	f.Add("abcd-ef-ghTij:kl:mnZ")
	f.Add("1990-12-31T00:00:00+0000")
	f.Add("1990-12-31T00:00:00+25:00")
	f.Add("1990-12-31T00:00:00.0Z")
	f.Add("1990-12-31T00:00:00.000Z")
	f.Add(" 1990-12-31T00:00:00Z")
	f.Add("1990-12-31T00:00:00Z ")
	f.Add(strings.Repeat("a", 100))

	p := NewParser()

	f.Fuzz(func(t *testing.T, input string) {
		rawResult, err := p.Parse(ColumnConfig{
			BaseColumnConfig: BaseColumnConfig{
				Type: externaldata.ColumnTypeDateTime,
			},
		}, input)

		result, ok := rawResult.(int64)
		if !ok {
			t.Errorf("result is not an int64")
		}

		if err != nil && result != 0 {
			t.Errorf("parseDatetime returned non-zero result %v with error %v for input %q", result, err, input)
		}
	})
}

func TestParseTimestamp(t *testing.T) {
	testCases := []struct {
		name                 string
		input                string
		expectedResult       int64
		expectedErrorMessage string
	}{
		{
			name:           "valid timestamp",
			input:          "1672531200",
			expectedResult: 1672531200,
		},
		{
			name:           "large timestamp",
			input:          "1672531200000",
			expectedResult: 1672531200000,
		},
		{
			name:           "zero",
			input:          "0",
			expectedResult: 0,
		},
		{
			name:           "one",
			input:          "1",
			expectedResult: 1,
		},
		{
			name:                 "leading zeros",
			input:                "0001672531200",
			expectedErrorMessage: errStringIsNotAValidTimestamp.Error(),
		},
		{
			name:                 "empty string",
			input:                "",
			expectedErrorMessage: errStringIsNotAValidTimestamp.Error(),
		},
		{
			name:                 "negative",
			input:                "-1",
			expectedErrorMessage: errStringIsNotAValidTimestamp.Error(),
		},
		{
			name:                 "non-numeric",
			input:                "abc",
			expectedErrorMessage: errStringIsNotAValidTimestamp.Error(),
		},
		{
			name:                 "mixed alphanumeric",
			input:                "123abc",
			expectedErrorMessage: errStringIsNotAValidTimestamp.Error(),
		},
		{
			name:                 "decimal",
			input:                "1672531200.5",
			expectedErrorMessage: errStringIsNotAValidTimestamp.Error(),
		},
		{
			name:                 "leading space",
			input:                " 1672531200",
			expectedErrorMessage: errStringIsNotAValidTimestamp.Error(),
		},
		{
			name:                 "trailing space",
			input:                "1672531200 ",
			expectedErrorMessage: errStringIsNotAValidTimestamp.Error(),
		},
		{
			name:                 "plus sign",
			input:                "+1672531200",
			expectedErrorMessage: errStringIsNotAValidTimestamp.Error(),
		},
		{
			name:                 "overflow",
			input:                "99999999999999999999",
			expectedErrorMessage: errStringIsNotAValidTimestamp.Error(),
		},
		{
			name:                 "minus sign only",
			input:                "-",
			expectedErrorMessage: errStringIsNotAValidTimestamp.Error(),
		},
		{
			name:                 "plus sign only",
			input:                "+",
			expectedErrorMessage: errStringIsNotAValidTimestamp.Error(),
		},
	}

	p := NewParser()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rawResult, err := p.Parse(ColumnConfig{
				BaseColumnConfig: BaseColumnConfig{
					Type: externaldata.ColumnTypeTimestamp,
				},
			}, tc.input)
			if err != nil {
				if tc.expectedErrorMessage == "" {
					t.Errorf("error is not expected: %v", err)
				} else if err.Error() != tc.expectedErrorMessage {
					t.Errorf("expected error containing %q but got %q", tc.expectedErrorMessage, err.Error())
				}
			} else if tc.expectedErrorMessage != "" {
				t.Error("expected error but got none")
			}

			result, ok := rawResult.(int64)
			if !ok {
				t.Errorf("result is not an int64")
			}

			if result != tc.expectedResult {
				t.Errorf("for input %q: expected %v, got %v", tc.input, tc.expectedResult, result)
			}
		})
	}
}

func FuzzParseTimestamp(f *testing.F) {
	f.Add("1672531200")
	f.Add("0")
	f.Add("")
	f.Add("-1")
	f.Add("abc")
	f.Add("123abc")
	f.Add("+1")
	f.Add("99999999999999999999")

	p := NewParser()

	f.Fuzz(func(t *testing.T, input string) {
		rawResult, err := p.Parse(ColumnConfig{
			BaseColumnConfig: BaseColumnConfig{
				Type: externaldata.ColumnTypeTimestamp,
			},
		}, input)

		result, ok := rawResult.(int64)
		if !ok {
			t.Errorf("result is not an int64")
		}

		if err != nil && result != 0 {
			t.Errorf("parseTimestamp returned non-zero result %v with error %v for input %q", result, err, input)
		}
	})
}

func TestParseRegexp(t *testing.T) {
	testCases := []struct {
		name           string
		input          string
		expectedRegexp string
		expectedScore  int
		expectedError  string
	}{
		{
			name:           "given empty string should be score 0",
			input:          "",
			expectedRegexp: "",
			expectedScore:  0,
		},
		{
			name:           "given dot star should be score 0",
			input:          ".*",
			expectedRegexp: ".*",
			expectedScore:  0,
		},
		{
			name:           "given dot plus should be score 0",
			input:          ".+",
			expectedRegexp: ".+",
			expectedScore:  0,
		},
		{
			name:           "given pattern with dot star should be score 1 with anchors added",
			input:          ".*test",
			expectedRegexp: "^.*test$",
			expectedScore:  1,
		},
		{
			name:           "given pattern with dot plus should be score 1 with anchors added",
			input:          "test.+",
			expectedRegexp: "^test.+$",
			expectedScore:  1,
		},
		{
			name:           "given pattern with both dot star and dot plus should be score 1 with anchors added",
			input:          ".*test.+",
			expectedRegexp: "^.*test.+$",
			expectedScore:  1,
		},
		{
			name:           "given fully anchored pattern should be score 3",
			input:          "^test$",
			expectedRegexp: "^test$",
			expectedScore:  3,
		},
		{
			name:           "given complex fully anchored pattern should be score 2",
			input:          "^[a-z]+\\d{2,4}$",
			expectedRegexp: "^[a-z]+\\d{2,4}$",
			expectedScore:  2,
		},
		{
			name:           "given simple pattern should be score 3 with anchors added",
			input:          "test",
			expectedRegexp: "^test$",
			expectedScore:  3,
		},
		{
			name:           "given pattern with character classes should be score 2 with anchors added",
			input:          "[a-z]+",
			expectedRegexp: "^[a-z]+$",
			expectedScore:  2,
		},
		{
			name:           "given pattern with quantifiers should be score 2 with anchors added",
			input:          "\\d{2,4}",
			expectedRegexp: "^\\d{2,4}$",
			expectedScore:  2,
		},
		{
			name:           "given pattern starting with anchor should be score 2",
			input:          "^test",
			expectedRegexp: "^test",
			expectedScore:  2,
		},
		{
			name:           "given pattern ending with anchor should be score 2",
			input:          "test$",
			expectedRegexp: "test$",
			expectedScore:  2,
		},
		{
			name:           "given pattern with wildcard starting with anchor should be score 1",
			input:          "^.*test",
			expectedRegexp: "^.*test",
			expectedScore:  1,
		},
		{
			name:           "given pattern with wildcard ending with anchor should be score 1",
			input:          "test.*$",
			expectedRegexp: "test.*$",
			expectedScore:  1,
		},
		{
			name:           "given pattern with anchor in middle should be score 3 with anchors added",
			input:          "te^st",
			expectedRegexp: "^te^st$",
			expectedScore:  3,
		},
		{
			name:           "given pattern with dollar in middle should be score 3 with anchors added",
			input:          "te$st",
			expectedRegexp: "^te$st$",
			expectedScore:  3,
		},
		{
			name:          "given invalid regex should return error",
			input:         "[unclosed",
			expectedError: "unable to parse regex: error parsing regexp: unterminated [] set in `[unclosed`",
		},
		{
			name:          "given invalid regex with unmatched parentheses should return error",
			input:         "test(",
			expectedError: "unable to parse regex: error parsing regexp: missing closing ) in `test(`",
		},
		{
			name:           "given pattern with nested wildcards in group should be score 1 with anchors added",
			input:          "(abc|.*def)",
			expectedRegexp: "^(abc|.*def)$",
			expectedScore:  1,
		},
		{
			name:           "given pattern with nested plus wildcard in group should be score 1 with anchors added",
			input:          "(start|end.+)",
			expectedRegexp: "^(start|end.+)$",
			expectedScore:  1,
		},
		{
			name:           "given begin text only",
			input:          "^",
			expectedRegexp: "^",
			expectedScore:  2,
		},
		{
			name:           "given end text only",
			input:          "$",
			expectedRegexp: "$",
			expectedScore:  2,
		},
		{
			name:           "given begin and end text only",
			input:          "^$",
			expectedRegexp: "^$",
			expectedScore:  3,
		},
		{
			name:           "given concat regexp anchors shouldn't be added",
			input:          "(?:(^ABC)|(XYZ$))",
			expectedRegexp: "(?:(^ABC)|(XYZ$))",
			expectedScore:  2,
		},
		{
			name:           "given literals concat regexp anchors should be added",
			input:          "abc|def",
			expectedRegexp: "^(abc|def)$",
			expectedScore:  3,
		},
		{
			name:           "given wildcards concat regexp anchors should be added",
			input:          ".*abc|.*def",
			expectedRegexp: "^(.*abc|.*def)$",
			expectedScore:  1,
		},
		{
			name:           "given regexp2 kind of regexp should be ok",
			input:          "(?!resource_CPU).*",
			expectedRegexp: "(?!resource_CPU).*",
			expectedScore:  2,
		},
	}

	p := NewParser()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rawResult, err := p.Parse(ColumnConfig{
				BaseColumnConfig: BaseColumnConfig{
					Type: externaldata.ColumnTypeRegexp,
				},
			}, tc.input)

			if tc.expectedError != "" {
				if err == nil {
					t.Error("expected error but got none")
				} else if err.Error() != tc.expectedError {
					t.Errorf("expected error %q but got %q", tc.expectedError, err.Error())
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			result, ok := rawResult.(parsedRegexp)
			if !ok {
				t.Error("result is not a parsedRegexp")
				return
			}

			if result.regexp != tc.expectedRegexp {
				t.Errorf("expected regexp %q, got %q", tc.expectedRegexp, result.regexp)
			}

			if result.score != tc.expectedScore {
				t.Errorf("expected score %d, got %d", tc.expectedScore, result.score)
			}
		})
	}
}
