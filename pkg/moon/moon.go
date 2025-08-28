package moon

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

func CalcMoonNumber(yearGiven int) int {
	d := 6 // 2000
	mod := 19
	yearCurrent := 2000

	if yearCurrent < yearGiven {
		for yearCurrent+mod < yearGiven {
			yearCurrent += mod
		}
		for yearCurrent < yearGiven {
			yearCurrent++
			d++
			if d > 19 {
				d = 1
			}
		}
	} else {
		for yearCurrent-mod > yearGiven {
			yearCurrent -= mod
		}
		for yearCurrent > yearGiven {
			yearCurrent--
			d--
			if d < 1 {
				d = 19
			}
		}
	}
	return d
}

func jyear(td float64) (int, int, int) {
	td += 0.5
	z := math.Floor(td)
	f := td - z

	var a float64
	if z < 2299161.0 {
		a = z
	} else {
		alpha := math.Floor((z - 1867216.25) / 36524.25)
		a = z + 1 + alpha - math.Floor(alpha/4)
	}

	b := a + 1524
	c := math.Floor((b - 122.1) / 365.25)
	d := math.Floor(365.25 * c)
	e := math.Floor((b - d) / 30.6001)

	mm := int(math.Floor(e))
	if mm >= 14 {
		mm -= 13
	} else {
		mm -= 1
	}

	year := int(math.Floor(c))
	if mm > 2 {
		year -= 4716
	} else {
		year -= 4715
	}

	day := int(math.Floor(b - d - math.Floor(30.6001*e) + f))

	return year, mm, day
}

type MoonTableElement struct {
	TNew  time.Time
	TFull time.Time
	T1    float64
	T2    float64
}

/*
JHMS  --  Convert Julian time to hour, minutes, and seconds,

	returned as three separate values.
*/
func jhms(j float64) (int, int, int) {
	j += 0.5 // Astronomical to civil
	ij := (j - math.Floor(j)) * 86400.0
	hours := math.Floor(ij / 3600)
	minutes := math.Floor((ij / 60))
	seconds := math.Floor(ij)
	return int(hours), int(math.Mod(minutes, 60)), int(math.Mod(seconds, 60))
}

/*  DTR  --  Degrees to radians.  */
func dtr(d float64) float64 {
	return (d * math.Pi) / 180.0
}

/*  FIXANGLE  --  Range reduce angle in degrees.  */
func fixangle(a float64) float64 {
	return a - 360.0*(math.Floor(a/360.0))
}

/*  SUMSER  --  Sum the series of periodic terms.  */
func sumser(trig func(float64) float64, D, M, F, T float64, argtab []float64, coeff []float64, tfix []int, tfixc []float64) float64 {
	D = dtr(fixangle(D))
	M = dtr(fixangle(M))
	F = dtr(fixangle(F))

	j := 0
	n := 0
	sum := 0.0

	for i := 0; i < len(coeff) && coeff[i] != 0.0; i++ {
		arg := (D * argtab[j]) + (M * argtab[j+1]) + (F * argtab[j+2])
		j += 3
		coef := coeff[i]

		if n < len(tfix) && i == tfix[n] {
			coef += T * tfixc[n]
			n++
		}
		sum += coef * trig(arg)
	}

	return sum
}

