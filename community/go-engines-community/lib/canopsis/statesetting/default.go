package statesetting

func DefaultWorst() StateSetting {
	return StateSetting{
		Method:          "worst",
		JUnitThresholds: nil,
	}
}
