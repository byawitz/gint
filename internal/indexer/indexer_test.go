package indexer

import (
	"fmt"
	Copy "github.com/otiai10/copy"
	"math/rand/v2"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

func TestIndexerAllFiles(t *testing.T) {
	tmpEnvFolderName := before(t)

	files := GetFiles([]string{tmpEnvFolderName}, false, nil)

	allFiles := []string{
		filepath.Join(tmpEnvFolderName, "app/index.php"),
		filepath.Join(tmpEnvFolderName, "conf/config.php"),
	}

	if len(files) != len(allFiles) {
		t.Fatalf("expected %v PHP files, but got %v files", len(allFiles), len(files))
	}

	if !slices.Equal(files, allFiles) {
		t.Fatalf("expected to get the files %v\n, but got %v files", allFiles, files)
	}

	after(t, tmpEnvFolderName)
}

func TestIndexerNotGitRepo(t *testing.T) {
	tmpEnvFolderName := before(t)

	defer func(folder string) {
		after(t, folder)

		if err := recover(); err != nil {
			expect := fmt.Sprintf("error while checking if there are changed PHP files at: %s", folder)
			err, ok := err.(string)

			if !ok || strings.Contains(expect, err) {
				t.Fatalf("expected to get %v but got %v files", err, expect)
			}
		}
	}(tmpEnvFolderName)

	GetFiles([]string{tmpEnvFolderName}, true, nil)
}

func TestIndexerGitDirty(t *testing.T) {
	tmpEnvFolderName := before(t)
	initAndCommit(t, tmpEnvFolderName)

	files := GetFiles([]string{tmpEnvFolderName}, true, nil)

	allFiles := []string{
		filepath.Join(tmpEnvFolderName, "conf/config.php"),
	}

	if len(files) != len(allFiles) {
		t.Fatalf("expected %v PHP files, but got %v files", len(allFiles), len(files))
	}

	if !slices.Equal(files, allFiles) {
		t.Fatalf("expected to get the files \n%v\nbut got \n%v\n", allFiles, files)
	}

	after(t, tmpEnvFolderName)
}
func before(t *testing.T) string {
	tmpEnvName := fmt.Sprintf("gint_%f", rand.Float64())
	tmpEnvFolderName := filepath.Join(os.TempDir(), tmpEnvName)

	err := os.Mkdir(tmpEnvFolderName, 0777)

	if err != nil {
		t.Fatal("unable to create tmp environment")
	}
	currentPath, err := os.Getwd()

	if err != nil {
		t.Fatal("unable to get current path")
	}

	err = Copy.Copy(filepath.Join(currentPath, "../../tests_assets/demo-project"), tmpEnvFolderName)

	if err != nil {
		t.Fatalf("unable to copy tmp environment %v", err)
	}
	return tmpEnvFolderName
}

func after(t *testing.T, folderName string) {
	err := os.RemoveAll(folderName)
	if err != nil {
		t.Fatal("unable to delete tmp environment")
	}
}

func initAndCommit(t *testing.T, folder string) {
	runCommand(t, folder, "git", "unable to initialize git repository", []string{"init"})
	runCommand(t, folder, "git", "unable to set git name", []string{"config", "user.name", "\"Test name\""})
	runCommand(t, folder, "git", "unable to set git email", []string{"config", "user.email", "\"test@test.com\""})
	runCommand(t, folder, "git", "unable to add file to git", []string{"add", "./app/index.php"})
	runCommand(t, folder, "git", "unable to commit file to git", []string{"commit", "-m", "\"Initial commit\""})
}

func runCommand(t *testing.T, folder, command, errorString string, args []string) {
	cmd := exec.Command(command, args...)
	cmd.Dir = folder

	err := cmd.Run()
	if err != nil {
		t.Fatalf("%v, Error:%v", errorString, err)
	}
}
