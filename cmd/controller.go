package cmd

import (
	"path/filepath"
	"os"
	"strings"
	log "github.com/sirupsen/logrus"
)

type Controller struct {
	Name string
	Package string
	IsAdmin bool
	absPath string
	controllerPath string
	templatesPath string
}

// NewController returns Controller with specified absolute path to
// package.
func NewController(absPath string, name string, pack string, isAdmin bool) *Controller {
	if absPath == "" {
		log.Fatal("can't create controller: absPath can't be blank")
	}
	if !filepath.IsAbs(absPath) {
		log.Fatal("can't create controller: absPath is not absolute")
	}

	// If absPath is symlink, use its destination.
	fi, err := os.Lstat(absPath)
	if err != nil {
		log.Fatalf("can't read path info: %s", err.Error())
	}
	if fi.Mode()&os.ModeSymlink != 0 {
		path, err := os.Readlink(absPath)
		if err != nil {
			log.Fatalf("can't read the destination of symlink: %s", err.Error())
		}
		absPath = path
	}

	p := new(Controller)
	p.absPath = strings.TrimSuffix(absPath, findControllerDir(absPath))
	p.Name = name
	p.Package = name
	if pack != "" {
		p.Package = pack
	}
	p.IsAdmin = isAdmin
	return p
}

// findControllerDir checks if base of absPath is controller dir and returns it or
// looks for existing controller dir in absPath.
func findControllerDir(absPath string) string {
	if !exists(absPath) || isEmpty(absPath) {
		return "controllers"
	}

	if isControllerDir(absPath) {
		return filepath.Base(absPath)
	}

	files, _ := filepath.Glob(filepath.Join(absPath, "c*"))
	for _, file := range files {
		if isControllerDir(file) {
			return filepath.Base(file)
		}
	}

	return "controllers"
}
// findTemplatesDir checks if base of absPath is templates dir and returns it or
// looks for existing templates dir in absPath.
func findTemplatesDir(absPath string, isAdmin bool) string {
	suffix := "site"
	if isAdmin {
		suffix = "admin"
	}
	if !exists(absPath) || isEmpty(absPath) {
		return filepath.Join("templates", suffix)
	}

	if isTemplateDir(absPath) {
		return filepath.Join(filepath.Base(absPath), suffix)
	}

	files, _ := filepath.Glob(filepath.Join(absPath, "t*"))
	for _, file := range files {
		if isTemplateDir(file) {
			return filepath.Join(filepath.Base(file), suffix)
		}
	}

	return filepath.Join("templates", suffix)
}
// isControllerDir checks if base of name is one of controllerDir.
func isControllerDir(name string) bool {
	name = filepath.Base(name)
	for _, controllerDir := range []string{"controllers"} {
		if name == controllerDir {
			return true
		}
	}
	return false
}

// isTemplateDir checks if base of name is one of templateDir.
func isTemplateDir(name string) bool {
	name = filepath.Base(name)
	for _, templateDir := range []string{"templates"} {
		if name == templateDir {
			return true
		}
	}
	return false
}
// AbsPath returns absolute path of controllers.
func (c *Controller) AbsPath() string {
	return c.absPath
}

// ControllerPath returns absolute path to directory, where all controllers are located.
func (c *Controller) ControllerPath() string {
	if c.absPath == "" {
		return ""
	}
	if c.controllerPath == "" {
		c.controllerPath = filepath.Join(c.absPath, findControllerDir(c.absPath))
	}
	return c.controllerPath
}


// TemplatesPath returns absolute path to directory, where all controllers are located.
func (c *Controller) TemplatesPath() string {
	if c.absPath == "" {
		return ""
	}
	if c.templatesPath == "" {
		c.templatesPath = filepath.Join(c.absPath, findTemplatesDir(c.absPath, c.IsAdmin))
	}
	return c.templatesPath
}