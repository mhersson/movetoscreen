package main

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

type monitor struct {
	xres int
	yres int
	xmin int
	xmax int
	ymin int
	ymax int
}

func newMonitor(res, pos string) monitor {
	xyres := strings.Split(res, "x")
	xypos := strings.Split(pos, "+")
	xres, _ := strconv.Atoi(xyres[0])
	yres, _ := strconv.Atoi(xyres[1])
	xpos, _ := strconv.Atoi(xypos[0])
	ypos, _ := strconv.Atoi(xypos[1])
	xmin := xpos
	xmax := xpos + xres
	ymin := ypos
	ymax := ypos + yres
	return monitor{xres, yres, xmin, xmax, ymin, ymax}
}

func (m monitor) IsInsideRange(x, y int) bool {
	if m.xmin < x && x < m.xmax && m.ymin < y && y < m.ymax {
		return true
	}
	return false
}

func (m monitor) GetResolution() (int, int) {
	return m.xres, m.yres
}

func getMonitors() []monitor {
	var monitors []monitor
	output, _ := exec.Command("xrandr").CombinedOutput()

	re := regexp.MustCompile(`\sconnected.+?(([0-9]+x[0-9]+)\+?([0-9]+\+[0-9]+))`)
	match := re.FindAllStringSubmatch(string(output), -1)
	for i := range match {
		monitors = append(monitors, newMonitor(match[i][2], match[i][3]))
	}

	return monitors
}

func getActiveWindowPosition() (int, int) {
	args := []string{"getactivewindow", "getwindowgeometry"}
	output, _ := exec.Command("xdotool", args...).CombinedOutput()
	re := regexp.MustCompile(`Position:\s([0-9]+,[0-9]+)`)
	match := re.FindStringSubmatch(string(output))
	xypos := strings.Split(match[1], ",")
	xpos, _ := strconv.Atoi(xypos[0])
	ypos, _ := strconv.Atoi(xypos[1])
	return xpos, ypos
}

func getActiveMonitorRes(monitors []monitor, xpos, ypos int) (int, int) {
	for _, monitor := range monitors {
		if monitor.IsInsideRange(xpos, ypos) {
			return monitor.GetResolution()
		}
	}
	return 0, 0
}

func getNewPos(dir string, xres, yres, xpos, ypos int) (int, int) {
	switch dir {
	case "left":
		return xpos - xres, ypos
	case "down":
		return xpos, ypos + yres
	case "up":
		return xpos, ypos - yres
	case "right":
		return xpos + xres, ypos
	}
	return 0, 0
}

func isInViewport(monitors []monitor, xpos, ypos int) bool {
	for _, monitor := range monitors {
		if monitor.IsInsideRange(xpos, ypos) {
			return true
		}
	}
	return false
}

func moveActiveWindow(xpos, ypos int) {
	args := []string{"getactivewindow", "windowraise", "windowmove",
		strconv.Itoa(xpos), strconv.Itoa(ypos)}
	err := exec.Command("xdotool", args...).Run()
	if err != nil {
		fmt.Println(err)
	}
}

func getMousePosition() (int, int) {
	args := []string{"getmouselocation"}
	output, _ := exec.Command("xdotool", args...).CombinedOutput()
	strout := strings.Split(string(output), " ")
	xpos, _ := strconv.Atoi(strout[0][2:])
	ypos, _ := strconv.Atoi(strout[1][2:])
	return xpos, ypos
}

func moveMousePointer(xpos, ypos int) {
	args := []string{"mousemove", strconv.Itoa(xpos), strconv.Itoa(ypos)}
	err := exec.Command("xdotool", args...).Run()
	if err != nil {
		fmt.Println(err)
	}
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func run(direction string, moveWindow bool) {
	monitors := getMonitors()
	if moveWindow {
		xpos, ypos := getActiveWindowPosition()
		xres, yres := getActiveMonitorRes(monitors, xpos, ypos)
		newxpos, newypos := getNewPos(direction, xres, yres, xpos, ypos)
		if isInViewport(monitors, newxpos, newypos) {
			moveActiveWindow(newxpos, newypos)
		}
	}
	xpos, ypos := getMousePosition()
	xres, yres := getActiveMonitorRes(monitors, xpos, ypos)
	newxpos, newypos := getNewPos(direction, xres, yres, xpos, ypos)
	if isInViewport(monitors, newxpos, newypos) {
		moveMousePointer(newxpos, newypos)
	}

}

func main() {
	var directions = []string{"left", "down", "up", "right"}

	if len(os.Args) == 2 && stringInSlice(os.Args[1], directions) {

		run(os.Args[1], false)

	} else if len(os.Args) == 3 &&
		stringInSlice(os.Args[1], directions) &&
		os.Args[2] == "window" {

		run(os.Args[1], true)
	}
}
