// Package git 负责获取 Git 仓库信息
package git

import (
	"os/exec"
	"strings"
)

// Info 表示 Git 仓库状态
type Info struct {
	Branch     string
	HasChanges bool
}

// GetInfo 获取当前目录的 Git 分支和状态
// Return:
//   - *Info: Git 信息, 非 git 仓库时返回 nil
//   - error: 执行错误
func GetInfo() (*Info, error) {
	// 检查是否在 git 仓库中
	if err := exec.Command("git", "rev-parse", "--git-dir").Run(); err != nil {
		return nil, nil
	}

	// 获取当前分支
	out, err := exec.Command("git", "branch", "--show-current").Output()
	if err != nil {
		return nil, err
	}
	branch := strings.TrimSpace(string(out))

	// detached HEAD 时 branch 为空, 尝试获取短 commit hash
	if branch == "" {
		hashOut, err := exec.Command("git", "rev-parse", "--short", "HEAD").Output()
		if err != nil {
			return nil, nil
		}
		branch = strings.TrimSpace(string(hashOut))
		if branch == "" {
			return nil, nil
		}
	}

	// 检查是否有未提交的更改
	statusOut, err := exec.Command("git", "status", "--porcelain").Output()
	if err != nil {
		return &Info{Branch: branch}, nil
	}

	return &Info{
		Branch:     branch,
		HasChanges: strings.TrimSpace(string(statusOut)) != "",
	}, nil
}
