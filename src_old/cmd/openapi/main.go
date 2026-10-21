// Command openapi genererer dokumentationen uden at starte serveren eller åbne databasen.
// Kør fra repository-roden: go -C src run ./cmd/openapi
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

const generator = "github.com/swaggo/swag/v2/cmd/swag@v2.0.0-rc5"

func main() {
	if err := generate(); err != nil {
		fmt.Fprintln(os.Stderr, "OpenAPI:", err)
		os.Exit(1)
	}
}

func generate() error {
	if _, err := os.Stat("main.go"); err != nil {
		return fmt.Errorf("kør fra roden med go -C src run ./cmd/openapi: %w", err)
	}
	// Rå output er midlertidigt; eksisterende spec bevares ved generatorfejl.
	dir, err := os.MkdirTemp("", "whoknows-openapi-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	cmd := exec.Command("go", "run", generator, "init", "--v3.1", "--parseInternal",
		"--outputTypes", "json", "--generalInfo", "main.go", "--exclude", "cmd", "--output", dir)
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	raw, err := os.ReadFile(filepath.Join(dir, "swagger.json"))
	if err != nil {
		return err
	}
	var doc map[string]any
	if err := json.Unmarshal(raw, &doc); err != nil {
		return err
	}
	if err := normalize(doc); err != nil {
		return err
	}
	data, err := json.MarshalIndent(doc, "", "    ")
	if err != nil {
		return err
	}
	output := filepath.Join("..", "docs", "openapi", "swagger.json")
	if err := os.MkdirAll(filepath.Dir(output), 0755); err != nil {
		return err
	}
	if err := os.WriteFile(output, append(data, '\n'), 0644); err != nil {
		return err
	}
	fmt.Println("Genereret:", output)
	return nil
}

// normalize udfylder kontraktens union-typer, som Swag rc5 ikke udleder af Go-typerne.
// Runtime sender konkrete værdier og string-lokationer; kontrakten tillader
// desuden null og numeriske array-indekser. Handlernes adfærd ændres ikke.
func normalize(doc map[string]any) error {
	if doc["openapi"] != "3.1.0" {
		return fmt.Errorf("forventede OpenAPI 3.1.0, fik %v", doc["openapi"])
	}
	for _, field := range []struct{ schema, name, kind string }{
		{"delivery.AuthResponse", "statusCode", "integer"},
		{"delivery.AuthResponse", "message", "string"},
		{"delivery.RequestValidationError", "message", "string"},
	} {
		property, err := object(doc, "components", "schemas", field.schema, "properties", field.name)
		if err != nil {
			return err
		}
		property["type"] = []string{field.kind, "null"}
	}
	items, err := object(doc, "components", "schemas", "delivery.ValidationErrorDetail", "properties", "loc", "items")
	if err != nil {
		return err
	}
	items["type"] = []string{"string", "integer"}
	search, err := object(doc, "paths", "/api/search", "get")
	if err != nil {
		return err
	}
	parameters, ok := search["parameters"].([]any)
	if !ok {
		return fmt.Errorf("søgeparametre mangler")
	}
	found := false
	for _, value := range parameters {
		parameter, ok := value.(map[string]any)
		if !ok || parameter["name"] != "language" {
			continue
		}
		schema, err := object(parameter, "schema")
		if err != nil {
			return err
		}
		schema["type"] = []string{"string", "null"}
		// Default(en) bevares som dokumentation af vores faktiske implementation.
		found = true
	}
	if !found {
		return fmt.Errorf("language-parameteren mangler")
	}
	// rc5 opretter externalDocs uden en reel dokumentationsadresse.
	if external, ok := doc["externalDocs"].(map[string]any); ok && external["url"] == "" {
		delete(doc, "externalDocs")
	}
	// Formularnavne kommer fra json-tags; rc5 lækker også form-tags ind i schemaet.
	login, err := object(doc, "components", "schemas", "delivery.LoginRequest", "properties")
	if err != nil {
		return err
	}
	for _, value := range login {
		if field, ok := value.(map[string]any); ok {
			delete(field, "form")
		}
	}
	return nil
}

// Fejl ved ændrede modelnavne skal være synlige, ikke give en delvist rettet fil.
func object(root map[string]any, keys ...string) (map[string]any, error) {
	current := root
	for _, key := range keys {
		next, ok := current[key].(map[string]any)
		if !ok {
			return nil, fmt.Errorf("mangler objekt %q i %v", key, keys)
		}
		current = next
	}
	return current, nil
}
