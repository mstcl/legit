package routes

import (
	"os"
	"path/filepath"

	"git.icyphox.sh/legit/git"
)

func isGoModule(gr *git.GitRepo) bool {
	_, err := gr.FileContent("go.mod")
	return err == nil
}

func getDescription(path string) (desc string) {
	db, err := os.ReadFile(filepath.Join(path, "description"))
	if err == nil {
		desc = defaultDescription(string(db))
	} else {
		desc = ""
	}
	return
}

func defaultDescription(desc string) string {
	if desc == "Unnamed repository; edit this file 'description' to name the repository.\n" {
		return "No description."
	}
	return desc
}

func (d *deps) isIgnored(name string) bool {
	for _, i := range d.c.Repo.Ignore {
		if name == i {
			return true
		}
	}

	return false
}