func Gen(year int, month int, day int) (string, string, []*MoonTableElement, time.Duration) {
	var moonDays time.Duration
	moonTable := []*MoonTableElement{}
	tGiven := time.Date(year, getMonth(month), day, 0, 0, 0, 0, time.UTC)

	var /*v,*/ s string
	var /*sk,*/ kr []float64
	var l int
	var perigee bool
	var dat []float64
	var evt [][]float64
	var m int
	var epad, pchar, phnear string
	pmin := math.MaxFloat64
	var pminx int
	pmax := -math.MaxFloat64
	var pmaxx int
	var yrange, centile float64
	const TOLERANCE = 0.01
	var k1, mtime float64
	var minx int
	var phaset []float64
	const Itemlen = 36
	const Pitemlen = 25

	s = ""

	skVal := math.Floor((float64(year) - 1999.97) * 13.2555)
	dat = make([]float64, 0)
	evt = make([][]float64, 0)
	phaset = make([]float64, 0)

	// Tabulate perigees and apogees for the year
	for l = 0; ; l++ {
		kr = moonpa(skVal)
		datey, _, _ := jyear(kr[0])
		perigee = (skVal - math.Floor(skVal)) < 0.25

		if datey == year {
			if kr[2] < pmin {
				pmin = kr[2]
				pminx = m
			} else if kr[2] > pmax {
				pmax = kr[2]
				pmaxx = m
			}
			dat = append(dat, skVal)
			evt = append(evt, kr)
			m++
		}
		if datey > year {
			break
		}
		skVal += 0.5
	}
	yrange = pmax - pmin

	// Tabulate new and full moons surrounding the year
	k1 = math.Floor((float64(year) - 1900) * 12.3685) // - 4
	minx = 0
	isNext := true
	for l = 0; ; l++ {
		mtime = truephase(k1, float64(l&1)*0.5)
		datey, _, _ := jyear(mtime)
		if datey >= year {
			minx++
		}
		phaseSign := mtime
		if (l & 1) == 0 {
			phaseSign = -mtime
		}
		phaset = append(phaset, phaseSign)
		if !isNext {
			break
		}
		if datey > year {
			minx++
			isNext = false
		}
		if (l & 1) != 0 {
			k1 += 1
		}
	}

	// Generate perigee and apogee table
	perigeeApogeeTable := ""
	for l = 0; l < m; l++ {
		skVal = dat[l]
		kr = evt[l]

		perigee = (skVal - math.Floor(skVal)) < 0.25
		if !perigee && s == "" {
			s = pad("", Itemlen, " ")
		}
		phnear = nearphase(kr[0], phaset)
		if strings.HasPrefix(phnear, "F") {
			pchar = "+"
		} else {
			pchar = "-"
		}
		if l == pminx || l == pmaxx {
			epad = pchar + pchar
		} else {
			centile = (kr[2] - pmin) / yrange
			if centile <= TOLERANCE || centile >= (1-TOLERANCE) {
				epad = pchar + " "
			} else {
				epad = "  "
			}
		}
		s += edate(kr[0]) + " " + fmt.Sprintf("%f", math.Round(kr[2])) + " km " + epad + " " + phnear
		if len(s) < Itemlen {
			s = pad(s, Itemlen, " ")
		} else {
			perigeeApogeeTable += s + "\n"
			s = ""
		}
	}
	if s != "" {
		perigeeApogeeTable += s + "\n"
	}

	s = ""
	var lastnew float64
	//var lastfull float64
	phaseTable := ""
	for l = 0; l < minx; l++ {

		elem := &MoonTableElement{}

		mp := phaset[l]
		var data string
		if mp < 0 {
			mp = -mp
			if lastnew != 0 {
				dataVal := cuzcoNight(mp) - cuzcoNight(lastnew)
				if dataVal == 30 {
					data = "*" + fmt.Sprintf("%f", dataVal)
				} else if dataVal != 29 {
					data = "@" + fmt.Sprintf("%f", dataVal)
				}
			} else {
				data = " "
			}
			elem.T1 = mp
			elem.T2 = lastnew
			//
			s += pad(data, 3, " ")
			lastnew = mp
		} else {
			if s == "" {
				s = pad(s, Pitemlen, " ")
			}
		}

		elem.T1 = mp
		elem.T2 = lastnew

		datey, _, _ := jyear(mp)
		elem.TNew = cuzcoDateTime(lastnew)
		elem.TFull = cuzcoDateTime(mp)

		s += "   " + strconv.Itoa(datey) + " " + cuzcoDate(mp)

		if elem.T1 != elem.T2 {
			moonTable = append(moonTable, elem)
			if tGiven.After(elem.TNew) && tGiven.Before(elem.TFull) {
				moonDays = tGiven.Sub(elem.TNew)
			}
		}

		if len(s) < Pitemlen {
			s = pad(s, Pitemlen, " ")
		} else {
			phaseTable += s + "\n"
			s = ""
		}
	}
	if s != "" {
		phaseTable += s + "\n"
	}
	return perigeeApogeeTable, phaseTable, moonTable, moonDays
}

