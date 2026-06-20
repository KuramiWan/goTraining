# Space Shooter - Go Game Testing Sample

A 2D space shooter built with Go + Ebiten, designed as a game testing learning case.

## Controls

| Key | Action |
|-----|--------|
| WASD / Arrow Keys | Move |
| Mouse Left / Right | Rotate ship |
| Space | Shoot |
| P | Pause / Resume |
| Enter / Space | Start game / Restart after game over |

## Game Features

- **3 Lives** — lose one per meteor hit, with 2-second invulnerability after each hit
- **Difficulty Scaling** — meteor spawn rate increases every 10 points
- **Screen Boundary Clamping** — player cannot fly off-screen
- **Off-screen Cleanup** — bullets and meteors removed when they leave the screen
- **Game States** — Start screen → Playing → Pause → Game Over → Restart

## Exported Fields (for testing)

| Field | Type | Description |
|-------|------|-------------|
| `Game.Score` | `int` | Current score |
| `Game.Lives` | `int` | Remaining lives |
| `Game.State` | `GameState` | Current game state (StateStart/StatePlaying/StatePaused/StateGameOver) |
| `Game.DifficultyLevel` | `int` | Current difficulty level (score / 10) |

## Known Test Points

- Collision detection uses AABB (no rotation-aware hitbox)
- Bullets inherit player velocity on spawn
- Meteor spawn position is random around a circle centered on screen
- Invulnerability after hit lasts 2 seconds with flashing visual
- Game resets to StateStart on first launch; subsequent restarts go to StatePlaying

## Build & Run

```bash
go mod tidy
go run .
```
