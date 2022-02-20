package main

// user configs here
var Config = UserConfig{
	// Program settings
	Market:          "LOOKS-PERP",
	SizeGranularity: 0,
	PricePrecision:  3,
	Aggregate:       true, // compile multi-transactions into singles

	// TUI
	PriceTrim:    2,     // how many digits to the right of price are hidden
	VolumeSymbol: "█",   // recommended: '#' or '█'
	PriceMarker:  " <<", // tailing last-price marker
}