func pad(str string, length int, padChar string) string {
	for len(str) < length {
		str = padChar + str
	}
	return str
}

func moonpa(sk float64) []float64 {
	var t, t2, t3, t4, JDE, D, M, F, par float64
	var apogee bool
	EarthRad := 6378.14

	k := sk
	t = k - math.Floor(k)
	if t > 0.499 && t < 0.501 {
		apogee = true
	} else if t > 0.999 || t < 0.001 {
		apogee = false
	} else {
		return nil
	}

	t = k / 1325.55
	t2 = t * t
	t3 = t2 * t
	t4 = t3 * t

	/* Mean time of perigee or apogee */
	JDE = 2451534.6698 + 27.55454989*k -
		0.0006691*t2 -
		0.000001098*t3 +
		0.0000000052*t4

	/* Mean elongation of the Moon */
	D = 171.9179 + 335.9106046*k -
		0.0100383*t2 -
		0.00001156*t3 +
		0.000000055*t4

	/* Mean anomaly of the Sun */
	M = 347.3477 + 27.1577721*k -
		0.0008130*t2 -
		0.0000010*t3

	/* Moon's argument of latitude */
	F = 316.6109 + 364.5287911*k -
		0.0125053*t2 -
		0.0000148*t3

	// Determine which coefficients to use based on apogee flag
	var argtab, coeff []float64
	var tfix []int
	var tfixc []float64

	if apogee {
		argtab = apoarg
		coeff = apocoeff
		tfix = apotft
		tfixc = apotfc
	} else {
		argtab = periarg
		coeff = pericoeff
		tfix = peritft
		tfixc = peritfc
	}

	JDE += sumser(math.Sin, D, M, F, t, argtab, coeff, tfix, tfixc)

	// Use different coefficients for the second sumser call
	if apogee {
		argtab = apoparg
		coeff = apopcoeff
		tfix = apoptft
		tfixc = apoptfc
	} else {
		argtab = periparg
		coeff = peripcoeff
		tfix = periptft
		tfixc = periptfc
	}

	par = sumser(math.Cos, D, M, F, t, argtab, coeff, tfix, tfixc)

	par = dtr(par / 3600.0)
	return []float64{JDE, par, EarthRad / math.Sin(par)}
}

