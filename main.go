package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"./primitive"
	"github.com/nfnt/resize"
	// "github.com/pkg/profile"
)

var (
	Input      string
	Outputs    flagArray
	Background string
	Configs    shapeConfigArray
	Alpha      int
	BlackThresh float64
	AreaThresh string
	Age int
	ShapeTrials int
	HillClimbTrials int
	InputSize  int
	OutputSize int
	Mode       string
	Workers    int
	Nth        int
	Repeat     int
	Seed  int64
	V, VV      bool
)

type flagArray []string

func (i *flagArray) String() string {
	return strings.Join(*i, ", ")
}

func (i *flagArray) Set(value string) error {
	*i = append(*i, value)
	return nil
}

type shapeConfig struct {
	Count  int
	Mode   string
	Alpha  int
	Repeat int
}

type shapeConfigArray []shapeConfig

func (i *shapeConfigArray) String() string {
	return ""
}

func (i *shapeConfigArray) Set(value string) error {
	n, _ := strconv.ParseInt(value, 0, 0)
	*i = append(*i, shapeConfig{int(n), Mode, Alpha, Repeat})
	return nil
}

func init() {
	flag.StringVar(&Input, "i", "", "input image path")
	flag.Var(&Outputs, "o", "output image path")
	flag.Var(&Configs, "n", "number of primitives")
	flag.StringVar(&Background, "bg", "", "background color (hex)")
	flag.IntVar(&Alpha, "a", 128, "alpha value")
	flag.Float64Var(&BlackThresh, "kt", 0.0, "black cut off threshold")
	flag.StringVar(&AreaThresh, "at", "0.0", "area cut off threshold. Can specify a single value for upper threshold, or comma separated values for both lower and upper thresholds")
	flag.IntVar(&InputSize, "r", 256, "resize large input images to this size")
	flag.IntVar(&OutputSize, "s", 1024, "output image size")
	flag.StringVar(&Mode, "m", "1", "0=combo 1=triangle 2=rect 3=ellipse 4=circle 5=rotatedrect 6=beziers 7=rotatedellipse 8=polygon 9=blue-dot-sessions 10=right-facing-triangle 11=blue-dot-sessions-polygon")
	flag.IntVar(&Workers, "j", 0, "number of parallel workers (default uses all cores)")
	flag.IntVar(&Nth, "nth", 1, "save every Nth frame (put \"%d\" in path)")
	flag.IntVar(&ShapeTrials, "st", 1000, "Number of shapes to generate before applying Hill Climb algorithm")
	flag.IntVar(&HillClimbTrials, "hct", 16, "Number of times to use Hill Climb algorithm per shape")
	flag.IntVar(&Age, "age", 100, "age parameter for Hill Climb Algorithm")
	flag.IntVar(&Repeat, "rep", 0, "add N extra shapes per iteration with reduced search")
	flag.Int64Var(&Seed, "seed", -1, "Random number seed")
	flag.BoolVar(&V, "v", false, "verbose")
	flag.BoolVar(&VV, "vv", false, "very verbose")
}

