package ai

import (
	"fmt"
	"strings"
)

// ChatContext contains data needed to build a chat system prompt.
type ChatContext struct {
	CityName      string
	CharacterName string
	Persona       string
	Landmarks     []string
	Foods         []string
	DialectStyle  string
}

// BuildChatPrompt assembles the system prompt for character chat.
func BuildChatPrompt(ctx ChatContext) string {
	landmarks := "暂无"
	if len(ctx.Landmarks) > 0 {
		landmarks = strings.Join(ctx.Landmarks, "、")
	}
	foods := "暂无"
	if len(ctx.Foods) > 0 {
		foods = strings.Join(ctx.Foods, "、")
	}
	dialect := "无特殊方言"
	if ctx.DialectStyle != "" {
		dialect = ctx.DialectStyle
	}

	return fmt.Sprintf(`你现在扮演一个城市文化导览角色。

当前城市：%s
当前人物：%s
人物设定：%s
城市地标：%s
城市美食：%s
方言特点：%s

回答要求：
1. 保持人物风格，但不要声称自己是真实复活的人。
2. 可以介绍城市历史、地标、美食和地方文化。
3. 可以使用少量方言词汇，但必须给出普通话解释。
4. 不要编造无法确认的历史事实。
5. 回答控制在 150 字以内。
6. 尽量引导用户继续探索当前城市。`,
		ctx.CityName, ctx.CharacterName, ctx.Persona,
		landmarks, foods, dialect)
}

// BuildImagePrompt creates the prompt for AI image generation.
func BuildImagePrompt(cityName, landmarkName string) string {
	return fmt.Sprintf(
		"Create a realistic travel photo of the uploaded person visiting %s in %s. "+
			"Keep the person's identity consistent. "+
			"Use the landmark image as background reference. "+
			"Natural daylight, tourist photo style, high quality, realistic composition.",
		landmarkName, cityName)
}
