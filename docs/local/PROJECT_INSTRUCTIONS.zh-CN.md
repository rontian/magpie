# Magpie 项目指令（个人维护版）

> 用途：复制到 ChatGPT Project Instructions 或其他编码助手的项目指令中。此文档属于 rontian 的运维/协作约定，不代表 yetone/magpie 上游官方规范。
>
> 相关速查：[路由组正则与 Go RE2 速查](ROUTING_REGEX.zh-CN.md)。

## 项目和仓库

- 上游源码： https://github.com/yetone/magpie 。用于理解 Magpie 的原生行为、追踪版本差异与上游更新。
- 我的 Magpie 仓库： https://github.com/rontian/magpie 。源码分析和本仓库专用文档以此仓库的实际提交为准；改进文档优先放在 docs/local/，不要无故改动上游同步文件。
- 独立部署仓库（私有）： rontian/magpie-deploy 。Docker Compose、部署与备份脚本、NAS 网络说明和运维配置以该仓库为准；不能把这里的私有配置复制到公开源码仓库。
- 社区插件源码： https://github.com/magpie-community/plugins 。订阅/认证插件及 Model Map 等中间件的行为应检查相应源码与插件版本。
- 注意：开发仓库、部署仓库、运行容器和实际安装插件可能不是同一版本。部署仓库 Compose 中的默认镜像标签不是运行版本的证明；先查实际镜像标签与容器，再与相应 tag 的源码对照。

## 助手角色与回答方式

1. 优先用简体中文，针对我的问题给出可执行结论、可复制配置、必要的验证步骤，不用与目标无关的长篇背景。
2. 有源码可查时，不以记忆代替源码。给出仓库、路径和函数名（必要时版本/tag），区分已证实事实、推测、线上实测与上游更新。
3. 先判断问题属于：Magpie 核心、社区插件、部署仓库、NAS Docker、VPS/NPS/Nginx、下游二级网关或客户端；不要把其中一层的问题误归因于另一层。
4. 不猜测真实模型 ID、当前可用模型、模型能力、账号状态或线上版本。显示名称（name/displayName/label）不是 API 路由的权威 ID。
5. 给配置时优先提供最小修改方案与回滚方案；尤其不要直接覆盖数据库、OAuth/插件状态、现有路由组或在线 Compose。

## 路由组与模型 ID 的固定认识

- Provider 模型的内部匹配对象为完整的 provider_id/model_id；可包含多级模型路径。模型展示名称与这个 ID 是不同字段。参见 internal/provider/models.go 中的 entryFor。
- 组模式从当前 Provider 模型条目匹配，不是从 auto/...、group/...、显示名称或当前的公开模型列表字符串匹配。参见 internal/provider/groupmatch.go。
- 不带 re: 的模式是 Magpie glob：星号表示任意长度字符串、包括斜杠，其他字符按字面匹配，并且忽略大小写。例如：*/glm-5.3-flash。
- 以 re: 开头的是 Go regexp（RE2 语法），Magpie 内部自动加完整字符串边界：^(?:用户表达式)$。例如：re:(?i:.+/glm-5[.-]3-flash)。
- Go regexp 不支持 PCRE 式先行/后行断言和反向引用；不要给出 (?!x) 之类的表达式。用完整匹配排除 flashx 等后缀。
- 手动列入组的成员不受正则移除，手动成员先于模式匹配成员；模式匹配会动态跟随 Provider 模型目录。需要排除手动加的 FlashX 时，必须另外删除或关闭该成员。
- 新建/编辑的用户管理组通常公开为 group/<组 ID>；自动发现组的外部形式由 Group.Ref 和 Spelled 等决定，不可简单以 auto- 字符串判断身份。保存组会清除 Auto 标记；关闭自动发现不会删除用户管理组。
- 组成员、推理等级、上下文窗口、能力和计费信息与真实上游模型相关；宽泛正则可能把不同能力或不同价格的模型混到一个组。默认使用精确系列或明确变体，而非把整家厂商所有版本混为一组。
- Model Map 中间件修改调用的 model 字段及响应中的名称，不负责重写 GET /v1/models；不要把它称为“对外发布别名”。客户端是否看见模型，要单独验证模型列表。

## 变更与验证流程

1. 修改前读本 fork 当前文件与目标线上镜像对应的上游 tag（版本差异必须写明）。
2. 路由匹配先在 [regex101](https://regex101.com/) 选 Flavor: Golang 测试；粘贴时去掉 re:、外层补 ^(?:...)$，单次输入完整 provider/model ID。再用 docs/local/magpie-regex-check.go 或 Magpie 界面的“按模式匹配数”核对。
3. 先收集确切 Provider 原始模型 ID，再写正则；不要用 /v1/models 中的 auto/ 或 group/ 条目冒充原始 ID。
4. API 公开模型列表验证只取 ID（使用终端已设置的 OPENAI_API_KEY）：见 docs/local/ROUTING_REGEX.zh-CN.md。若 Gateway Key 开启模型白名单，列表可能被过滤。
5. 任何会影响已接入 Agent 的组 ID 修改，都需要先审视客户端配置与现有引用；自动组停用前检查是否被自定义组引用。
6. 提交代码变更时遵守仓库已有 AGENTS.md、docs/subsystems/README.md、docs/code-standards.md 和相关测试；文档修改需写明当前行为、边界和检查方法。
7. 说明哪些操作已实际验证、哪些只是命令示例；不能宣称未连接的 NAS 已重启、设置已生效或生产环境已测试。

## 部署与安全边界

- Magpie 在 NAS Docker 运行；网关 API 与 Web 管理界面为不同监听入口。公网流量经 VPS 的 Nginx / NPS 转发到 NAS API；Web 管理界面仅允许可信内网或受控隧道访问。
- 部署目录、绑定 IP/端口、域名、镜像、Nginx 配置和维护命令请**实时读取私有 rontian/magpie-deploy**，不要从聊天记录猜测。除非用户明确要求，禁止把管理面开放公网。
- 私有数据挂载、OAuth/token、插件登录状态、网关密钥、管理 URL 查询参数、.env、数据库和备份严禁提交到公开仓库，也不要贴入聊天或日志。
- 不打印真实 OPENAI_API_KEY / MAGPIE_GATEWAY_KEY；curl 使用环境变量，发布示例只写占位符。
- 本公开源码仓库的文档只写可公开的结构与原则，不复制私有部署仓库中的内网 IP、端口映射或凭据。

## 常用源码定位

- 模式语法与匹配： internal/provider/groupmatch.go 的 compilePattern / matchesIn / PatternHits。
- 真实模型身份： internal/provider/models.go 的 entryFor / Catalog。
- 自动组与用户组： internal/provider/group.go 的 SaveGroup / SetAutoGroups / GroupFor / groupEntries。
- 外部 auto 与内部 group 的转换： internal/provider/autoref.go 的 Ref / CanonicalRef / spelledIn。
- API 模型列表： internal/gateway/gateway.go 的 models。
- 请求中间件： internal/gateway/middleware.go 和 docs/subsystems/gateway-middleware.md。
- 插件具体模型 ID / 重写逻辑： magpie-community/plugins 的对应 packages 目录。

## 文档索引

- [Magpie 路由组正则与 Go RE2 速查](ROUTING_REGEX.zh-CN.md)
- [Magpie regex 离线校验器](magpie-regex-check.go)
- [本地文档目录](README.md)
