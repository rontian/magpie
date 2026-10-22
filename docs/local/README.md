# rontian/magpie 本地补充文档

> 这部分是 rontian 维护的使用笔记和协作指令，不属于 yetone/magpie 上游官方手册。为便于追踪上游更新，个人内容集中放在本目录。

## 入口

- [项目指令（简体中文）](PROJECT_INSTRUCTIONS.zh-CN.md)：可复制进 ChatGPT Project 或其他编码助手的项目指令；说明上游、fork、私有部署仓库、插件源码以及变更与安全边界。
- [路由组正则 / Go RE2 速查](ROUTING_REGEX.zh-CN.md)：regex101 Golang 测试方法、Magpie glob 与 re: 区别、GLM/GPT/Claude/Grok/DeepSeek/Qwen/Gemini 等规则、手动成员和别名注意事项。
- [独立 Go 规则校验器](magpie-regex-check.go)：按 Magpie 的 compilePattern 规则测试完整 provider/model 字符串；可直接 go run，无第三方依赖。

## 其他仓库

- 上游： https://github.com/yetone/magpie
- 本地维护的 fork： https://github.com/rontian/magpie
- 独立部署配置（私有）： rontian/magpie-deploy
- 社区插件： https://github.com/magpie-community/plugins

私有部署仓库中的 IP、端口、账号、令牌、备份和配置值不可复制到本公开源码仓库。部署状态请以实时容器与部署仓库为准，不以本目录中示例推断。

## 版本验证

本目录的规则按 2026-10-22 读取的本 fork 源码整理；实际运行环境可使用不同的镜像 tag 或插件版本。模型 ID 与正则命中结果必须在实际 Provider 模型目录中确认。
