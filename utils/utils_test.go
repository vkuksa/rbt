package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type UtilsSuite struct {
	suite.Suite
}

const intSize = 32 << (^uint(0) >> 63)

func (suite *UtilsSuite) TestMaxValue() {
	i := MaxValue[int]()
	assert.Equal(suite.T(), i, int(1<<(intSize-1)-1))

	i8 := MaxValue[int8]()
	assert.Equal(suite.T(), i8, int8(1<<7-1))

	i16 := MaxValue[int16]()
	assert.Equal(suite.T(), i16, int16(1<<15-1))

	i32 := MaxValue[int32]()
	assert.Equal(suite.T(), i32, int32(1<<31-1))

	i64 := MaxValue[int64]()
	assert.Equal(suite.T(), i64, int64(1<<63-1))

	f32 := MaxValue[float32]()
	assert.Equal(suite.T(), f32, float32(0x1p127*(1+(1-0x1p-23))))

	f64 := MaxValue[float64]()
	assert.Equal(suite.T(), f64, float64(0x1p1023*(1+(1-0x1p-52))))

	u := MaxValue[uint]()
	assert.Equal(suite.T(), u, uint(1<<intSize-1))

	u8 := MaxValue[uint8]()
	assert.Equal(suite.T(), u8, uint8(1<<8-1))

	u64 := MaxValue[uint64]()
	assert.Equal(suite.T(), u64, uint64(1<<64-1))

	str := MaxValue[string]()
	assert.Greater(suite.T(), str, "zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz")
}

func (suite *UtilsSuite) TestMinValue() {
	i := MinValue[int]()
	assert.Equal(suite.T(), i, int(-1<<(intSize-1)))

	i8 := MinValue[int8]()
	assert.Equal(suite.T(), i8, int8(-1<<7))

	i16 := MinValue[int16]()
	assert.Equal(suite.T(), i16, int16(-1<<15))

	i32 := MinValue[int32]()
	assert.Equal(suite.T(), i32, int32(-1<<31))

	i64 := MinValue[int64]()
	assert.Equal(suite.T(), i64, int64(-1<<63))

	f32 := MinValue[float32]()
	assert.Equal(suite.T(), f32, -float32(0x1p127*(1+(1-0x1p-23))))

	f64 := MinValue[float64]()
	assert.Equal(suite.T(), f64, -float64(0x1p1023*(1+(1-0x1p-52))))

	u := MinValue[uint]()
	assert.Equal(suite.T(), u, uint(0))

	u8 := MinValue[uint8]()
	assert.Equal(suite.T(), u8, uint8(0))

	u64 := MinValue[uint64]()
	assert.Equal(suite.T(), u64, uint64(0))

	str := MinValue[string]()
	assert.Equal(suite.T(), str, "")
}

func TestUtilsSuite(t *testing.T) {
	suite.Run(t, new(UtilsSuite))
}
