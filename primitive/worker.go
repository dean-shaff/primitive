package primitive

import (
	"image"
	"math/rand"
	"time"
	// "fmt"

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
	LowerAreaThresh float64
	UpperAreaThresh float64
	Counter    int
}

func NewWorker(target *image.RGBA, blackThresh, lowerAreaThresh, upperAreaThresh float64, seed int64) *Worker {
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
	if seed == -1 {
		seed = time.Now().UnixNano()
	}
	vv("NewWorker: seed=%d\n", seed)
	worker.Rnd = rand.New(rand.NewSource(seed))
	worker.BlackThresh = blackThresh
	worker.UpperAreaThresh = upperAreaThresh
	worker.LowerAreaThresh = lowerAreaThresh
	vv("NewWorker: BlackThresh=%.2f, LowerAreaThresh=%.2f, UpperAreaThresh=%.2f\n", worker.BlackThresh, worker.LowerAreaThresh, worker.UpperAreaThresh)
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
	if worker.UpperAreaThresh > 0.0 {
		area := shape.Area()
		if area != -1 {
			total_area := float64(worker.H * worker.W)
			frac := area / total_area
			if frac > worker.UpperAreaThresh || frac < worker.LowerAreaThresh {
				// vv("Energy: area=%.2f, total area=%.2f, frac=%.2f\n", area, total_area, frac)
				return 1.0
			} else {
				vv("Energy: area=%.2f, total area=%.2f, frac=%.2f\n", area, total_area, frac)
			}
		}
	}

	copyLines(worker.Buffer, worker.Current, lines)
	drawLines(worker.Buffer, color, lines)
	return differencePartial(worker.Target, worker.Current, worker.Buffer, worker.Score, lines)
}

func (worker *Worker) BestHillClimbState(t ShapeType, a, n, age, m, idx int, fn NewShapeFunc, rand_val float64) *State {
	var bestEnergy float64
	var bestState *State
	// rand_val := worker.Rnd.Float64()
	v("BestHillClimbState: n=%d, m=%d, r=%f\n", n, m, rand_val)
	for i := 0; i < m; i++ {
		state := worker.BestRandomState(t, a, n, idx, fn, rand_val)
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

func (worker *Worker) BestRandomState(t ShapeType, a, n, idx int, fn NewShapeFunc, rand_val float64) *State {
	var bestEnergy float64
	var bestState *State

	for i := 0; i < n; i++ {
		state := worker.RandomState(t, a, idx, fn, rand_val)
		energy := state.Energy()
		// vv("BestRandomState: i=%d energy=%.2f, bestEnergy=%.2f\n", i, energy, bestEnergy)
		if i == 0 || energy < bestEnergy {
			bestEnergy = energy
			bestState = state
		}
	}
	return bestState
}

func NewBlueDotSessionsShapeFactory (modes []int, percs []float64) NewShapeFunc {
	var cum_percs []float64
	cur_cum_val := 0.0
	for idx, _ := range percs {
		cur_cum_val += percs[idx]
		cum_percs = append(cum_percs, cur_cum_val)
	}
	// fmt.Println(cum_percs)
	return func (worker *Worker, a, idx int, rand_val float64) Shape {
		// vv("NewBlueDotSessionsShapeFactory")
		// rnd := worker.Rnd
		// rand_val := rnd.Float64()
		for idy, val := range cum_percs {
			if rand_val <= val {
				return worker.SimpleRandomShape(ShapeType(modes[idy]))
			}
		}
		return worker.SimpleRandomShape(ShapeType(modes[len(modes) - 1]))
	}
}

type NewShapeFunc func(worker* Worker, a, idx int, rand_val float64) Shape

func (worker *Worker) SimpleRandomShape(t ShapeType) Shape {
	switch t {
	default:
		return worker.SimpleRandomShape(ShapeType(worker.Rnd.Intn(8)+1))
	case ShapeTypeTriangle:
		return NewRandomTriangle(worker)
	case ShapeTypeRectangle:
		return NewRandomRectangle(worker)
	case ShapeTypeEllipse:
		return NewRandomEllipse(worker)
	case ShapeTypeCircle:
		return NewRandomCircle(worker)
	case ShapeTypeRotatedRectangle:
		return NewRandomRotatedRectangle(worker)
	case ShapeTypeQuadratic:
		return NewRandomQuadratic(worker)
	case ShapeTypeRotatedEllipse:
		return NewRandomRotatedEllipse(worker)
	case ShapeTypePolygon:
		return NewRandomPolygon(worker, 4, true, 15, 40, 0)
	case ShapeTypeRightFacingTriangle:
		return NewRandomRFTriangle(worker)
	case ShapeTypeDiamond:
		return NewRandomDiamond(worker, 4, true, 15, 20, 0)
	}
}


func (worker *Worker) RandomState(t ShapeType, a, idx int, fn NewShapeFunc, rand_val float64) *State {
	vv("RandomState: idx=%d\n", idx)
	if t == ShapeTypeBlueDotSessions {
		return NewState(worker, fn(worker, a, idx, rand_val), a)
	} else {
		return NewState(worker, worker.SimpleRandomShape(t), a)
	}
}
