package ECC

import (
	// "fmt"
	"fmt"
	. "math/big"
	"sync"
	. "github.com/zarismine/golang_Crypto/base"
)

type Point struct {
	X, Y *Int
	E    *ECC
}
type ECC struct {
	P, A, B *Int
}

func Curve(p, a, b *Int) *ECC {
	return &ECC{p, a, b}
}

func (e *ECC) SetPoint(x, y *Int) *Point {
	t1 := PowMod(x, NewInt(3), e.P)
	t1 = Mod(Add(Add(t1, Mul(e.A, x)), e.B), e.P)
	t2 := PowMod(y, NewInt(2), e.P)
	if t2.Cmp(t1) == 0 {
		return &Point{x, y, e}
	} else {
		fmt.Println(t1, t2)
		panic("点不在曲线上")
	}
}

func Point_Add(Q, P *Point) *Point {
	if P.E != Q.E {
		panic("不在一个曲线上")
	}
	if P.X.Cmp(NewInt(0)) == 0 && P.Y.Cmp(NewInt(1)) == 0 {
		return Q
	}
	if Q.X.Cmp(NewInt(0)) == 0 && Q.Y.Cmp(NewInt(1)) == 0 {
		return P
	}
	if P.X.Cmp(Q.X) == 0 && P.X.Cmp(Q.X) != 0 {
		return &Point{NewInt(0), NewInt(1), P.E}
	}
	p, a := P.E.P, P.E.A
	var m *Int
	if P.X.Cmp(Q.X) == 0 && P.Y.Cmp(Q.Y) == 0 {
		m = Add(Mul(NewInt(3), Mul(P.X, P.X)), a)
		m = Mul(m, ModInv(Mul(NewInt(2), P.Y), p))
		m = Mod(m, p)
	} else {
		m = Sub(Q.Y, P.Y)
		m = Mul(m, ModInv(Sub(Q.X, P.X), p))
		m = Mod(m, p)
	}
	x := Sub(Mul(m, m), Add(P.X, Q.X))
	y := Sub(Mul(m, Sub(P.X, x)), P.Y)
	x, y = Mod(x, p), Mod(y, p)
	return &Point{x, y, P.E}
}

func Point_Mul(G *Point, s *Int) *Point {
	R := &Point{NewInt(0), NewInt(1), G.E}
	N := new(Int).Set(s)
	for N.Sign() != 0 {
		if new(Int).And(N, NewInt(1)).Cmp(NewInt(0)) != 0 {
			R = Point_Add(R, G)
		}
		G = Point_Add(G, G)
		N.Rsh(N, 1)
	}
	return R
}
func Bsgs(Q, P *Point, bound int) int {
	step := Iroot(bound, 2)
	h := make(map[string]int)
	temp_P := &Point{NewInt(0), NewInt(1), P.E}
	s_P := Point_Mul(P, NewInt(int64(step)))
	hs_P := Point_Mul(s_P, NewInt(int64(4096)))
	var wg sync.WaitGroup
	var mu sync.Mutex
	wg.Add(256)
	for i := 0; i < 256; i++ {
		// fmt.Println(i)
		go func(temp *Point, i int) {
			for j := 0; j < 4096; j++ {
				mu.Lock()
				h[temp.Y.String()] = ((i*4096 + j) * 1048576)
				mu.Unlock()
				temp = Point_Add(temp, s_P)
			}
			wg.Done()
		}(temp_P, i)
		temp_P = Point_Add(temp_P,hs_P)
	}
	wg.Wait()
	ch := make(chan int)
	fmt.Println("table done!!", step)
	temp_Q := Q
	hs_P = Point_Mul(P, NewInt(int64(4096)))
	s_P = P
	for i := 0; i < 256; i++ {
		// fmt.Println(i)
		go func(temp *Point, i int) {
			for j := 0; j < 4096; j++ {
				r, ok := h[temp.Y.String()]
				if ok {
					ch <- (r - (4096*i + j))
				}
				temp = Point_Add(temp, s_P)
			}
		}(temp_Q, i)
		temp_Q = Point_Add(temp_Q, hs_P)
	}
	return <-ch
}
