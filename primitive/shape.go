package primitive

import "github.com/fogleman/gg"

type Shape interface {
	Rasterize() []Scanline
	Copy() Shape
	Mutate()
	Draw(dc *gg.Context, scale float64)
	SVG(attrs string) string
}

type BlueDotSessionsShapeConfig struct {
	QuadPercent float64
	StartTriangles int
}

type ShapeType int

const (
	ShapeTypeAny ShapeType = iota
	ShapeTypeTriangle
	ShapeTypeRectangle
	ShapeTypeEllipse
	ShapeTypeCircle
	ShapeTypeRotatedRectangle
	ShapeTypeQuadratic
	ShapeTypeRotatedEllipse
	ShapeTypePolygon
	ShapeTypeBlueDotSessions
)
