package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func CheckErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func Copy(src, dest string) {
	srcFile, err := os.Open(src)
	CheckErr(err)
	defer srcFile.Close()

	destFile, err := os.Create(dest)
	CheckErr(err)
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	CheckErr(err)

	err = destFile.Sync()
	CheckErr(err)
}

func main() {
	interval := 5
	maxBackups := 5

	fmt.Println("ALICE MADNESS RETURNS CHECKPOINT BACKUPPER")
	fmt.Println("------------------------------------------")
	fmt.Println("Interval:", interval, "minutes")
	fmt.Println("Maximum backups:", maxBackups)
	fmt.Println("------------------------------------------")

	home, err := os.UserHomeDir()
	CheckErr(err)

	saveFolder := home + "/Documents/My Games/Alice Madness Returns/AliceGame/CheckPoint"
	files, err := os.ReadDir(saveFolder)
	CheckErr(err)

	var profiles []string
	for _, f := range files {
		if f.IsDir() {
			profiles = append(profiles, f.Name())
		}
	}

	profile := profiles[0]

	for len(profiles) > 1 {
		fmt.Println("Multiple profile folders detected")
		for i, p := range profiles {
			fmt.Println(i, p)
		}

		fmt.Print("Choose profile (enter a number): ")

		var t int
		_, err := fmt.Scan(&t)
		CheckErr(err)

		if t >= 0 && t < len(profiles) {
			profile = profiles[t]
			break
		}
	}

	checkpoint := saveFolder + "/" + profile + "/Alice2Checkpoint.sav"
	fileInfo, err := os.Stat(checkpoint)
	CheckErr(err)

	files, err = os.ReadDir(saveFolder + "/" + profile)
	CheckErr(err)

	var backupRoll []string
	for _, f := range files {
		if strings.Contains(f.Name(), ".backup") {
			backupRoll = append(backupRoll, f.Name())
		}
	}

	lastModTime := fileInfo.ModTime()
	fmt.Println("Now watching save for changes")
	fmt.Println(checkpoint)
	fmt.Println("Last modified", lastModTime.String())
	fmt.Println("Just close shell or press ctrl+c to exit")

	for {
		fileInfo, err := os.Stat(checkpoint)
		CheckErr(err)

		if fileInfo.ModTime() != lastModTime {
			lastModTime = fileInfo.ModTime()
			fmt.Println("Change detected, copying new backup...")
			newFileName := "Alice2Checkpoint.sav." + strconv.FormatInt(lastModTime.Unix(), 10) + ".backup"
			newFilePath := saveFolder + "/" + profile + "/" + newFileName
			Copy(checkpoint, newFilePath)
			fmt.Println("Copied backup to", newFileName)
			backupRoll = append(backupRoll, newFileName)
		}

		for len(backupRoll) > maxBackups {
			fmt.Println("Deleting old backup: ", backupRoll[0])
			err := os.Remove(saveFolder + "/" + profile + "/" + backupRoll[0])
			CheckErr(err)
			backupRoll = backupRoll[1:]
		}

		time.Sleep(time.Duration(interval) * time.Minute)
	}
}
