// Copyright © 2018 akuma06 <contact@akuma06.tk>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"github.com/spf13/cobra"
	log "github.com/sirupsen/logrus"
	"unicode"
	"github.com/spf13/viper"
	"time"
	"os"
	"path/filepath"
	"bufio"
	"strings"
)

var parentName, urlController string
var isAdmin, generateTemplate bool

// newControllerCmd represents the newController command
var newControllerCmd = &cobra.Command{
	Use:   "newController",
	Short: "Create a new controller",
	Long: `This command can be used to create a new controller.
It removes a part of the tedious and robotic way when adding a 
route to the website. Appending flags to it can allow to generate 
the view and define the main path.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			log.Fatal("newController needs a name for the controller")
		}
		wd, err := os.Getwd()
		if err != nil {
			log.Fatal(err.Error())
		}
		wd = filepath.Join(wd, "app")
		if !exists(wd) {
			log.Fatal("this command needs to be executed from DokoKai main directory")
		}
		controllerName := validateControllerName(args[0])
		controller := NewController(wd, controllerName, parentName, isAdmin)
		controllerPath := filepath.Join(controller.ControllerPath(), controller.Package, controller.Name+".go")
		createControllerFile(controller, controllerPath)
		if generateTemplate {
			templatesPath := filepath.Join(controller.TemplatesPath(), controller.Package, controller.Name+".jet.html")
			createTemplateFile(controller, templatesPath)
		}
	},
}

func init() {
	rootCmd.AddCommand(newControllerCmd)

	newControllerCmd.Flags().BoolVarP(&isAdmin, "isAdmin", "a", false, "Defines if a controller is for Admin (Default: false)")
	newControllerCmd.Flags().BoolVarP(&generateTemplate, "template", "t", true, "Should generate a basic template (Default: true)")
	newControllerCmd.Flags().StringVarP(&parentName, "parent", "p", "", "Parent package (Default: controller name)")
	newControllerCmd.Flags().StringVarP(&urlController, "url", "u", "", "URL (Default: /controller name)")
}

// validateControllerName returns source without any dashes and underscore.
// If there will be dash or underscore, next letter will be uppered.
// It supports only ASCII (1-byte character) strings.
// https://github.com/spf13/cobra/issues/269
func validateControllerName(source string) string {
		i := 0
		l := len(source)
		// The output is initialized on demand, then first dash or underscore
		// occurs.
		var output string

		for i < l {
		if source[i] == '-' || source[i] == '_' {
		if output == "" {
		output = source[:i]
	}

		// If it's last rune and it's dash or underscore,
		// don't add it output and break the loop.
		if i == l-1 {
		break
	}

		// If next character is dash or underscore,
		// just skip the current character.
		if source[i+1] == '-' || source[i+1] == '_' {
		i++
		continue
	}

		// If the current character is dash or underscore,
		// upper next letter and add to output.
		output += string(unicode.ToUpper(rune(source[i+1])))
		// We know, what source[i] is dash or underscore and source[i+1] is
		// uppered character, so make i = i+2.
		i += 2
		continue
	}

		// If the current character isn't dash or underscore,
		// just add it.
		if output != "" {
		output += string(source[i])
	}
		i++
	}

		if output == "" {
		return source // source is initially valid name.
	}
	return output
}

func createControllerFile(controller *Controller, path string) {
	template := `{{comment .copyright}}
{{if .license}}{{comment .license}}{{end}}
package {{.parentName}}Controller
import (
	{{ if .generateTemplate }}
	"github.com/akuma06/DokoKai/app/templates"
	{{ end }}
	"github.com/gin-gonic/gin"
)

// {{.controllerName}}Handler : Show {{.controllerName}}
func {{.controllerName}}Handler(c *gin.Context) {
	{{ if .generateTemplate }}
	templates.Static(c, "{{ if .isAdmin }}admin{{ else }}site{{ end }}/{{.parentName}}/{{.controllerName}}.jet.html")
	{{ else }}
	
	{{ end }}
}
`
	data := make(map[string]interface{})
	data["copyright"] = copyrightLine()
	data["license"] = MIT.Header
	data["parentName"] = controller.Package
	data["isAdmin"] = controller.IsAdmin
	data["generateTemplate"] = generateTemplate
	data["controllerName"] = controller.Name

	controllerScript, err := executeTemplate(template, data)
	if err != nil {
		log.Fatal(err.Error())
	}
	err = writeStringToFile(path, controllerScript)
	if err != nil {
		log.Fatal(err.Error())
	}
	routerPath := filepath.Join(controller.ControllerPath(), controller.Package, "router.go")
	if !exists(routerPath) {
		createRoute(controller, routerPath)
		routesPath := filepath.Join(controller.AbsPath(), "router.go")
		addTo(routesPath, "import", `    _ "github.com/akuma06/DokoKai/app/controllers/` +  controller.Package + `" // ` +  controller.Package + `controller`)
		return
	}
	addTo(routerPath, "init()", `    controllers.Get().Any("/` + urlController + `", ` + controller.Name + `Handler)`)
}

func createTemplateFile(controller *Controller, path string) {
	template := `{{"{{"}} extends "layouts/index_{{ if .isAdmin }}admin{{ else }}site{{ end }}" {{"}}"}}
{{"{{"}}block title(){{"}}"}}{{ .controllerName }}{{"{{"}}end{{"}}"}}
{{"{{"}}block content_body(){{"}}"}}
Your content goes here.
{{"{{"}}end{{"}}"}}
`
	data := make(map[string]interface{})
	data["isAdmin"] = controller.IsAdmin
	data["parentName"] = controller.Package
	data["controllerName"] = controller.Name

	controllerScript, err := executeTemplate(template, data)
	if err != nil {
		log.Fatal(err.Error())
	}
	err = writeStringToFile(path, controllerScript)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func addTo(path string, lookfor string, toappend string) {
	if (!exists(path)) {
		log.Fatalf("couldn't add route: cannot find %s", path)
	}
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var modifiedFile []string
	for scanner.Scan() {
		line := scanner.Text()
		if pos := strings.Index(line, lookfor); pos > -1 {
			modifiedFile = append(modifiedFile, line)
			modifiedFile = append(modifiedFile, toappend)
			continue
		}
		modifiedFile = append(modifiedFile, line)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	err = writeStringToFile(path, strings.Join(modifiedFile, "\n"))
	if err != nil {
		log.Fatal(err)
	}
}

func copyrightLine() string {
	author := viper.GetString("author")

	year := viper.GetString("year") // For tests.
	if year == "" {
		year = time.Now().Format("2006")
	}

	return "Copyright © " + year + " " + author
}

func createRoute(controller *Controller, path string) {
	template := `{{comment .copyright}}
{{if .license}}{{comment .license}}{{end}}
package {{.parentName}}Controller
import (
	"github.com/akuma06/DokoKai/app/controllers"
)

func init() {
	controllers.Get().Any("{{.url}}", {{.controllerName}}Handler)
}
`
	data := make(map[string]interface{})
	data["copyright"] = copyrightLine()
	data["license"] = MIT.Header
	data["parentName"] = controller.Package
	data["controllerName"] = controller.Name
	if urlController == "" {
		urlController = "/" + controller.Name
	}
	data["url"] = urlController

	controllerScript, err := executeTemplate(template, data)
	if err != nil {
		log.Fatal(err.Error())
	}
	err = writeStringToFile(path, controllerScript)
	if err != nil {
		log.Fatal(err.Error())
	}
}
