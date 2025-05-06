package jupiter

import (
	"github.com/shopspring/decimal"
)

const (
	SwapModeExactIn  = "ExactIn"
	SwapModeExactOut = "ExactOut"
)

// QuoteParams are the parameters for a quote request.
type QuoteParams struct {
	InputMint  string `url:"inputMint"`  // required
	OutputMint string `url:"outputMint"` // required
	Amount     uint64 `url:"amount"`     // required

	MaxAccounts         uint64 `url:"maxAccounts,omitempty"` // Max number of accounts to return
	SwapMode            string `url:"swapMode,omitempty"`    // Swap mode, default is ExactIn; Available values : ExactIn, ExactOut.
	SlippageBps         uint64 `url:"slippageBps,omitempty"`
	FeeBps              uint64 `url:"feeBps,omitempty"`              // Fee BPS (only pass in if you want to charge a fee on this swap)
	OnlyDirectRoutes    bool   `url:"onlyDirectRoutes,omitempty"`    // Only return direct routes (no hoppings and split trade)
	AsLegacyTransaction bool   `url:"asLegacyTransaction,omitempty"` // Only return routes that can be done in a single legacy transaction. (Routes might be limited)
	UserPublicKey       string `url:"userPublicKey,omitempty"`       // Public key of the user (only pass in if you want deposit and fee being returned, might slow down query)
	Dexes               string `url:"dexes,omitempty"`               // Dexes to query, separated by comma. If not provided, all dexes will be queried.
	ExcludeDexes        string `url:"excludeDexes,omitempty"`        // Dexes to exclude, separated by comma. If not provided, all dexes will be queried.
}

type QuoteResponse struct {
	InputMint            string          `json:"inputMint"`
	InAmount             decimal.Decimal `json:"inAmount"`
	OutputMint           string          `json:"outputMint"`
	OutAmount            decimal.Decimal `json:"outAmount"`
	OtherAmountThreshold decimal.Decimal `json:"otherAmountThreshold"`
	SwapMode             string          `json:"swapMode"`
	SlippageBps          int             `json:"slippageBps"`
	PlatformFee          decimal.Decimal `json:"platformFee"`
	PriceImpactPct       decimal.Decimal `json:"priceImpactPct"`
	RoutePlan            []Route         `json:"routePlan"`
	ContextSlot          int             `json:"contextSlot"`
	TimeTaken            decimal.Decimal `json:"timeTaken"`
}

type Route struct {
	SwapInfo SwapInfo        `json:"swapInfo"`
	Percent  decimal.Decimal `json:"percent"`
}

type SwapInfo struct {
	AmmKey     string          `json:"ammKey"`
	Label      string          `json:"label"`
	InputMint  string          `json:"inputMint"`
	OutputMint string          `json:"outputMint"`
	InAmount   decimal.Decimal `json:"inAmount"`
	OutAmount  decimal.Decimal `json:"outAmount"`
	FeeAmount  decimal.Decimal `json:"feeAmount"`
	FeeMint    string          `json:"feeMint"`
}

type QuoteForSwapResponse struct {
	InputMint            string      `json:"inputMint"`
	InAmount             string      `json:"inAmount"`
	OutputMint           string      `json:"outputMint"`
	OutAmount            string      `json:"outAmount"`
	OtherAmountThreshold string      `json:"otherAmountThreshold"`
	SwapMode             string      `json:"swapMode"`
	SlippageBps          int         `json:"slippageBps"`
	PlatformFee          interface{} `json:"platformFee"`
	PriceImpactPct       string      `json:"priceImpactPct"`
	RoutePlan            []struct {
		SwapInfo struct {
			AmmKey     string `json:"ammKey"`
			Label      string `json:"label"`
			InputMint  string `json:"inputMint"`
			OutputMint string `json:"outputMint"`
			InAmount   string `json:"inAmount"`
			OutAmount  string `json:"outAmount"`
			FeeAmount  string `json:"feeAmount"`
			FeeMint    string `json:"feeMint"`
		} `json:"swapInfo"`
		Percent int `json:"percent"`
	} `json:"routePlan"`
	ContextSlot int     `json:"contextSlot"`
	TimeTaken   float64 `json:"timeTaken"`
}

// SwapParams are the parameters for a swap request.
type SwapParams struct {
	UserPublicKey           string               `json:"userPublicKey"`
	DestinationTokenAccount string               `json:"destinationTokenAccount,omitempty"`
	WrapAndUnwrapSol        bool                 `json:"wrapAndUnwrapSol"`
	QuoteResponse           QuoteForSwapResponse `json:"quoteResponse"`
	DynamicSlippage         bool                 `json:"dynamicSlippage"`
}

// SwapResponse is the response from a swap request.
type SwapResponse struct {
	SwapTransaction string `json:"swapTransaction"` // base64 encoded transaction string
}

type SwapInstructionsResponse struct {
	TokenLedgerInstruction      *Instruction  `json:"tokenLedgerInstruction"`
	ComputeBudgetInstructions   []Instruction `json:"computeBudgetInstructions"`
	SetupInstructions           []Instruction `json:"setupInstructions"`
	SwapInstruction             *Instruction  `json:"swapInstruction"`
	CleanupInstruction          *Instruction  `json:"cleanupInstruction"`
	OtherInstructions           []Instruction `json:"otherInstructions"`
	AddressLookupTableAddresses []string      `json:"addressLookupTableAddresses"`
	PrioritizationFeeLamports   int           `json:"prioritizationFeeLamports"`
	ComputeUnitLimit            int           `json:"computeUnitLimit"`
	PrioritizationType          struct {
		ComputeBudget struct {
			MicroLamports          int `json:"microLamports"`
			EstimatedMicroLamports int `json:"estimatedMicroLamports"`
		} `json:"computeBudget"`
	} `json:"prioritizationType"`
	DynamicSlippageReport struct {
		SlippageBps                  int         `json:"slippageBps"`
		OtherAmount                  interface{} `json:"otherAmount"`
		SimulatedIncurredSlippageBps interface{} `json:"simulatedIncurredSlippageBps"`
		AmplificationRatio           interface{} `json:"amplificationRatio"`
		CategoryName                 string      `json:"categoryName"`
		HeuristicMaxSlippageBps      int         `json:"heuristicMaxSlippageBps"`
	} `json:"dynamicSlippageReport"`
	SimulationError struct {
		ErrorCode string `json:"errorCode"`
		Error     string `json:"error"`
	} `json:"simulationError"`
}

type Instruction struct {
	ProgramId string    `json:"programId"`
	Accounts  []Account `json:"accounts"`
	Data      string    `json:"data"`
}

type Account struct {
	Pubkey     string `json:"pubkey"`
	IsSigner   bool   `json:"isSigner"`
	IsWritable bool   `json:"isWritable"`
}
