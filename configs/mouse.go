package configs

import "time"

type AccelerationConfig struct {
	VelocitySensitivity float64
	VelocityThreshold   float64
	MaxVelocityFactor   float64

	DurationSensitivity     float64
	MinConsistentDurationMs time.Duration
	MaxTrackedDurationMs    time.Duration
	MaxDurationFactor       float64

	ScrollSensitivity float64

	MaxTotalFactor float64
	ResetTimeGapMs time.Duration
}

func DefaultAccelerationConfig() AccelerationConfig {
	return AccelerationConfig{
		VelocitySensitivity:     0.03,
		VelocityThreshold:       1.0,
		MaxVelocityFactor:       3.0,
		DurationSensitivity:     0.0005,
		MinConsistentDurationMs: 150 * time.Millisecond,
		MaxTrackedDurationMs:    1000 * time.Millisecond,
		MaxDurationFactor:       2.0,

		ScrollSensitivity: 4.0,

		MaxTotalFactor: 4.0, ResetTimeGapMs: 100 * time.Millisecond,
	}
}

const PlatformWheelDelta = 120
