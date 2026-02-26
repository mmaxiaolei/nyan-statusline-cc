package state

import (
	"os"
	"path/filepath"
	"testing"
)

// TestIsProcessing_NoStateFile 无状态文件时应返回 true (默认处理中)
func TestIsProcessing_NoStateFile(t *testing.T) {
	dir := t.TempDir()
	if !IsProcessing(dir) {
		t.Error("IsProcessing() without state file should return true")
	}
}

// TestIsProcessing_Processing hook 写入 processing 状态时应返回 true
func TestIsProcessing_Processing(t *testing.T) {
	dir := t.TempDir()
	writeState(t, dir, `{"status":"processing"}`)

	if !IsProcessing(dir) {
		t.Error("IsProcessing() with processing status should return true")
	}
}

// TestIsProcessing_Completed hook 写入 completed 状态时应返回 false
func TestIsProcessing_Completed(t *testing.T) {
	dir := t.TempDir()
	writeState(t, dir, `{"status":"completed"}`)

	if IsProcessing(dir) {
		t.Error("IsProcessing() with completed status should return false")
	}
}

// TestIsProcessing_CorruptedFile 状态文件损坏时应返回 true (安全降级)
func TestIsProcessing_CorruptedFile(t *testing.T) {
	dir := t.TempDir()
	writeState(t, dir, "invalid json")

	if !IsProcessing(dir) {
		t.Error("IsProcessing() with corrupted file should return true")
	}
}

// TestIsProcessing_EmptyStatus status 为空字符串时应返回 true
func TestIsProcessing_EmptyStatus(t *testing.T) {
	dir := t.TempDir()
	writeState(t, dir, `{"status":""}`)

	if !IsProcessing(dir) {
		t.Error("IsProcessing() with empty status should return true")
	}
}

// TestIsProcessing_UnknownStatus 未知 status 值应返回 true
func TestIsProcessing_UnknownStatus(t *testing.T) {
	dir := t.TempDir()
	writeState(t, dir, `{"status":"unknown"}`)

	if !IsProcessing(dir) {
		t.Error("IsProcessing() with unknown status should return true")
	}
}

// TestIsProcessing_LegacyFormat 旧格式 (output_tokens) 应返回 true (兼容降级)
func TestIsProcessing_LegacyFormat(t *testing.T) {
	dir := t.TempDir()
	writeState(t, dir, `{"output_tokens":1234}`)

	if !IsProcessing(dir) {
		t.Error("IsProcessing() with legacy format should return true (status field missing)")
	}
}

// TestSetStatus_Processing 写入 processing 状态后可正确读取
func TestSetStatus_Processing(t *testing.T) {
	dir := t.TempDir()
	if err := SetStatus(dir, StatusProcessing); err != nil {
		t.Fatalf("SetStatus() error: %v", err)
	}
	if !IsProcessing(dir) {
		t.Error("IsProcessing() should return true after SetStatus(processing)")
	}
}

// TestSetStatus_Completed 写入 completed 状态后可正确读取
func TestSetStatus_Completed(t *testing.T) {
	dir := t.TempDir()
	if err := SetStatus(dir, StatusCompleted); err != nil {
		t.Fatalf("SetStatus() error: %v", err)
	}
	if IsProcessing(dir) {
		t.Error("IsProcessing() should return false after SetStatus(completed)")
	}
}

// TestSetStatus_Overwrite 多次写入应以最后一次为准
func TestSetStatus_Overwrite(t *testing.T) {
	dir := t.TempDir()
	_ = SetStatus(dir, StatusCompleted)
	_ = SetStatus(dir, StatusProcessing)

	if !IsProcessing(dir) {
		t.Error("IsProcessing() should return true after overwriting to processing")
	}
}

// TestSetStatus_CreatesFile 状态文件不存在时应自动创建
func TestSetStatus_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	if err := SetStatus(dir, StatusCompleted); err != nil {
		t.Fatalf("SetStatus() error: %v", err)
	}
	path := filepath.Join(dir, stateFileName)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Error("SetStatus() should create state file if not exists")
	}
}

func writeState(t *testing.T, dir, content string) {
	t.Helper()
	path := filepath.Join(dir, stateFileName)
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write state file: %v", err)
	}
}
