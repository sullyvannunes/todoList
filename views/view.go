package views

import (
	"fmt"
	"html/template"
	"io"
	"os"
	"strings"
)

type ViewConfig struct {
	Namespace             string
	ComponentDir          string
	IsProductionMode      bool
	LayoutDefault         string
	TemplateFileExtension string
}

type View struct {
	namespace             string
	componentDir          string
	isProductionMode      bool
	layoutDefault         string
	templateFileExtension string
}

func New(config ViewConfig) *View {
	namespace := config.Namespace + "/"
	templateFileExtension := ".html"
	if config.TemplateFileExtension != "" {
		templateFileExtension = config.TemplateFileExtension
	}

	return &View{
		namespace:             namespace,
		componentDir:          namespace + config.ComponentDir,
		isProductionMode:      config.IsProductionMode,
		layoutDefault:         config.LayoutDefault,
		templateFileExtension: templateFileExtension,
	}
}

func (v *View) Render(w io.Writer, tmplName string, data any) error {
	return v.RenderWithLayout(w, v.layoutDefault, tmplName, data)
}

func (v *View) RenderWithLayout(w io.Writer, layoutName, tmplName string, data any) error {
	fs := os.DirFS(v.componentDir)
	comTmpls, err := template.ParseFS(fs, "*"+v.templateFileExtension)
	if err != nil {
		return fmt.Errorf("parse component templates: %w", err)
	}

	tmplPath := strings.Split(tmplName, "/")
	fs = os.DirFS(v.namespace + tmplPath[0])

	partialTmpls, err := template.ParseFS(fs, "_*"+v.templateFileExtension)
	if err != nil {
		return fmt.Errorf("parse component templates: %w", err)
	}

	t, err := template.ParseFiles(v.namespace + tmplName + v.templateFileExtension)
	if err != nil {
		return fmt.Errorf("parse template file: %w", err)
	}

	newTmpl, err := v.mergeTemplates(t, comTmpls)
	if err != nil {
		return fmt.Errorf("merge template and components: %w", err)
	}

	newTmpl, err = v.mergeTemplates(newTmpl, partialTmpls)
	if err != nil {
		return fmt.Errorf("merge template and partials: %w", err)
	}

	layout, err := template.ParseFiles(v.namespace + layoutName + v.templateFileExtension)
	if err != nil {
		return fmt.Errorf("parse layout template: %w", err)
	}

	finalTmpl, err := v.mergeTemplates(layout, newTmpl)
	if err != nil {
		return fmt.Errorf("merge template and layout: %w", err)
	}

	if err := finalTmpl.Execute(w, data); err != nil {
		return fmt.Errorf("execute final template: %w", err)
	}

	return nil
}

func (v *View) mergeTemplates(dst, src *template.Template) (*template.Template, error) {
	dstClone, err := dst.Clone()
	if err != nil {
		return nil, fmt.Errorf("clone dst template: %w", err)
	}

	for _, t := range src.Templates() {
		_, err := dstClone.AddParseTree(t.Name(), t.Tree)
		if err != nil {
			return nil, fmt.Errorf("add parse tree: %w", err)
		}
	}

	return dstClone, nil
}
