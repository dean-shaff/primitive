package primitive

import "github.com/fogleman/gg"

type Shape interface {
	Rasterize() []Scanline
	Copy() Shape
	Mutate()
	Draw(dc *gg.Context, scale float64)
	SVG(attrs string) string
	Area() float64
}

type BlueDotSessionsShapeConfig struct {
	StartRect int
	EndRect int
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
	ShapeTypeRightFacingTriangle
	ShapeTypeDiamond
	ShapeTypeBlueDotSessions
)
