#!/usr/bin/env python3
"""Synthetischer Rettungshubschrauber (EC135/H145-Charakter) für die Landeplatz-Übung.

Erzeugt zwei Dateien:
  1_anflug_und_schwebeflug.mp3     10:30 min. Leiser Anflug aus der Ferne (45 s Crescendo),
                                   danach konstanter Schwebeflug bis zum Ende.
  2_landung_und_abschaltung.mp3    1:25 min. Startet im Schwebeflug (8 s) (gleiche Lautstärke wie
                                   Track 1), Landung/Bodennähe, Turbine auf Leerlauf,
                                   Abschaltung, Rotor läuft aus, Stille.

Der Schwebeflug in Track 1 ist eine 60-s-Schleife, deren periodische Anteile exakt auf
60 s ausgelegt sind (kein Klick an der Naht).

Benötigt: numpy, imageio-ffmpeg (für MP3). Aufruf: python make_heli.py
"""
import subprocess
import tempfile
import wave
from pathlib import Path

import numpy as np

SR = 44100
LOOP_S = 60
FR0 = 6.5            # Rotordrehzahl in Hz (390 U/min), x4 Blätter = 26 Hz Blattfolge
OUT = Path(__file__).parent
rng = np.random.default_rng(135)
HOVER_PEAK = 0.70    # Spitzenpegel des Schwebeflugs (Bezug für beide Tracks)
GROUND_GAIN = 1.25
W_LOW, W_SLAP, W_HIGH = 1.6, 0.85, 0.10   # Mischung: dumpfe Rotorschläge vorn, Turbine/Heulen nur dezent   # Bodennähe/Downwash ist lauter als Schwebeflug in Entfernung

THETA = rng.uniform(0, 0.4, 260)   # kleine Blatt-zu-Blatt-Phasenunterschiede


def fftfilt(x, lo=None, hi=None, order=4):
    spec = np.fft.rfft(x)
    f = np.fft.rfftfreq(len(x), 1 / SR)
    h = np.ones_like(f)
    if lo:
        h *= 1 / np.sqrt(1 + (lo / np.maximum(f, 1e-9)) ** (2 * order))
    if hi:
        h *= 1 / np.sqrt(1 + (f / hi) ** (2 * order))
    return np.fft.irfft(spec * h, len(x))


def noise(n, lo, hi):
    x = fftfilt(rng.standard_normal(n), lo, hi)
    return x / (x.std() * 3.7)


_rotor_peak = None


def render(n, fr, sw):
    """Liefert pro Kanal (low, slap, high). fr: Rotor-Hz, sw: Turbinendrehzahl relativ (1 = Flug).
    fr/sw sind Skalare oder Arrays der Länge n. Bei konstanten Werten, deren Produkt mit der
    Dauer ganzzahlig ist, ist das Ergebnis exakt periodisch."""
    global _rotor_peak
    t = np.arange(n) / SR
    fr = np.broadcast_to(np.asarray(fr, float), (n,))
    sw = np.broadcast_to(np.asarray(sw, float), (n,))
    phi = np.cumsum(2 * np.pi * fr / SR)
    rotor = np.zeros(n)
    for k in range(1, 260):
        a = k ** -0.85 * np.exp(-k * fr / 600)
        if k % 4:
            a = a * 0.35
            th = THETA[k]
        else:
            th = 0.0
        rotor += a * np.cos((k * phi + th) % (2 * np.pi))
    if _rotor_peak is None:
        _rotor_peak = np.max(np.abs(rotor))
    rotor /= _rotor_peak
    env = np.clip(rotor, 0, None) ** 2
    amp = (fr / FR0) ** 2                                   # Schlaglautstärke ~ Drehzahl²
    wob = 1 + 0.12 * np.sin(2 * np.pi * (12 / LOOP_S) * t + 1.0) + 0.06 * np.sin(2 * np.pi * (34 / LOOP_S) * t)
    low = fftfilt(rotor, None, 220) * wob * amp
    # Turbinen-/Fenestron-Heulen mit Drehzahl-abhängiger Frequenz
    fw = 2790 * sw * (1 + 0.002 * np.sin(2 * np.pi * (3 / LOOP_S) * t))
    pw = np.cumsum(2 * np.pi * fw / SR)
    whine = (0.10 * np.sin(pw % (2 * np.pi)) + 0.30 * np.sin((pw / 2) % (2 * np.pi))) \
        * (0.7 + 0.3 * env) * sw ** 2
    chans = []
    for _ in range(2):
        slap = noise(n, 120, 1100) * (0.25 + env) * wob * amp
        turb = noise(n, 1500, 4500) * 0.35 * sw ** 3
        chans.append((low.astype(np.float32), slap.astype(np.float32),
                      (turb * 0.6 + whine * 0.35).astype(np.float32)))
    return chans


