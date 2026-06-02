**概览**

- 工具：ctdcnctl（源代码在本仓库）
- 作用：运维/管理 CLI，支持 Kubernetes、Kafka、Redis、MySQL、Harbor、Elasticsearch、MinIO、Nacos、StreamPark 等平台的常用管理与检查功能。

**快速开始**

- 运行本地代码：

```bash
go run main.go <command> [subcommand] [flags]
```

- 使用已编译二进制（在仓库根目录或上传到服务器）：

```bash
./ctdcnctl <command> [subcommand] [flags]
```

**Kubernetes（示例）**

- 列出命名空间（默认使用 `~/.kube/config`，可用 `--kubeconfig` 指定）：

```bash
# 使用默认 kubeconfig
./ctdcnctl k8s namespace list

# 指定 kubeconfig
./ctdcnctl k8s namespace list --kubeconfig /root/.kube/config
```

- 列出 Pod（可指定命名空间，留空表示所有命名空间）：

```bash
# 指定命名空间
./ctdcnctl k8s pod check -n default

# 所有命名空间
./ctdcnctl k8s pod check
```

- 注意：程序通过 `client-go` 直接访问 Kubernetes API，不依赖系统 `kubectl`。如果 kubeconfig 使用 `exec` 插件（例如 gcloud、aws-iam-authenticator），请确保相应的可执行程序也在服务器上可用。

**常用全局 flag**

- `--kubeconfig, -c`：指定 kubeconfig 文件路径（覆盖默认）
- `--namespace, -n`：用于 k8s pod 等子命令指定命名空间

**部署建议**

- 把 `ctdcnctl` 可执行文件上传到服务器（例如 `/usr/local/bin/`），并赋予执行权限：

```bash
scp ctdcnctl user@server:/usr/local/bin/
ssh user@server "chmod +x /usr/local/bin/ctdcnctl"
```

- 将 kubeconfig 放到 `/root/.kube/config` 或使用环境变量/配置文件指定路径。也可将本仓库生成的 `ctdcnctl` 配置文件（仓库根目录的 `ctdcnctl`）复制到 `/etc/profile.d/` 或直接在 shell 中 `source`。

**其他平台命令**

- 项目包含若干管理命令，位于 `cmd/` 目录。示例：
  - Kafka：`./ctdcnctl kafka ...`
  - Redis：`./ctdcnctl redis ...`
  - MySQL：`./ctdcnctl mysql ...`

**调试与构建**

- 本地构建：

```bash
# 本机默认构建
go build -o ctdcnctl .

# 交叉编译 Linux amd64
GOOS=linux GOARCH=amd64 go build -o ctdcnctl .
```

**安全提示**

- 不要在仓库中提交包含密钥（例如 DEEPSEEK API Key）的文件。建议使用环境变量（`DEEPSEEK_API_KEY`）或安全的配置管理方案。

**文件位置**

- 入口：`main.go`
- 命令实现：`cmd/` 目录
- 服务实现：`internal/` 目录
- 工具库：`pkg/` 目录

---

如需我将这份文档追加到 README 或生成更详细的命令手册（按 `cmd/` 自动扫描并生成子命令清单），告诉我要生成的格式（Markdown、HTML 或 plain text）。
