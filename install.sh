#!/bin/bash
set -euo pipefail

BINARY_NAME="nyan-statusline"
INSTALL_DIR="$HOME/.claude"
SETTINGS_FILE="$INSTALL_DIR/settings.json"

echo "=== Nyan StatusLine Installer ==="

# 检测架构
ARCH=$(uname -m)
echo "Detected architecture: $ARCH"

# 构建
echo "Building..."
go build -ldflags="-s -w" -o "$BINARY_NAME" .
echo "Build complete."

# 安装二进制
mkdir -p "$INSTALL_DIR"
cp "$BINARY_NAME" "$INSTALL_DIR/$BINARY_NAME"
chmod +x "$INSTALL_DIR/$BINARY_NAME"
echo "Installed to $INSTALL_DIR/$BINARY_NAME"

# 更新 settings.json (statusLine + hooks)
# 使用 python3 处理 JSON 合并 (macOS 自带)
if [ -f "$SETTINGS_FILE" ]; then
    cp "$SETTINGS_FILE" "$SETTINGS_FILE.bak"
    echo "Backed up $SETTINGS_FILE to $SETTINGS_FILE.bak"
fi

python3 << 'PYEOF'
import json, os

settings_file = os.path.expanduser("~/.claude/settings.json")
binary = os.path.expanduser("~/.claude/nyan-statusline")

# 读取已有配置或创建空配置
cfg = {}
if os.path.exists(settings_file):
    with open(settings_file, "r") as f:
        cfg = json.load(f)

# statusLine 配置
cfg["statusLine"] = {
    "type": "command",
    "command": binary,
    "padding": 0,
}

# hooks 配置: UserPromptSubmit → processing, Stop → completed
nyan_hook_processing = {
    "matcher": "",
    "hooks": [
        {
            "type": "command",
            "command": f"{binary} --state processing",
        }
    ],
}
nyan_hook_completed = {
    "matcher": "",
    "hooks": [
        {
            "type": "command",
            "command": f"{binary} --state completed",
        }
    ],
}

# 合并 hooks (保留用户已有的 hooks)
if "hooks" not in cfg:
    cfg["hooks"] = {}
hooks = cfg["hooks"]

for event, hook_entry in [
    ("UserPromptSubmit", nyan_hook_processing),
    ("Stop", nyan_hook_completed),
]:
    if event not in hooks:
        hooks[event] = []
    # 移除旧的 nyan hook (通过检测 nyan-statusline 关键字)
    hooks[event] = [
        h for h in hooks[event]
        if not any(
            "nyan-statusline" in hk.get("command", "")
            for hk in h.get("hooks", [])
        )
    ]
    hooks[event].append(hook_entry)

with open(settings_file, "w") as f:
    json.dump(cfg, f, indent=2, ensure_ascii=False)

print("Updated settings.json with statusLine + hooks config.")
PYEOF

# 清理构建产物
rm -f "$BINARY_NAME"

echo ""
echo "=== Installation complete ==="
echo "Restart Claude Code to see the Nyan StatusLine!"
