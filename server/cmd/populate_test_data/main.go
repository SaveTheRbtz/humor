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
		{"Animals", "What do you call a bear with no teeth? A gummy bear."},
		{"Animals", "Why did the squirrel swim on its back? To keep its nuts dry."},
		{"Technology", "Why do scientists trust atoms? Because they make up everything."},
		{"Technology", "Why do programmers prefer dark mode? Because light attracts bugs."},
		{"Technology", "How many programmers does it take to change a light bulb? None, it's a hardware problem."},
		{"Technology", "Why do programmers always mix up Christmas and Halloween? Because Oct 31 == Dec 25."},
		{"Technology", "Why do Java developers wear glasses? Because they can't C#."},
		{"Food", "Why did the cookie go to the doctor? Because it was feeling crumbly."},
		{"Food", "I'm on a seafood diet. I see food and I eat it."},
		{"Food", "Why did the tomato turn red? Because it saw the salad dressing."},
		{"Food", "What do you call cheese that isn't yours? Nacho cheese."},
		{"Food", "Why did the coffee file a police report? It got mugged."},
	}

	themes := make(map[string]struct{})
	for _, joke := range jokes {
		themes[joke.Theme] = struct{}{}
	}

	themeToID := make(map[string]string)
	for themeName := range themes {
		rnd := rand.Float64()
		ref, err := addTheme(ctx, client, themeName, rnd)
		if err != nil {
			log.Fatalf("Failed to add theme %s: %v", themeName, err)
		}
		themeToID[themeName] = ref.ID
		fmt.Printf("Added theme: %f: %s\n", rnd, themeName)
	}

	for _, joke := range jokes {
		rnd := rand.Float64()
		err := addJoke(ctx, client, themeToID[joke.Theme], joke.Theme, joke.Text, rnd)
		if err != nil {
			log.Fatalf("Failed to add joke: %v", err)
		}
		fmt.Printf("Added joke to theme %f: %s (%s): %s\n",
			rnd, joke.Theme, themeToID[joke.Theme], joke.Text)
	}

	fmt.Println("Data population completed.")
}

func addTheme(
	ctx context.Context,
	client *firestore.Client,
	themeName string,
	rnd float64,
) (*firestore.DocumentRef, error) {
	theme := serverImpl.Theme{
		Text:   themeName,
		Random: rnd,
		Active: true,
	}
	docRef, _, err := client.Collection("themes").Add(ctx, theme)
	return docRef, err
}

func addJoke(ctx context.Context,
	client *firestore.Client,
	themeID string,
	themeName string,
	text string,
	rnd float64,
) error {
	joke := serverImpl.Joke{
		ThemeID: themeID,
		Theme:   themeName,
		Text:    text,
		Model:   fmt.Sprintf("dad-%d", rand.Intn(3)),
		Random:  rnd,
		Active:  true,
	}
	_, _, err := client.Collection("jokes").Add(ctx, joke)
	return err
}
