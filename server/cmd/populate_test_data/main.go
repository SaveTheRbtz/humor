package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"cloud.google.com/go/firestore"
	serverImpl "github.com/SaveTheRbtz/humor/server/internal/server"
)

func main() {
	ctx := context.Background()

	// Ensure FIRESTORE_EMULATOR_HOST is set
	if os.Getenv("FIRESTORE_EMULATOR_HOST") == "" {
		log.Fatal("FIRESTORE_EMULATOR_HOST is not set")
	}

	// Connect to the emulator
	client, err := firestore.NewClient(ctx, "humor-arena")
	if err != nil {
		log.Fatalf("Failed to create Firestore client: %v", err)
	}
	defer client.Close()

	rand.Seed(time.Now().UnixNano())

	jokes := []struct {
		Theme string
		Text  string
	}{
		{"Animals", "Why do cows have hooves instead of feet? Because they lactose."},
		{"Animals", "What do you call a fish wearing a bowtie? Sofishticated."},
		{"Animals", "Why did the chicken join a band? Because it had the drumsticks."},
		{"Technology", "Why do programmers prefer dark mode? Because light attracts bugs."},
		{"Technology", "How many programmers does it take to change a light bulb? None, it's a hardware problem."},
		{"Food", "I'm on a seafood diet. I see food and I eat it."},
		{"Food", "Why did the tomato turn red? Because it saw the salad dressing."},
	}

	themes := make(map[string]struct{})
	for _, joke := range jokes {
		themes[joke.Theme] = struct{}{}
	}

	for themeName := range themes {
		rnd := rand.Float64()
		err := addTheme(ctx, client, themeName, rnd)
		if err != nil {
			log.Fatalf("Failed to add theme %s: %v", themeName, err)
		}
		fmt.Printf("Added theme: %f: %s\n", rnd, themeName)
	}

	for _, joke := range jokes {
		rnd := rand.Float64()
		err := addJoke(ctx, client, joke.Theme, joke.Text, rnd)
		if err != nil {
			log.Fatalf("Failed to add joke: %v", err)
		}
		fmt.Printf("Added joke to theme %f: %s: %s\n", rnd, joke.Theme, joke.Text)
	}

	fmt.Println("Data population completed.")
}

func addTheme(ctx context.Context, client *firestore.Client, themeName string, rnd float64) error {
	theme := serverImpl.Theme{
		Text:   themeName,
		Random: rnd,
	}
	_, err := client.Collection("themes").Doc(themeName).Set(ctx, theme)
	return err
}

func addJoke(ctx context.Context, client *firestore.Client, theme string, text string, rnd float64) error {

	joke := serverImpl.Joke{
		Theme:  theme,
		Text:   text,
		Random: rnd,
	}
	_, _, err := client.Collection("jokes").Add(ctx, joke)
	return err
}
