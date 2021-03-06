package primitive

import (
	"fmt"
	"strings"
	"math"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/raster"
)

type Polygon struct {
	Worker *Worker
	Order  int
	Convex bool
	MinAngle float64
	SizeFactor float64
	BoundsFactor int
	X, Y   []float64
}

// const m = 0
// const minAngle = 15


func NewRandomPolygon(worker *Worker, order int, convex bool, minAngle float64, sizeFactor float64, boundsFactor int) *Polygon {

	w := worker.W
	h := worker.H
	rnd := worker.Rnd
	x := make([]float64, order)
	y := make([]float64, order)
	x[0] = rnd.Float64() * float64(w)
	y[0] = rnd.Float64() * float64(h)
	for i := 1; i < order; i++ {
		x[i] = clamp(x[0] + rnd.Float64()*sizeFactor - sizeFactor/2, float64(-boundsFactor), float64(w-1+boundsFactor))
		y[i] = clamp(y[0] + rnd.Float64()*sizeFactor - sizeFactor/2, float64(-boundsFactor), float64(h-1+boundsFactor))
	}
	p := &Polygon{worker, order, convex, minAngle, sizeFactor, boundsFactor, x, y}
	p.Mutate()
	return p
}

func (p *Polygon) Draw(dc *gg.Context, scale float64) {
	dc.NewSubPath()
	for i := 0; i < p.Order; i++ {
		dc.LineTo(p.X[i], p.Y[i])
	}
	dc.ClosePath()
	dc.Fill()
}

func (p *Polygon) SVG(attrs string) string {
	ret := fmt.Sprintf(
		"<polygon %s points=\"",
		attrs)
	points := make([]string, len(p.X))
	for i := 0; i < len(p.X); i++ {
		points[i] = fmt.Sprintf("%f,%f", p.X[i], p.Y[i])
	}

	return ret + strings.Join(points, ",") + "\" />"
}

func (p *Polygon) Copy() Shape {
	a := *p
	a.X = make([]float64, p.Order)
	a.Y = make([]float64, p.Order)
	copy(a.X, p.X)
	copy(a.Y, p.Y)
	return &a
}

func (p *Polygon) Mutate() {
	// vv("Polygon.Mutate")
	w := p.Worker.W
	h := p.Worker.H
	boundsFactor := p.BoundsFactor
	sizeFactor := p.SizeFactor
	rnd := p.Worker.Rnd
	for {
		if rnd.Float64() < 0.25 {
			i := rnd.Intn(p.Order)
			j := rnd.Intn(p.Order)
			p.X[i], p.Y[i], p.X[j], p.Y[j] = p.X[j], p.Y[j], p.X[i], p.Y[i]
		} else {
			i := rnd.Intn(p.Order)
			p.X[i] = clamp(p.X[i]+rnd.NormFloat64()*sizeFactor/2, float64(-boundsFactor), float64(w-1+boundsFactor))
			p.Y[i] = clamp(p.Y[i]+rnd.NormFloat64()*sizeFactor/2, float64(-boundsFactor), float64(h-1+boundsFactor))
		}

		if p.Valid() {
			break
		}
	}
}

func (p *Polygon) Valid() bool {
	if !p.Convex {
		return true
	}
	var sign bool
	for a := 0; a < p.Order; a++ {
		i := (a + 0) % p.Order
		j := (a + 1) % p.Order
		k := (a + 2) % p.Order
		c := cross3(p.X[i], p.Y[i], p.X[j], p.Y[j], p.X[k], p.Y[k])
		if a == 0 {
			sign = c > 0
		} else if c > 0 != sign {
			return false
		}
	}
	// check internal angles of polygon.
	for a := 0; a < p.Order; a++ {
		i := (a + 0) % p.Order
		j := (a + 1) % p.Order
		k := (a + 2) % p.Order
		a := angleBetween(p.X[i], p.Y[i], p.X[j], p.Y[j], p.X[k], p.Y[k])
		if a < p.MinAngle {
			return false
		}
	}

	w, h := float64(p.Worker.W), float64(p.Worker.H)
	boundsFactor := float64(p.BoundsFactor)
	for idx := 0; idx < p.Order; idx++ {
		if p.X[idx] < -boundsFactor || p.X[idx] > w + boundsFactor {
			return false
		}
		if p.Y[idx] < -boundsFactor || p.Y[idx] > h + boundsFactor {
			return false
		}
	}


	return true
}

func cross3(x1, y1, x2, y2, x3, y3 float64) float64 {
	dx1 := x2 - x1
	dy1 := y2 - y1
	dx2 := x3 - x2
	dy2 := y3 - y2
	return dx1*dy2 - dy1*dx2
}

func (p *Polygon) Rasterize() []Scanline {
	var path raster.Path
	for i := 0; i <= p.Order; i++ {
		f := fixp(p.X[i%p.Order], p.Y[i%p.Order])
		if i == 0 {
			path.Start(f)
		} else {
			path.Add1(f)
		}
	}
	return fillPath(p.Worker, path)
}


func (p *Polygon) Area() float64 {

	x := p.X
	y := p.Y
	x = append(x, x[0])
	y = append(y, y[0])

	order := len(x)
	// var area float64 = p.X[order - 1]*p.X[0] - p.X[0]*p.Y[order - 1]
	var area float64 = 0.0

	for idx := 0; idx < order - 1; idx++ {
		area += x[idx]*y[idx + 1] - x[idx + 1]*y[idx]
	}

	area = math.Abs(area)/2.0

	return area
}

func mag (x, y float64) float64 {
	return math.Sqrt(math.Pow(x, 2) + math.Pow(y, 2))
}

func dot (x0, y0, x1, y1 float64) float64 {
	return x0*x1 + y0*y1
}

// Get the angle created by vector between (x1, y1) -> (x0, y0) and (x1, y1) -> (x2, y2). Returns angle in degrees
func angleBetween (x0, y0, x1, y1, x2, y2 float64) float64 {
	x0 -= x1
	y0 -= y1
	x2 -= x1
	y2 -= y1

	mag0 := mag(x0, y0)
	mag2 := mag(x2, y2)
	dot02 := dot(x0, y0, x2, y2)

	angle := degrees(math.Acos(dot02 / (mag0 * mag2)))

	return angle
}
