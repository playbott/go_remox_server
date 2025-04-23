package server

import (
	"log"
	"math"
	"remox/configs"
	"remox/platform"
	"time"
)

func handleMoveCommand(
	controller platform.InputController,
	dx, dy int,
	config configs.AccelerationConfig,
	moveState *clientMoveState,
	deltaTime time.Duration,
) {

	var finalDx, finalDy int

	if dx == 0 && dy == 0 && moveState.consistentDuration == 0 {
		finalDx, finalDy = 0, 0

		moveState.lastDx = 0
		moveState.lastDy = 0
	} else {

		currentSignX := math.Copysign(1, float64(dx))
		if dx == 0 {
			currentSignX = 0
		}
		currentSignY := math.Copysign(1, float64(dy))
		if dy == 0 {
			currentSignY = 0
		}

		resetDuration := false
		if deltaTime > config.ResetTimeGapMs {
			resetDuration = true

		} else if dx == 0 && dy == 0 {
			resetDuration = true

			moveState.consistentSignX = 0
			moveState.consistentSignY = 0
		} else {

			changedX := currentSignX != 0 && moveState.consistentSignX != 0 && currentSignX != moveState.consistentSignX
			changedY := currentSignY != 0 && moveState.consistentSignY != 0 && currentSignY != moveState.consistentSignY

			startedX := currentSignX != 0 && moveState.consistentSignX == 0
			startedY := currentSignY != 0 && moveState.consistentSignY == 0
			if changedX || changedY || (changedX && startedY) || (changedY && startedX) {
				resetDuration = true

				moveState.consistentSignX = currentSignX
				moveState.consistentSignY = currentSignY
			}
		}

		if resetDuration {
			moveState.consistentDuration = 0

			if !(dx == 0 && dy == 0) {
				moveState.consistentSignX = currentSignX
				moveState.consistentSignY = currentSignY
			}
		} else if dx != 0 || dy != 0 {

			moveState.consistentDuration += deltaTime

			if moveState.consistentDuration > config.MaxTrackedDurationMs {
				moveState.consistentDuration = config.MaxTrackedDurationMs
			}

			if moveState.consistentSignX == 0 && currentSignX != 0 {
				moveState.consistentSignX = currentSignX
			}
			if moveState.consistentSignY == 0 && currentSignY != 0 {
				moveState.consistentSignY = currentSignY
			}
		}

		velocityFactor, durationFactor := 1.0, 1.0
		if config.VelocitySensitivity > 0 {
			velocity := math.Sqrt(float64(dx*dx + dy*dy))
			if velocity > config.VelocityThreshold {
				velocityFactor = 1.0 + velocity*config.VelocitySensitivity
				if config.MaxVelocityFactor > 1.0 && velocityFactor > config.MaxVelocityFactor {
					velocityFactor = config.MaxVelocityFactor
				}
			}
		}
		if config.DurationSensitivity > 0 && moveState.consistentDuration >= config.MinConsistentDurationMs {
			relevantDuration := moveState.consistentDuration - config.MinConsistentDurationMs
			durationFactor = 1.0 + (float64(relevantDuration)/float64(time.Millisecond))*config.DurationSensitivity
			if config.MaxDurationFactor > 1.0 && durationFactor > config.MaxDurationFactor {
				durationFactor = config.MaxDurationFactor
			}
		}
		totalFactor := velocityFactor * durationFactor
		if config.MaxTotalFactor > 1.0 && totalFactor > config.MaxTotalFactor {
			totalFactor = config.MaxTotalFactor
		}

		if totalFactor == 1.0 {
			finalDx, finalDy = dx, dy
		} else {
			finalDx = int(math.Round(float64(dx) * totalFactor))
			finalDy = int(math.Round(float64(dy) * totalFactor))

			if dx != 0 && finalDx == 0 {
				finalDx = int(math.Copysign(1, float64(dx)))
			}
			if dy != 0 && finalDy == 0 {
				finalDy = int(math.Copysign(1, float64(dy)))
			}
		}

	}

	moveState.lastDx = dx
	moveState.lastDy = dy

	if finalDx != 0 || finalDy != 0 {
		err := controller.MoveMouseBy(finalDx, finalDy)
		if err != nil {
			log.Printf("Error executing MoveMouseBy(%d, %d): %v", finalDx, finalDy, err)

			moveState.consistentDuration = 0
		}
	}
}

func handleButtonState(
	controller platform.InputController,
	reportedButtons ButtonStateMap,
	lastKnownButtons ButtonStateMap,
) {

	for button, isPressed := range reportedButtons {
		lastPressed, exists := lastKnownButtons[button]

		if !exists || lastPressed != isPressed {
			err := controller.SetButtonState(button, isPressed)
			if err != nil {
				log.Printf("Error setting button '%s' state to %v: %v", button, isPressed, err)

			}
			lastKnownButtons[button] = isPressed
		}
	}

	for button, wasPressed := range lastKnownButtons {
		if wasPressed {
			if _, stillReported := reportedButtons[button]; !stillReported {
				log.Printf("Button '%s' presumed released (missing from report), attempting release.", button)
				err := controller.SetButtonState(button, false)
				if err != nil {
					log.Printf("Error auto-releasing button '%s': %v", button, err)
				}
				lastKnownButtons[button] = false
			}
		}
	}
}

func handleScroll(
	controller platform.InputController,
	deltaY int,
	state *clientState,
	config configs.AccelerationConfig,
) {

	if deltaY == 0 && math.Abs(state.accumulatedScroll) < 1.0 {
		return
	}

	sensitivity := config.ScrollSensitivity
	if sensitivity <= 0.01 {
		sensitivity = 0.01
	}

	threshold := float64(configs.PlatformWheelDelta) / sensitivity

	state.accumulatedScroll += float64(deltaY)

	if math.Abs(state.accumulatedScroll) >= threshold {

		numTicks := int(state.accumulatedScroll / threshold)

		var amountPerTick int
		if numTicks > 0 {
			amountPerTick = -configs.PlatformWheelDelta
		} else if numTicks < 0 {
			amountPerTick = configs.PlatformWheelDelta
		} else {
			return
		}

		successTicks := 0
		for i := 0; i < int(math.Abs(float64(numTicks))); i++ {
			err := controller.Scroll(amountPerTick)
			if err != nil {
				log.Printf("Error executing scroll tick %d/%d: %v", i+1, int(math.Abs(float64(numTicks))), err)
				break
			}
			successTicks++
		}

		if successTicks > 0 {
			accumSign := 1.0
			if state.accumulatedScroll < 0 {
				accumSign = -1.0
			}
			processedAmount := accumSign * float64(successTicks) * threshold
			state.accumulatedScroll -= processedAmount
		}
	}

}
