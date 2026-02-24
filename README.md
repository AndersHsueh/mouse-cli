# mouse-cli

Virtual mouse CLI - 通过命令行模拟鼠标输入

## 功能

- **move**: 移动鼠标 `mouse-cli move 100` 或 `mouse-cli move 100,200`
- **scroll**: 滚轮滚动 `mouse-cli scroll 3`
- **click**: 点击按钮 `mouse-cli click left`
- **press**: 按住按钮 `mouse-cli press left`
- **release**: 释放按钮 `mouse-cli release left`

## 使用方法

```bash
# 移动鼠标 (相对移动)
mouse-cli move 100              # 向右移动 100 像素
mouse-cli move -100,-50         # 向左 100, 向上 50

# 滚轮滚动
mouse-cli scroll 3              # 向上滚动 3 行
mouse-cli scroll -3             # 向下滚动 3 行

# 点击
mouse-cli click left           # 左键点击
mouse-cli click right          # 右键点击
mouse-cli click middle         # 中键点击

# 按住/释放 (用于拖拽)
mouse-cli press left           # 按住左键
mouse-cli move 100,100         # 移动
mouse-cli release left         # 释放左键
```

## 权限

需要 root 权限或加入 `input` 用户组来访问 `/dev/uinput`：

```bash
# 方法1: 使用 sudo
sudo ./mouse-cli click left

# 方法2: 加入 input 组
sudo usermod -aG input $USER
# 然后重新登录
```

## 技术

- 使用 Linux uinput 内核模块创建虚拟鼠标设备
- Go 1.21+

## License

MIT
