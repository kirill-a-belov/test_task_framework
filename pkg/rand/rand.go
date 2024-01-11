package rand

import (
	"crypto/rand"
	"math/big"

	"github.com/pkg/errors"
)

func Rand(size int, max int64) ([]int64, error) {
	if size < 0 {
		return nil, errors.New("invalid size")
	}
	if max < 0 {
		return nil, errors.New("invalid max")
	}

	result := make([]int64, size)
	for i := 0; i < size; i++ {
		rndNum, err := rand.Int(rand.Reader, big.NewInt(max))
		if err != nil {
			return nil, err
		}

		result[i] = rndNum.Int64()
	}

	return result, nil
}
