package dyaic

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"time"
)

func genPatch(old, new, patchName string) {
	fmt.Println(patchName, ": diff begin")
	beginTime := time.Now()
	patch, _ := exec.Command("diff", old, new).Output()
	fmt.Println(string(patch))
	err := ioutil.WriteFile(patchName, patch, 0644)
	if err != nil {
		log.Panic(err)
	}
	endTime := time.Now()
	fmt.Println(patchName, ": diff finished in", endTime.Sub(beginTime))
}

func genPatchForDirectory(old, new, patchName string) {
	beginTime := time.Now()
	patch, _ := exec.Command("diff", "-ruN", old, new).Output()
	err := ioutil.WriteFile(patchName, patch, 0644)
	if err != nil {
		log.Panic(err)
	}
	endTime := time.Now()
	fmt.Println("Gen patch for directory in", endTime.Sub(beginTime))
}

func patch(old, new, patchName string, clean bool) {
	cmd := exec.Command("patch", old, "-i", patchName, "-o", new)
	err := cmd.Run()
	if err != nil {
		log.Panic(err)
	}
	if clean {
		err = os.Remove(patchName)
		if err != nil {
			log.Panic(err)
		}
	}
}

// NOTICE: currently this function is only designed to apply patches in "~/.dyaic/patches/". To support patches in any directories, genPatchForDirectory needs to be modified together.
func patchForDirectory(old, patchName string, clean bool) {
	cmdString := fmt.Sprint("patch -d ", old, " -s -p5 < ", patchName)
	cmd := exec.Command("bash", "-c", cmdString)
	fmt.Println(cmd.String())
	err := cmd.Run()
	if err != nil {
		log.Panic(err)
	}
	if clean {
		err = os.Remove(patchName)
		if err != nil {
			log.Panic(err)
		}
	}
}

func genBSPatch(old, new, patchName string) {
	cmd := exec.Command("bsdiff", old, new, patchName)
	beginTime := time.Now()
	fmt.Println(patchName, ": bsdiff begin")
	err := cmd.Run()
	if err != nil { // bsdiff returns 0 on success and -1 on failure
		log.Panic(err)
	}
	endTime := time.Now()
	fmt.Println(patchName, ": bsdiff finished in", endTime.Sub(beginTime))
}

func bsPatch(old, new, patchName string, clean bool) {
	cmd := exec.Command("bspatch", old, new, patchName)
	err := cmd.Run()
	if err != nil {
		log.Panic(err)
	}
	if clean {
		err = os.Remove(patchName)
		if err != nil {
			log.Panic(err)
		}
	}
}
