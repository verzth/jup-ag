package tests

import (
	jupiter "git.verzth.work/go/jup-ag"
	"testing"
)

func TestHclient(t *testing.T) {
	cl := jupiter.NewJupag("https://lite-api.jup.ag/swap/v1", true)

	quote, err := cl.QuoteForSwap(jupiter.QuoteParams{
		InputMint:           "2b1kV6DkPAnxd5ixfnxCpjxmKwqjjaYmCZfHsFu24GXo",
		OutputMint:          "DezXAZ8z7PnrnRJjz3wXBoRgixCa6xjnB7YaB1pPB263",
		Amount:              1344716,
		SwapMode:            jupiter.SwapModeExactIn,
		SlippageBps:         50,
		FeeBps:              0,
		OnlyDirectRoutes:    false,
		AsLegacyTransaction: false,
		UserPublicKey:       "",
	})

	if err != nil {
		t.Fatal(err)
		return
	}

	t.Logf("quote: %+v", quote)

	swap, err := cl.SwapInstructions(jupiter.SwapParams{
		UserPublicKey:    "",
		WrapAndUnwrapSol: false,
		QuoteResponse:    *quote,
		DynamicSlippage:  false,
	})

	if err != nil {
		t.Fatal(err)
		return
	}

	t.Logf("swap: %+v", swap)
}
