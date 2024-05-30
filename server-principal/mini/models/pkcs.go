package models

import (
	"math/big"
)

const p = 31
const R = 4
const ru = 5
const ry = 7

type PprmS_ struct {
	N   int64
	G   *big.Int
	Y   *big.Int
	H   int64
	CtU [2]*big.Int
	CtY [2]*big.Int
}

type PrivprmS_ struct {
	U int64
	Y int64
	X int64
}

var PprmS PprmS_
var PrivprmS PrivprmS_

func Enc(pprmS PprmS_, gx int64, yx int64, hx int64) [2]*big.Int {
	var r [2]*big.Int

	r[0] = new(big.Int).Exp(pprmS.G, big.NewInt(gx), nil)             //	G**gx
	r11 := new(big.Int).Exp(pprmS.Y, big.NewInt(yx), nil)             //	Y**yx
	r12 := new(big.Int).Exp(big.NewInt(pprmS.H), big.NewInt(hx), nil) //	H**hx
	r[1] = new(big.Int).Mul(r11, r12)

	return r
}

func init() {
	PrivprmS.U = 7
	PrivprmS.Y = 9
	PrivprmS.X = 7

	PprmS.N = p
	PprmS.G = new(big.Int).Exp(big.NewInt(R), big.NewInt(2*PprmS.N), nil) // G = R ** 2N
	PprmS.Y = new(big.Int).Exp(PprmS.G, big.NewInt(PrivprmS.X), nil)      // Y = G ** X
	PprmS.H = (1 + PprmS.N) % (PprmS.N * PprmS.N)
	ctu := Enc(PprmS, ru, ru, PrivprmS.U)
	cty := Enc(PprmS, ry, ry, PrivprmS.Y)
	PprmS.CtU[0] = ctu[0]
	PprmS.CtU[1] = ctu[1]
	PprmS.CtY[0] = cty[0]
	PprmS.CtY[1] = cty[1]
}
