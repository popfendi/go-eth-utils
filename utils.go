package goethutils

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

const (
	Wei   = 1
	GWei  = 1e9
	Ether = 1e18
)

// formats wei as ether
func WeiToEther(wei *big.Int) *big.Float {
	return new(big.Float).Quo(new(big.Float).SetInt(wei), big.NewFloat(Ether))
}

// formats ether as wei
func EtherToWei(eth *big.Float) *big.Int {
	truncInt, _ := eth.Int(nil)
	truncInt = new(big.Int).Mul(truncInt, big.NewInt(Ether))
	fracStr := strings.Split(fmt.Sprintf("%.18f", eth), ".")[1]
	fracStr += strings.Repeat("0", 18-len(fracStr))
	fracInt, _ := new(big.Int).SetString(fracStr, 10)
	wei := new(big.Int).Add(truncInt, fracInt)
	return wei
}

// go implimentation of ethers parse units (returns big.Int as big.Float to whichever decimal precision you like)
func ParseUnits(num *big.Int, decimals uint8) *big.Float {
	if num.Cmp(big.NewInt(0)) == 0 {
		return big.NewFloat(0)
	}

	// Create a big.Float variable with the desired decimal places
	factor := new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil))

	// Convert the big.Int to a big.Float
	result := new(big.Float).SetInt(num)

	// Divide the big.Float by the factor to shift the decimal places
	result.Quo(result, factor)

	return result
}

// Removes i amount of decimals from num (returning num as *big.Int)
func RemoveZerosFromEnd(num *big.Int, i int) *big.Int {
	// Create a big.Int representing 10^decimals
	factor := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(i)), nil)

	// Multiply the original num by the factor
	result := new(big.Int).Div(num, factor)

	return result
}

// Adds i amount of decimals from num (returning num as *big.Int)
func AddZerosToEnd(num *big.Int, i int) *big.Int {
	// Create a big.Int representing 10^decimals
	factor := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(i)), nil)

	// Multiply the original num by the factor
	result := new(big.Int).Mul(num, factor)

	return result
}

func SortTokens(tokenA, tokenB common.Address) (common.Address, common.Address) {
	if new(big.Int).SetBytes(tokenA.Bytes()).Cmp(new(big.Int).SetBytes(tokenB.Bytes())) < 0 {
		return tokenA, tokenB
	}
	return tokenB, tokenA
}

// GeneratePairAddress generates a pair address for the given tokens (for UniswapV2 or clones, requires pairAddressSuffic and factory address of clone instance)
func GeneratePairAddress(token0, token1, FactoryAddress common.Address, pairAddressSuffix string) common.Address {
	// addresses need to be sorted in an ascending order for proper behaviour
	token0, token1 = SortTokens(token0, token1)

	// 255 is required as a prefix for this to work
	message := []byte{255}

	message = append(message, FactoryAddress.Bytes()...)

	addrSum := token0.Bytes()
	addrSum = append(addrSum, token1.Bytes()...)

	message = append(message, crypto.Keccak256(addrSum)...)

	b, _ := hex.DecodeString(pairAddressSuffix)
	message = append(message, b...)
	hashed := crypto.Keccak256(message)
	addressBytes := big.NewInt(0).SetBytes(hashed)
	addressBytes = addressBytes.Abs(addressBytes)
	return common.BytesToAddress(addressBytes.Bytes())
}

/*
	returns the storage key slot for a holder.

Position is the index slot (storage index of amount balances map).
You can use github.com/vocdoni/storage-proofs-eth-go/token useful library to get the memory slot of the map its self
*/
func GetBalancesMapSlot(holder common.Address, position int) [32]byte {
	return crypto.Keccak256Hash(
		common.LeftPadBytes(holder[:], 32),
		common.LeftPadBytes(big.NewInt(int64(position)).Bytes(), 32),
	)
}

// implimentation of solidity abi.Encode
func SolidityEncode(types []string, values []interface{}) (data []byte, err error) {
	var args abi.Arguments
	for _, t := range types {
		ty, err := abi.NewType(t, t, nil)
		if err != nil {
			return nil, err
		}

		args = append(args, abi.Argument{Type: ty})
	}

	data, err = args.Pack(
		values...,
	)
	return
}
