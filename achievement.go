package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

// AchievementID uniquely identifies each achievement.
type AchievementID int

const (
	AchFirstKill AchievementID = iota
	AchScore100
	AchScore500
	AchScore1000
	AchScore5000
	AchSurvive1Min
	AchSurvive5Min
	AchKill50Enemies
	AchKill100Enemies
	AchKillBoss
	AchCollectAllPowerUps
	AchNoDamage1Min
	AchMaxDifficulty
	AchBronzeMilestone
	AchPiercingMaster
)

// Achievement holds the state of one achievement.
type Achievement struct {
	ID          AchievementID `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Unlocked    bool          `json:"unlocked"`
	UnlockedAt  time.Time     `json:"unlocked_at,omitempty"`
}

var allAchievements = []Achievement{
	{ID: AchFirstKill, Name: "初次击杀", Description: "摧毁第一个敌人"},
	{ID: AchScore100, Name: "小试牛刀", Description: "得分达到 100"},
	{ID: AchScore500, Name: "渐入佳境", Description: "得分达到 500"},
	{ID: AchScore1000, Name: "弹幕高手", Description: "得分达到 1000"},
	{ID: AchScore5000, Name: "传奇飞行员", Description: "得分达到 5000"},
	{ID: AchSurvive1Min, Name: "生存者", Description: "存活 1 分钟"},
	{ID: AchSurvive5Min, Name: "太空老兵", Description: "存活 5 分钟"},
	{ID: AchKill50Enemies, Name: "猎手", Description: "摧毁 50 个敌人"},
	{ID: AchKill100Enemies, Name: "清道夫", Description: "摧毁 100 个敌人"},
	{ID: AchKillBoss, Name: "Boss杀手", Description: "击败一个 Boss"},
	{ID: AchCollectAllPowerUps, Name: "收集狂", Description: "收集所有类型道具"},
	{ID: AchNoDamage1Min, Name: "完美闪避", Description: "连续 1 分钟不受伤害"},
	{ID: AchMaxDifficulty, Name: "极限挑战", Description: "难度等级达到 20"},
	{ID: AchBronzeMilestone, Name: "青铜里程碑", Description: "收集 5 个青铜道具"},
	{ID: AchPiercingMaster, Name: "穿透大师", Description: "一次穿透射击击中 5 个目标"},
}

// AchievementTracker manages achievement state and triggers.
type AchievementTracker struct {
	Achievements map[AchievementID]*Achievement `json:"achievements"`
	filePath     string
	// Transient tracking state (not persisted)
	lastDamageTime   time.Time
	noDamageStarted  bool
	collectedTypes   map[PowerUpType]bool
	bronzeCount      int
}

// NewAchievementTracker creates a new tracker and loads persisted state.
func NewAchievementTracker() *AchievementTracker {
	dir, err := os.UserConfigDir()
	if err != nil {
		dir = "."
	}
	appDir := filepath.Join(dir, "space-shooter")
	os.MkdirAll(appDir, 0755)

	t := &AchievementTracker{
		Achievements:   make(map[AchievementID]*Achievement),
		filePath:       filepath.Join(appDir, "achievements.json"),
		collectedTypes: make(map[PowerUpType]bool),
	}

	// Initialize all achievements as locked
	for i := range allAchievements {
		ach := allAchievements[i]
		t.Achievements[ach.ID] = &ach
	}

	return t
}

// Load reads achievement state from disk.
func (t *AchievementTracker) Load() error {
	data, err := os.ReadFile(t.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	var saved []Achievement
	if err := json.Unmarshal(data, &saved); err != nil {
		return err
	}
	for i := range saved {
		if ach, ok := t.Achievements[saved[i].ID]; ok {
			ach.Unlocked = saved[i].Unlocked
			ach.UnlockedAt = saved[i].UnlockedAt
		}
	}
	return nil
}

// Save writes achievement state to disk.
func (t *AchievementTracker) Save() error {
	var list []Achievement
	for _, ach := range t.Achievements {
		if ach.Unlocked {
			list = append(list, *ach)
		}
	}
	data, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(t.filePath, data, 0644)
}

// Unlock marks an achievement as unlocked and returns whether it was newly unlocked.
func (t *AchievementTracker) Unlock(id AchievementID) bool {
	ach, ok := t.Achievements[id]
	if !ok || ach.Unlocked {
		return false
	}
	ach.Unlocked = true
	ach.UnlockedAt = time.Now()
	t.Save()
	return true
}

// RecentlyUnlocked returns an achievement if one was just unlocked (for popup display).
// Caller should track the last displayed ID to avoid showing the same one repeatedly.
func (t *AchievementTracker) IsUnlocked(id AchievementID) bool {
	ach, ok := t.Achievements[id]
	if !ok {
		return false
	}
	return ach.Unlocked
}

// RecordDamage notifies the tracker that the player took damage.
func (t *AchievementTracker) RecordDamage() {
	t.lastDamageTime = time.Now()
}

// RecordKill notifies the tracker of an enemy kill.
func (t *AchievementTracker) RecordKill(kills int) {
	if kills >= 1 {
		t.Unlock(AchFirstKill)
	}
	if kills >= 50 {
		t.Unlock(AchKill50Enemies)
	}
	if kills >= 100 {
		t.Unlock(AchKill100Enemies)
	}
}

// RecordPowerUpCollected records a collected power-up type.
func (t *AchievementTracker) RecordPowerUpCollected(pType PowerUpType) {
	t.collectedTypes[pType] = true
	if pType == PowerUpBronze {
		t.bronzeCount++
		if t.bronzeCount >= 5 {
			t.Unlock(AchBronzeMilestone)
		}
	}
	// Check if we've collected all types
	allTypes := []PowerUpType{PowerUpSilver, PowerUpBronze, PowerUpGold, PowerUpShield, PowerUpHealth, PowerUpSpeed}
	for _, typ := range allTypes {
		if !t.collectedTypes[typ] {
			return
		}
	}
	t.Unlock(AchCollectAllPowerUps)
}

// CheckStats evaluates stat-based achievements.
func (t *AchievementTracker) CheckStats(stats *GameStats, score int, difficultyLevel int, elapsedTime time.Duration) {
	if score >= 100 {
		t.Unlock(AchScore100)
	}
	if score >= 500 {
		t.Unlock(AchScore500)
	}
	if score >= 1000 {
		t.Unlock(AchScore1000)
	}
	if score >= 5000 {
		t.Unlock(AchScore5000)
	}
	if elapsedTime >= 1*time.Minute {
		t.Unlock(AchSurvive1Min)
	}
	if elapsedTime >= 5*time.Minute {
		t.Unlock(AchSurvive5Min)
	}
	if difficultyLevel >= 20 {
		t.Unlock(AchMaxDifficulty)
	}
	if stats.PiercingHits >= 5 {
		t.Unlock(AchPiercingMaster)
	}
}

// TimeSinceLastDamage returns the duration since the last damage event.
func (t *AchievementTracker) TimeSinceLastDamage() time.Duration {
	if t.lastDamageTime.IsZero() {
		return 0
	}
	return time.Since(t.lastDamageTime)
}
