package main

var Config = UserConfig{
	// Profile settings
	Market:          "SLP-PERP",
	SizeGranularity: 0,
	PricePrecision:  4,
	Aggregate:       true,
	PriceTrim:       2,

	// TUI
	VolumeSymbol: "█", // recommended: '#' or '█'
	PriceMarker:  "<",
}
