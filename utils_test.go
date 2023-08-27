package goethutils

import (
	"math/big"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

func TestWeiToEther(t *testing.T) {
	w, _ := new(big.Float).SetString("0.0000001")
	type args struct {
		wei *big.Int
	}
	tests := []struct {
		name string
		args args
		want *big.Float
	}{
		{name: "Pass Case", args: args{wei: big.NewInt(100000000000)}, want: w},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WeiToEther(tt.args.wei); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WeiToEther() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEtherToWei(t *testing.T) {
	a, _ := new(big.Float).SetString("0.0000001")
	type args struct {
		eth *big.Float
	}
	tests := []struct {
		name string
		args args
		want *big.Int
	}{
		{name: "Pass Case", args: args{eth: a}, want: big.NewInt(100000000000)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EtherToWei(tt.args.eth); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EtherToWei() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRemoveZerosFromEnd(t *testing.T) {
	type args struct {
		num *big.Int
		i   int
	}
	tests := []struct {
		name string
		args args
		want *big.Int
	}{
		{name: "Pass Case", args: args{num: big.NewInt(1000000), i: 3}, want: big.NewInt(1000)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RemoveZerosFromEnd(tt.args.num, tt.args.i); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RemoveZerosFromEnd() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAddZerosToEnd(t *testing.T) {
	type args struct {
		num *big.Int
		i   int
	}
	tests := []struct {
		name string
		args args
		want *big.Int
	}{
		{name: "Pass Case", args: args{num: big.NewInt(1000), i: 3}, want: big.NewInt(1000000)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AddZerosToEnd(tt.args.num, tt.args.i); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AddZerosToEnd() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseUnits(t *testing.T) {
	w, _ := new(big.Float).SetString("10.01")
	type args struct {
		num      *big.Int
		decimals uint8
	}
	tests := []struct {
		name string
		args args
		want *big.Float
	}{
		{name: "Pass Case", args: args{num: big.NewInt(1001), decimals: 2}, want: w},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseUnits(tt.args.num, tt.args.decimals); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseUnits() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSortTokens(t *testing.T) {
	weth := common.HexToAddress("0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2")
	usdt := common.HexToAddress("0xdac17f958d2ee523a2206206994597c13d831ec7")
	type args struct {
		tokenA common.Address
		tokenB common.Address
	}
	tests := []struct {
		name  string
		args  args
		want  common.Address
		want1 common.Address
	}{
		{name: "Pass  Case", args: args{tokenA: usdt, tokenB: weth}, want: weth, want1: usdt},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := SortTokens(tt.args.tokenA, tt.args.tokenB)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SortTokens() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("SortTokens() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestGeneratePairAddress(t *testing.T) {
	weth := common.HexToAddress("0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2")
	usdt := common.HexToAddress("0xdac17f958d2ee523a2206206994597c13d831ec7")
	factory := common.HexToAddress("0x5C69bEe701ef814a2B6a3EDD4B1652CB9cc5aA6f")
	suffix := "96e8ac4277198ff8b6f785478aa9a39f403cb768dd02cbee326c3e7da348845f"
	type args struct {
		token0            common.Address
		token1            common.Address
		FactoryAddress    common.Address
		pairAddressSuffix string
	}
	tests := []struct {
		name string
		args args
		want common.Address
	}{
		{name: "Pass Case", args: args{token0: weth, token1: usdt, FactoryAddress: factory, pairAddressSuffix: suffix}, want: common.HexToAddress("0x0d4a11d5eeaac28ec3f61d100daf4d40471f1852")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GeneratePairAddress(tt.args.token0, tt.args.token1, tt.args.FactoryAddress, tt.args.pairAddressSuffix); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GeneratePairAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSolidityEncode(t *testing.T) {
	w := common.Hex2Bytes("000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc200000000000000000000000000000000000000000000000000000002540be4000000000000000000000000000000000000000000000000000000000000000000")
	type args struct {
		types  []string
		values []interface{}
	}
	tests := []struct {
		name     string
		args     args
		wantData []byte
		wantErr  bool
	}{
		{name: "Pass Case", args: args{types: []string{"address", "uint256", "bytes32"}, values: []interface{}{common.HexToAddress("0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2"), big.NewInt(10000000000), [32]byte{}}}, wantData: w, wantErr: false},
		{name: "Fail Case", args: args{types: []string{"not", "a", "real", "type"}, values: []interface{}{common.HexToAddress("0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2"), common.HexToAddress("0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2"), common.HexToAddress("0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2"), common.HexToAddress("0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2")}}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := SolidityEncode(tt.args.types, tt.args.values)
			if (err != nil) != tt.wantErr {
				t.Errorf("SolidityEncode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}

}
