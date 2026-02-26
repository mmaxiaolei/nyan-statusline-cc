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

# 更新 settings.json
if [ -f "$SETTINGS_FILE" ]; then
    # 备份原有配置
    cp "$SETTINGS_FILE" "$SETTINGS_FILE.bak"
    echo "Backed up $SETTINGS_FILE to $SETTINGS_FILE.bak"

    # 检查是否已有 statusLine 配置
    if grep -q '"statusLine"' "$SETTINGS_FILE"; then
        echo "statusLine config already exists in settings.json, skipping auto-config."
        echo "Please manually update if needed."
    else
        # 在 JSON 末尾的 } 前插入 statusLine 配置
        # 使用 python3 处理 JSON (macOS 自带)
        python3 -c "
import json, sys
with open('$SETTINGS_FILE', 'r') as f:
    cfg = json.load(f)
cfg['statusLine'] = {
    'type': 'command',
    'command': '$INSTALL_DIR/$BINARY_NAME',
    'padding': 0
}
with open('$SETTINGS_FILE', 'w') as f:
    json.dump(cfg, f, indent=2)
"
        echo "Updated $SETTINGS_FILE with statusLine config."
    fi
else
    # 创建新的 settings.json
    cat > "$SETTINGS_FILE" << EOF
{
  "statusLine": {
    "type": "command",
    "command": "$INSTALL_DIR/$BINARY_NAME",
    "padding": 0
  }
}
EOF
    echo "Created $SETTINGS_FILE with statusLine config."
fi

# 清理构建产物
rm -f "$BINARY_NAME"

echo ""
echo "=== Installation complete ==="
echo "Restart Claude Code to see the Nyan StatusLine!"
