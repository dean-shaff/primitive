package primitive

import (
	"image"
	"math/rand"
	"time"

	"github.com/golang/freetype/raster"
)

type Worker struct {
	W, H       int
	Target     *image.RGBA
	Current    *image.RGBA
	Buffer     *image.RGBA
	Rasterizer *raster.Rasterizer
	Lines      []Scanline
	Heatmap    *Heatmap
	Rnd        *rand.Rand
	Score      float64
	BlackThresh float64
	AreaThresh float64
	Counter    int
}

func NewWorker(target *image.RGBA, blackThresh, areaThresh float64) *Worker {
	w := target.Bounds().Size().X
	h := target.Bounds().Size().Y
	worker := Worker{}
	worker.W = w
	worker.H = h
	worker.Target = target
	worker.Buffer = image.NewRGBA(target.Bounds())
	worker.Rasterizer = raster.NewRasterizer(w, h)
	worker.Lines = make([]Scanline, 0, 4096) // TODO: based on height
	worker.Heatmap = NewHeatmap(w, h)
	worker.Rnd = rand.New(rand.NewSource(time.Now().UnixNano()))
	worker.BlackThresh = blackThresh
	worker.AreaThresh = areaThresh
	vv("NewWorker: BlackThresh=%.2f, AreaThresh=%.2f\n", worker.BlackThresh, worker.AreaThresh)
	return &worker
}

func (worker *Worker) Init(current *image.RGBA, score float64) {
	worker.Current = current
	worker.Score = score
	worker.Counter = 0
	worker.Heatmap.Clear()
}

func (worker *Worker) Energy(shape Shape, alpha int) float64 {
	black := Color{0, 0, 0, alpha}
	worker.Counter++
	lines := shape.Rasterize()
	// worker.Heatmap.Add(lines)
	color := computeColor(worker.Target, worker.Current, lines, alpha)
	diff := RGBADiffColor(color, black)
	if diff < worker.BlackThresh {
		return 1.0
	}
	if worker.AreaThresh > 0.0 {
		area := shape.Area()
		if area != -1 {
			total_area := float64(worker.H * worker.W)
			frac := area / total_area
			if frac > worker.AreaThresh {
				// vv("Energy: area=%.2f, total area=%.2f, frac=%.2f\n", area, total_area, frac)
				return 1.0
			}
		}
	}

	copyLines(worker.Buffer, worker.Current, lines)
	drawLines(worker.Buffer, color, lines)
	return differencePartial(worker.Target, worker.Current, worker.Buffer, worker.Score, lines)
}

func (worker *Worker) BestHillClimbState(t ShapeType, a, n, age, m, idx int, fn NewShapeFunc) *State {
	var bestEnergy float64
	var bestState *State
	vv("BestHillClimbState: n=%d, m=%d\n", n, m)
	for i := 0; i < m; i++ {
		state := worker.BestRandomState(t, a, n, idx, fn)
		before := state.Energy()
		area_before := state.Shape.Area()
		state = HillClimb(state, age).(*State)
		energy := state.Energy()
		area_after := state.Shape.Area()
		vv("%dx random: %.6f -> %dx hill climb: %.6f (area %.1f -> %.1f)\n", n, before, age, energy, area_before, area_after)
		if i == 0 || energy < bestEnergy {
			bestEnergy = energy
			bestState = state
		}
	}
	return bestState
}

func (worker *Worker) BestRandomState(t ShapeType, a, n, idx int, fn NewShapeFunc) *State {
	var bestEnergy float64
	var bestState *State
	for i := 0; i < n; i++ {
		state := worker.RandomState(t, a, idx, fn)
		energy := state.Energy()
		// vv("BestRandomState: energy=%.2f, bestEnergy=%.2f\n", energy, bestEnergy)
		if i == 0 || energy < bestEnergy {
			bestEnergy = energy
			bestState = state
		}
	}
	return bestState
}

func NewBlueDotSessionsShapeFactory (quadPercent float64, startTriangles int) NewShapeFunc {
	return func (worker *Worker, a, idx int) Shape {
		// vv("NewBlueDotSessionsShapeFactory")
		if idx < startTriangles {
			return NewRandomTriangle(worker)
		} else {
			var rflt float64 = worker.Rnd.Float64()
			if (rflt < quadPercent) {
				return NewRandomPolygon(worker, 4, false)
			} else {
				return NewRandomTriangle(worker)
			}
		}
	}
}

type NewShapeFunc func(worker* Worker, a, idx int) Shape


func (worker *Worker) RandomState(t ShapeType, a, idx int, fn NewShapeFunc) *State {
	// vv("RandomState: a=%d\n", a)
	switch t {
	default:
		return worker.RandomState(ShapeType(worker.Rnd.Intn(8)+1), a, idx, fn)
	case ShapeTypeTriangle:
		return NewState(worker, NewRandomTriangle(worker), a)
	case ShapeTypeRectangle:
		return NewState(worker, NewRandomRectangle(worker), a)
	case ShapeTypeEllipse:
		return NewState(worker, NewRandomEllipse(worker), a)
	case ShapeTypeCircle:
		return NewState(worker, NewRandomCircle(worker), a)
	case ShapeTypeRotatedRectangle:
		return NewState(worker, NewRandomRotatedRectangle(worker), a)
	case ShapeTypeQuadratic:
		return NewState(worker, NewRandomQuadratic(worker), a)
	case ShapeTypeRotatedEllipse:
		return NewState(worker, NewRandomRotatedEllipse(worker), a)
	case ShapeTypePolygon:
		return NewState(worker, NewRandomPolygon(worker, 4, true), a)
	case ShapeTypeBlueDotSessions:
		return NewState(worker, fn(worker, a, idx), a)
	}
}
