package idg

import "testing"

func TestLoad(t *testing.T) {
	t.Run("progress", func(t *testing.T) {
		progress := NewProgress("./progress")
		data := map[string]int64{"file1": 0, "file2": 100, "file3": 200}
		progress.Refresh(data)
		progress.Load()
	})
}
