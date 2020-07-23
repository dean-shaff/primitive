package primitive

import (
	"github.com/fogleman/gg"
)

type BDSPolygon struct {
  polygon Polygon
}

const incre = 40
// const m = 0

func NewRandomBDSPolygon(worker *Worker, order int, convex bool) *BDSPolygon {
	var p *Polygon

	if order == 4 {
		w := worker.W
		h := worker.H
		rnd := worker.Rnd
		x := make([]float64, order)
		y := make([]float64, order)
		x[0] = rnd.Float64() * float64(w)
		y[0] = rnd.Float64() * float64(h)

		x[1] = x[0] + rnd.Float64() * incre
		y[1] = y[0] + rnd.Float64() * incre

		x[2] = x[0]
		y[2] = y[1] + rnd.Float64() * incre

		x[3] = x[0] - rnd.Float64() * 0.5*incre
		y[3] = y[1]
		p = &Polygon{worker, order, convex, x, y}
	}
	bdsp := &BDSPolygon{*p}
	bdsp.Mutate()
	return bdsp


}

func (bdsp *BDSPolygon) Draw(dc *gg.Context, scale float64) {
  bdsp.polygon.Draw(dc, scale)
}

func (bdsp *BDSPolygon) SVG(attrs string) string {
  return bdsp.polygon.SVG(attrs)
}

func (bdsp *BDSPolygon) Copy() Shape {
  a := bdsp.polygon
  a.X = make([]float64, bdsp.polygon.Order)
  a.Y = make([]float64, bdsp.polygon.Order)
  copy(a.X, bdsp.polygon.X)
  copy(a.Y, bdsp.polygon.Y)
  return &BDSPolygon{a}
  // return a
}

func (bdsp *BDSPolygon) Mutate() {
  p := bdsp.polygon
	w := p.Worker.W
	h := p.Worker.H
	rnd := p.Worker.Rnd
	for {
		i := rnd.Intn(p.Order)

		p.X[i] = clamp(p.X[i]+rnd.NormFloat64()*incre, -m, float64(w-1+m))
		p.Y[i] = clamp(p.Y[i]+rnd.NormFloat64()*incre, -m, float64(h-1+m))

		for idx := i; idx < p.Order; idx++ {
			switch idx {
			case 1:
				p.X[1] = clamp(p.X[0] + rnd.Float64() * incre, -m, float64(w-1+m))
				p.Y[1] = clamp(p.Y[0] + rnd.Float64() * incre, -m, float64(h-1+m))
			case 2:
				p.X[2] = p.X[0] //+ rnd.NormFloat64()*0.1*incre
				p.Y[2] = p.Y[1] + rnd.Float64() * incre
			case 3:
				p.X[3] = p.X[2] - rnd.Float64()*0.5*incre
				p.Y[3] = p.Y[1] //+ rnd.NormFloat64()*0.1*incre
			}

			// switch idx {
			// case 1:
			// 	p.X[1] = p.X[1] + p.X[0]
			// 	p.Y[1] = p.Y[1] + p.Y[0]
			// case 2:
			// 	p.X[2] = p.X[0] //+ rnd.NormFloat64()*0.1*incre
			// 	p.Y[2] = p.Y[2] + p.Y[1]
			// case 3:
			// 	p.X[3] = p.X[2]
			// 	p.Y[3] = p.Y[1] //+ rnd.NormFloat64()*0.1*incre
			// }
		}

		if p.Valid() {
			break
		}
	}
}

func (bdsp *BDSPolygon) Valid() bool {
  p := bdsp.polygon
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
	return true
}


func (bdsp *BDSPolygon) Rasterize() []Scanline {
  return bdsp.polygon.Rasterize()
}


func (bdsp *BDSPolygon) Area() float64 {
  return bdsp.polygon.Area()
}