func truephase(k, phase float64) float64 {
	var t, t2, t3, pt, m, mprime, f float64
	SynMonth := 29.53058868 // Synodic month (mean time from new to next new Moon)

	k += phase           // Add phase to new moon time
	t = k / 1236.85      // Time in Julian centuries from 1900 January 0.5
	t2 = t * t           // Square for frequent use
	t3 = t2 * t          // Cube for frequent use
	pt = 2415020.75933 + // Mean time of phase
		SynMonth*k +
		0.0001178*t2 -
		0.000000155*t3 +
		0.00033*dsin(166.56+132.87*t-0.009173*t2)

	m = 359.2242 + // Sun's mean anomaly
		29.10535608*k -
		0.0000333*t2 -
		0.00000347*t3
	mprime = 306.0253 + // Moon's mean anomaly
		385.81691806*k +
		0.0107306*t2 +
		0.00001236*t3
	f = 21.2964 + // Moon's argument of latitude
		390.67050646*k -
		0.0016528*t2 -
		0.00000239*t3

	if (phase < 0.01) || (math.Abs(phase-0.5) < 0.01) {
		// Corrections for New and Full Moon
		pt += (0.1734-0.000393*t)*dsin(m) +
			0.0021*dsin(2*m) -
			0.4068*dsin(mprime) +
			0.0161*dsin(2*mprime) -
			0.0004*dsin(3*mprime) +
			0.0104*dsin(2*f) -
			0.0051*dsin(m+mprime) -
			0.0074*dsin(m-mprime) +
			0.0004*dsin(2*f+m) -
			0.0004*dsin(2*f-m) -
			0.0006*dsin(2*f+mprime) +
			0.0010*dsin(2*f-mprime) +
			0.0005*dsin(m+2*mprime)
	} else if (math.Abs(phase-0.25) < 0.01) || (math.Abs(phase-0.75) < 0.01) {
		pt += (0.1721-0.0004*t)*dsin(m) +
			0.0021*dsin(2*m) -
			0.6280*dsin(mprime) +
			0.0089*dsin(2*mprime) -
			0.0004*dsin(3*mprime) +
			0.0079*dsin(2*f) -
			0.0119*dsin(m+mprime) -
			0.0047*dsin(m-mprime) +
			0.0003*dsin(2*f+m) -
			0.0004*dsin(2*f-m) -
			0.0006*dsin(2*f+mprime) +
			0.0021*dsin(2*f-mprime) +
			0.0003*dsin(m+2*mprime) +
			0.0004*dsin(m-2*mprime) -
			0.0003*dsin(2*m+mprime)

		if phase < 0.5 {
			// First quarter correction
			pt += 0.0028 - 0.0004*dcos(m) + 0.0003*dcos(mprime)
		} else {
			// Last quarter correction
			pt += -0.0028 + 0.0004*dcos(m) - 0.0003*dcos(mprime)
		}
	}
	return pt
}

func nearphase(jd float64, phaset []float64) string {
	var closest int
	dt := math.MaxFloat64
	var rs string

	for i := 0; i < len(phaset); i++ {
		absPhase := math.Abs(phaset[i])
		currentDt := math.Abs(jd - absPhase)
		if currentDt < dt {
			dt = currentDt
			closest = i
		}
	}

	if phaset[closest] < 0 {
		rs = "N"
	} else {
		rs = "F"
	}

	if jd > math.Abs(phaset[closest]) {
		rs += "+"
	} else {
		rs += "-"
	}

	if dt >= 1 {
		days := int(math.Floor(dt))
		rs += fmt.Sprintf("%dd", days)
		dt -= float64(days)
	} else {
		rs += "  "
	}

	hours := int(math.Floor((dt * 86400) / 3600))
	if hours < 10 {
		rs += " "
	}
	rs += fmt.Sprintf("%dh", hours)

	return rs
}

func edate(j float64) string {
	j += (30.0 / (24 * 60 * 60)) // Round to nearest minute
	_, datem, dated := jyear(j)
	timeh, timem, _ := jhms(j)

	return months[datem-1] + " " + pad(strconv.Itoa(dated), 2, " ") + " " +
		pad(strconv.Itoa(timeh), 2, " ") + ":" + pad(strconv.Itoa(timem), 2, "0")
}

func cuzcoDate(j float64) string {
	j -= 5.0 / 24.0 // 5 timezones west of UTC
	return edate(j)
}

func getMonth(datem int) time.Month {
	datem = datem - 1
	if datem < 0 {
		datem = 0
	}
	if datem > 11 {
		datem = 11
	}

	return monthsGo[datem]
}

func cuzcoDateTime(j float64) time.Time {
	datey, datem, dated := jyear(j)
	//t.AddDate(datey, datem, dated)

	j1 := j
	j1 -= 5.0 / 24.0
	j1 += (30.0 / (24 * 60 * 60))

	timeh, timem, times := jhms(j1)

	t := time.Date(datey, getMonth(datem), dated, timeh, timem, times, 0, time.UTC)
	return t
}

func cuzcoNight(j float64) float64 {
	j -= 5.0 / 24.0      // 5 timezones west of UTC
	j -= 6.0 / 24.0      // anything up to 6am is considered previous night
	j += 0.5             // Astronomical to civil
	return math.Floor(j) // round to days
}

func dsin(x float64) float64 {
	return math.Sin(dtr(x))
}