def mix(chans, g_low, g_mid, g_high):
    return np.stack([l * W_LOW * g_low + s * W_SLAP * g_mid + h * W_HIGH * g_high for l, s, h in chans], 1)


def smooth(t, t0, t1, v0, v1):
    """Weicher Übergang (Cosinus) von v0 nach v1 zwischen t0 und t1."""
    x = np.clip((t - t0) / (t1 - t0), 0, 1)
    return v0 + (v1 - v0) * (1 - np.cos(np.pi * x)) / 2


def write(name, y, fade_in=0.0, fade_out=0.0):
    y = np.tanh(y * SCALE / 0.95) * 0.95           # sanfter Limiter
    if fade_in:
        y[: int(fade_in * SR)] *= np.linspace(0, 1, int(fade_in * SR))[:, None]
    if fade_out:
        y[-int(fade_out * SR):] *= np.linspace(1, 0, int(fade_out * SR))[:, None]
    import imageio_ffmpeg
    with tempfile.TemporaryDirectory() as d:
        wav = Path(d) / "x.wav"
        with wave.open(str(wav), "wb") as w:
            w.setnchannels(2); w.setsampwidth(2); w.setframerate(SR)
            w.writeframes((y * 32767).astype("<i2").tobytes())
        subprocess.run([imageio_ffmpeg.get_ffmpeg_exe(), "-y", "-loglevel", "error", "-i", str(wav),
                        "-codec:a", "libmp3lame", "-q:a", "4", str(OUT / f"{name}.mp3")], check=True)
    print("geschrieben:", name)


# ---- Referenz: stationärer Schwebeflug (60 s, periodisch) -------------------------------
hover = render(LOOP_S * SR, FR0, 1.0)
SCALE = 1.0
SCALE = HOVER_PEAK / np.max(np.abs(mix(hover, 1, 1, 1)))

# ---- Track 1: Anflug + Schwebeflug, 10:30 min -------------------------------------------
TOTAL1, RAMP = 630, 45
n1 = TOTAL1 * SR
t1 = np.arange(n1) / SR
x = np.clip(t1 / RAMP, 0, 1)
# Entfernung: Tiefe tragen weit, Höhen werden stark gedämpft
g = [x ** 1.6, x ** 2.4, x ** 3.6]
reps = -(-n1 // len(hover[0][0]))
y = np.zeros((n1, 2), np.float32)
for c in range(2):
    comps = [np.tile(a, reps)[:n1] for a in hover[c]]
    y[:, c] = comps[0] * W_LOW * g[0] + comps[1] * W_SLAP * g[1] + comps[2] * W_HIGH * g[2]
write("1_anflug_und_schwebeflug", y, fade_in=1.0, fade_out=8.0)
del y

# ---- Track 2: Landung und Abschaltung, 1:25 min -----------------------------------------
TOTAL2 = 85
n2 = TOTAL2 * SR
t2 = np.arange(n2) / SR
#  0-8 s   Schwebeflug (wie Track 1)
#  8-20 s  Sinkflug, Bodennähe: lauter (Downwash, Abstand kleiner)
# 20-22 s  Aufsetzen
# 22-36 s  Turbine auf Bodenleerlauf, Rotor etwas langsamer
# 36-75 s  Abschaltung: Turbine läuft aus, Rotor bremst aus, Stille
fr = FR0 * np.ones(n2)
fr = smooth(t2, 22, 32, FR0, 5.2)
shut = t2 > 36
fr = np.where(shut, 5.2 * np.exp(-(t2 - 36) / 10), fr)
sw = smooth(t2, 22, 32, 1.0, 0.78)
sw = np.where(shut, 0.78 * np.exp(-(t2 - 36) / 6), sw)
G = smooth(t2, 8, 20, 1.0, GROUND_GAIN)
G = np.where(t2 > 22, smooth(t2, 22, 34, GROUND_GAIN, 1.0), G)
chans = render(n2, fr, sw)
y = mix(chans, G, G, G ** 1.3)
write("2_landung_und_abschaltung", y, fade_in=0.05, fade_out=3.0)
print("fertig")
