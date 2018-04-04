package templates

import (
	"github.com/CloudyKit/jet"
	"github.com/spf13/viper"
	"github.com/gin-gonic/gin"
	"github.com/justinas/nosurf"
	"net/http"
	"path"
	log "github.com/sirupsen/logrus"
)

const (
	templatesDir = "app/templates"
	// AdminDir points to the admin template sub directory
	AdminDir = "admin"
	// SiteDir points to the website template sub directory
	SiteDir = "site"
	// ErrorsDir points to the errors page templates sub directory
	ErrorsDir = "errors"
)
// View : Jet Template Renderer
var View = jet.NewHTMLSet(GetDir("."))

// GetDir return the templates directory path
func GetDir(absPath string) string {
	if absPath == "" {
		absPath = "."
	}
	templatePath := absPath + "/" + templatesDir
	log.Infof("Templates directory: %s", templatePath)
	return  templatePath
}

// Commonvariables return a jet.VarMap variable containing the necessary variables to run index layouts
func Commonvariables(c *gin.Context) jet.VarMap {
	token := nosurf.Token(c.Request)
	variables := templateFunctions(make(jet.VarMap))

	variables.Set("URL", c.Request.URL)
	variables.Set("CsrfToken", token)
	variables.Set("Config", viper.AllSettings())
	// TODO: Errors, Infos, User, Nav, Search
	return variables
}

// Render is a function rendering a template
func Render(c *gin.Context, templateName string, variables jet.VarMap) {
	t, err := View.GetTemplate(templateName)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if err = t.Execute(c.Writer, variables, nil); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
}

// HttpError render an error template
func HttpError(c *gin.Context, errorCode int) {
	switch errorCode {
	case http.StatusNotFound:
		Static(c, path.Join(ErrorsDir, "404.jet.html"))
		c.AbortWithStatus(errorCode)
		return
	case http.StatusBadRequest:
		Static(c, path.Join(ErrorsDir, "400.jet.html"))
		c.AbortWithStatus(errorCode)
		return
	case http.StatusInternalServerError:
		Static(c, path.Join(ErrorsDir, "500.jet.html"))
		c.AbortWithStatus(errorCode)
		return
	}
}

// Static render static templates
func Static(c *gin.Context, templateName string) {
	var variables jet.VarMap
	if isAdminTemplate(templateName) {
		variables = NewPanelCommonvariables(c)
	} else {
		variables = Commonvariables(c)
	}
	Render(c, templateName, variables)
}

// ModelList render list models templates
func ModelList(c *gin.Context, templateName string, models interface{}, nav Navigation, search SearchForm) {
	var variables jet.VarMap
	if isAdminTemplate(templateName) {
		variables = NewPanelCommonvariables(c)
	} else {
		variables = Commonvariables(c)
	}
	variables.Set("Models", models)
	variables.Set("Navigation", nav)
	variables.Set("Search", search)
	Render(c, templateName, variables)
}

// Form render a template form
func Form(c *gin.Context, templateName string, form interface{}) {
	var variables jet.VarMap
	if isAdminTemplate(templateName) {
		variables = NewPanelCommonvariables(c)
	} else {
		variables = Commonvariables(c)
	}
	variables.Set("Form", form)
	Render(c, templateName, variables)
}
// NewPanelSearchForm : Helper that creates a search form without items/page field
// these need to be used when the templateVariables don't include `navigation`
func NewPanelSearchForm() SearchForm {
	form := NewSearchForm()
	form.ShowItemsPerPage = false
	return form
}

// NewPanelCommonvariables return a jet.VarMap variable containing the necessary variables to run index admin layouts
func NewPanelCommonvariables(c *gin.Context) jet.VarMap {
	common := Commonvariables(c)
	common.Set("Search", NewPanelSearchForm())
	return common
}

func isAdminTemplate(templateName string) bool {
	if templateName != "" && len(templateName) > len(AdminDir) {
		return templateName[:5] == AdminDir
	}
	return false
}