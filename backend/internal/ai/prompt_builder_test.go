package ai

import (
	"strings"
	"testing"
)

func TestBuildChatPromptUsesShortRoleplayContext(t *testing.T) {
	prompt := BuildChatPrompt(ChatContext{
		CityName:      "北京",
		CharacterName: "朱棣",
	})

	for _, want := range []string{
		"你在和用户进行角色扮演的游戏，你扮演的人物是朱棣",
		"当前城市是北京",
		"不要声称自己是真实复活的人",
		"不编造无法确认的确定性史实",
		"回答控制在150字以内",
	} {
		if !strings.Contains(prompt, want) {
			t.Fatalf("prompt missing %q: %s", want, prompt)
		}
	}

	for _, banned := range []string{"人物设定", "城市地标", "城市美食"} {
		if strings.Contains(prompt, banned) {
			t.Fatalf("prompt should not include %q: %s", banned, prompt)
		}
	}
}

func TestBuildChatPromptIncludesOptionalDialectRule(t *testing.T) {
	prompt := BuildChatPrompt(ChatContext{
		CityName:      "西安",
		CharacterName: "李白",
		DialectStyle:  "关中话",
	})

	want := "可少量使用关中话风格表达，但必须附普通话解释。"
	if !strings.Contains(prompt, want) {
		t.Fatalf("prompt missing dialect rule %q: %s", want, prompt)
	}
}
