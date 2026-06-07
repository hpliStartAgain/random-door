# Soundscape Assets

地标声景音频放在本目录，供 `landmarks.soundscape_url` 通过 `/static/soundscapes/...` 引用。

约束：

- 文件必须是项目可合法使用的授权音频，禁止直接搬运未知版权素材。
- 推荐格式：`wav` 或压缩后的 `mp3`，单文件建议小于 500KB。
- 前端只允许用户点击后播放，不自动播放。
- 命名建议：`<city>_<landmark>.wav`，例如 `beijing_forbidden_city.wav`。

当前仓库包含的短音频是轻量样例资源，后续可替换为更高质量的授权环境声。
