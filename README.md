# typocat

```
 /\_/\
( ^.^ )
 / >☕
/|_|_|\
```

<video src="https://raw.githubusercontent.com/davitostes/typocat/master/video.mp4" controls width="600" poster="https://raw.githubusercontent.com/davitostes/typocat/master/thumbnail.png"></video>

Tiny terminal typing trainer. Cat watches. Combo grows. Mistype and the kitty wilts.

## Run

```sh
go run .
```

Or build it.

```sh
go build -o typocat
./typocat
```

## Combo tiers

| Streak | Cat says |
|-------:|----------|
|     1+ | Yay! |
|    10+ | Nice. |
|    20+ | Amazing! |
|    50+ | Incredible!! |
|   100+ | SPLENDID!!! |

## Layout

Single file, `main.go`. Frames, palette, and text live at the top. Swap them, recompile, vibe shifts.
