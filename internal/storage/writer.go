package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/WarnetBes/cursor-tool/internal/backup"
	"github.com/WarnetBes/cursor-tool/internal/uuid"
)

var TelemetryFields = []string{
	"telemetry.machineId",
	"telemetry.macMachineId",
	"telemetry.devDeviceId",
	"telemetry.sqmId",
}

type Result struct {
	Before     map[string]string
	After      map[string]string
	BackupPath string
}

func ModifyStorageIDs(path string, bm *backup.Manager) (*Result, error) {
	// Symlink guard
	info, err := os.Lstat(path)
	if err != nil {
		return nil, fmt.Errorf("lstat storage.json: %w", err)
	}
	if info.Mode()&os.ModeSymlink != 0 {
		return nil, fmt.Errorf("security: storage.json is a symlink")
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read: %w", err)
	}
	var content map[string]interface{}
	if err = json.Unmarshal(data, &content); err != nil {
		return nil, fmt.Errorf("parse json: %w", err)
	}

	bakPath, err := bm.Create(path)
	if err != nil {
		return nil, fmt.Errorf("backup: %w", err)
	}

	res := &Result{
		Before:     make(map[string]string),
		After:      make(map[string]string),
		BackupPath: bakPath,
	}

	for _, field := range TelemetryFields {
		if v, ok := content[field].(string); ok {
			res.Before[field] = v
		} else {
			res.Before[field] = "<not set>"
		}
		newID, err := uuid.Generate()
		if err != nil {
			return nil, rollback(bakPath, path, err)
		}
		content[field] = newID
		res.After[field] = newID
	}

	newData, err := json.MarshalIndent(content, "", "  ")
	if err != nil {
		return nil, rollback(bakPath, path, err)
	}
	if err = atomicWrite(path, newData); err != nil {
		return nil, rollback(bakPath, path, err)
	}

	verify, err := os.ReadFile(path)
	if err != nil || !json.Valid(verify) {
		return nil, rollback(bakPath, path,
			fmt.Errorf("post-write verification failed"))
	}

	return res, nil
}

func ReadCurrentIDs(path string) (map[string]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var content map[string]interface{}
	if err = json.Unmarshal(data, &content); err != nil {
		return nil, err
	}
	ids := make(map[string]string)
	for _, f := range TelemetryFields {
		if v, ok := content[f].(string); ok {
			ids[f] = v
		} else {
			ids[f] = "<not set>"
		}
	}
	return ids, nil
}

func atomicWrite(dst string, data []byte) error {
	dir := filepath.Dir(dst)
	tmp := filepath.Join(dir, filepath.Base(dst)+".tmp")
	f, err := os.OpenFile(tmp, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil { return err }
	defer func() { _ = os.Remove(tmp) }()
	if _, err = f.Write(data); err != nil { f.Close(); return err }
	if err = f.Sync(); err != nil { f.Close(); return err }
	f.Close()
	if err = os.Rename(tmp, dst); err != nil { return err }
	return os.Chmod(dst, 0600)
}

func rollback(bakPath, originalPath string, cause error) error {
	data, readErr := os.ReadFile(bakPath)
	if readErr != nil {
		return fmt.Errorf("failed (%w) + rollback read failed (%v)", cause, readErr)
	}
	if writeErr := atomicWrite(originalPath, data); writeErr != nil {
		return fmt.Errorf("failed (%w) + rollback write failed (%v)", cause, writeErr)
	}
	return fmt.Errorf("failed (rolled back): %w", cause)
}
