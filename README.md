# huachuang-note 命令行笔记管理工具

## 1. 项目概述

`huachuang-note` 是一个轻量级的命令行笔记管理工具，支持笔记的增删改查、搜索、导出等功能。

**技术栈**：
- **Go** — 编程语言
- **Cobra** — 命令行框架（[github.com/spf13/cobra](https://github.com/spf13/cobra)）
- **BoltDB (bbolt)** — 嵌入式 Key-Value 本地数据库（[go.etcd.io/bbolt](https://go.etcd.io/bbolt)）

---

## 2. 项目结构

```
huachuang_note-cli/
├── main.go                  # 程序入口，调用 cmd.Execute()
├── go.mod                   # Go Module 定义
├── go.sum                   # 依赖版本校验
├── cmd/                     # 命令行层（Cobra 命令定义）
│   ├── root.go              # 根命令 & 生命周期钩子 & DB 初始化
│   ├── add.go               # note add — 添加笔记
│   ├── list.go              # note list / note delete — 列出 & 删除笔记
│   ├── search.go            # note search — 搜索笔记
│   ├── export.go            # note export — 导出 JSON
│   ├── edit.go              # note edit — 编辑笔记
│   └── info.go              # note info — 查看笔记详情
└── storage/                 # 存储层（数据库操作封装）
    └── bolt.go              # BoltDB 封装 & Note 数据结构定义
```

**文件职责说明**：

| 文件 | 职责 |
|------|------|
| [main.go](file:///Users/lindan/Desktop/solo/task_project/huachuang_note-cli/main.go) | 程序入口，仅调用 `cmd.Execute()` 启动 Cobra |
| [cmd/root.go](file:///Users/lindan/Desktop/solo/task_project/huachuang_note-cli/cmd/root.go) | 定义根命令、生命周期钩子（DB 连接管理）、环境变量读取 |
| [cmd/add.go](file:///Users/lindan/Desktop/solo/task_project/huachuang_note-cli/cmd/add.go) | `note add` 子命令 |
| [cmd/list.go](file:///Users/lindan/Desktop/solo/task_project/huachuang_note-cli/cmd/list.go) | `note list` 和 `note delete` 子命令，含 Unicode 安全截断函数 |
| [cmd/search.go](file:///Users/lindan/Desktop/solo/task_project/huachuang_note-cli/cmd/search.go) | `note search` 子命令 |
| [cmd/export.go](file:///Users/lindan/Desktop/solo/task_project/huachuang_note-cli/cmd/export.go) | `note export` 子命令 |
| [cmd/edit.go](file:///Users/lindan/Desktop/solo/task_project/huachuang_note-cli/cmd/edit.go) | `note edit` 子命令 |
| [cmd/info.go](file:///Users/lindan/Desktop/solo/task_project/huachuang_note-cli/cmd/info.go) | `note info` 子命令 |
| [storage/bolt.go](file:///Users/lindan/Desktop/solo/task_project/huachuang_note-cli/storage/bolt.go) | Note 结构定义、BoltDB 打开/关闭、所有 CRUD 操作 |

---

## 3. 命令完整列表

本工具共提供 **7 个子命令**：

### 3.1 note add — 添加笔记

| 项 | 说明 |
|----|------|
| **命令格式** | `note add "笔记内容"` |
| **功能** | 将笔记保存到本地数据库，自动记录创建时间 |
| **参数** | 一个或多个字符串，自动用空格拼接为笔记内容 |
| **示例** | `note add "今天学习了 Go 语言的并发编程"` |
| **错误处理** | 参数不足时 Cobra 自动提示用法 |
| **输出** | `笔记已添加，ID: 1` |

---

### 3.2 note list — 列出所有笔记

| 项 | 说明 |
|----|------|
| **命令格式** | `note list` |
| **功能** | 列出所有笔记，显示 ID、内容摘要（前 50 字）和创建时间 |
| **参数** | 无 |
| **示例** | `note list` |
| **错误处理** | 无 |
| **输出** | 表格形式，空数据时显示 `暂无笔记` |

```
ID      内容摘要                                                     创建时间
------  -------------------------------------------------------  -------------------
1       今天学习了 Go 语言的并发编程，goroutine 和 channel 非常强大                2026-06-22T10:58:41+08:00
```

---

### 3.3 note delete — 删除笔记

| 项 | 说明 |
|----|------|
| **命令格式** | `note delete [ID]` |
| **功能** | 删除指定 ID 的笔记 |
| **参数** | 一个整数 ID |
| **示例** | `note delete 1` |
| **错误处理** | ID 不存在时：`笔记 ID 999 不存在`，退出码 1；ID 格式无效时：`无效的 ID: xxx` |
| **输出** | 成功时：`笔记 1 已删除` |

---

### 3.4 note search — 搜索笔记

| 项 | 说明 |
|----|------|
| **命令格式** | `note search "关键词"` |
| **功能** | 忽略大小写搜索笔记内容，显示匹配笔记的完整内容 |
| **参数** | 一个或多个字符串作为关键词 |
| **示例** | `note search "go"` |
| **错误处理** | 无 |
| **输出** | 匹配结果列表；无结果时显示 `未找到包含'xxx'的笔记` |

---

### 3.5 note export — 导出 JSON

| 项 | 说明 |
|----|------|
| **命令格式** | `note export` |
| **功能** | 将所有笔记导出为格式化的 JSON 文件到当前目录 |
| **参数** | 无 |
| **示例** | `note export` |
| **错误处理** | 写入文件失败时提示 |
| **输出** | `已导出 4 条笔记到 notes_export_20260622.json` |
| **文件名** | `notes_export_YYYYMMDD.json`（如 `notes_export_20260622.json`） |

---

### 3.6 note edit — 编辑笔记

| 项 | 说明 |
|----|------|
| **命令格式** | `note edit [ID] "新内容"` |
| **功能** | 更新指定 ID 的笔记内容，并记录修改时间（`updated_at`） |
| **参数** | 第一个参数为整数 ID，后续参数拼接为新内容 |
| **示例** | `note edit 1 "修改后的新内容"` |
| **错误处理** | ID 不存在时：`笔记 ID 999 不存在`，退出码 1 |
| **输出** | 成功时：`笔记 1 已更新` |

---

### 3.7 note info — 查看笔记详情

| 项 | 说明 |
|----|------|
| **命令格式** | `note info [ID]` |
| **功能** | 查看指定 ID 笔记的完整内容、创建时间、修改时间（如有） |
| **参数** | 一个整数 ID |
| **示例** | `note info 1` |
| **错误处理** | ID 不存在时：`笔记 ID 999 不存在`，退出码 1 |
| **输出** | ID、创建时间、修改时间（仅编辑过的笔记显示）、完整内容 |

---

## 4. 数据库设计

### 4.1 BoltDB 存储结构

本工具使用 **BoltDB**（bbolt）作为嵌入式本地数据库，采用单 Bucket、单表结构：

| 项 | 说明 |
|----|------|
| **Bucket 名称** | `notes`（字节序列 `[]byte("notes")`），定义于 [storage/bolt.go#L22](file:///Users/lindan/Desktop/solo/task_project/huachuang_note-cli/storage/bolt.go#L22) |
| **Key** | 8 字节大端序 `uint64`，即笔记 ID，通过 `binary.BigEndian.PutUint64()` 编码 |
| **Value** | JSON 序列化后的 `Note` 结构体 |

Key 的编码方式确保了 BoltDB 的 B+ 树按 ID 升序排列，遍历输出时保持插入顺序。

### 4.2 Note 结构体

定义于 [storage/bolt.go#L15-L20](file:///Users/lindan/Desktop/solo/task_project/huachuang_note-cli/storage/bolt.go#L15-L20)：

```go
type Note struct {
    ID        uint64     `json:"id"`
    Content   string     `json:"content"`
    CreatedAt time.Time  `json:"created_at"`
    UpdatedAt *time.Time `json:"updated_at,omitempty"`
}
```

| 字段 | 类型 | JSON Tag | 说明 |
|------|------|----------|------|
| `ID` | `uint64` | `id` | 笔记唯一标识，由 BoltDB Bucket 的 `NextSequence()` 自增生成 |
| `Content` | `string` | `content` | 笔记完整内容 |
| `CreatedAt` | `time.Time` | `created_at` | 创建时间，添加笔记时自动写入 `time.Now()` |
| `UpdatedAt` | `*time.Time` | `updated_at,omitempty` | 修改时间，仅在编辑笔记时写入；`omitempty` 确保未编辑时 JSON 不输出该字段 |

### 4.3 为什么使用 BoltDB？

1. **零配置嵌入式**：单文件数据库，不需要独立的数据库服务进程，编译后可执行文件即独立运行
2. **ACID 事务**：所有读写操作都在事务中进行，保证数据一致性
3. **Go 原生实现**：纯 Go 编写，无 CGO 依赖，交叉编译方便
4. **适合小规模数据**：对于个人笔记场景（数千至数万条）性能绰绰有余
5. **B+ 树存储**：键有序存储，遍历效率高

---

## 5. 核心逻辑说明

### 5.1 Cobra 命令注册与生命周期钩子

**命令注册机制**：
每个子命令文件（如 `add.go`、`list.go`）通过 Go 的 `init()` 函数在包初始化时自动调用 `rootCmd.AddCommand()` 将自己注册到根命令。这种模式下：
- 新增命令只需创建新文件，无需修改其他文件
- 导入 `cmd` 包即触发所有 `init()` 执行

**生命周期钩子**（定义于 [cmd/root.go#L19-L37](file:///Users/lindan/Desktop/solo/task_project/huachuang_note-cli/cmd/root.go#L19-L37)）：

| 钩子 | 作用 |
|------|------|
| `PersistentPreRunE` | 在任何子命令执行前运行，负责读取环境变量、打开数据库连接，赋值给包级变量 `store` |
| `PersistentPostRunE` | 在任何子命令执行后运行，负责关闭数据库连接 |

两个钩子配合实现了**资源自动管理**：开发者编写子命令时只需直接使用 `store`，无需关心 DB 的打开和关闭。

### 5.2 NOTE_DB_PATH 环境变量

逻辑位于 [cmd/root.go#L20-L28](file:///Users/lindan/Desktop/solo/task_project/huachuang_note-cli/cmd/root.go#L20-L28)：

```
读取 NOTE_DB_PATH 环境变量
    │
    ├─ 已设置 → 使用该路径作为数据库文件
    │
    └─ 未设置 → 获取用户主目录 (os.UserHomeDir())
                  → 使用 ~/.note.db 作为数据库文件
```

**设计目的**：
- 默认行为符合 Unix 工具习惯（配置文件放在用户主目录）
- 通过环境变量可方便地进行多环境切换、测试、数据隔离

### 5.3 Unicode 安全截断实现

函数位于 [cmd/list.go#L58-L76](file:///Users/lindan/Desktop/solo/task_project/huachuang_note-cli/cmd/list.go#L58-L76)：

```go
func truncate(s string, n int) string {
    runes := []rune(s)               // 转为 rune 切片，按 Unicode 码点而非字节计数
    if len(runes) <= n {
        return s
    }
    cut := n
    for cut > 0 {
        r := runes[cut-1]
        // 检查最后一个字符是否为：组合标记 / 变体选择符 / 零宽连接符
        if unicode.IsMark(r) || unicode.Is(unicode.Variation_Selector, r) || r == '\u200D' {
            cut--  // 往前回退
        } else {
            break  // 找到完整字符边界
        }
    }
    return string(runes[:cut]) + "..."
}
```

**处理的三类问题字符**：

| 字符类型 | 示例 | 说明 |
|----------|------|------|
| Unicode Mark (Mn/Mc/Me) | `é` = `e` + `́`（组合重音）、中文拼音声调 | 组合标记必须依附于前一个字符，不能独立显示 |
| Variation Selector | `VS-1` ~ `VS-256`、emoji 肤色选择符 | 用于改变前一个字符的显示样式 |
| ZWJ（零宽连接符 `U+200D`） | 👨‍👩‍👧‍👦 = 👨 + ZWJ + 👩 + ZWJ + 👧 + ZWJ + 👦 | 连接多个 emoji 形成复合图形字符 |

如果截断位置恰好落在这些字符上，算法会不断回退，直到找到一个可以独立显示的完整字符。

### 5.4 搜索功能的实现

位于 [storage/bolt.go#L117-L134](file:///Users/lindan/Desktop/solo/task_project/huachuang_note-cli/storage/bolt.go#L117-L134)：

```go
func (s *Store) SearchNotes(keyword string) ([]Note, error) {
    keywordLower := strings.ToLower(keyword)        // 关键词转小写
    // ...
    if strings.Contains(strings.ToLower(note.Content), keywordLower) {
        notes = append(notes, note)                 // 笔记内容转小写后匹配
    }
}
```

**忽略大小写的实现方式**：不修改原始数据，只在匹配时将关键词和笔记内容都转为小写再进行 `strings.Contains` 子串匹配。

**设计权衡**：
- ✅ 实现简单，无需额外索引
- ✅ 支持部分匹配
- ❌ 数据量大时性能一般（全表扫描），对于个人笔记场景可接受

### 5.5 JSON 序列化与 omitempty 标签

导出功能使用 `json.MarshalIndent` 生成格式化 JSON（定义于 [cmd/export.go#L21](file:///Users/lindan/Desktop/solo/task_project/huachuang_note-cli/cmd/export.go#L21)）。

`omitempty` 标签作用于 [storage/bolt.go#L19](file:///Users/lindan/Desktop/solo/task_project/huachuang_note-cli/storage/bolt.go#L19) 的 `UpdatedAt` 字段：

```go
UpdatedAt *time.Time `json:"updated_at,omitempty"`
```

**效果对比**：

| 场景 | JSON 输出 |
|------|-----------|
| 笔记未编辑过（`UpdatedAt == nil`） | `{"id": 2, "content": "...", "created_at": "..."}` — 不包含 `updated_at` |
| 笔记已编辑过（`UpdatedAt != nil`） | `{"id": 1, "content": "...", "created_at": "...", "updated_at": "..."}` |

如果不用 `omitempty`，`nil` 指针会被序列化为 `null`，不够语义化。

### 5.6 UpdatedAt 字段的指针设计

为什么使用 `*time.Time`（指针类型）而不是 `time.Time`（值类型）？

**原因 1：区分"零值"和"未设置"**
- `time.Time{}` 的零值是 `0001-01-01 00:00:00 UTC`，是一个合法的时间值
- 如果使用值类型，无法区分"笔记从未被修改"和"笔记在公元元年被修改"
- 使用指针类型后，`nil` 明确表示"未修改过"，非 `nil` 表示"有修改时间"

**原因 2：配合 omitempty**
- 只有指针、接口、切片、map 等引用类型的 `nil` 值才会被 `omitempty` 识别
- 值类型 `time.Time` 的零值 `time.Time{}` 不会被 `omitempty` 跳过

---

## 6. 代码架构

本项目采用经典的**两层分层架构**：

```
┌─────────────────────────────────┐
│   cmd/  (命令行层)               │
│   ┌─────────────────────────┐   │
│   │  root.go                │   │
│   │  add.go  list.go        │   │  依赖 ↓
│   │  search.go  export.go   │   │
│   │  edit.go  info.go       │   │
│   └─────────────────────────┘   │
└─────────────────────────────────┘
               │
               │ 使用 store 全局变量调用
               ↓
┌─────────────────────────────────┐
│   storage/  (存储层)             │
│   ┌─────────────────────────┐   │
│   │  bolt.go                │   │
│   │  - Note 结构体          │   │
│   │  - Store (DB 封装)      │   │
│   │  - CRUD 方法            │   │
│   └─────────────────────────┘   │
└─────────────────────────────────┘
               │
               │ 封装调用
               ↓
        BoltDB (bbolt)
```

**分层原则**：

| 层 | 职责 | 依赖方向 |
|----|------|----------|
| `cmd/` 命令行层 | 解析命令行参数、格式化输出、错误提示、调用存储层 | 依赖 `storage/` 包 |
| `storage/` 存储层 | 数据结构定义、数据库操作封装、事务管理 | **不依赖** `cmd/` 包，只依赖标准库和 bbolt |

**各层之间的依赖关系**：
- `cmd/` → `storage/`：通过包级变量 `store *storage.Store` 持有数据库实例，各子命令直接调用 `store.AddNote()`、`store.ListNotes()` 等方法
- `storage/` → `cmd/`：无任何依赖，存储层保持纯粹，可被其他上层（如 HTTP API、GUI）复用
- 各 `cmd/*.go` 文件之间：通过共享的 `rootCmd` 和 `store` 变量协作，但彼此互不调用

---

## 7. 运行方式

### 7.1 编译

```bash
# 在项目根目录下执行
go build -o note .
```

编译完成后会在当前目录生成可执行文件 `note`（macOS/Linux）或 `note.exe`（Windows）。

### 7.2 安装到系统

```bash
go install .
```

编译后的二进制会被安装到 `$GOPATH/bin`（或 `$HOME/go/bin`），确保该目录在 `PATH` 中即可全局使用 `note` 命令。

### 7.3 基本使用

```bash
# 添加笔记
./note add "今天学习了 Go 语言"

# 列出所有笔记
./note list

# 查看笔记详情
./note info 1

# 编辑笔记
./note edit 1 "修改后的内容"

# 搜索笔记（忽略大小写）
./note search "go"

# 删除笔记
./note delete 1

# 导出所有笔记为 JSON
./note export
```

### 7.4 使用 NOTE_DB_PATH 环境变量

**默认路径**：`~/.note.db`（用户主目录下的隐藏文件）

**自定义路径**：

```bash
# 临时设置（仅当前命令生效）
NOTE_DB_PATH=/tmp/my_notes.db ./note add "测试笔记"

# 长期设置（当前 Shell 会话）
export NOTE_DB_PATH=/path/to/custom.db
./note list
./note add "另一条笔记"
```

**使用场景**：
- 测试时使用临时目录，避免污染用户数据
- 不同项目/场景使用独立的笔记库
- 便携运行：将数据库文件放在 U 盘或同步目录
