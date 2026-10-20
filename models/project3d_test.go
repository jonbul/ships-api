package models

import (
	"math"
	"strings"
	"testing"
)

func validProject() Project3d {
	return Project3d{
		Name: "  Falcon  ",
		Layers: []Layer3d{{Name: "Hull", Visible: true, Parts: []Part{
			{Type: PartBox, Position: Vec3{0, 0, 0}, Scale: Vec3{1, 1, 1}, Color: "#aabbcc"},
		}}},
	}
}

func TestValidateAcceptsAndTrims(t *testing.T) {
	p := validProject()
	if err := p.Validate(); err != nil {
		t.Fatalf("valid project rejected: %v", err)
	}
	if p.Name != "Falcon" {
		t.Errorf("name not trimmed: %q", p.Name)
	}
}

func TestValidateAcceptsEveryPartType(t *testing.T) {
	for _, partType := range []string{
		PartBox, PartSphere, PartCylinder, PartCone, PartWedge, PartTorus,
		PartDome, PartCapsule, PartPyramid, PartTube, PartPrism, PartFrustum,
		PartLight,
	} {
		p := validProject()
		p.Layers[0].Parts[0].Type = partType
		p.Layers[0].Parts[0].Intensity = 2
		if err := p.Validate(); err != nil {
			t.Errorf("%s rejected: %v", partType, err)
		}
	}
}

func TestValidateAcceptsEveryTexture(t *testing.T) {
	for _, texture := range []string{
		"", "panels", "grid", "stripes", "carbon", "hexagons", "worn", "camo", "windows",
	} {
		p := validProject()
		p.Layers[0].Parts[0].Texture = texture
		if err := p.Validate(); err != nil {
			t.Errorf("texture %q rejected: %v", texture, err)
		}
	}
}

func TestValidateAcceptsImageLink(t *testing.T) {
	p := validProject()
	p.Layers[0].Parts[0].TextureURL = "https://i.imgur.com/abc123.png?size=2"
	if err := p.Validate(); err != nil {
		t.Fatalf("image link rejected: %v", err)
	}
}

func TestValidateRejects(t *testing.T) {
	cases := map[string]func(p *Project3d){
		"empty name":        func(p *Project3d) { p.Name = "   " },
		"long name":         func(p *Project3d) { p.Name = strings.Repeat("x", MaxProjectNameLength+1) },
		"no layers":         func(p *Project3d) { p.Layers = nil },
		"unknown type":      func(p *Project3d) { p.Layers[0].Parts[0].Type = "teapot" },
		"bad color":         func(p *Project3d) { p.Layers[0].Parts[0].Color = "red" },
		"unknown texture":   func(p *Project3d) { p.Layers[0].Parts[0].Texture = "marble" },
		"http image link":   func(p *Project3d) { p.Layers[0].Parts[0].TextureURL = "http://a.com/a.png" },
		"spaced image link": func(p *Project3d) { p.Layers[0].Parts[0].TextureURL = "https://a.com/a b" },
		"long image link": func(p *Project3d) {
			p.Layers[0].Parts[0].TextureURL = "https://a.com/" + strings.Repeat("a", MaxTextureURLLength)
		},
		"texture and image link": func(p *Project3d) {
			p.Layers[0].Parts[0].Texture = "panels"
			p.Layers[0].Parts[0].TextureURL = "https://a.com/a.png"
		},
		"far position":   func(p *Project3d) { p.Layers[0].Parts[0].Position[1] = MaxCoordinate + 1 },
		"zero scale":     func(p *Project3d) { p.Layers[0].Parts[0].Scale[2] = 0 },
		"NaN rotation":   func(p *Project3d) { p.Layers[0].Parts[0].Rotation[0] = math.NaN() },
		"too many parts": func(p *Project3d) { p.Layers[0].Parts = make([]Part, MaxParts+1) },
	}
	for name, mutate := range cases {
		t.Run(name, func(t *testing.T) {
			p := validProject()
			mutate(&p)
			if err := p.Validate(); err == nil {
				t.Error("invalid project accepted")
			}
		})
	}
}

