package order

type TradingGenesis struct {
	ExpireFee              int64 `json:"expire_fee"`
	IOCExpireFee           int64 `json:"ioc_expire_fee"`
	FeeRateWithNativeToken int64 `json:"fee_rate_native"`
	FeeRate                int64 `json:"fee_rate"`
}

var DefaultTradingGenesis = TradingGenesis{
	ExpireFee:              100,
	IOCExpireFee:           50,
	FeeRateWithNativeToken: 500,
	FeeRate:                1000,
}
