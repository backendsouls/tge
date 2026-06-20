//go:build functional

// Package functional contains black-box functional tests that exercise the real
// tge CLI binary inside a container, provisioned with testcontainers.
//
// SQLite is embedded, so there is no database server to containerize; instead we
// containerize the shipped artifact itself and exec real commands against it,
// verifying the binary, the SQLite driver and on-disk persistence end to end.
//
// Run with:  go test -tags functional ./test/functional/...
package functional

import (
	"context"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	tcexec "github.com/testcontainers/testcontainers-go/exec"
	"github.com/testcontainers/testcontainers-go/wait"
)

// startCLIContainer builds the CLI image from the repo Dockerfile and starts a
// long-lived container so multiple commands can share one SQLite database.
func startCLIContainer(t *testing.T) testcontainers.Container {
	t.Helper()
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		FromDockerfile: testcontainers.FromDockerfile{
			Context:    "../..", // repo root holds the Dockerfile and sources
			Dockerfile: "Dockerfile",
			KeepImage:  true,
		},
		WaitingFor: wait.ForExec([]string{"true"}).WithStartupTimeout(3 * time.Minute),
	}
	c, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("start container: %v", err)
	}
	t.Cleanup(func() { _ = testcontainers.TerminateContainer(c) })
	return c
}

// execCLI runs `tge <args...>` inside the container, returning the exit code and
// the combined stdout/stderr output.
func execCLI(t *testing.T, c testcontainers.Container, args ...string) (int, string) {
	t.Helper()
	code, reader, err := c.Exec(context.Background(), append([]string{"tge"}, args...), tcexec.Multiplexed())
	if err != nil {
		t.Fatalf("exec %v: %v", args, err)
	}
	out, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("read output for %v: %v", args, err)
	}
	return code, string(out)
}

func TestFunctional_MultiCharacterFlow(t *testing.T) {
	c := startCLIContainer(t)

	// Two power systems, one with a nested power tree.
	for _, sys := range []string{"Universe A Cultivation", "Universe B Sorcery"} {
		if code, out := execCLI(t, c, "powersystem", "add", "--name", sys); code != 0 {
			t.Fatalf("add system %q exit = %d, out:\n%s", sys, code, out)
		}
	}
	if code, out := execCLI(t, c, "powersystem", "add-power", "--system", "Universe A Cultivation", "--name", "Body"); code != 0 {
		t.Fatalf("add-power exit = %d, out:\n%s", code, out)
	}

	// A female MainCharacter in two systems (repeatable --system).
	if code, out := execCLI(t, c, "character", "create", "--name", "Lin Feng", "--type", "MainCharacter", "--gender", "Female",
		"--system", "Universe A Cultivation", "--system", "Universe B Sorcery"); code != 0 {
		t.Fatalf("create MC exit = %d, out:\n%s", code, out)
	}
	// A Hero is allowed because the MC is female.
	if code, out := execCLI(t, c, "character", "create", "--name", "Yun Hai", "--type", "Hero", "--gender", "Male",
		"--system", "Universe A Cultivation"); code != 0 {
		t.Fatalf("create Hero exit = %d, out:\n%s", code, out)
	}

	// status shows the main character only, with both systems (proves persistence).
	code, out := execCLI(t, c, "status")
	if code != 0 {
		t.Fatalf("status exit = %d, out:\n%s", code, out)
	}
	for _, want := range []string{"Lin Feng", "MainCharacter", "Female", "Universe A Cultivation", "Universe B Sorcery"} {
		if !strings.Contains(out, want) {
			t.Errorf("status missing %q:\n%s", want, out)
		}
	}

	// character list shows the full roster.
	code, out = execCLI(t, c, "character", "list")
	if code != 0 {
		t.Fatalf("list exit = %d, out:\n%s", code, out)
	}
	if !strings.Contains(out, "Lin Feng") || !strings.Contains(out, "Yun Hai") {
		t.Errorf("roster missing a character:\n%s", out)
	}
}

func TestFunctional_HeroineRejectedForFemaleMainCharacter(t *testing.T) {
	c := startCLIContainer(t)

	if code, out := execCLI(t, c, "powersystem", "add", "--name", "Universe A Cultivation"); code != 0 {
		t.Fatalf("add system exit = %d, out:\n%s", code, out)
	}
	if code, out := execCLI(t, c, "character", "create", "--name", "Lin Feng", "--type", "MainCharacter", "--gender", "Female",
		"--system", "Universe A Cultivation"); code != 0 {
		t.Fatalf("create MC exit = %d, out:\n%s", code, out)
	}
	// Heroine requires a male MC, but the MC is female.
	code, out := execCLI(t, c, "character", "create", "--name", "Bai Li", "--type", "Heroine", "--gender", "Female",
		"--system", "Universe A Cultivation")
	if code == 0 {
		t.Fatalf("expected non-zero exit for Heroine with female MC, out:\n%s", out)
	}
	if !strings.Contains(out, "Heroine requires") {
		t.Errorf("expected role error, got:\n%s", out)
	}
}

