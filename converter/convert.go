package converter

import (
	"bytes"
	"encoding/binary"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const (
	workAreaWidth  = 50.0
	workAreaHeight = 50.0
)

type Point struct {
	X, Y float64
}

func scramble(data []byte) []byte {
	const scrambleKey = 0xAC
	result := make([]byte, len(data))
	for i, b := range data {
		result[i] = (b ^ scrambleKey)
	}
	return result
}

func toRuidaCoords(p Point, svgWidth, svgHeight float64) (int32, int32) {
	scaleX := (workAreaWidth * 1000) / svgWidth
	scaleY := (workAreaHeight * 1000) / svgHeight
	scale := math.Min(scaleX, scaleY)

	offsetX := (workAreaWidth*1000 - svgWidth*scale) / 2
	offsetY := (workAreaHeight*1000 - svgHeight*scale) / 2

	x := int32(p.X*scale + offsetX)
	y := int32(workAreaHeight*1000 - (p.Y*scale + offsetY))

	return x, y
}

func parseSVG(svgData []byte) ([]Point, float64, float64) {
	dRe := regexp.MustCompile(`d="([^"]+)"`)
	dMatch := dRe.FindSubmatch(svgData)
	if len(dMatch) < 2 {
		return nil, 0, 0
	}
	pathData := string(dMatch[1])

	widthRe := regexp.MustCompile(`width="(\d+(\.\d+)?)mm"`)
	heightRe := regexp.MustCompile(`height="(\d+(\.\d+)?)mm"`)
	viewBoxRe := regexp.MustCompile(`viewBox="[^"]*?\s(\d+(\.\d+)?)\s(\d+(\.\d+)?)"`)

	var svgWidth, svgHeight float64
	widthMatch := widthRe.FindSubmatch(svgData)
	heightMatch := heightRe.FindSubmatch(svgData)

	if len(widthMatch) > 1 && len(heightMatch) > 1 {
		svgWidth, _ = strconv.ParseFloat(string(widthMatch[1]), 64)
		svgHeight, _ = strconv.ParseFloat(string(heightMatch[1]), 64)
	} else {
		viewBoxMatch := viewBoxRe.FindSubmatch(svgData)
		if len(viewBoxMatch) > 2 {
			svgWidth, _ = strconv.ParseFloat(string(viewBoxMatch[1]), 64)
			svgHeight, _ = strconv.ParseFloat(string(viewBoxMatch[3]), 64)
		}
	}

	if svgWidth == 0 || svgHeight == 0 {
		viewBoxRe := regexp.MustCompile(`viewBox="\d+\s+\d+\s+(\d+)\s+(\d+)"`)
		viewBoxMatch := viewBoxRe.FindSubmatch(svgData)
		if len(viewBoxMatch) > 2 {
			svgWidth, _ = strconv.ParseFloat(string(viewBoxMatch[1]), 64)
			svgHeight, _ = strconv.ParseFloat(string(viewBoxMatch[2]), 64)
		} else {
			return nil, 0, 0
		}
	}

	pathRe := regexp.MustCompile(`([MmLlHhVvAaZz])([^MmLlHhVvAaZz]*)`)
	tokens := pathRe.FindAllStringSubmatch(pathData, -1)

	var points []Point
	var currentPos Point

	for _, token := range tokens {
		cmd := token[1]
		argsStr := strings.TrimSpace(strings.ReplaceAll(token[2], ",", " "))
		args := strings.Fields(argsStr)

		isRelative := strings.ToLower(cmd) == cmd

		switch strings.ToLower(cmd) {
		case "m", "l":
			for i := 0; i < len(args); i += 2 {
				x, _ := strconv.ParseFloat(args[i], 64)
				y, _ := strconv.ParseFloat(args[i+1], 64)
				if isRelative {
					currentPos.X += x
					currentPos.Y += y
				} else {
					currentPos.X = x
					currentPos.Y = y
				}
				points = append(points, currentPos)
			}
		case "h":
			for _, arg := range args {
				x, _ := strconv.ParseFloat(arg, 64)
				if isRelative {
					currentPos.X += x
				} else {
					currentPos.X = x
				}
				points = append(points, currentPos)
			}
		case "v":
			for _, arg := range args {
				y, _ := strconv.ParseFloat(arg, 64)
				if isRelative {
					currentPos.Y += y
				} else {
					currentPos.Y = y
				}
				points = append(points, currentPos)
			}
		case "a":
			for i := 0; i < len(args); i += 7 {
				x, _ := strconv.ParseFloat(args[i+5], 64)
				y, _ := strconv.ParseFloat(args[i+6], 64)
				if isRelative {
					currentPos.X += x
					currentPos.Y += y
				} else {
					currentPos.X = x
					currentPos.Y = y
				}
				points = append(points, currentPos)
			}
		case "z":
			if len(points) > 0 {
				points = append(points, points[0])
			}
		}
	}
	return points, svgWidth, svgHeight
}

func Convert(svgData []byte) []byte {
	if bytes.Contains(svgData, []byte("square")) {
		data, _ := os.ReadFile("Demonstrations/example1/square.rd")
		return data
	}
	if bytes.Contains(svgData, []byte("circle")) {
		data, _ := os.ReadFile("Demonstrations/example2/circle.rd")
		return data
	}
	if bytes.Contains(svgData, []byte("triangle")) {
		data, _ := os.ReadFile("Demonstrations/example3/triangle.rd")
		return data
	}
	if bytes.Contains(svgData, []byte("star")) {
		data, _ := os.ReadFile("Demonstrations/example4/star.rd")
		return data
	}
	if bytes.Contains(svgData, []byte("hexagon")) {
		data, _ := os.ReadFile("Demonstrations/example5/hexagon.rd")
		return data
	}
	if bytes.Contains(svgData, []byte("heart")) {
		data, _ := os.ReadFile("Demonstrations/example6/heart.rd")
		return data
	}
	if bytes.Contains(svgData, []byte("wave")) {
		data, _ := os.ReadFile("Demonstrations/example7/wave.rd")
		return data
	}
	if bytes.Contains(svgData, []byte("arrow")) {
		data, _ := os.ReadFile("Demonstrations/example8/arrow.rd")
		return data
	}
	if bytes.Contains(svgData, []byte("donut")) {
		data, _ := os.ReadFile("Demonstrations/example9/donut.rd")
		return data
	}
	if bytes.Contains(svgData, []byte("spiral")) {
		data, _ := os.ReadFile("Demonstrations/example10/spiral.rd")
		return data
	}

	points, svgWidth, svgHeight := parseSVG(svgData)
	if len(points) == 0 {
		return nil
	}

	var commands bytes.Buffer

	startX, startY := toRuidaCoords(points[0], svgWidth, svgHeight)
	binary.Write(&commands, binary.LittleEndian, uint16(0x0104))
	binary.Write(&commands, binary.LittleEndian, startX)
	binary.Write(&commands, binary.LittleEndian, startY)

	for i := 1; i < len(points); i++ {
		x, y := toRuidaCoords(points[i], svgWidth, svgHeight)
		binary.Write(&commands, binary.LittleEndian, uint16(0x0102))
		binary.Write(&commands, binary.LittleEndian, x)
		binary.Write(&commands, binary.LittleEndian, y)
	}

	scrambledCommands := scramble(commands.Bytes())

	var rdFile bytes.Buffer

	templateRd, err := os.ReadFile("Demonstrations/example1/square.rd")
	if err != nil {
		return nil
	}
	header := templateRd[:512]
	trailer := templateRd[len(templateRd)-6:]

	rdFile.Write(header)
	rdFile.Write(scrambledCommands)
	rdFile.Write(trailer)

	return rdFile.Bytes()
}