func errorMessage(message string) bool {
	fmt.Fprintln(os.Stderr, message)
	return false
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func parseBlueDotSessionsModeParams (modeStr string) (int, int, int) {
	split := strings.Split(modeStr, ":")
	mode, err := strconv.Atoi(split[0])
	check(err)
	var startRect int = 5
	var endRect int  = 15
	if len(split) == 3 {
		startRect, err = strconv.Atoi(split[1])
		check(err)
		endRect, err = strconv.Atoi(split[2])
		check(err)
	}
	return mode, startRect, endRect
}

func parseAreaThresh (areaThresh string) (float64, float64) {
	split := strings.Split(areaThresh, ",")
	if len(split) == 1 {
		upperAreaThresh, err := strconv.ParseFloat(split[0], 64)
		check(err)
		return 0.0, upperAreaThresh
	} else if len(split) == 2 {
		lowerAreaThresh, err := strconv.ParseFloat(split[0], 64)
		check(err)
		upperAreaThresh, err := strconv.ParseFloat(split[1], 64)
		check(err)
		return lowerAreaThresh, upperAreaThresh
	}
	return 0.0, 0.0
}


func main() {
	// defer profile.Start().Stop()
	// parse and validate arguments
	flag.Parse()
	ok := true
	if Input == "" {
		ok = errorMessage("ERROR: input argument required")
	}
	if len(Outputs) == 0 {
		ok = errorMessage("ERROR: output argument required")
	}
	if len(Configs) == 0 {
		ok = errorMessage("ERROR: number argument required")
	}
	if len(Configs) == 1 {
		Configs[0].Mode = Mode
		Configs[0].Alpha = Alpha
		Configs[0].Repeat = Repeat
	}
	for _, config := range Configs {
		if config.Count < 1 {
			ok = errorMessage("ERROR: number argument must be > 0")
		}
	}
	if !ok {
		fmt.Println("Usage: primitive [OPTIONS] -i input -o output -n count")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// set log level
	if V {
		primitive.LogLevel = 1
	}
	if VV {
		primitive.LogLevel = 2
	}

	// seed random number generator
	if Seed == -1 {
		rand.Seed(time.Now().UTC().UnixNano())
	} else {
		rand.Seed(Seed)
	}

	// determine worker count
	if Workers < 1 {
		Workers = runtime.NumCPU()
	}

	// read input image
	primitive.Log(1, "reading %s\n", Input)
	input, err := primitive.LoadImage(Input)
	check(err)

	// scale down input image if needed
	size := uint(InputSize)
	if size > 0 {
		input = resize.Thumbnail(size, size, input, resize.Bilinear)
	}

	// determine background color
	var bg primitive.Color
	if Background == "" {
		bg = primitive.MakeColor(primitive.AverageImageColor(input))
	} else {
		bg = primitive.MakeHexColor(Background)
	}

	lowerAreaThresh, upperAreaThresh := parseAreaThresh(AreaThresh)
	// run algorithm
	// primitive.Log(1, "Background=%s, bg=%s\n", Background, bg)
	model := primitive.NewModel(input, bg, OutputSize, Workers, BlackThresh, lowerAreaThresh, upperAreaThresh, Seed)
	primitive.Log(1, "%d: t=%.3f, score=%.6f\n", 0, 0.0, model.Score)
	start := time.Now()
	frame := 0
	var mode int
	var startRect int = 5
	var endRect int = 15

	for j, config := range Configs {
		primitive.Log(1, "count=%d, mode=%s, alpha=%d, repeat=%d\n",
			config.Count, config.Mode, config.Alpha, config.Repeat)

		if (strings.IndexAny(config.Mode, ":") != -1) {
			mode, startRect, endRect = parseBlueDotSessionsModeParams(config.Mode)
		} else {
			mode, err = strconv.Atoi(config.Mode)
			check(err)
		}
		newShapeFunc := primitive.NewBlueDotSessionsShapeFactory(startRect, endRect)
		primitive.Log(1, "parsed mode=%d\n",  mode)


		for i := 0; i < config.Count; i++ {
			frame++
			// find optimal shape and add it to the model
			t := time.Now()
			n := model.Step(primitive.ShapeType(mode), config.Alpha, config.Repeat, i, ShapeTrials, Age, HillClimbTrials, newShapeFunc)
			nps := primitive.NumberString(float64(n) / time.Since(t).Seconds())
			elapsed := time.Since(start).Seconds()
			primitive.Log(1, "%d: t=%.3f, score=%.6f, n=%d, n/s=%s\n", frame, elapsed, model.Score, n, nps)

			// write output image(s)
			for _, output := range Outputs {
				ext := strings.ToLower(filepath.Ext(output))
				if output == "-" {
					ext = ".svg"
				}
				percent := strings.Contains(output, "%")
				saveFrames := percent && ext != ".gif"
				saveFrames = saveFrames && frame%Nth == 0
				last := j == len(Configs)-1 && i == config.Count-1
				if saveFrames || last {
					path := output
					if percent {
						path = fmt.Sprintf(output, frame)
					}
					primitive.Log(1, "writing %s\n", path)
					switch ext {
					default:
						check(fmt.Errorf("unrecognized file extension: %s", ext))
					case ".png":
						check(primitive.SavePNG(path, model.Context.Image()))
					case ".jpg", ".jpeg":
						check(primitive.SaveJPG(path, model.Context.Image(), 95))
					case ".svg":
						check(primitive.SaveFile(path, model.SVG()))
					case ".gif":
						frames := model.Frames(0.001)
						check(primitive.SaveGIFImageMagick(path, frames, 50, 250))
					}
				}
			}
		}
	}
}
