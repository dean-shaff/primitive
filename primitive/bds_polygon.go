package primitive

import (
	"github.com/fogleman/gg"
)

type BDSPolygon struct {
  polygon Polygon
}

func NewRandomBDSPolygon(worker *Worker, order int, convex bool) *BDSPolygon {
  p := NewRandomPolygon(worker, order, convex)
  bdsp := &BDSPolygon{*p}
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
	const m = 16
  p := bdsp.polygon
	w := p.Worker.W
	h := p.Worker.H
	rnd := p.Worker.Rnd
	for {
		if rnd.Float64() < 0.25 {
			i := rnd.Intn(p.Order)
			j := rnd.Intn(p.Order)
			p.X[i], p.Y[i], p.X[j], p.Y[j] = p.X[j], p.Y[j], p.X[i], p.Y[i]
		} else {
			i := rnd.Intn(p.Order)
			p.X[i] = clamp(p.X[i]+rnd.NormFloat64()*16, -m, float64(w-1+m))
			p.Y[i] = clamp(p.Y[i]+rnd.NormFloat64()*16, -m, float64(h-1+m))
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
