package main

import (
	"fmt"
	"log"

	"golang.org/x/image/math/f64"
	"gonum.org/v1/gonum/mat"
)

func findH(p1, p2 [4]f64.Vec2) mat.VecDense {
	a := mat.NewDense(9, 9, prepareP(p1, p2))

	var svd mat.SVD
	if ok := svd.Factorize(a, mat.SVDFull); !ok {
		log.Fatal("failed to factorize P")
	}

	const rcond = 1e-15
	rank := svd.Rank(rcond)
	if rank == 0 {
		log.Fatal("zero rank system")
	}

	b := mat.NewVecDense(9, []float64{0, 0, 0, 0, 0, 0, 0, 0, 1})

	var x mat.VecDense
	svd.SolveVecTo(&x, b, rank)

	fmt.Printf("singular values = %v\nrank = %d\nx = %.15f\n",
		format(svd.Values(nil), 4, rcond), rank, mat.Formatted(&x, mat.Prefix("    ")))

	return x
}

func main() {
	h1 := findH(
		[4]f64.Vec2{{0, 0}, {100, 0}, {0, 100}, {100, 100}},
		// [4]f64.Vec2{{0, 0}, {100, 0}, {0, 100}, {100, 100}},
		[4]f64.Vec2{{20, 20}, {80, 40}, {20, 80}, {80, 60}},
	)
	h := mat.NewDense(3, 3, h1.RawVector().Data)
	p1 := mat.NewVecDense(3, []float64{0, 0, 1})
	var p2 mat.VecDense
	p2.MulVec(h, p1)
	p2.ScaleVec(1/p2.AtVec(2), &p2)
	fmt.Println(mat.Formatted(&p2))
	// var m mat.Dense
	// m.Mul(mat.NewVe)
	// fmt.Println(
	// 	mat.Formatted(),
	// )
}

func prepareP(p1, p2 [4]f64.Vec2) []float64 {
	// From https://math.stackexchange.com/questions/494238/how-to-compute-homography-matrix-h-from-corresponding-points-2d-2d-planar-homog
	p := make([]float64, 0, 9*2*4+9)
	for i := 0; i < 4; i++ {
		xi1, yi1 := p1[i][0], p1[i][1]
		xi2, yi2 := p2[i][0], p2[i][1]
		p = append(p,
			-xi1, -yi1, -1, 0, 0, 0, xi1*xi2, yi1*xi2, xi2,
			0, 0, 0, -xi1, -yi1, -1, xi1*yi2, yi1*yi2, yi2,
		)
	}
	return append(p, 0, 0, 0, 0, 0, 0, 0, 0, 1)
}

func format(vals []float64, prec int, eps float64) []string {
	s := make([]string, len(vals))
	for i, v := range vals {
		if v < eps {
			s[i] = fmt.Sprintf("<%.*g", prec, eps)
			continue
		}
		s[i] = fmt.Sprintf("%.*g", prec, v)
	}
	return s
}
