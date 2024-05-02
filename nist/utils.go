package nist

import "math"

var (
	// MAXLOG is the maximum log value to prevent underflow.
	MAXLOG float64 = 7.09782712893383996732e2
	// big is a large number used to stabilize the computation of the continued fraction.
	big float64 = 4.503599627370496e15
	// biginv is the inverse of 'big', used to rescale variables in the continued fraction.
	biginv float64 = 2.22044604925031308085e-16
	// MACHEP is the machine epsilon, which is the smallest number such that
	// 1.0 + MACHEP > 1.0. It is used as a convergence criterion.
	MACHEP float64 = 1.38777878078144567553e-17
)

// ref: https://nvlpubs.nist.gov/nistpubs/Legacy/SP/nistspecialpublication800-22r1a.pdf (section 5.5.3, p.99)

func igam(a, x float64) float64 {
	var res, ax, c, r float64
	if x <= 0 || a <= 0 {
		return 0.0
	}

	if x >= 1.0 && x > a {
		return 1.0 - igamc(a, x)
	}

	tempLgam, _ := math.Lgamma(a)
	ax = a*math.Log(x) - x - tempLgam
	if ax < -MAXLOG {
		panic("igam: Underflow")
	}
	ax = math.Exp(ax)

	// power series
	r = a
	c = 1.0
	res = 1.0

	for {
		r += 1.0
		c *= x / r
		res += c

		if c/res <= MACHEP {
			break
		}
	}

	return res * ax / a
}

func igamc(a, x float64) float64 {
	if x <= 0 || a <= 0 {
		return 1.0
	}
	if x < 1.0 || x < a {
		return 1.0 - igam(a, x)
	}

	tempLgam, _ := math.Lgamma(a)
	ax := a*math.Log(x) - x - tempLgam
	if ax < -MAXLOG {
		panic("igamc: Underflow")
	}
	ax = math.Exp(ax)

	// continued fraction
	y := 1.0 - a
	z := x + y + 1.0
	c := 0.0
	pkm2 := 1.0
	qkm2 := x
	pkm1 := x + 1.0
	qkm1 := z * x
	ans := pkm1 / qkm1

	var yc, pk, qk, r, t float64
	for {
		c += 1.0
		y += 1.0
		z += 2.0
		yc = y * c
		pk = pkm1*z - pkm2*yc
		qk = qkm1*z - qkm2*yc
		if qk != 0 {
			r = pk / qk
			t = math.Abs((ans - r) / r)
			ans = r
		} else {
			t = 1.0
		}
		pkm2 = pkm1
		pkm1 = pk
		qkm2 = qkm1
		qkm1 = qk
		if math.Abs(pk) > big {
			pkm2 *= biginv
			pkm1 *= biginv
			qkm2 *= biginv
			qkm1 *= biginv
		}
		if t <= MACHEP {
			break
		}
	}

	return ans * ax
}
