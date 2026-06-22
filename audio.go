package main

import (
	"bytes"
	"encoding/binary"
	"io"
	"log"
	"math"
	"path"
	"sync"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/jfreymuth/oggvorbis"
)

// AudioManager handles sound effects and background music.
type AudioManager struct {
	context *audio.Context
	sfx     map[string]*soundEntry
	muted   bool
	sfxVol  float64
	bgmVol  float64
	mu      sync.Mutex
}

type soundEntry struct {
	data []byte // decoded PCM (16-bit little-endian, 2 channels)
}

const audioSampleRate = 44100

// NewAudioManager creates a new audio manager and preloads all sound effects.
func NewAudioManager(sfxVolume, bgmVolume float64) *AudioManager {
	ctx := audio.NewContext(audioSampleRate)
	am := &AudioManager{
		context: ctx,
		sfx:     make(map[string]*soundEntry),
		muted:   false,
		sfxVol:  clampVolume(sfxVolume),
		bgmVol:  clampVolume(bgmVolume),
	}

	// Preload all bundled OGG files
	files := []string{
		"sfx_laser1.ogg",
		"sfx_laser2.ogg",
		"sfx_lose.ogg",
		"sfx_shieldDown.ogg",
		"sfx_shieldUp.ogg",
		"sfx_twoTone.ogg",
		"sfx_zap.ogg",
	}

	for _, file := range files {
		if err := am.loadSFX(file); err != nil {
			log.Printf("audio: failed to load %s: %v", file, err)
		}
	}

	return am
}

func clampVolume(v float64) float64 {
	if v < 0 {
		return 0
	}
	if v > 1 {
		return 1
	}
	return v
}

// loadSFX reads and decodes an OGG file from the embedded assets.
func (am *AudioManager) loadSFX(filename string) error {
	fullPath := path.Join("assets/Bonus", filename)
	file, err := data.Open(fullPath)
	if err != nil {
		return err
	}
	defer file.Close()

	raw, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	reader := bytes.NewReader(raw)
	decoder, err := oggvorbis.NewReader(reader)
	if err != nil {
		return err
	}

	sampleRate := decoder.SampleRate()
	channels := decoder.Channels()

	// Read all float32 samples
	var samples []float32
	buf := make([]float32, 4096)
	for {
		n, err := decoder.Read(buf)
		if n > 0 {
			samples = append(samples, buf[:n]...)
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
	}

	// If the original sample rate differs, resample (simple nearest-neighbor)
	if sampleRate != audioSampleRate {
		samples = simpleResampleF32(samples, channels, sampleRate, audioSampleRate)
	}

	// Convert float32[-1.0, 1.0] to int16 little-endian
	pcm := float32ToInt16LE(samples)

	am.sfx[filename] = &soundEntry{
		data: pcm,
	}
	return nil
}

// float32ToInt16LE converts float32 samples (range -1 to 1) to int16 little-endian bytes.
func float32ToInt16LE(samples []float32) []byte {
	out := make([]byte, len(samples)*2)
	for i, s := range samples {
		// Clamp
		if s > 1.0 {
			s = 1.0
		} else if s < -1.0 {
			s = -1.0
		}
		val := int16(s * math.MaxInt16)
		binary.LittleEndian.PutUint16(out[i*2:], uint16(val))
	}
	return out
}

// simpleResampleF32 does nearest-neighbor resampling.
func simpleResampleF32(samples []float32, channels, fromRate, toRate int) []float32 {
	ratio := float64(fromRate) / float64(toRate)
	sampleCount := len(samples) / channels
	newCount := int(float64(sampleCount) / ratio)
	out := make([]float32, newCount*channels)
	for i := 0; i < newCount; i++ {
		srcIdx := int(float64(i) * ratio)
		if srcIdx >= sampleCount {
			srcIdx = sampleCount - 1
		}
		copy(out[i*channels:], samples[srcIdx*channels:(srcIdx+1)*channels])
	}
	return out
}

// PlaySFX plays a sound effect by filename (non-blocking).
func (am *AudioManager) PlaySFX(filename string) {
	am.mu.Lock()
	defer am.mu.Unlock()

	if am.muted || am.sfxVol <= 0 {
		return
	}

	entry, ok := am.sfx[filename]
	if !ok || entry == nil {
		return
	}

	player, err := am.context.NewPlayer(bytes.NewReader(entry.data))
	if err != nil {
		log.Printf("audio: failed to create player for %s: %v", filename, err)
		return
	}
	player.SetVolume(am.sfxVol)
	player.Play()
}

// SetSFXVolume sets the SFX volume (0.0 to 1.0).
func (am *AudioManager) SetSFXVolume(v float64) {
	am.mu.Lock()
	defer am.mu.Unlock()
	am.sfxVol = clampVolume(v)
}

// SetBGMVolume sets the BGM volume (0.0 to 1.0).
func (am *AudioManager) SetBGMVolume(v float64) {
	am.mu.Lock()
	defer am.mu.Unlock()
	am.bgmVol = clampVolume(v)
}

// ToggleMute toggles master mute.
func (am *AudioManager) ToggleMute() {
	am.mu.Lock()
	defer am.mu.Unlock()
	am.muted = !am.muted
}

// IsMuted returns whether audio is muted.
func (am *AudioManager) IsMuted() bool {
	am.mu.Lock()
	defer am.mu.Unlock()
	return am.muted
}

// SFXVolume returns the current SFX volume.
func (am *AudioManager) SFXVolume() float64 {
	am.mu.Lock()
	defer am.mu.Unlock()
	return am.sfxVol
}

// BGMVolume returns the current BGM volume.
func (am *AudioManager) BGMVolume() float64 {
	am.mu.Lock()
	defer am.mu.Unlock()
	return am.bgmVol
}

// PlayLaser1 plays the player laser sound.
func (am *AudioManager) PlayLaser1() {
	am.PlaySFX("sfx_laser1.ogg")
}

// PlayLaser2 plays the enemy laser sound.
func (am *AudioManager) PlayLaser2() {
	am.PlaySFX("sfx_laser2.ogg")
}

// PlayLose plays the game over / death sound.
func (am *AudioManager) PlayLose() {
	am.PlaySFX("sfx_lose.ogg")
}

// PlayShieldUp plays the shield activation sound.
func (am *AudioManager) PlayShieldUp() {
	am.PlaySFX("sfx_shieldUp.ogg")
}

// PlayShieldDown plays the shield deactivation sound.
func (am *AudioManager) PlayShieldDown() {
	am.PlaySFX("sfx_shieldDown.ogg")
}

// PlayPowerUpPickup plays the power-up collection sound.
func (am *AudioManager) PlayPowerUpPickup() {
	am.PlaySFX("sfx_twoTone.ogg")
}

// PlayExplosion plays the explosion sound.
func (am *AudioManager) PlayExplosion() {
	am.PlaySFX("sfx_zap.ogg")
}
