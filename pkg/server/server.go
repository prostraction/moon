package server

import (
	"errors"
	"math"
	"moon/pkg/moon"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

type Server struct {
	Days         string
	JD           float64
	L            float64
	Illumination string
	Test         string
	Test2        string
	Test3        []*moon.MoonTableElement
	Test4        float64
}

func (s *Server) NewRouter() *fiber.App {
	app := fiber.New(fiber.Config{
		Prefork:       true,
		ServerHeader:  "Fiber",
		CaseSensitive: true,
		StrictRouting: true,
	})

	app.Get("/v1.1/getMoonPhase", s.getMoonPhaseV1)
	app.Get("/v1.2/getMoonPhase", s.getMoonPhaseV2)
	app.Get("/v1.3/getMoonPhase", s.getMoonPhaseV3)
	return app
}

func (s *Server) getMoonPhaseV1(c *fiber.Ctx) error {
	D := c.Query("d", "default")
	DInt, err := strconv.Atoi(D)
	if err != nil {
		return err
	}

	M := c.Query("m", "default")
	MInt, err := strconv.Atoi(M)
	if err != nil {
		return err
	}

	Y := c.Query("y", "default")
	yInt, err := strconv.Atoi(Y)
	if err != nil {
		return err
	}

	L := moon.CalcMoonNumber(yInt)
	LCalc := ((L * 11) - 14) % 30

	s.Days = strconv.Itoa((LCalc + DInt + MInt) % 30)
	//s.Illumination = strconv.Itoa((((LCalc + DInt + MInt) % 30) / 4) % 4)
	return c.JSON(s)
}

func (s *Server) getMoonPhaseV2(c *fiber.Ctx) error {
	D := c.Query("d", "default")
	DInt, err := strconv.Atoi(D)
	if err != nil {
		return err
	}

	M := c.Query("m", "default")
	MInt, err := strconv.Atoi(M)
	if err != nil {
		return err
	}

	Y := c.Query("y", "default")
	yInt, err := strconv.Atoi(Y)
	if err != nil {
		return err
	}

	referenceDate := time.Date(2000, time.January, 6, 0, 0, 0, 0, time.UTC)
	givenDate := time.Date(yInt, time.Month(MInt), DInt, 0, 0, 0, 0, time.UTC)
	days := givenDate.Sub(referenceDate).Hours() / 24

	moons := days / 29.53058770576

	// Step 3: Return the fractional part
	s.Days = strconv.FormatFloat(29.53058770576*(moons-float64(int(moons))), 'f', -3, 64)
	//return moons - float64(int(moons)), nil

	if MInt < 3 {
		yInt--
		MInt += 12
	}
	A := float64(yInt) / 100
	B := A / 4
	C := 2 - A + B
	E := 365.25 * float64(yInt+4716)
	F := 30.6001 * (float64(MInt) + 1)
	JD := C + float64(float64(DInt)) + E + F - 1524.5

	if (JD - 2451545.0) == 0 {
		return errors.New("JD is 0")
	}
	T := (JD - 2451545.0) / 36525
	L := 29.5305888531 + 0.00000021621*T - 3.64*math.Pow(10, -12)*(T*T)

	s.L = L
	s.JD = JD
	//s.Illumination = strconv.Itoa((((LCalc + DInt + MInt) % 30) / 4) % 4)
	return c.JSON(s)
}

func (s *Server) getMoonPhaseV3(c *fiber.Ctx) error {
	D := c.Query("d", "default")
	DInt, err := strconv.Atoi(D)
	if err != nil {
		return err
	}

	M := c.Query("m", "default")
	MInt, err := strconv.Atoi(M)
	if err != nil {
		return err
	}

	Y := c.Query("y", "default")
	yInt, err := strconv.Atoi(Y)
	if err != nil {
		return err
	}

	L := moon.CalcMoonNumber(yInt)
	LCalc := ((L * 11) - 14) % 30

	s.Days = strconv.Itoa((LCalc + DInt + MInt) % 30)

	///testingData := []*moon.MoonTableElement{}
	var test4 time.Duration
	_, _, _, test4 = moon.Gen(yInt, MInt, DInt)
	s.Test4 = test4.Hours()
	s.Test4 = s.Test4 / 24

	//for i := range testingData {
	//	s.Test3 += strconv.FormatFloat(testingData[i].T1, 'E', -1, 64)
	//}

	//s.Illumination = strconv.Itoa((((LCalc + DInt + MInt) % 30) / 4) % 4)
	return c.JSON(s)
}
