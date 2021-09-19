package main

import (
	"fmt"
	"math"
	"os"

	"github.com/Konstantin8105/sm"
)

func calc(expr string) string {
	var err error
	for iter := 0; iter < 3; iter++ {
		expr, err = sm.Sexpr(nil, expr)
		if err != nil {
			panic(err)
		}
	}
	return expr
}

func simath() {
	sm.FloatFormat = 10
	// [ K11 K12 ] [ X ]    [ F1 ]
	// [ K21 K22 ] [ Y ]  = [ F2 ]
	F1 := "-x*(x-2)*(-y)*(y-2)*y+x*6-8*y"
	K11 := calc("d(" + F1 + ",x); variable(x); variable(y)")
	K12 := calc("d(" + F1 + ",y); variable(x); variable(y)")
	// calculate F2 = integral(K12,x)
	F2 := "-8*x - 4*(x*(x*x/3.0*y)) - 6*(x*x/2.0*(y*y)) + 3*(x*(x*x/3.0*(y*y))) + 8*(x*x/2.0*y)+1/(y+1)-1"
	K21 := calc("d(" + F2 + ",x); variable(x); variable(y)")
	K22 := calc("d(" + F2 + ",y); variable(x); variable(y)")

	sm.FloatFormat = 5
	F1 = calc("(" + F1 + ")*1.00000001")
	F2 = calc("(" + F2 + ")*1.00000001")
	K11 = calc("(" + K11 + ")*1.00000001")
	K12 = calc("(" + K12 + ")*1.00000001")
	K21 = calc("(" + K21 + ")*1.00000001")
	K22 = calc("(" + K22 + ")*1.00000001")

	fmt.Fprintf(os.Stdout, "F1  : %s\n", F1)
	fmt.Fprintf(os.Stdout, "F2  : %s\n", F2)
	fmt.Fprintf(os.Stdout, "K11 : %s\n", K11)
	fmt.Fprintf(os.Stdout, "K12 : %s\n", K12)
	fmt.Fprintf(os.Stdout, "K21 : %s\n", K21)
	fmt.Fprintf(os.Stdout, "K22 : %s\n", K22)
	fmt.Fprintf(os.Stdout, "K12 == K21 - %v\n", calc(K12+"-("+K21+")") == "0.00000")
}

func Kstiff(d [2]float64) (K [2][2]float64) {
	x := d[0]
	y := d[1]
	K[0][0] = 6.0000000000 - 2.0000000000*(y*(y*y)) - 4.0000000000*(x*(y*y)) + 4.0000000000*(y*y) + 2.0000000000*(x*(y*(y*y)))
	K[0][1] = -8.0000000000 - 6.0000000000*(x*(y*y)) - 4.0000000000*(x*(x*y)) + 8.0000000000*(x*y) + 3.0000000000*(x*(x*(y*y)))
	K[1][0] = -8.0000000000 - 6.0000000000*(x*(y*y)) - 3.9999999996*(x*(x*y)) + 8.0000000000*(x*y) + 2.9999999997*(x*(x*(y*y)))
	K[1][1] = -1.3333333332*(x*(x*x)) - 6.0000000000*(x*(x*y)) + (1.9999999998*(x*(x*(x*y))) + 4.0000000000*(x*x) + -1.0000000000/(1.0000000000+2.0000000000*y+y*y))
	return
}

func force(d [2]float64) (F [2]float64) {
	x := d[0]
	y := d[1]
	F[0] = -x*(x-2)*(-y)*(y-2)*y + x*6 - 8*y
	F[1] = -8*x - 4*(x*(x*x/3.0*y)) - 6*(x*x/2.0*(y*y)) + 3*(x*(x*x/3.0*(y*y))) + 8*(x*x/2.0*y) + 1/(y+1) - 1
	return
}

func main() {

	// simath()

	d := [2]float64{6.00000, 10.00000}

	//K := Kstiff(d)
	F := force(d)
	fmt.Fprintf(os.Stdout, "d   : %f\n", d)
	//fmt.Fprintf(os.Stdout, "K   : %f\n", K)
	fmt.Fprintf(os.Stdout, "F   : %f\n", F)
	// 	{
	// 		e1 := F[0] - (K[0][0]*d[0] + K[0][1]*d[1])
	// 		e2 := F[1] - (K[1][0]*d[0] + K[1][1]*d[1])
	// 		fmt.Println(	">:::", F[0] , (K[0][0]*d[0] + K[0][1]*d[1]) )
	// 		error := math.Sqrt(e1*e1 + e2*e2)
	// 		fmt.Fprintf(os.Stdout, "err : %f\n", error)
	// 	}

	// 	fmt.Println(	">>>>>>>>>>>>>>>>>>>>")
	// 	for i := 0; i < 10; i++ {
	// 		x := x * float64(i) / float64(10-1)
	// 		y := y * float64(i) / float64(10-1)
	// 		f := force(x,y)
	// 		fmt.Fprintf(os.Stdout,"%f , %f\n", f[0], f[1])
	// 	}
	// 	fmt.Println(	">>>>>>>>>>>>>>>>>>>>")

	steps(F, [2]float64{0, 0})
	// 	for iter := 0; iter < 25; iter++ {
	// 		x, y = steps(F1, F2, x, y)
	// 	fmt.Println(
	// 		x, ",",
	// 		y, ",",
	// 		F1-K[0][0]*x-K[0][1]*y, ",",
	// 		F2-K[1][0]*x-K[1][1]*y,
	// 	)
	// 	}
}

func steps(F, di [2]float64) (d [2]float64) {

	distance := 0.1

	d = di

	K := Kstiff(d)
	e1 := F[0] - (K[0][0]*d[0] + K[0][1]*d[1])
	e2 := F[1] - (K[1][0]*d[0] + K[1][1]*d[1])
	error := math.Sqrt(e1*e1 + e2*e2)
	//fmt.Println(d, ",", error)

	for iter := 0; iter < 200; iter++ {
		K = Kstiff(d)
		fmt.Println(":", d, K)

		vs := []struct {
			coord [2]float64
			err   float64
		}{
			{coord: [2]float64{d[0] + distance, d[1]}},
			{coord: [2]float64{d[0] - distance, d[1]}},
			{coord: [2]float64{d[0], d[1] + distance}},
			{coord: [2]float64{d[0], d[1] - distance}},
			{coord: [2]float64{d[0] + distance, d[1] + distance}},
			{coord: [2]float64{d[0] - distance, d[1] - distance}},
			{coord: [2]float64{d[0] - distance, d[1] + distance}},
			{coord: [2]float64{d[0] + distance, d[1] - distance}},
		}
		for i := range vs {
			dd := vs[i].coord
			edd1 := F[0] - (K[0][0]*dd[0] + K[0][1]*dd[1])
			edd2 := F[1] - (K[1][0]*dd[0] + K[1][1]*dd[1])
			vs[i].err = math.Sqrt(edd1*edd1 + edd2*edd2)
		}
		pos := -1
		e := error
		for i := range vs {
			if vs[i].err < e {
				pos = i
				e = vs[i].err
			}
		}
		// fmt.Println(">", pos, vs)
		if pos < 0 {
			fmt.Println("break")
			break
		}
		error = vs[pos].err
		d = vs[pos].coord
		fmt.Println(error)
	}

	// check
	// fmt.Println("Error: ", math.Abs(F1-K11*x-K12*y)<1e-9)
	// fmt.Println("Error: ", math.Abs(F2-K21*x-K22*y)<1e-9)

	return
}
