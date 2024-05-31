package base

import (
	"crypto/rand"
	"math"
	. "math/big"
	// "math/rand"
)
func Mod(a,b *Int) *Int{
	return new(Int).Mod(a,b)
}
func Mul(a,b *Int) *Int{
	return new(Int).Mul(a,b)
}
func Div(a,b *Int) *Int{
	return new(Int).Div(a, b)
}
func Add(a,b *Int) *Int{
	return new(Int).Add(a, b)
}
func Sub(a,b *Int) *Int{
	return new(Int).Sub(a, b)
}
func PowMod(base, exp, mod *Int) *Int {
	if mod == NewInt(1) {
		return NewInt(0)
	}
	result := NewInt(1)
	base = Mod(base,mod)
	for exp.Cmp(NewInt(0)) == 1 {
		if Mod(exp,NewInt(2)).Cmp(NewInt(1)) == 0 {
			result = Mod(Mul(result,base),mod)
		}
		exp.Rsh(exp, 1)
		base = Mod(Mul(base,base),mod)
	}
	return result
}
func Pow(base *Int,exp int) *Int {
	e := NewInt(int64(exp))
	return new(Int).Exp(base, e, nil)
}

func randomBigInt(n *Int) *Int {
	// 生成一个 [0, n) 范围内的随机整数
	random, err := rand.Int(rand.Reader, n)
	if err != nil {
		panic(err)
	}

	return random
}
func IsPrime(n *Int) bool {
	if n.Cmp(NewInt(2))==0 || n.Cmp(NewInt(3))==0 {
		return true
	}
	k := 100
	if Mod(n,NewInt(2)) == NewInt(0) {
		return false
	}
	r, d := 0, Sub(n,NewInt(1))
	for Mod(d,NewInt(2)).Cmp(NewInt(0)) ==0 {
		r++
		d = Div(d,NewInt(2))
	}
	for i := 0; i < k; i++ {
		a := randomBigInt(Sub(n,NewInt(1)))
		x := PowMod(a, d, n)
		j := 0
		if x.Cmp(NewInt(1)) == 0 || x.Cmp(Sub(n,NewInt(1))) == 0 {
			continue
		}
		for ; j < r; j++ {
			x = PowMod(x, NewInt(2), n)
			if x.Cmp(Sub(n,NewInt(1))) == 0 {
				break
			}
		}
		if j == r {
			return false
		}
	}
	return true
}


func XGCD(a, b *Int) (*Int, *Int, *Int) {
    if a.Cmp(NewInt(0))==0 {
        return b, NewInt(0), NewInt(1)
    }
    gcd, x1, y1 := XGCD(Mod(b,a), a)
    x := Sub(y1,Mul(Div(b,a),x1))
    y := x1
    return gcd, x, y
}

func ModInv(a, p *Int) *Int {
	gcd, x, _ := XGCD(a, p)
    if gcd.Cmp(NewInt(1)) != 0 {
        panic("Modular inverse does not exist")
    }
    return Mod(Add(Mod(x,p),p),p)
}

func Iroot(n int, exp int) int {
    // 首先确定合理的初始值x0
    x0 := float64(n) / float64(exp)
	
    // 设置精度要求
    const epsilon = 1e-15

    // 迭代计算
    for {
        x1 := (float64(exp-1)*x0 + float64(n)/math.Pow(x0, float64(exp-1))) / float64(exp)
        if math.Abs(x1-x0) < epsilon {
            return int(x1)
        }
        x0 = x1
    }
}