# Debugging State Transitions

When debugging UI issues, **PNGs only capture a single moment** and miss bugs that occur during state changes. For example, a sidebar might render correctly initially but break after data loads - a temporal bug invisible in static screenshots.

## The GIF Limitation

GIF animations capture state transitions, but Claude's image viewer only displays the first frame. To see what's happening across the animation, extract frames with ffmpeg.

## Extract Frames with ffmpeg

```bash
# Extract specific frames (0, 10, 20, 30, 40)
ffmpeg -i testdata/vhs/demo.gif \
  -vf "select=eq(n\,0)+eq(n\,10)+eq(n\,20)+eq(n\,30)+eq(n\,40)" \
  -vsync vfr /tmp/frame%d.png

# Extract every 10th frame
ffmpeg -i testdata/vhs/demo.gif \
  -vf "select=not(mod(n\,10))" \
  -vsync vfr /tmp/frame%d.png
```

## Debugging Workflow

1. **Record with VHS**: `vhs testdata/vhs/startup.tape`
2. **Extract key frames**: `ffmpeg -i startup.gif -vf "select=eq(n\,30)+eq(n\,40)" -vsync vfr /tmp/frame%d.png`
3. **Read frames**: Compare before/after states with `Read /tmp/frame1.png` then `Read /tmp/frame2.png`
4. **Identify transition**: Find which state change introduces the bug

## Real Example

This technique exposed a bug where the sidebar rendered correctly at frame 30 (before conversations loaded) but was broken at frame 40 (after conversations loaded). The static screenshot showed only the broken state, making the temporal nature of the bug invisible.
