package runner

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/kanister10l/GoCraft/db"
	"github.com/kanister10l/GoCraft/eventmanager"
	"github.com/kanister10l/GoCraft/logger"
	"github.com/ztrue/tracerr"
)

// Runner ...
type Runner struct {
	Connection   *db.Connection
	ServerID     int
	Name         string
	Port         string
	Image        string
	EventManager *eventmanager.Manager
}

// NewRunner ...
func NewRunner(conn *db.Connection, id int, port, name, image string, manager *eventmanager.Manager) *Runner {
	runner := &Runner{
		Connection:   conn,
		ServerID:     id,
		Name:         name,
		Port:         port,
		Image:        image,
		EventManager: manager,
	}

	return runner
}

// StartServers ...
func StartServers(c *db.Connection, manager *eventmanager.Manager) error {
	rows, err := c.Database.Query(db.SelectServers)
	if err != nil {
		tracerr.PrintSourceColor(tracerr.Wrap(err))
		return err
	}

	var id int
	var image string
	var name string
	var port string
	for rows.Next() {
		rows.Scan(&id, &image, &name, &port)
		NewRunner(c, id, port, name, image, manager).Run()
	}

	return nil
}

// GetVolume ...
func (r *Runner) GetVolume() (string, error) {
	rows, err := r.Connection.Database.Query(fmt.Sprintf(db.SelectVolumesByServer, r.ServerID))
	if err != nil {
		tracerr.PrintSourceColor(tracerr.Wrap(err))
		return "", err
	}

	var id int
	var name string
	var location string
	rows.Next()
	rows.Scan(&id, &name, &location)

	return location, nil
}

// CopyFile ...
func (r *Runner) CopyFile(source, target string) error {
	input, err := ioutil.ReadFile(source)
	if err != nil {
		tracerr.PrintSourceColor(tracerr.Wrap(err))
		return err
	}

	err = ioutil.WriteFile(target, input, 0644)
	if err != nil {
		tracerr.PrintSourceColor(tracerr.Wrap(err))
		return err
	}

	return nil
}

// PrepareEnv ...
func (r *Runner) PrepareEnv() (string, error) {
	dir, err := ioutil.TempDir("", "gocraft")
	logger.Logger.Infow("new tmp dir created", "location", dir)
	if err != nil {
		tracerr.PrintSourceColor(tracerr.Wrap(err))
		return "", err
	}

	if err := r.CopyFile("./_generic/eula.txt", filepath.Join(dir, "eula.txt")); err != nil {
		return "", err
	}

	if err := r.CopyFile("./_generic/server.properties", filepath.Join(dir, "server.properties")); err != nil {
		return "", err
	}

	if err := r.CopyFile("./_generic/startscript.sh", filepath.Join(dir, "startscript.sh")); err != nil {
		return "", err
	}

	if err := r.CopyFile("./_images/"+r.Image, filepath.Join(dir, "Dockerfile")); err != nil {
		return "", err
	}

	return dir, nil
}

// Build ...
func (r *Runner) Build(dir string) error {
	cmd := exec.Command("docker", "build", "-t", r.Name, dir)
	err := cmd.Run()
	if err != nil {
		tracerr.PrintSourceColor(tracerr.Wrap(err))
		return err
	}

	logger.Logger.Infow("successful docker build", "image", r.Name)

	return nil
}

// GetContainerParameterString ...
func (r *Runner) GetContainerParameterString() (string, error) {
	rows, err := r.Connection.Database.Query(fmt.Sprintf(db.SelectContainerParametersByServer, r.ServerID))
	if err != nil {
		tracerr.PrintSourceColor(tracerr.Wrap(err))
		return "", err
	}

	pstr := ""

	var id int
	var key string
	var value string
	for rows.Next() {
		rows.Scan(&id, &key, &value)
		switch key {
		case "cpu":
			pstr += "--cpus=" + value + " "
		case "memory":
			pstr += "--memory=" + value + " "
		}
	}
	return strings.TrimRight(pstr, " "), nil
}

// Run ...
func (r *Runner) Run() error {
	dir, err := r.PrepareEnv()
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)

	if r.Build(dir) != nil {
		return err
	}

	vloc, err := r.GetVolume()
	if err != nil {
		return err
	}

	logger.Logger.Infow("docker run", "volume", vloc)

	settings, err := r.GetContainerParameterString()
	if err != nil {
		return err
	}

	variadicSettings := []string{"run", "-dit", "--name", r.Name, "--rm", "-v", vloc + ":/mc", "-p", r.Port + ":25565"}
	variadicSettings = append(variadicSettings, strings.Split(settings, " ")...)
	variadicSettings = append(variadicSettings, r.Name)

	log.Println(variadicSettings)

	cmd := exec.Command("docker", variadicSettings...)
	err = cmd.Run()
	if err != nil {
		tracerr.PrintSourceColor(tracerr.Wrap(err))
		return err
	}

	logger.Logger.Infow("successful docker run", "container", r.Name)

	cname := r.Name
	err = r.EventManager.Register("stop", func() {
		cmd := exec.Command("docker", "kill", cname)
		cmd.Run()
	})

	if err != nil {
		tracerr.PrintSourceColor(tracerr.Wrap(err))
		return err
	}

	return nil
}
