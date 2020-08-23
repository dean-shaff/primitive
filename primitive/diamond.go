package primitive

import (
	"github.com/fogleman/gg"
)

type Diamond struct {
  polygon Polygon
}


func NewRandomDiamond(worker *Worker, order int, convex bool, minAngle, sizeFactor float64, boundsFactor int) *Diamond {
	var p *Polygon

	if order == 4 {
		w := worker.W
		h := worker.H
		rnd := worker.Rnd

		x0 := rnd.Float64() * float64(w)
		y0 := rnd.Float64() * float64(h)
		x, y := Init(worker, x0, y0, sizeFactor)
		// vv("NewRandomDiamond: x=%f, y=%f\n", x, y)
		p = &Polygon{worker, order, convex, minAngle, sizeFactor, boundsFactor, x, y}
	}
	diam := &Diamond{*p}
	diam.Mutate()
	return diam
}

func Init(worker *Worker, x0, y0, sizeFactor float64) ([]float64, []float64) {
	rnd := worker.Rnd

	x := make([]float64, 4)
	y := make([]float64, 4)
	x[0] = x0
	y[0] = y0

	x[1] = x[0] + rnd.Float64() * sizeFactor + 0.1*sizeFactor
	y[1] = y[0] + rnd.Float64() * sizeFactor + 0.1*sizeFactor

	x[2] = x[0]
	y[2] = y[1] + rnd.Float64() * sizeFactor

	x[3] = x[0] - (rnd.Float64()*0.2*sizeFactor) - 0.1*sizeFactor
	y[3] = y[1]
	return x, y
}

func (diam *Diamond) Draw(dc *gg.Context, scale float64) {
  diam.polygon.Draw(dc, scale)
}

func (diam *Diamond) SVG(attrs string) string {
  return diam.polygon.SVG(attrs)
}

func (diam *Diamond) Copy() Shape {
  a := diam.polygon
  a.X = make([]float64, diam.polygon.Order)
  a.Y = make([]float64, diam.polygon.Order)
  copy(a.X, diam.polygon.X)
  copy(a.Y, diam.polygon.Y)
  return &Diamond{a}
  // return a
}

func (diam *Diamond) Mutate() {
  p := diam.polygon
	w := p.Worker.W
	h := p.Worker.H
	rnd := p.Worker.Rnd
	sizeFactor := p.SizeFactor
	boundsFactor := p.BoundsFactor
	boundsFactorf := float64(boundsFactor)
	var delta_x1 float64
	for {
		idx := rnd.Intn(p.Order)
		if idx % 2 == 0 {
			p.X[idx] = clamp(p.X[idx]+rnd.NormFloat64()*sizeFactor, p.X[3], p.X[1])
			p.Y[idx] = clamp(p.Y[idx]+rnd.NormFloat64()*sizeFactor, -boundsFactorf, float64(h-1+boundsFactor))
			p.X[(idx + 2) % p.Order] = p.X[idx]
		} else {
			if idx == 1 {
				p.X[idx] = clamp(p.X[idx]+rnd.NormFloat64()*sizeFactor, -boundsFactorf, float64(w-1+boundsFactor))
			} else if idx == 3 {
				delta_x1 = p.X[1] - p.X[0]
				low_bound_3 := p.X[0] - (delta_x1 * 0.5)
				if low_bound_3 < -boundsFactorf {
					low_bound_3 = -boundsFactorf
				}
				p.X[idx] = clamp(p.X[idx]+rnd.NormFloat64()*sizeFactor, low_bound_3, p.X[0] - 0.1*sizeFactor)
			}
			p.Y[idx] = clamp(p.Y[idx]+rnd.NormFloat64()*sizeFactor, p.Y[0] + 0.3*sizeFactor, p.Y[2] - 0.3*sizeFactor)
			p.Y[(idx + 2) % p.Order] = p.Y[idx]
		}
		if diam.Valid() {
			break
		}
	}
}

func (diam *Diamond) Valid() bool {
	p := diam.polygon
	delta_x1 := p.X[1] - p.X[0]
	delta_x3 := p.X[0] - p.X[3]
	if delta_x3 > delta_x1 {
		return false
	}
	return diam.polygon.Valid()
}


func (diam *Diamond) Rasterize() []Scanline {
  return diam.polygon.Rasterize()
}


func (diam *Diamond) Area() float64 {
  return diam.polygon.Area()
}
