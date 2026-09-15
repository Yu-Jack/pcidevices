package modules

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExists(t *testing.T) {
	tests := []struct {
		name      string
		contents  string
		requested []string
		want      bool
	}{
		{
			name:      "all modules loaded",
			contents:  "vfio_pci 1 0 - Live 0x0\nvfio_iommu_type1 2 0 - Live 0x0\n",
			requested: []string{"vfio_pci", "vfio_iommu_type1"},
			want:      true,
		},
		{
			name:      "module missing",
			contents:  "vfio_pci 1 0 - Live 0x0\n",
			requested: []string{"vfio_pci", "vfio_iommu_type1"},
			want:      false,
		},
		{
			name:      "similar module name does not match",
			contents:  "vfio_pci_extra 1 0 - Live 0x0\n",
			requested: []string{"vfio_pci"},
			want:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filename := filepath.Join(t.TempDir(), "modules")
			if err := os.WriteFile(filename, []byte(tt.contents), 0600); err != nil {
				t.Fatal(err)
			}

			got, err := Exists(filename, tt.requested)
			if err != nil {
				t.Fatal(err)
			}
			if got != tt.want {
				t.Fatalf("Exists() = %v, want %v", got, tt.want)
			}
		})
	}
}