func TestFunctional_Novels(t *testing.T) {
	c := startCLIContainer(t)

	if code, out := execCLI(t, c, "powersystem", "add", "--name", "Universe A Cultivation"); code != 0 {
		t.Fatalf("add system exit = %d, out:\n%s", code, out)
	}
	// Two main characters, each able to lead a novel.
	if code, out := execCLI(t, c, "character", "create", "--name", "Lin Feng", "--type", "MainCharacter", "--gender", "Female",
		"--system", "Universe A Cultivation"); code != 0 {
		t.Fatalf("create MC1 exit = %d, out:\n%s", code, out)
	}
	if code, out := execCLI(t, c, "character", "create", "--name", "Mu Chen", "--type", "MainCharacter", "--gender", "Male",
		"--system", "Universe A Cultivation"); code != 0 {
		t.Fatalf("create MC2 exit = %d, out:\n%s", code, out)
	}

	if code, out := execCLI(t, c, "novel", "add", "--title", "Ascension", "--main-character", "Lin Feng"); code != 0 {
		t.Fatalf("add novel exit = %d, out:\n%s", code, out)
	}
	// The same main character can't lead a second novel.
	if code, out := execCLI(t, c, "novel", "add", "--title", "Another", "--main-character", "Lin Feng"); code == 0 {
		t.Fatalf("expected rejection for taken main character, out:\n%s", out)
	}

	if code, out := execCLI(t, c, "novel", "add", "--title", "Sorcerer's Path", "--main-character", "Mu Chen"); code != 0 {
		t.Fatalf("add second novel exit = %d, out:\n%s", code, out)
	}

	code, out := execCLI(t, c, "novel", "list")
	if code != 0 {
		t.Fatalf("novel list exit = %d, out:\n%s", code, out)
	}
	for _, want := range []string{"Ascension", "Lin Feng", "Sorcerer's Path", "Mu Chen"} {
		if !strings.Contains(out, want) {
			t.Errorf("novel list missing %q:\n%s", want, out)
		}
	}

	// Volumes and chapters, surviving across processes (SQLite persistence).
	if code, out := execCLI(t, c, "novel", "add-volume", "--novel", "Ascension", "--number", "1", "--title", "Beginnings"); code != 0 {
		t.Fatalf("add-volume exit = %d, out:\n%s", code, out)
	}
	if code, out := execCLI(t, c, "novel", "add-chapter", "--novel", "Ascension", "--volume", "1", "--number", "1", "--title", "A Mortal's Dream"); code != 0 {
		t.Fatalf("add-chapter exit = %d, out:\n%s", code, out)
	}
	code, out = execCLI(t, c, "novel", "show", "--title", "Ascension")
	if code != 0 {
		t.Fatalf("novel show exit = %d, out:\n%s", code, out)
	}
	for _, want := range []string{"Volume 1: Beginnings", "Chapter 1: A Mortal's Dream"} {
		if !strings.Contains(out, want) {
			t.Errorf("novel show missing %q:\n%s", want, out)
		}
	}
}

func TestFunctional_Universes(t *testing.T) {
	c := startCLIContainer(t)

	for _, sys := range []string{"Cultivation", "Sorcery"} {
		if code, out := execCLI(t, c, "powersystem", "add", "--name", sys); code != 0 {
			t.Fatalf("add system %q exit = %d, out:\n%s", sys, code, out)
		}
	}
	if code, out := execCLI(t, c, "universe", "add", "--name", "Universe A"); code != 0 {
		t.Fatalf("add universe exit = %d, out:\n%s", code, out)
	}
	if code, out := execCLI(t, c, "universe", "add-system", "--universe", "Universe A", "--system", "Cultivation"); code != 0 {
		t.Fatalf("add-system exit = %d, out:\n%s", code, out)
	}

	// A second universe can't claim a system that already belongs to another.
	if code, out := execCLI(t, c, "universe", "add", "--name", "Universe B"); code != 0 {
		t.Fatalf("add universe B exit = %d, out:\n%s", code, out)
	}
	if code, out := execCLI(t, c, "universe", "add-system", "--universe", "Universe B", "--system", "Cultivation"); code == 0 {
		t.Fatalf("expected rejection for taken system, out:\n%s", out)
	}

	// Membership survives across processes (persistence).
	code, out := execCLI(t, c, "universe", "show", "--name", "Universe A")
	if code != 0 {
		t.Fatalf("show exit = %d, out:\n%s", code, out)
	}
	if !strings.Contains(out, "Cultivation") {
		t.Errorf("universe show missing system:\n%s", out)
	}
}

func TestFunctional_CharacterRequiresExistingSystem(t *testing.T) {
	c := startCLIContainer(t)

	code, out := execCLI(t, c, "character", "create", "--name", "Lin Feng", "--type", "SideCharacter", "--gender", "Male", "--system", "Nonexistent")
	if code == 0 {
		t.Fatalf("expected non-zero exit for missing system, out:\n%s", out)
	}
	if !strings.Contains(out, "not found") {
		t.Errorf("expected not-found error, got:\n%s", out)
	}
}
