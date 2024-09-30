package jupiter

import (
	"context"
	"fmt"
	"git.verzth.work/go/jup-ag/program/perpetuals"
	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/shopspring/decimal"
)

type Jupiter struct {
	client *rpc.Client
}

func NewJupiter(rpcEndpoint string) *Jupiter {
	return &Jupiter{
		client: rpc.New(rpcEndpoint),
	}
}

func (m *Jupiter) GetApr() (decimal.Decimal, error) {
	pubKey := solana.MustPublicKeyFromBase58(JupiterJlpPool)
	// basic usage
	resp, err := m.client.GetAccountInfo(
		context.TODO(),
		pubKey,
	)
	if err != nil {
		return decimal.Zero, err
	}

	jlp := perpetuals.Pool{}
	err = jlp.UnmarshalWithDecoder(bin.NewBorshDecoder(resp.GetBinary()))
	if err != nil {
		return decimal.Zero, err
	}

	fmt.Println("feeAprBps: ", jlp.PoolApr.FeeAprBps)

	apr := decimal.NewFromUint64(jlp.PoolApr.FeeAprBps)
	apr = apr.Div(decimal.NewFromInt(100))

	fmt.Println("APR: ", apr)

	return apr, nil
}