func dcos(x float64) float64 {
	return math.Cos(dtr(x))
}

var months = []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"}
var monthsGo = []time.Month{time.January, time.February, time.March, time.April, time.May, time.June, time.July, time.August, time.September, time.October, time.November, time.December}

var periarg = []float64{
	/*  D,  M,  F   */
	2, 0, 0,
	4, 0, 0,
	6, 0, 0,
	8, 0, 0,
	2, -1, 0,
	0, 1, 0,
	10, 0, 0,
	4, -1, 0,
	6, -1, 0,
	12, 0, 0,
	1, 0, 0,
	8, -1, 0,
	14, 0, 0,
	0, 0, 2,
	3, 0, 0,
	10, -1, 0,
	16, 0, 0,
	12, -1, 0,
	5, 0, 0,
	2, 0, 2,
	18, 0, 0,
	14, -1, 0,
	7, 0, 0,
	2, 1, 0,
	20, 0, 0,
	1, 1, 0,
	16, -1, 0,
	4, 1, 0,
	9, 0, 0,
	4, 0, 2,

	2, -2, 0,
	4, -2, 0,
	6, -2, 0,
	22, 0, 0,
	18, -1, 0,
	6, 1, 0,
	11, 0, 0,
	8, 1, 0,
	4, 0, -2,
	6, 0, 2,
	3, 1, 0,
	5, 1, 0,
	13, 0, 0,
	20, -1, 0,
	3, 2, 0,
	4, -2, 2,
	1, 2, 0,
	22, -1, 0,
	0, 0, 4,
	6, 0, -2,
	2, 1, -2,
	0, 2, 0,
	0, -1, 2,
	2, 0, 4,
	0, -2, 2,
	2, 2, -2,
	24, 0, 0,
	4, 0, -4,
	2, 2, 0,
	1, -1, 0,
}

var pericoeff = []float64{
	-1.6769,
	0.4589,
	-0.1856,
	0.0883,
	-0.0773,
	0.0502,
	-0.0460,
	0.0422,
	-0.0256,
	0.0253,
	0.0237,
	0.0162,
	-0.0145,
	0.0129,
	-0.0112,
	-0.0104,
	0.0086,
	0.0069,
	0.0066,
	-0.0053,
	-0.0052,
	-0.0046,
	-0.0041,
	0.0040,
	0.0032,
	-0.0032,
	0.0031,
	-0.0029,
	0.0027,
	0.0027,

	-0.0027,
	0.0024,
	-0.0021,
	-0.0021,
	-0.0021,
	0.0019,
	-0.0018,
	-0.0014,
	-0.0014,
	-0.0014,
	0.0014,
	-0.0014,
	0.0013,
	0.0013,
	0.0011,
	-0.0011,
	-0.0010,
	-0.0009,
	-0.0008,
	0.0008,
	0.0008,
	0.0007,
	0.0007,
	0.0007,
	-0.0006,
	-0.0006,
	0.0006,
	0.0005,
	0.0005,
	-0.0004,
}

var peritft = []int{
	4,
	5,
	7,
	-1,
}

var peritfc = []float64{
	0.00019,
	-0.00013,
	-0.00011,
}

var apoarg = []float64{
	/*  D,  M,  F   */
	2, 0, 0,
	4, 0, 0,
	0, 1, 0,
	2, -1, 0,
	0, 0, 2,
	1, 0, 0,
	6, 0, 0,
	4, -1, 0,
	2, 0, 2,
	1, 1, 0,
	8, 0, 0,
	6, -1, 0,
	2, 0, -2,
	2, -2, 0,
	3, 0, 0,
	4, 0, 2,

	8, -1, 0,
	4, -2, 0,
	10, 0, 0,
	3, 1, 0,
	0, 2, 0,
	2, 1, 0,
	2, 2, 0,
	6, 0, 2,
	6, -2, 0,
	10, -1, 0,
	5, 0, 0,
	4, 0, -2,
	0, 1, 2,
	12, 0, 0,
	2, -1, 2,
	1, -1, 0,
}

