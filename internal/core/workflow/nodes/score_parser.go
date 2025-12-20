package nodes

import (
	"encoding/json"
	"fmt"
	"regexp"
)

// StructuredScore represents the parsed output from Adjudicator.
// This is used by the Workflow Engine to drive loop exit conditions.
type StructuredScore struct {
	Score struct {
		StrategicAlignment int `json:"strategic_alignment"`
		PracticalValue     int `json:"practical_value"`
		LogicalConsistency int `json:"logical_consistency"`
		WeightedTotal      int `json:"weighted_total"`
	} `json:"score"`
	Verdict            string `json:"verdict"`
	ExitRecommendation bool   `json:"exit_recommendation"`
}

// ParseAdjudicatorOutput extracts structured score from Adjudicator's markdown output.
// The JSON block is expected to be wrapped in ```json ... ``` code fences.
func ParseAdjudicatorOutput(content string) (*StructuredScore, error) {
	// Match JSON block within markdown code fences
	re := regexp.MustCompile("(?s)```json\\s*(\\{.*?\\})\\s*```")
	matches := re.FindStringSubmatch(content)
	if len(matches) < 2 {
		return nil, fmt.Errorf("no structured score JSON block found in adjudicator output")
	}

	var score StructuredScore
	if err := json.Unmarshal([]byte(matches[1]), &score); err != nil {
		return nil, fmt.Errorf("failed to parse structured score JSON: %w", err)
	}

	return &score, nil
}

// GetWeightedScore returns the weighted total score from parsed output.
// Falls back to calculating from individual scores if weighted_total is 0.
func (s *StructuredScore) GetWeightedScore() int {
	if s.Score.WeightedTotal > 0 {
		return s.Score.WeightedTotal
	}
	// Calculate weighted total: 40% strategic + 30% practical + 30% logical
	return (s.Score.StrategicAlignment*40 + s.Score.PracticalValue*30 + s.Score.LogicalConsistency*30) / 100
}

// ShouldExit determines if the workflow should exit based on score and recommendation.
func (s *StructuredScore) ShouldExit(threshold int) bool {
	if s.ExitRecommendation {
		return true
	}
	return s.GetWeightedScore() >= threshold
}
