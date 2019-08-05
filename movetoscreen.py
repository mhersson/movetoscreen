#!/usr/bin/env python
'''
Move active window and/or mouse to adjacent monitor,
and make sure it's not outside of the viewport before moving

Dependencies: xrandr, xdotool

Usage:
movestocreen.py <left|down|up|right> [window]

Written by Morten Hersson / 2019 - Public domain
Inspired by Antoine Calando's https://github.com/calandoa/movescreen
'''
import re
import sys
import subprocess


class Monitor():
    def __init__(self, res, pos):
        self.xres = int(res[0])
        self.yres = int(res[1])
        self.xpos = int(pos[0])
        self.ypos = int(pos[1])
        self.xmin = self.xpos
        self.xmax = self.xpos + self.xres
        self.ymin = self.ypos
        self.ymax = self.ypos + self.yres

    def is_inside_range(self, x, y):
        if self.xmin < x < self.xmax and self.ymin < y < self.ymax:
            return True
        return False

    def get_resolution(self):
        return self.xres, self.yres


def get_monitors():
    monitors = []
    out = subprocess.check_output(['xrandr'])
    pattern = r"\sconnected.+?(([0-9]+x[0-9]+)\+?([0-9]+\+[0-9]+))"
    for line in out.decode().splitlines():
        match = re.search(pattern, line)
        if match:
            monitors.append(Monitor(match.group(2).split("x"),
                                    match.group(3).split("+")))
    return monitors


def get_active_window_pos():
    out = subprocess.check_output(
        ['xdotool', 'getactivewindow', 'getwindowgeometry'])

    match = re.search(r"Position:\s([0-9]+,[0-9]+)", out.decode())
    if match:
        xpos, ypos = match.group(1).split(",")
        return int(xpos), int(ypos)
    return None


def get_current_resolution(monitors, xpos, ypos):
    # Get resolution of current monitor
    for s in monitors:
        if s.is_inside_range(xpos, ypos):
            return s.xres, s.yres
    return None


def get_new_position(direction, xres, yres, xpos, ypos):
    # Direction: direction to move
    # xres and yres and the resolution of the current monitor
    # xpos and ypos are the current position of the active window/mouse pointer
    if direction == "left":
        return xpos - xres, ypos
    if direction == "down":
        return xpos, ypos + yres
    if direction == "up":
        return xpos, ypos - yres
    if direction == "right":
        return xpos + xres, ypos
    return None


def is_in_viewport(monitors, xpos, ypos):
    for m in monitors:
        if m.is_inside_range(xpos, ypos):
            return True
    return False


def move_window(newxpos, newypos):
    subprocess.call(['xdotool', 'getactivewindow', 'windowmove',
                     str(newxpos), str(newypos), 'windowraise'])


def get_mouse_position():
    out = subprocess.check_output(['xdotool', 'getmouselocation'])
    pos = dict([w.split(':') for w in out.decode().split()])
    return int(pos['x']), int(pos['y'])


def move_mouse(newxpos, newypos):
    subprocess.call(['xdotool', 'mousemove', str(newxpos), str(newypos)])


def run(direction, focus_only=True):
    monitors = get_monitors()

    if not focus_only:
        xpos, ypos = get_active_window_pos()

        xres, yres = get_current_resolution(monitors, xpos, ypos)

        newxpos, newypos = get_new_position(direction, xres, yres, xpos, ypos)

        if is_in_viewport(monitors, newxpos, newypos):
            move_window(newxpos, newypos)

    xpos, ypos = get_mouse_position()

    xres, yres = get_current_resolution(monitors, xpos, ypos)

    newxpos, newypos = get_new_position(direction, xres, yres, xpos, ypos)

    if is_in_viewport(monitors, newxpos, newypos):
        move_mouse(newxpos, newypos)


if __name__ == "__main__":
    if len(sys.argv) == 2 and sys.argv[1] in ['left', 'down', 'up', 'right']:
        run(sys.argv[1])
    elif (len(sys.argv) == 3 and
          sys.argv[1] in ['left', 'down', 'up', 'right'] and
          sys.argv[2] == "window"):
        run(sys.argv[1], focus_only=False)
