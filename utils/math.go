package utils

import "math"

// Stats tracks a sample population and can be used for online calculation
// of mean, variance and standard deviation using Welford's method
// https://en.wikipedia.org/wiki/Algorithms_for_calculating_variance#Online_algorithm
type Stats struct {
	sampleCount  uint64
	mean, mValue float64
}

func NewStats() *Stats {
	return &Stats{}
}

func (s *Stats) Add(v float64) {
	oldSampleCount := s.sampleCount
	newSampleCount := oldSampleCount + 1

	// Update mean
	oldMean := s.mean
	newMean := oldMean + (v-oldMean)/float64(newSampleCount)

	// Update mValue
	oldMValue := s.mValue
	newMValue := oldMValue + (v-oldMean)*(v-newMean)

	s.sampleCount = newSampleCount
	s.mean = newMean
	s.mValue = newMValue
}

func (s *Stats) Count() uint64 {
	return s.sampleCount
}

func (s *Stats) Mean() float64 {
	return s.mean
}

func (s *Stats) Variance() float64 {
	return s.mValue / float64(s.sampleCount-1)
}

func (s *Stats) Std() float64 {
	return math.Sqrt(s.Variance())
}

// RecomputeAverage recalculates an average given the old average,
// the old samples count and the new sample value
func RecomputeAverage(newSample, oldAvg float64, oldSampleCount uint64) float64 {
	newTotal := (oldAvg * float64(oldSampleCount)) + newSample
	return newTotal / float64((oldSampleCount + 1))
}
