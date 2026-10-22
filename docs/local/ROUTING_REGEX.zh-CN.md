# Magpie 路由组模式 / Go RE2 正则速查

> 适用范围：rontian/magpie（代码源自 yetone/magpie）；编写于 2026-10-22。此文档是用户维护的规则库，不代表上游承诺支持某个模型。**所有表达式只对当前 Gateway 实际注册的 Provider 模型 ID 生效**。
>
> 配套：[项目指令](PROJECT_INSTRUCTIONS.zh-CN.md) · [离线校验器](magpie-regex-check.go) · [源码 groupmatch.go](../../internal/provider/groupmatch.go)。

## 一、在线测试 Go 正则

**首选：[regex101.com](https://regex101.com/) → Flavor 选 Golang。** 其官方 [引擎说明](https://docs.regex101.com/editor/engines/) 已列出 Go regexp。其他参考：[Go 官方 regexp 语法](https://go.dev/src/regexp/syntax/doc.go)、[Go Playground](https://go.dev/play/)。

Magpie 编辑框中的表达式：

~~~text
re:(?i:.+/(?:glm|gml)[.-]?5[.-]?3-flash)
~~~

在 regex101 的 **Regular Expression** 中填（去掉 re:，补上 Magpie 自动加的完整匹配边界）：

~~~regex
^(?:(?i:.+/(?:glm|gml)[.-]?5[.-]?3-flash))$
~~~

在 **Test String** 单独输入一条真实完整 ID，例如 zhipu/glm-5.3-flash。对下一条 ID 再单独验证。不要随意打开多行 m 模式测试整批文本，因为 Magpie 是对**每一条 ID 单独完整匹配**。批量测试请使用文末的 Go 校验器。

注意：在 regex101 的选择器中要明确切换到 **Golang**，而不是默认的 PCRE2 或 JavaScript；测试环境只是语法验证，不等于线上 Provider 模型列表。

## 二、Magpie 的两套模式语法

| Magpie 模式 | 含义 | 备注 |
| --- | --- | --- |
| \`*/glm-5.3-flash\` | 任意 Provider 的精确 glm-5.3-flash | Magpie glob；自动忽略大小写，* 可跨 / |
| \`re:(?i:.+/glm-5\.3-flash)\` | 同一目标的 Go 正则 | 正则必须使用 .+、.* 等；点要写 \. |
| \`re:(?i:.+/glm-5[.-]3-flash)\` | 允许 glm-5.3-flash、glm-5-3-flash | 不匹配 flashx |
| \`re:(?i:.+/(?:glm|gml)[.-]?5[.-]?3-flash)\` | 容忍 glm/gml 和多种版本分隔写法 | 误拼 gml 仅在确认需要时加入 |
| \`re:(?i:(?:zhipu|trae-cn)/glm-5[.-]3-flash)\` | 只匹配两个指定 Provider | 需要 Provider 的真实 ID 正确 |

源码中正则编译为：

~~~go
regexp.Compile("^(?:" + re + ")$")
~~~

而不带 re: 的 glob 会先把普通部分 regexp.QuoteMeta，再把星号转换成 .*，整体用 (?i)^...$ 匹配。因此：

- **正则的 * 不能单独出现**；错误示例：re:(*/glm-5.3-flash)。
- **正则默认区分大小写**；统一写 (?i:...) 才能忽略大小写。
- . 是“任意一个字符”，要表示小数点需写 \.，或者用字符组 [.] / [.-]。
- .+/ 允许任意 Provider 前缀，也容纳形如 openrouter/anthropic/claude-... 的多层模型 ID。
- Go RE2 **没有** (?!x)、(?=...)、(?<=...) 或反向引用 \1。排除 flashx 的最简方法是只写 ...-flash 并使用完整字符串匹配。
- 匹配范围是 **Entry.ID = providerID + "/" + modelID**，不是展示名称、group/...、auto/...。所有匹配结果都是 Magpie 当前实际注册的 Provider 模型。
- 模式是“动态吸纳”而非覆盖名单：手动添加的成员会保留，且排在按模式匹配的成员之前。截图中的 FlashX 若是手动成员，换正则也不会把它移除。
- 编辑框“匹配 N 个模型”只代表当前模型条目的命中情况。真实路由前还要考虑 Provider 是否启用、组是否禁用、账号就绪情况等。

## 三、精确匹配优先：常用模型变体

下表中的模式可直接复制到 Magpie「路由组 → 模式」。若未明确需要，请不要把一个厂商的不同版本、Flash / Pro / Thinking 等变体混成同一个故障转移组。

| 用途 | 可复制的 Magpie 模式 | 说明 |
| --- | --- | --- |
| GLM-5.3-Flash，仅点或横杠 | \`re:(?i:.+/glm-5[.-]3-flash)\` | 不含 FlashX、Preview |
| GLM Flash，容忍 5.3 / 5-3 / 53 和 glm/gml | \`re:(?i:.+/(?:glm|gml)[.-]?5[.-]?3-flash)\` | 命中 glm-5.3-flash、glm-53-flash、gml53-flash；**不含 x** |
| GLM-5.3-FlashX（单独一组） | \`re:(?i:.+/glm-5[.-]3-flashx)\` | 不含 Flash |
| GLM-5.3 基础版 | \`re:(?i:.+/glm-5[.-]3)\` | 不含 Flash / FlashX |
| GPT-5.6 / GPT-5.6-Codex | \`re:(?i:.+/gpt-5[.-]6(?:-codex)?)\` | 仅指定版本，示例命名不表示平台必然提供 |
| Claude Opus 系列 | \`re:(?i:.+/claude-opus-[0-9]+(?:[.-][0-9]+)*(?:-[0-9]{8})?)\` | 支持版本点、短横线和 8 位日期后缀 |
| DeepSeek V4.1 Flash | \`re:(?i:.+/deepseek-v4[.-]1-flash)\` | 不含 V4 Pro、FlashX |
| DeepSeek V4 Pro | \`re:(?i:.+/deepseek-v4-pro)\` | 精确变体 |
| Grok 4 Fast 系列 | \`re:(?i:.+/grok-4(?:[.-][0-9]+)?-fast(?:-reasoning)?)\` | 允许可选小版本和 reasoning 后缀 |
| Qwen3.8-Flash（模型 ID 直接就是它时） | \`re:(?i:.+/qwen3[.-]?8-flash)\` | **不匹配** qoder/qfmodel 等别名 ID |
| Qwen3.8-Flash（含已核对的 Qoder qfmodel） | \`re:(?i:(?:qoder/qfmodel|.+/qwen3[.-]?8-flash))\` | Qoder CN 的实际 key 可能不同，不要未经核实照搬 |
| 指定 Provider 的 GLM-5.3-Flash | \`re:(?i:(?:zhipu|trae-cn)/glm-5[.-]3-flash)\` | 限制 Provider 范围，避免混入其他来源 |
| 某个无法辨认的真实模型 ID | \`re:(?i:.+/dfmodel)\` | 仅示意“按真实 ID 精确匹配”，不要仅凭 ID 推断背后型号 |

### GLM 容错模式的验证样本

| 完整原始模型 ID | 应命中？ |
| --- | --- |
| zhipu/glm-5.3-flash | ✅ |
| trae-cn/GLM-5-3-Flash | ✅ |
| workbuddy/glm-53-flash | ✅ |
| workbuddy/gml53-flash | ✅（启用 gml 兼容时） |
| trae-cn/glm-5.3-flashx | ❌ |
| zhipu/glm-5.3-flash-preview | ❌ |
| zhipu/glm-5.3 | ❌ |

## 四、厂商级宽匹配：只适合模型浏览或经人工审核的路由组

下列模式会包含较多变体。**不要直接把全家族放进一个生产故障转移组**：推理、上下文窗口、图像、输出限制、价格和接口可能不同。特别是一个“名称相似”的 Qwen、GPT、Claude 可能并非同一版本或同一能力。

| 家族 | Go RE2 模式（直接贴到 Magpie） | 示例（仅用于说明字符串形态） |
| --- | --- | --- |
| GPT | \`re:(?i:.+/gpt[._-][a-z0-9]+(?:[._-][a-z0-9]+)*)\` | gpt-4o、gpt-5.6、gpt-5.6-codex |
| OpenAI o 系列 | \`re:(?i:.+/o[0-9]+(?:[-._](?:mini|pro|high|medium))?)\` | o3、o4-mini |
| Claude（Opus/Sonnet/Haiku） | \`re:(?i:.+/claude[-._](?:opus|sonnet|haiku)[-._][0-9]+(?:[.-][0-9]+)*(?:-[0-9]{8})?)\` | claude-opus-5.5、claude-sonnet-4-5 |
| DeepSeek | \`re:(?i:.+/deepseek[._-][a-z0-9]+(?:[._-][a-z0-9]+)*)\` | deepseek-v4-pro、deepseek-r1 |
| GLM | \`re:(?i:.+/glm(?:[._-]?[a-z0-9]+)+)\` | glm-5.3、glm-5.3-flash、glm-5.3-flashx |
| Grok | \`re:(?i:.+/grok(?:[._-]?[a-z0-9]+)+)\` | grok-4、grok-code-fast-1 |
| Gemini | \`re:(?i:.+/gemini[._-][a-z0-9]+(?:[._-][a-z0-9]+)*)\` | gemini-2.5-pro |
| Qwen | \`re:(?i:.+/qwen[0-9a-z._-]+)\` | qwen3.8-flash、qwen3-235b |
| Kimi | \`re:(?i:.+/kimi(?:[._-][a-z0-9]+)+)\` | kimi-k3、kimi-k2.5 |
| MiniMax | \`re:(?i:.+/minimax[._-]?[a-z0-9]+(?:[._-][a-z0-9]+)*)\` | minimax-m2.7 |
| Step | \`re:(?i:.+/step[._-][0-9]+(?:[._-][a-z0-9]+)*)\` | step-5-preview |
| Llama | \`re:(?i:.+/llama[._-]?[0-9]+(?:[._-][a-z0-9]+)*)\` | llama-4-maverick |
| Mistral | \`re:(?i:.+/mistral(?:[._-][a-z0-9]+)+)\` | mistral-large |
| ERNIE | \`re:(?i:.+/ernie(?:[._-][a-z0-9]+)+)\` | ernie-4.5 |
| Doubao | \`re:(?i:.+/doubao(?:[._-][a-z0-9]+)+)\` | doubao-seed-1.6 |

**以上只是一组字符匹配模板**；模型是否存在、是否可用、是否同能力，必须以当前 Provider 返回的实际 ID 和元数据为准。有的供应商把真实模型叫作 qfmodel、dfmodel 或其他内部键，厂商家族正则根本不会命中它；请添加已验证的精确 ID 或人工成员，不要按 displayName 猜测。

## 五、排除 / 组合 / 限定 Provider

| 目的 | 模式 |
| --- | --- |
| 只要 glm-5.3-flash，排除 flashx | \`re:(?i:.+/glm-5[.-]3-flash)\` |
| Flash 与 FlashX **都要**（确实希望混用时） | \`re:(?i:.+/glm-5[.-]3-flashx?)\` |
| GLM Flash 或 DeepSeek V4 Pro | \`re:(?i:.+/(?:glm-5[.-]3-flash|deepseek-v4-pro))\` |
| 仅 WorkBuddy 的 DeepSeek V4 Pro | \`re:(?i:workbuddy/deepseek-v4-pro)\` |
| 只匹配最后一段模型名，允许多级路径 | \`re:(?i:.+/glm-5[.-]3-flash)\` |
| 不用正则、直接精确名称 | \`*/glm-5.3-flash\` |

Go RE2 没有否定先行断言。严格写明允许的后缀，比试图“匹配所有但排除 x”更安全。

## 六、从实际运行的网关查看 ID

仅获取当前 API key 可见的模型 ID（不要把 key 本身粘贴进命令或文档）：

~~~bash
curl -fsS 'https://api.magpie.rontian.top/v1/models' \
  -H "Authorization: Bearer $OPENAI_API_KEY" \
  | jq -r '.data[].id' | sort -u
~~~

这个接口可能只列出 auto/...、group/...：它展示的是**公开目录**，并不一定列出所有原始 Provider 模型。若 Provider 模型被设置“仅通过组显示”、网关密钥有模型白名单或 Provider 没有开放该模型，结果尤其不完整。写模式时应在 Magpie Providers 页面核对该模型的真实 ID，或通过受信任的管理途径读取 Provider 模型清单；不要把上面输出的 group/... 当作模式匹配对象。

## 七、用真正的 Go 运行批量测试（推荐）

本仓库随文档提供不依赖外部模块的 [magpie-regex-check.go](magpie-regex-check.go)，完整模拟 Magpie compilePattern 的**模式语法和整串匹配**。在仓库根目录执行：

~~~bash
go run docs/local/magpie-regex-check.go \
  're:(?i:.+/(?:glm|gml)[.-]?5[.-]?3-flash)' <<'MODELS'
zhipu/glm-5.3-flash
trae-cn/glm-5.3-flashx
workbuddy/glm-53-flash
workbuddy/gml53-flash
MODELS
~~~

预期输出：

~~~text
MATCH   zhipu/glm-5.3-flash
----    trae-cn/glm-5.3-flashx
MATCH   workbuddy/glm-53-flash
MATCH   workbuddy/gml53-flash
~~~

另外可用精确 glob 测试：

~~~bash
printf '%s\n' 'zhipu/glm-5.3-flash' 'zhipu/glm-5.3-flashx' \
  | go run docs/local/magpie-regex-check.go '*/glm-5.3-flash'
~~~

**注意：** 这个 CLI 检查器只验证“字符串是否按 Magpie 规则匹配”，不读取线上 Provider 注册信息，不验证模型可调用、账号/密钥/禁用状态，也不能代替生产路由测试。

## 八、组管理的稳定性建议

1. 如果你已关闭 Magpie 自动发现组，可以自己维护 group/<固定 ID>。组 ID 尽量稳定，成员可精确指定或用审慎的模式动态吸纳。
2. 每个能力/价格/上下文窗口差异明显的变体独立成组，如 glm-5.3-flash 与 glm-5.3-flashx **不要随意混组**。
3. 始终查看编辑页的匹配数量；手动成员与模式结果分别检查。注意“模型手动关闭”不同于“从组移除”。
4. Group ID 的 group/ 前缀不可通过普通组名称修改；Model Map 只重写调用和回复，不会替换 /v1/models 列表中的 ID。
5. 二级网关要求简短别名时，由它的模型发布/映射机制（必须同时涵盖列表和请求）处理；Magpie 组内的正式路由身份仍保留 group/...。

### 源码依据

- [internal/provider/groupmatch.go](../../internal/provider/groupmatch.go)：regexp 与 glob 语法、整串锚点、动态匹配和 PatternHits。
- [internal/provider/models.go](../../internal/provider/models.go)：模型 Entry.ID 的构成。
- [internal/provider/group.go](../../internal/provider/group.go)：组成员排序、保存、Auto 状态与 GroupFor。
- [internal/provider/autoref.go](../../internal/provider/autoref.go)：auto/... 的映射与公开引用。
- [internal/gateway/gateway.go](../../internal/gateway/gateway.go)：/v1/models 的列表逻辑。

维护规则：如果源码升级改变上述行为，请先修改或标注本速查，不要直接把某次聊天中的正则建议当成新的事实。
