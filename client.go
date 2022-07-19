package grpctmpl

import (
	"bytes"
	"context"
	"fmt"
	grpctypes "github.com/cosmos/cosmos-sdk/types/grpc"
	"github.com/cosmos/cosmos-sdk/types/query"
	val "github.com/cosmos/cosmos-sdk/x/staking/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"log"
	"math/big"
	"os"
)

const (
	proposalAt = 417_738
	fiveDays   = 5 * 24 * 60 * 6 // 6 second blocks
	step       = fiveDays / 24
)

func Run() {
	log.SetOutput(os.Stderr)

	client, err := grpc.Dial(GRPCHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()
	vClient := val.NewQueryClient(client)

	buf := bytes.NewBufferString(`"height","validator","percentage"` + "\n")
	for i := proposalAt; i < proposalAt+(fiveDays*4); i += step {
		key := make([]byte, 0)
		validators := make([]val.Validator, 0)
		for {
			bonded, e := vClient.Validators(
				metadata.AppendToOutgoingContext(context.Background(), grpctypes.GRPCBlockHeightHeader, fmt.Sprintf("%d", i)),
				&val.QueryValidatorsRequest{
					Status: "BOND_STATUS_BONDED",
					Pagination: &query.PageRequest{
						Limit: 100,
						Key:   key,
					},
				},
			)
			if e != nil {
				log.Fatal(e)
			}
			validators = append(validators, bonded.Validators...)
			if bonded.Pagination.NextKey == nil || len(bonded.Pagination.NextKey) == 0 {
				break
			}
			key = bonded.Pagination.NextKey
		}

		// figure out total bonded
		totalBonded := new(big.Float)
		for _, v := range validators {
			totalBonded = totalBonded.Add(totalBonded, new(big.Float).SetUint64(v.BondedTokens().Uint64()))
		}

		// collect each validators percentage
		for _, v := range validators {
			p, _ := new(big.Float).Quo(new(big.Float).SetUint64(v.BondedTokens().Uint64()), totalBonded).Float64()
			buf.WriteString(fmt.Sprintf(`%d,"%s",%f%s`, i, v.Description.Moniker,
				p*100, "\n",
			))
		}
	}
	fmt.Println(buf.String())
}
