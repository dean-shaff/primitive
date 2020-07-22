package primitive

import (
	"github.com/fogleman/gg"
)

// Right Facing Triangle
type RFTriangle struct {
  triangle Triangle
  MutateFactor int
  MutateYTol float64
}

func NewRandomRFTriangle(worker *Worker) *RFTriangle {
  tol := 0.2
  rnd := worker.Rnd
  x1 := rnd.Intn(worker.W - 31)
  y1 := rnd.Intn(worker.H/2)
  x2 := x1 // + rnd.Intn(31) - 15
  y2 := y1 + rnd.Intn(worker.H/2) + 1
  x3 := x1 + rnd.Intn(31)
  v := y2 - y1
  bottom := int(float64(y1) + tol*float64(v))
  y3 := bottom + rnd.Intn(v)
  // vv("(%d, %d), (%d, %d), (%d, %d)\n", x1, y1, x2, y2, x3, y3)
	t := Triangle{worker, x1, y1, x2, y2, x3, y3}
  rft := &RFTriangle{t, 10, tol}
  rft.Mutate()
	return rft
}

func (t *RFTriangle) Draw(dc *gg.Context, scale float64) {
  t.triangle.Draw(dc, scale)
}

func (t *RFTriangle) SVG(attrs string) string {
  return t.triangle.SVG(attrs)
}

func (t *RFTriangle) Copy() Shape {
  a := *t
  return &a
}

func (rft *RFTriangle) Mutate() {
  // rft.triangle.Mutate()
  // return
  // vv("RFTriangle Mutate\n")
  t := rft.triangle
  worker := t.Worker
  w := worker.W
  h := worker.H
  rnd := worker.Rnd
  mfloat := float64(rft.MutateFactor)
  // mint := rft.MutateFactor
  for {
    switch rnd.Intn(3) {
    case 0:
      rft.triangle.X1 = clampInt(t.X1+int(rnd.NormFloat64()*mfloat), 0, w-1)
      rft.triangle.Y1 = clampInt(t.Y1+int(rnd.NormFloat64()*mfloat), 0, h-1)
      rft.triangle.X2 = rft.triangle.X1
      rft.triangle.Y2 = h-1 - rft.triangle.Y1
    case 1:
      rft.triangle.X2 = clampInt(t.X2+int(rnd.NormFloat64()*mfloat), 0, w-1)
      rft.triangle.Y2 = clampInt(t.Y2+int(rnd.NormFloat64()*mfloat), 0, h-1)
      rft.triangle.X1 = rft.triangle.X2
      rft.triangle.Y1 = h-1 - rft.triangle.Y2
    case 2:
      rft.triangle.X3 = clampInt(t.X3+int(rnd.NormFloat64()*mfloat), 0, w-1)
      v := t.Y2 - t.Y1 + 1
      bottom := int(float64(t.Y1) + rft.MutateYTol*float64(v))
      rft.triangle.Y3 = bottom + rnd.Intn(v)
    }

    if rft.Valid() {
      break
    }
  }
}

func (rft *RFTriangle) Valid() bool {
  t := rft.triangle
  if t.X3 < t.X1 {
    // vv("here 0\n")
    return false
  }
  if t.Y3 > t.Y2 || t.Y3 < t.Y1 {
    // vv("Y1=%d, Y2=%d, Y3=%d\n", t.Y1, t.Y1, t.Y3)
    return false
  }
  return rft.triangle.Valid()
}

func (t *RFTriangle) Rasterize() []Scanline {
  return t.triangle.Rasterize()
}

func (t *RFTriangle) Area() float64 {
  return t.triangle.Area()
}
