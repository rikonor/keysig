package utils

// RecomputeAverage recalculates an average given the old average,
// the old samples count and the new sample value
func RecomputeAverage(newSample, oldAvg float64, oldSampleCount uint64) float64 {
	newTotal := (oldAvg * float64(oldSampleCount)) + newSample
	return newTotal / float64((oldSampleCount + 1))
}
