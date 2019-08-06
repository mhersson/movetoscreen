# Movetoscreen
Switch focus (move mouse pointer), or move the active window, to an adjacent monitor. It is intendend for window managers lacking this shortcut. I use it with cwm

Repo contains versions for both Go and Python. I first wrote the Python script, then implemented the exact same functionality in Go. They are identical, but of course the Go version is faster.

## Usage

`movetoscreen <left|down|up|right> [window]`

If the window option is missing, it only moves the mouse pointer

## Dependencies
 - `xrandr`
 - `xdotool`

## From my ~/.cwmrc

```
bind-key CM-Left "movetoscreen left"
bind-key CM-Down "movetoscreen down"
bind-key CM-Up "movetoscreen up"
bind-key CM-Right "movetoscreen right"

bind-key CMS-Left "movetoscreen left window"
bind-key CMS-Down "movetoscreen down window"
bind-key CMS-Up "movetoscreen up window"
bind-key CMS-Right "movetoscreen right window"
```

## Author and licence

Morten Hersson - 2019 - Public domain

Inspired by Antoine Calando's [movescreen](https://github.com/calandoa/movescreen).
