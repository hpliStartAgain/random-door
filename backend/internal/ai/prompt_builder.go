package ai

import (
	"fmt"
)

// ChatContext contains data needed to build a chat system prompt.
type ChatContext struct {
	CityName      string
	CharacterName string
	DialectStyle  string
}

// BuildChatPrompt assembles the system prompt for character chat.
func BuildChatPrompt(ctx ChatContext) string {
	dialectRule := ""
	if ctx.DialectStyle != "" {
		dialectRule = fmt.Sprintf("可少量使用%s风格表达，但必须附普通话解释。", ctx.DialectStyle)
	}

	return fmt.Sprintf("你在和用户进行角色扮演的游戏，你扮演的人物是%s。当前城市是%s。不要声称自己是真实复活的人，不编造无法确认的确定性史实。回答控制在150字以内。%s",
		ctx.CharacterName, ctx.CityName, dialectRule)
}

// BuildImagePrompt creates the prompt for AI image generation.
func BuildImagePrompt(cityName, landmarkName string) string {
	return fmt.Sprintf(
		"请生成一张真实旅行打卡照：保留用户上传自拍中的人物身份、脸部特征和自然姿态，"+
			"将人物自然合成到%s的%s场景中。可参考地标图片作为背景风格和构图依据。"+
			"自然日光、游客照片风格、高质量、真实构图。禁止生成在世公众人物、色情、暴力或侮辱内容。",
		cityName, landmarkName)
}

// BuildGuessCaptionPrompt creates a short social caption prompt for a panorama screenshot.
func BuildGuessCaptionPrompt(cityName, targetName, sceneHint, platform string) (systemPrompt, userPrompt string) {
	if targetName == "" {
		targetName = cityName
	}
	if sceneHint == "" {
		sceneHint = "用户刚在全景视角中截取了一张城市文化场景。"
	}
	systemPrompt = "你是城市文化互动产品的社交文案助手。只写可直接发布的中文文案，不要解释生成过程，不要编造确定性史实。"
	userPrompt = fmt.Sprintf(
		"城市：%s\n场景：%s\n截图线索：%s\n平台：%s\n要求：围绕“猜猜我在哪”写一条自然、有悬念的短文案；微博80字以内，可带2个话题；朋友圈60字以内，不要话题。",
		cityName, targetName, sceneHint, platform,
	)
	return systemPrompt, userPrompt
}
