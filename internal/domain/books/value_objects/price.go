package valueobjects

import "errors"


type Price struct {
	minorAmount int
}


func NewPriceFromMajor(majorAmount int) (Price, error) {
	if majorAmount < 0 {
		return Price{}, errors.New("price amount cannot be negative")
	}
	return Price{minorAmount: majorAmount * 100}, nil
}


func NewPriceFromMinor(minorAmount int) (Price, error) {
	if minorAmount < 0 {
		return Price{}, errors.New("price amount cannot be negative")
	}
	return Price{minorAmount: minorAmount}, nil
}


func (p Price) MinorValue() int {
	return p.minorAmount
}


func (p Price) MajorValue() int {
	return p.minorAmount / 100
}

func (p Price) IsFree() bool {
	return p.minorAmount == 0
}