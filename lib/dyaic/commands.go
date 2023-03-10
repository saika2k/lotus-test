package dyaic

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"time"
)

func DyaicCommit(loc string, bs bool) {
	if loc == "" {
		loc = TempLocation
	}
	locLen := len(loc)
	err := filepath.Walk(loc, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rLoc := path[locLen:]
		repoLoc := RepoLocation + rLoc
		repoInfo, err := os.Stat(repoLoc)

		if exist(err) {
			if info.IsDir() {
				return nil
			}
			if info.ModTime().After(repoInfo.ModTime()) { // file has been modified, sync needed
				fmt.Println("File has been modified:", rLoc)
				patchName := repoLoc + ".patch"
				if bs {
					genBSPatch(repoLoc, path, patchName)
					bsPatch(repoLoc, repoLoc, patchName, true)
				} else {
					genPatch(repoLoc, path, patchName)
					patch(repoLoc, repoLoc, patchName, true)
				}
				fmt.Println("Updated.")
				// TODO: send changes tx
				// TODO: sync changes with other nodes
			}
		} else { // new file (or folder), creation needed
			if info.IsDir() {
				fmt.Println("Creating folder:", repoLoc)
				err = os.Mkdir(repoLoc, 0755)
				if err != nil {
					return err
				}
			} else {
				fmt.Println("New file:", rLoc)
				copy(path, repoLoc)
				fmt.Println("Copied.")
				// TODO: send file tx
				// TODO: sync file with other nodes
			}
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
}

func hashFile(loc string) {
	hashBegin := time.Now()
	if loc == "" {
		loc = TempLocation
	}
	fmt.Println(Md5File(loc))
	hashEnd := time.Now()
	fmt.Println(hashEnd.Sub(hashBegin))
}

func DyaicGitwalker() {
	homedir, err := os.UserHomeDir()
	if err != nil {
		log.Panic(err)
	}
	gitwalkerDir := homedir + "/.gitwalker/"
	patchesDir := homedir + "/.dyaic/patches/"
	versionNumber := countSubDirectories(gitwalkerDir)
	for d := 1; d < versionNumber; d++ {
		newDir := gitwalkerDir + fmt.Sprintf("%04d", d)
		oldDir := gitwalkerDir + fmt.Sprintf("%04d", d-1)
		patchName := patchesDir + fmt.Sprintf("%04d%04d.patch", d-1, d)
		fmt.Printf("Start Patching %04d~%04d\n", d-1, d)
		genPatchForDirectory(oldDir, newDir, patchName)
	}
}

// DyaicPatchGitwalker applies patches at ~/.dyaic/patches/ to repo versions at ~/.gitwalker/
func DyaicPatchGitwalker(bs bool) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		log.Panic(err)
	}
	gitwalkerDir := homedir + "/.gitwalker/"
	patchesDir := homedir + "/.dyaic/patches/"
	versionNumber := countSubDirectories(gitwalkerDir)
	for d := 1; d < versionNumber; d++ {
		oldDir := gitwalkerDir + fmt.Sprintf("%04d", d-1)
		patchName := patchesDir + fmt.Sprintf("%04d%04d.patch", d-1, d)
		fmt.Printf("Start Patching %04d~%04d\n", d-1, d)
		patchForDirectory(oldDir, patchName, bs)
	}
}

func DyaicPatch(loc string, bs bool) {
	if loc == "" {
		loc = TempLocation
	}
	locLen := len(loc)
	err := filepath.Walk(loc, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rLoc := path[locLen:]
		repoLoc := RepoLocation + rLoc
		//repoInfo, err := os.Stat(repoLoc)

		if exist(err) {
			if info.IsDir() {
				return nil
			}
			if !SameFile(path, repoLoc) { // file has been modified, sync needed
				fmt.Println("File has been modified:", rLoc, ", file size: ", info.Size())
				patchName := repoLoc + ".patch"
				if bs {
					genBSPatch(repoLoc, path, patchName)
				} else {
					genPatch(repoLoc, path, patchName)
				}
				fmt.Println("Generated patch file ", repoLoc, ".patch")
			}
		} else { // new file (or folder), creation needed
			if info.IsDir() {
				fmt.Println("New folder:", repoLoc)
			} else {
				fmt.Println("New file:", rLoc)
			}
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
}

func DyaicPrintDiff(loc string) {
	if loc == "" {
		loc = TempLocation
	}
	locLen := len(loc)
	err := filepath.Walk(loc, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rLoc := path[locLen:]
		repoLoc := RepoLocation + rLoc
		repoInfo, err := os.Stat(repoLoc)

		if exist(err) {
			if info.IsDir() {
				return nil
			}
			if info.ModTime().After(repoInfo.ModTime()) { // file has been modified, sync needed
				chs := GenerateChanges(repoLoc, path)
				if len(chs.Item) == 0 {
					return nil
				}
				fmt.Println("File has been modified:", rLoc)
				ShowDiff(repoLoc, path)
			}
		} else { // new file (or folder), creation needed
			if info.IsDir() {
				fmt.Println("New folder:", repoLoc)
			} else {
				fmt.Println("New file:", rLoc)
			}
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
}

func saveDiff(loc string) {
	if loc == "" {
		loc = TempLocation
	}
	locLen := len(loc)
	err := filepath.Walk(loc, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rLoc := path[locLen:]
		repoLoc := RepoLocation + rLoc
		repoInfo, err := os.Stat(repoLoc)

		if exist(err) {
			if info.IsDir() {
				return nil
			}
			if info.ModTime().After(repoInfo.ModTime()) { // file has been modified, sync needed
				chs := GenerateChanges(repoLoc, path)
				if len(chs.Item) == 0 {
					return nil
				}
				fmt.Println("File has been modified:", rLoc)
				SaveDyaicDiff(repoLoc, path)
			}
		} else { // new file (or folder), creation needed
			if info.IsDir() {
				fmt.Println("New folder:", repoLoc)
			} else {
				fmt.Println("New file:", rLoc)
			}
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
}

func DyaicPrintFolder(loc string) {
	if loc == "" {
		loc = TempLocation
	}
	err := filepath.Walk(loc, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		fmt.Println(path, info.ModTime(), info.Size())
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
}

func DyaicWatch(loc string) {
	watcher := Watch(loc)
	defer watcher.Close()
	select {}
}
