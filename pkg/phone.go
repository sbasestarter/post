package pkg

import (
	"github.com/ttacon/libphonenumber"
)

func ParsePhone(phone string) (number, region string, regionCode int, valid bool) {
	num, err := libphonenumber.Parse(phone, "")
	if err != nil {
		return
	}

	region = libphonenumber.GetRegionCodeForNumber(num)
	regionCode = libphonenumber.GetCountryCodeForRegion(region)

	valid = libphonenumber.IsValidNumberForRegion(num, region)

	if valid {
		number = libphonenumber.GetNationalSignificantNumber(num)
	}

	return
}