func TestValidateNormalisesRotation(t *testing.T) {
	p := validProject()
	p.Layers[0].Parts[0].Rotation = Vec3{5 * math.Pi, 0, -3 * math.Pi}
	if err := p.Validate(); err != nil {
		t.Fatal(err)
	}
	for _, r := range p.Layers[0].Parts[0].Rotation {
		if math.Abs(r) > math.Pi+1e-9 {
			t.Errorf("rotation %v not normalised", r)
		}
	}
}

func TestBoundingRadius(t *testing.T) {
	p := validProject()
	p.Layers[0].Parts = []Part{
		{Type: PartBox, Position: Vec3{3, 4, 0}, Scale: Vec3{2, 2, 2}, Color: "#000000"},
	}
	// |position| = 5, half diagonal of a 2x2x2 cube = sqrt(3).
	want := 5 + math.Sqrt(3)
	if got := p.BoundingRadius(); math.Abs(got-want) > 1e-9 {
		t.Errorf("BoundingRadius = %v, want %v", got, want)
	}
}

func TestBoundingRadiusIgnoresHiddenLayers(t *testing.T) {
	p := validProject()
	p.Layers = append(p.Layers, Layer3d{Visible: false, Parts: []Part{
		{Type: PartBox, Position: Vec3{40, 0, 0}, Scale: Vec3{1, 1, 1}, Color: "#000000"},
	}})
	if r := p.BoundingRadius(); r > 1 {
		t.Errorf("hidden layer counted: radius %v", r)
	}
	p.Layers[0].Visible = false
	p.Layers[1].Visible = false
	if r := p.BoundingRadius(); r != 0 {
		t.Errorf("a ship with nothing visible has radius %v, want 0", r)
	}
}

func TestValidateOpacity(t *testing.T) {
	for opacity, ok := range map[float64]bool{0: true, 0.05: true, 0.5: true, 1: true, 0.01: false, 1.5: false, math.NaN(): false} {
		p := validProject()
		p.Layers[0].Parts[0].Opacity = opacity
		if err := p.Validate(); (err == nil) != ok {
			t.Errorf("opacity %g: err = %v", opacity, err)
		}
	}
}

func TestValidateLights(t *testing.T) {
	light := Part{Type: PartLight, Scale: Vec3{0.4, 0.4, 0.4}, Color: "#fff1c4", Intensity: 2}
	for intensity, ok := range map[float64]bool{0.1: true, 10: true, 0: false, 11: false, math.Inf(1): false} {
		p := validProject()
		part := light
		part.Intensity = intensity
		p.Layers[0].Parts = append(p.Layers[0].Parts, part)
		if err := p.Validate(); (err == nil) != ok {
			t.Errorf("intensity %g: err = %v", intensity, err)
		}
	}

	// A mirrored light counts twice.
	p := validProject()
	mirrored := light
	mirrored.Mirror = true
	for i := 0; i < MaxLights/2; i++ {
		p.Layers[0].Parts = append(p.Layers[0].Parts, mirrored)
	}
	if err := p.Validate(); err != nil {
		t.Fatalf("%d lights rejected: %v", MaxLights, err)
	}
	p.Layers[0].Parts = append(p.Layers[0].Parts, light)
	if err := p.Validate(); err == nil || !strings.Contains(err.Error(), "lights") {
		t.Errorf("too many lights accepted: %v", err)
	}
}

func TestValidateDropsIntensityOffLights(t *testing.T) {
	p := validProject()
	p.Layers[0].Parts[0].Intensity = 3
	if err := p.Validate(); err != nil {
		t.Fatal(err)
	}
	if p.Layers[0].Parts[0].Intensity != 0 {
		t.Errorf("box kept intensity %g", p.Layers[0].Parts[0].Intensity)
	}
}
