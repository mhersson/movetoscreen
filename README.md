# Movetoscreen
Switch focus (move mouse pointer), or move the active window, to an adjacent monitor.
It is intendend for window managers lacking this shortcut. I use it with cwm

## Usage

`movestocreen.py <left|down|up|right> [window]`

If the window option is missing, it only moves the mouse pointer

## Dependencies
 - `xrandr`
 - `xdotool`

## Example from my ~/.cwmrc (make sure movetoscreen.py is in your path)

```
bind-key CM-Left "movetoscreen.py left"
bind-key CM-Down "movetoscreen.py down"
bind-key CM-Up "movetoscreen.py up"
bind-key CM-Right "movetoscreen.py right"

bind-key CMS-Left "movetoscreen.py left window"
bind-key CMS-Down "movetoscreen.py down window"
bind-key CMS-Up "movetoscreen.py up window"
bind-key CMS-Right "movetoscreen.py right window"
```

## Author and licence

Morten Hersson - 2019 - Public domain

Inspired by Antoine Calando's [movescreen](https://github.com/calandoa/movescreen).
