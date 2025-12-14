package prompt

import (
	"strings"
	"testing"
)

func TestBuildAffirmativeMessages(t *testing.T) {
	material := "测试材料内容"
	messages := BuildAffirmativeMessages(material)

	if len(messages) != 2 {
		t.Fatalf("BuildAffirmativeMessages() returned %d messages, want 2", len(messages))
	}

	// Check system message
	if messages[0].Role != "system" {
		t.Errorf("messages[0].Role = %v, want %v", messages[0].Role, "system")
	}
	if messages[0].Content != AffirmativeSystemPrompt {
		t.Error("messages[0].Content does not match AffirmativeSystemPrompt")
	}

	// Check user message
	if messages[1].Role != "user" {
		t.Errorf("messages[1].Role = %v, want %v", messages[1].Role, "user")
	}
	if !strings.Contains(messages[1].Content, material) {
		t.Errorf("messages[1].Content does not contain material: %v", messages[1].Content)
	}
	if !strings.Contains(messages[1].Content, "用户提供的材料") {
		t.Errorf("messages[1].Content does not contain expected prefix")
	}
}

func TestBuildNegativeMessages(t *testing.T) {
	material := "测试批判材料"
	messages := BuildNegativeMessages(material)

	if len(messages) != 2 {
		t.Fatalf("BuildNegativeMessages() returned %d messages, want 2", len(messages))
	}

	// Check system message
	if messages[0].Role != "system" {
		t.Errorf("messages[0].Role = %v, want %v", messages[0].Role, "system")
	}
	if messages[0].Content != NegativeSystemPrompt {
		t.Error("messages[0].Content does not match NegativeSystemPrompt")
	}

	// Check user message
	if messages[1].Role != "user" {
		t.Errorf("messages[1].Role = %v, want %v", messages[1].Role, "user")
	}
	if !strings.Contains(messages[1].Content, material) {
		t.Errorf("messages[1].Content does not contain material")
	}
}

func TestBuildAdjudicatorMessages(t *testing.T) {
	material := "原始材料"
	proArgument := "正方论点"
	conArgument := "反方论点"

	messages := BuildAdjudicatorMessages(material, proArgument, conArgument)

	if len(messages) != 2 {
		t.Fatalf("BuildAdjudicatorMessages() returned %d messages, want 2", len(messages))
	}

	// Check system message
	if messages[0].Role != "system" {
		t.Errorf("messages[0].Role = %v, want %v", messages[0].Role, "system")
	}
	if messages[0].Content != AdjudicatorSystemPrompt {
		t.Error("messages[0].Content does not match AdjudicatorSystemPrompt")
	}

	// Check user message contains all inputs
	if messages[1].Role != "user" {
		t.Errorf("messages[1].Role = %v, want %v", messages[1].Role, "user")
	}
	if !strings.Contains(messages[1].Content, material) {
		t.Error("messages[1].Content does not contain material")
	}
	if !strings.Contains(messages[1].Content, proArgument) {
		t.Error("messages[1].Content does not contain proArgument")
	}
	if !strings.Contains(messages[1].Content, conArgument) {
		t.Error("messages[1].Content does not contain conArgument")
	}
	if !strings.Contains(messages[1].Content, "原始材料") {
		t.Error("messages[1].Content does not contain expected section header")
	}
	if !strings.Contains(messages[1].Content, "正方观点") {
		t.Error("messages[1].Content does not contain pro section header")
	}
	if !strings.Contains(messages[1].Content, "反方观点") {
		t.Error("messages[1].Content does not contain con section header")
	}
}

func TestSystemPrompts(t *testing.T) {
	// Test AffirmativeSystemPrompt
	if AffirmativeSystemPrompt == "" {
		t.Error("AffirmativeSystemPrompt is empty")
	}
	if !strings.Contains(AffirmativeSystemPrompt, "战略支持者") {
		t.Error("AffirmativeSystemPrompt does not contain expected role description")
	}
	if !strings.Contains(AffirmativeSystemPrompt, "正方核心立场") {
		t.Error("AffirmativeSystemPrompt does not contain expected output format")
	}

	// Test NegativeSystemPrompt
	if NegativeSystemPrompt == "" {
		t.Error("NegativeSystemPrompt is empty")
	}
	if !strings.Contains(NegativeSystemPrompt, "批判性思维专家") {
		t.Error("NegativeSystemPrompt does not contain expected role description")
	}
	if !strings.Contains(NegativeSystemPrompt, "反方核心驳斥") {
		t.Error("NegativeSystemPrompt does not contain expected output format")
	}

	// Test AdjudicatorSystemPrompt
	if AdjudicatorSystemPrompt == "" {
		t.Error("AdjudicatorSystemPrompt is empty")
	}
	if !strings.Contains(AdjudicatorSystemPrompt, "首席裁决官") {
		t.Error("AdjudicatorSystemPrompt does not contain expected role description")
	}
	if !strings.Contains(AdjudicatorSystemPrompt, "综合裁决报告") {
		t.Error("AdjudicatorSystemPrompt does not contain expected output format")
	}
}

func TestBuildAffirmativeMessages_EmptyMaterial(t *testing.T) {
	messages := BuildAffirmativeMessages("")

	if len(messages) != 2 {
		t.Fatalf("BuildAffirmativeMessages() returned %d messages, want 2", len(messages))
	}
	// Should still work with empty material
	if messages[1].Role != "user" {
		t.Errorf("messages[1].Role = %v, want %v", messages[1].Role, "user")
	}
}

func TestBuildNegativeMessages_EmptyMaterial(t *testing.T) {
	messages := BuildNegativeMessages("")

	if len(messages) != 2 {
		t.Fatalf("BuildNegativeMessages() returned %d messages, want 2", len(messages))
	}
}

func TestBuildAdjudicatorMessages_EmptyInputs(t *testing.T) {
	messages := BuildAdjudicatorMessages("", "", "")

	if len(messages) != 2 {
		t.Fatalf("BuildAdjudicatorMessages() returned %d messages, want 2", len(messages))
	}
}

func TestBuildMessages_SpecialCharacters(t *testing.T) {
	material := "包含特殊字符: \n\t\"'<>&"

	affirmativeMessages := BuildAffirmativeMessages(material)
	if !strings.Contains(affirmativeMessages[1].Content, material) {
		t.Error("Special characters not preserved in affirmative messages")
	}

	negativeMessages := BuildNegativeMessages(material)
	if !strings.Contains(negativeMessages[1].Content, material) {
		t.Error("Special characters not preserved in negative messages")
	}

	adjudicatorMessages := BuildAdjudicatorMessages(material, "pro", "con")
	if !strings.Contains(adjudicatorMessages[1].Content, material) {
		t.Error("Special characters not preserved in adjudicator messages")
	}
}

func TestBuildMessages_LongContent(t *testing.T) {
	// Test with very long content
	longMaterial := strings.Repeat("这是一段很长的测试内容。", 1000)

	messages := BuildAffirmativeMessages(longMaterial)
	if !strings.Contains(messages[1].Content, longMaterial) {
		t.Error("Long content not preserved in messages")
	}
}
