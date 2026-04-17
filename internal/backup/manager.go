package backup

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

const DefaultMaxBackups = 5

type Manager struct{ MaxBackups int }

func New(max int) *Manager {
	if max <= 0 { max = DefaultMaxBackups }
	return &Manager{MaxBackups: max}
}

func (m *Manager) Create(originalPath string) (string, error) {
	info, err := os.Lstat(originalPath)
	if err != nil {
		return "", fmt.Errorf("lstat: %w", err)
	}
	if info.Mode()&os.ModeSymlink != 0 {
		return "", fmt.Errorf("security: %s is a symlink, aborting",
			filepath.Base(originalPath))
	}
	data, err := os.ReadFile(originalPath)
	if err != nil {
		return "", fmt.Errorf("read for backup: %w", err)
	}
	ts := time.Now().Format("20060102_150405")
	bakPath := fmt.Sprintf("%s.%s.bak", originalPath, ts)
	if err := atomicWrite(bakPath, data, 0600); err != nil {
		return "", fmt.Errorf("write backup: %w", err)
	}
	_ = m.Prune(originalPath)
	return bakPath, nil
}

func (m *Manager) Restore(originalPath string) (string, error) {
	baks, err := m.List(originalPath)
	if err != nil || len(baks) == 0 {
		return "", fmt.Errorf("no backups found for %s",
			filepath.Base(originalPath))
	}
	data, err := os.ReadFile(baks[0])
	if err != nil {
		return "", fmt.Errorf("read backup: %w", err)
	}
	if err := atomicWrite(originalPath, data, 0600); err != nil {
		return "", fmt.Errorf("restore write: %w", err)
	}
	return baks[0], nil
}

func (m *Manager) List(originalPath string) ([]string, error) {
	dir := filepath.Dir(originalPath)
	base := filepath.Base(originalPath)
	matches, err := filepath.Glob(filepath.Join(dir, base+".*.bak"))
	if err != nil {
		return nil, err
	}
	sort.Slice(matches, func(i, j int) bool {
		return extractTS(matches[i]) > extractTS(matches[j])
	})
	return matches, nil
}

func (m *Manager) Prune(originalPath string) error {
	baks, err := m.List(originalPath)
	if err != nil { return err }
	keep := m.MaxBackups
	if keep > len(baks) { keep = len(baks) }
	for _, old := range baks[keep:] {
		_ = os.Remove(old)
	}
	return nil
}

func extractTS(path string) string {
	base := filepath.Base(path)
	for _, p := range strings.Split(base, ".") {
		if len(p) == 15 && strings.Contains(p, "_") {
			return p
		}
	}
	return ""
}

func atomicWrite(dst string, data []byte, perm os.FileMode) error {
	tmp := dst + ".tmp"
	f, err := os.OpenFile(tmp, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil { return err }
	defer func() { _ = os.Remove(tmp) }()
	if _, err = f.Write(data); err != nil { f.Close(); return err }
	if err = f.Sync(); err != nil { f.Close(); return err }
	f.Close()
	return os.Rename(tmp, dst)
}
