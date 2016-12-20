package color

import (
	"math"
)

// Reference
// https://en.wikipedia.org/wiki/Color_difference#CIEDE2000

type Lab struct {
	L float64
	A float64
	B float64
}

var xyzConv = [][]float64{
	{0.4124, 0.3576, 0.1805},
	{0.2126, 0.7152, 0.0722},
	{0.0193, 0.1192, 0.9505},
}
var whiteXyz = []float64{
	0.9505, 1.0, 1.089,
}

func labFunc(t float64) float64 {
	if t > 0.008856 {
		return math.Pow(t, 1.0/3)
	} else {
		return (math.Pow(29.0/3, 3)*t + 16.0) / 116.0
	}
}

func FromRgb(r, g, b float64) Lab {
	x := r*xyzConv[0][0] + g*xyzConv[0][1] + b*xyzConv[0][2]
	y := r*xyzConv[1][0] + g*xyzConv[1][1] + b*xyzConv[1][2]
	z := r*xyzConv[2][0] + g*xyzConv[2][1] + b*xyzConv[2][2]
	xn := labFunc(x / whiteXyz[0])
	yn := labFunc(y / whiteXyz[1])
	zn := labFunc(z / whiteXyz[2])
	return Lab{
		L: 116*yn - 16,
		A: 500 * (xn - yn),
		B: 200 * (yn - zn),
	}
}

const (
	PI2 = math.Pi * 2
	RAD = PI2 / 360.0
)

func (self *Lab) Diff(other Lab) float64 {
	l1, a1, b1 := self.L, self.A, self.B
	l2, a2, b2 := other.L, other.A, other.B

	// 1. Calculate Ci', hi'
	cs1, cs2 := math.Hypot(a1, b1), math.Hypot(a2, b2)
	csb := (cs1 + cs2) / 2
	csb7 := math.Pow(csb, 7)
	g := 0.5 * (1.0 - math.Sqrt(csb7/(csb7+math.Pow(25.0, 7))))
	ad1, ad2 := (1+g)*a1, (1+g)*a2
	cd1, cd2 := math.Hypot(ad1, b1), math.Hypot(ad2, b2)
	hd1 := math.Mod(math.Atan2(b1, ad1)+PI2, PI2)
	hd2 := math.Mod(math.Atan2(b2, ad2)+PI2, PI2)

	// 2. Calculate DL', DC', DH'
	dLd := l2 - l1
	dCd := cd2 - cd1
	dhd := 0.0
	if cd1*cd2 != 0.0 {
		hdiff := hd2 - hd1
		dhd = hdiff
		if hdiff > math.Pi {
			dhd += PI2
		} else if hdiff < -math.Pi {
			dhd -= PI2
		}
	}
	dHd := 2.0 * math.Sqrt(cd1*cd2) * math.Sin(dhd/2)

	// 3. Calculate CIEDE2000
	lbd := (l1 + l2) / 2
	cbd := (cd1 + cd2) / 2
	hsum := hd1 + hd2
	hbd := hsum
	if cd1*cd2 != 0.0 {
		if math.Abs(hd2-hd1) <= math.Pi {
			hbd = hsum / 2
		} else if hsum < PI2 {
			hbd = (hsum + PI2) / 2
		} else {
			hbd = (hsum - PI2) / 2
		}
	}
	t := 1 - 0.17*math.Cos(hbd-30*RAD) + 0.24*math.Cos(2*hbd) + 0.32*math.Cos(3*hbd+6*RAD) - 0.20*math.Cos(4*hbd-63*RAD)
	dt := 30 * math.Exp(-math.Pow((hbd-275*RAD)/(25*RAD), 2))
	rc := 2 * math.Sqrt(math.Pow(cbd, 7)/(math.Pow(cbd, 7)+math.Pow(25, 7)))
	sl := 1.0 + (0.015*math.Pow(lbd-50, 2))/math.Sqrt(20+math.Pow(lbd-50, 2))
	sc := 1.0 + 0.045*cbd
	sh := 1.0 + 0.015*cbd*t
	rt := -math.Sin(2*dt) * rc
	dE2 := math.Pow(dLd/sl, 2) + math.Pow(dCd/sc, 2) + math.Pow(dHd/sh, 2) + rt*(dCd/sc)*(dHd/sh)
	return math.Sqrt(dE2)
}
