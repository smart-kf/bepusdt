package rate

import (
	"regexp"

	"github.com/shopspring/decimal"
	"github.com/spf13/cast"

	"github.com/v03413/bepusdt/app/help"
	"github.com/v03413/bepusdt/app/log"
)

func parseFloatRate(syntax string, rawVal float64) float64 {
	if syntax == "" {
		return rawVal
	}

	if help.IsNumber(syntax) {
		return cast.ToFloat64(syntax)
	}

	match, err := regexp.MatchString(`^[~+-]\d+(\.\d+)?$`, syntax)
	if !match || err != nil {
		log.Error("浮动语法解析错误", err)

		return 0
	}

	act := syntax[0:1]
	raw := decimal.NewFromFloat(rawVal)
	base := decimal.NewFromFloat(cast.ToFloat64(syntax[1:]))

	switch act {
	case "~":
		return raw.Mul(base).InexactFloat64()
	case "+":
		return raw.Add(base).InexactFloat64()
	case "-":
		return raw.Sub(base).InexactFloat64()
	}

	return 0
}
