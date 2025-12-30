package nodes

import (
	"testing"
)

func TestParseStructuredScore_Success(t *testing.T) {
	content := `## 💡 决策简报 (Executive Summary)

【评分: 85/100】 【结论：细节优化】

## 🎯 结构化评分输出 (Structured Score Output)

` + "```json" + `
{
  "score": {
    "strategic_alignment": 90,
    "practical_value": 80,
    "logical_consistency": 85,
    "weighted_total": 85
  },
  "verdict": "细节优化",
  "exit_recommendation": false
}
` + "```" + `

### 评分矩阵总结
...
`

	score, err := ParseStructuredScore(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if score.Score.StrategicAlignment != 90 {
		t.Errorf("expected strategic_alignment 90, got %d", score.Score.StrategicAlignment)
	}
	if score.Score.WeightedTotal != 85 {
		t.Errorf("expected weighted_total 85, got %d", score.Score.WeightedTotal)
	}
	if score.Verdict != "细节优化" {
		t.Errorf("expected verdict '细节优化', got '%s'", score.Verdict)
	}
	if score.ExitRecommendation {
		t.Error("expected exit_recommendation to be false")
	}
}

func TestParseStructuredScore_NoJSON(t *testing.T) {
	content := `## 决策简报

没有 JSON 块的输出内容
`

	_, err := ParseStructuredScore(content)
	if err == nil {
		t.Error("expected error for content without JSON block")
	}
}

func TestStructuredScore_GetWeightedScore(t *testing.T) {
	// Test with explicit weighted_total
	score1 := &StructuredScore{}
	score1.Score.WeightedTotal = 85
	if got := score1.GetWeightedScore(); got != 85 {
		t.Errorf("expected 85, got %d", got)
	}

	// Test calculated weighted score (40/30/30 weights)
	score2 := &StructuredScore{}
	score2.Score.StrategicAlignment = 100
	score2.Score.PracticalValue = 80
	score2.Score.LogicalConsistency = 60
	// (100*40 + 80*30 + 60*30) / 100 = (4000 + 2400 + 1800) / 100 = 82
	if got := score2.GetWeightedScore(); got != 82 {
		t.Errorf("expected 82, got %d", got)
	}
}

func TestStructuredScore_ShouldExit(t *testing.T) {
	// Exit recommendation takes precedence
	score1 := &StructuredScore{ExitRecommendation: true}
	score1.Score.WeightedTotal = 50 // Below threshold
	if !score1.ShouldExit(90) {
		t.Error("should exit when exit_recommendation is true")
	}

	// Score threshold check
	score2 := &StructuredScore{}
	score2.Score.WeightedTotal = 92
	if !score2.ShouldExit(90) {
		t.Error("should exit when score >= threshold")
	}

	// Below threshold
	score3 := &StructuredScore{}
	score3.Score.WeightedTotal = 85
	if score3.ShouldExit(90) {
		t.Error("should not exit when score < threshold")
	}
}
