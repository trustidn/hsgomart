# POS sounds (optional)

Place custom audio files here for POS feedback:

- **beep.mp3** – played when a product is added to cart (success). Short POS-style beep.
- **beep-error.mp3** – played when a barcode is not found (error). Different tone/length.

If these files are not present, the app falls back to Web Audio API beeps.

Serve path: `/beep.mp3` and `/beep-error.mp3` (files must be in `public/` so they are served at root).
