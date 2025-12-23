package lib

import (
	"bufio"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

func ScanAndUpdate(cfg *Config, db *ModsDB) {
	root := cfg.Root
	if root == "" {
		log.Println("scan skipped: workshop_root is empty")
		return
	}

	entries, err := os.ReadDir(root)
	if err != nil {
		log.Println("scan error:", err)
		return
	}

	type result struct {
		folder string
		entry  ModEntry
	}

	results := make(chan result, len(entries))
	var wg sync.WaitGroup
	sem := make(chan struct{}, 10)

	existingFolders := map[string]bool{}
	existingMods := map[string]ModEntry{}

	for _, m := range db.Mods {
		existingMods[m.Folder] = m
	}

	for _, e := range entries {
		if !e.IsDir() {
			continue
		}

		folder := e.Name()
		if strings.HasPrefix(folder, ".") {
			continue
		}

		existingFolders[folder] = true

		if _, ok := existingMods[folder]; ok {
			continue
		}

		wg.Add(1)
		go func(folder string) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			fullPath := filepath.Join(root, folder)
			codename, name := extractFromFolder(fullPath)

			results <- result{
				folder: folder,
				entry: ModEntry{
					Name:         name,
					CodeName:     codename,
					Folder:       folder,
					Enabled:      true,
					DiscoveredAt: time.Now().UTC(),
				},
			}
		}(folder)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	found := map[string]ModEntry{}
	for r := range results {
		found[r.folder] = r.entry
	}

	newList := make([]ModEntry, 0, len(existingFolders))

	for _, m := range db.Mods {
		if !existingFolders[m.Folder] {
			log.Println("removing missing mod:", m.Folder)
			continue
		}

		if fresh, ok := found[m.Folder]; ok {
			fresh.Enabled = m.Enabled
			fresh.DiscoveredAt = m.DiscoveredAt
			if m.Name != "" {
				fresh.Name = m.Name
			}
			newList = append(newList, fresh)
		} else {
			newList = append(newList, m)
		}
	}

	for folder, m := range found {
		if _, ok := existingMods[folder]; !ok {
			newList = append(newList, m)
		}
	}

	db.Mods = newList
}

func extractFromFolder(folder string) (codename, name string) {
	log.Println("Extracting from folder:", folder)
	codename, name = extractFromScripts(folder)

	if codename == "" || name == "" {
		id := filepath.Base(folder)
		if title, err := FetchSteamTitle(id); err == nil && title != "" {
			log.Println("Steam API success for", id)
			if codename == "" {
				codename = id
			}
			if name == "" {
				name = title
			}
		} else {
			log.Println("Steam API failed for", id, "using folder name")
		}
	}

	if codename == "" {
		codename = filepath.Base(folder)
	}
	if name == "" {
		name = filepath.Base(folder)
	}

	return codename, name
}

func parseRpyFile(path string) (codename, pretty string) {
	f, err := os.Open(path)
	if err != nil {
		return "", ""
	}
	defer f.Close()
	s := bufio.NewScanner(f)
	for s.Scan() {
		line := s.Text()
		m := modsLineRe.FindStringSubmatch(line)
		if len(m) == 3 {
			k := m[1]
			val := strings.TrimSpace(m[2])
			clean := braceRe.ReplaceAllString(val, "")
			clean = strings.ReplaceAll(clean, `\"`, `"`)
			clean = strings.Trim(clean, `"`)
			clean = strings.Join(strings.Fields(clean), " ")
			return k, clean
		}
	}
	return "", ""
}

func extractFromScripts(folder string) (codename, name string) {
	codename, name = "", ""
	filepath.WalkDir(folder, func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		low := strings.ToLower(d.Name())
		if !strings.HasSuffix(low, ".rpy") {
			return nil
		}
		c, n := parseRpyFile(p)
		if c != "" && codename == "" {
			codename = c
		}

		if n != "" && name == "" {
			name = n
		}

		if codename != "" &&
			name != "" {
			return fs.SkipDir
		}
		return nil
	})
	return
}
