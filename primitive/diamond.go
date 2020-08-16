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
		x := make([]float64, order)
		y := make([]float64, order)
		x[0] = rnd.Float64() * float64(w)
		y[0] = rnd.Float64() * float64(h)

		x[1] = x[0] + rnd.Float64() * sizeFactor
		y[1] = y[0] + rnd.Float64() * sizeFactor

		x[2] = x[0]
		y[2] = y[1] + rnd.Float64() * sizeFactor

		x[3] = x[0] - rnd.Float64() * 0.5*sizeFactor - 0.3*sizeFactor
		y[3] = y[1]
		p = &Polygon{worker, order, convex, minAngle, sizeFactor, boundsFactor, x, y}
	}
	diam := &Diamond{*p}
	diam.Mutate()
	return diam


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
	for {
		idx := rnd.Intn(p.Order)

		if idx % 2 == 0 {
			p.X[idx] = clamp(p.X[idx]+rnd.NormFloat64()*sizeFactor, p.X[3] + 0.5*sizeFactor, p.X[1] - 0.5*sizeFactor)
			p.Y[idx] = clamp(p.Y[idx]+rnd.NormFloat64()*sizeFactor, float64(-boundsFactor), float64(h-1+boundsFactor))
			p.X[(idx + 2) % p.Order] = p.X[idx]
		} else {
			if idx == 1 {
				p.X[idx] = clamp(p.X[idx]+rnd.NormFloat64()*sizeFactor, float64(-boundsFactor), float64(w-1+boundsFactor))
			} else if idx == 3 {
				p.X[idx] = clamp(p.X[idx]+rnd.NormFloat64()*sizeFactor, float64(-boundsFactor), p.X[0] - 0.5*sizeFactor)
			}
			p.Y[idx] = clamp(p.Y[idx]+rnd.NormFloat64()*sizeFactor, p.Y[0] + 0.5*sizeFactor, p.X[2] - 0.5*sizeFactor)
			p.Y[(idx + 2) % p.Order] = p.Y[idx]
		}

		// p.X[i] = clamp(p.X[i]+rnd.NormFloat64()*incre, -m, float64(w-1+m))
		// p.Y[i] = clamp(p.Y[i]+rnd.NormFloat64()*incre, -m, float64(h-1+m))
		// for idx := i; idx < p.Order; idx++ {
		// 	switch idx {
		// 	case 1:
		// 		p.X[1] = clamp(p.X[0] + rnd.Float64() * incre, -m, float64(w-1+m))
		// 		p.Y[1] = clamp(p.Y[0] + rnd.Float64() * incre, -m, float64(h-1+m))
		// 	case 2:
		// 		p.X[2] = p.X[0] //+ rnd.NormFloat64()*0.1*incre
		// 		p.Y[2] = p.Y[1] + rnd.Float64() * incre
		// 	case 3:
		// 		p.X[3] = p.X[2] - rnd.Float64()*0.5*incre
		// 		p.Y[3] = p.Y[1] //+ rnd.NormFloat64()*0.1*incre
		// 	}
		// }

		if diam.Valid() {
			break
		}
	}
}

func (diam *Diamond) Valid() bool {
	return diam.polygon.Valid()
}


func (diam *Diamond) Rasterize() []Scanline {
  return diam.polygon.Rasterize()
}


func (diam *Diamond) Area() float64 {
  return diam.polygon.Area()
}
