package modules

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Exists checks whether all requested kernel modules are present in a proc-modules-style file.
func Exists(filename string, requested []string) (bool, error) {
	if len(requested) == 0 {
		return true, nil
	}

	file, err := os.Open(filename)
	if err != nil {
		return false, fmt.Errorf("failed to open %s: %w", filename, err)
	}
	defer file.Close()

	found := make(map[string]struct{}, len(requested))
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) == 0 {
			continue
		}
		for _, module := range requested {
			if fields[0] == module {
				found[module] = struct{}{}
			}
		}
		if len(found) == len(requested) {
			return true, nil
		}
	}
	if err := scanner.Err(); err != nil {
		return false, fmt.Errorf("failed to read %s: %w", filename, err)
	}

	return false, nil
}