var apocoeff = []float64{
	0.4392,
	0.0684,
	0.0456,
	0.0426,
	0.0212,
	-0.0189,
	0.0144,
	0.0113,
	0.0047,
	0.0036,
	0.0035,
	0.0034,
	-0.0034,
	0.0022,
	-0.0017,
	0.0013,

	0.0011,
	0.0010,
	0.0009,
	0.0007,
	0.0006,
	0.0005,
	0.0005,
	0.0004,
	0.0004,
	0.0004,
	-0.0004,
	-0.0004,
	0.0003,
	0.0003,
	0.0003,
	-0.0003,
}

var apotft = []int{
	2,
	3,
	-1,
}

var apotfc = []float64{
	-0.00011,
	-0.00011,
}

var periparg = []float64{
	/*  D,  M,  F   */
	0, 0, 0,
	2, 0, 0,
	4, 0, 0,
	2, -1, 0,
	6, 0, 0,
	1, 0, 0,
	8, 0, 0,
	0, 1, 0,
	0, 0, 2,
	4, -1, 0,
	2, 0, -2,
	10, 0, 0,
	6, -1, 0,
	3, 0, 0,
	2, 1, 0,
	1, 1, 0,
	12, 0, 0,
	8, -1, 0,
	2, 0, 2,
	2, -2, 0,
	5, 0, 0,
	14, 0, 0,

	10, -1, 0,
	4, 1, 0,
	12, -1, 0,
	4, -2, 0,
	7, 0, 0,
	4, 0, 2,
	16, 0, 0,
	3, 1, 0,
	1, -1, 0,
	6, 1, 0,
	0, 2, 0,
	14, -1, 0,
	2, 2, 0,
	6, -2, 0,
	2, -1, -2,
	9, 0, 0,
	18, 0, 0,
	6, 0, 2,
	0, -1, 2,
	16, -1, 0,
	4, 0, -2,
	8, 1, 0,
	11, 0, 0,
	5, 1, 0,
	20, 0, 0,
}

var peripcoeff = []float64{
	3629.215,
	63.224,
	-6.990,
	2.834,
	1.927,
	-1.263,
	-0.702,
	0.696,
	-0.690,
	-0.629,
	-0.392,
	0.297,
	0.260,
	0.201,
	-0.161,
	0.157,
	-0.138,
	-0.127,
	0.104,
	0.104,
	-0.079,
	0.068,

	0.067,
	0.054,
	-0.038,
	-0.038,
	0.037,
	-0.037,
	-0.035,
	-0.030,
	0.029,
	-0.025,
	0.023,
	0.023,
	-0.023,
	0.022,
	-0.021,
	-0.020,
	0.019,
	0.017,
	0.014,
	-0.014,
	0.013,
	0.012,
	0.011,
	0.010,
	-0.010,
}

var periptft = []int{
	3,
	7,
	9,
	-1,
}

var periptfc = []float64{
	-0.0071,
	-0.0017,
	0.0016,
}

var apoparg = []float64{
	/*  D,  M,  F   */
	0, 0, 0,
	2, 0, 0,
	1, 0, 0,
	0, 0, 2,
	0, 1, 0,
	4, 0, 0,
	2, -1, 0,
	1, 1, 0,
	4, -1, 0,
	6, 0, 0,
	2, 1, 0,
	2, 0, 2,
	2, 0, -2,
	2, -2, 0,
	2, 2, 0,
	0, 2, 0,
	6, -1, 0,
	8, 0, 0,
}

var apopcoeff = []float64{
	3245.251,
	-9.147,
	-0.841,
	0.697,
	-0.656,
	0.355,
	0.159,
	0.127,
	0.065,

	0.052,
	0.043,
	0.031,
	-0.023,
	0.022,
	0.019,
	-0.016,
	0.014,
	0.010,
}

var apoptft = []int{
	4,
	-1,
}

var apoptfc = []float64{
	0.0016,
	-1,
}
