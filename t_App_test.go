package mango

import "testing"

func Test_NewApplication(t *testing.T) {
	app, err := NewApplication()
	// app.Print()

	if err != nil {
		t.Fatal(app, err)
	}

	// Check default paths
	if app.binPath == "" {
		t.Fatal("binPath is empty")
	}

	// Trim from binPath end to content path end
	if app.ContentPath[len(app.binPath):] != "/test-files/content" {
		t.Fatal("Default ContentPath must end with /test-files/content")
	}
	// Trim from binPath end to public path end
	if app.PublicPath[len(app.binPath):] != "/test-files/public" {
		t.Fatal("Default PublicPath must end with /test-files/public")
	}

}

func Test_NewApplicationFuncs(t *testing.T) {
	app, _ := NewApplication()

	// new virtual page
	// linked to app, but no parents
	// not listed anywhere
	p := app.NewPage("Virtual reality!")
	p.Set("Custom", "param")
	if p.Params["IsVirtual"] != "Yes" ||
		p.Params["Label"] != "Virtual reality!" ||
		p.Params["VirtualSlug"] != "virtual-reality" ||
		p.App == nil {
		p.Print()
		t.Fatal("Labeled NewPage")
	}

	// Empty label
	p = app.NewPage("")
	p.Set("Custom", "param")
	if p.Params["IsVirtual"] != "Yes" ||
		p.Params["VirtualSlug"] != "" ||
		p.App == nil {
		p.Print()
		t.Fatal("Empty labeled NewPage")
	}

}
