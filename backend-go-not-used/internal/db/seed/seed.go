package seed

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/jeevanions/italian-learning/internal/api/models"
	dbmodels "github.com/jeevanions/italian-learning/internal/db/models"
	"github.com/jeevanions/italian-learning/internal/domain/services"
)

func SeedBasicData(wordService services.WordService, groupService services.GroupService) error {
	ctx := context.Background()

	// Seed basic vocabulary groups
	groups := []*dbmodels.Group{
		{
			Name:            "Basic Greetings",
			Description:     "Essential Italian greetings and farewells",
			DifficultyLevel: 1,
			Category:        "vocabulary",
		},
		{
			Name:            "Numbers 1-10",
			Description:     "Basic Italian numbers",
			DifficultyLevel: 1,
			Category:        "vocabulary",
		},
	}

	for _, group := range groups {
		if err := groupService.CreateGroup(ctx, group); err != nil {
			return err
		}
	}

	// Seed basic words
	verbConjugation, _ := json.Marshal(map[string]string{
		"present_1s": "sono",
		"present_2s": "sei",
		"present_3s": "è",
		"present_1p": "siamo",
		"present_2p": "siete",
		"present_3p": "sono",
	})

	words := []*dbmodels.Word{
		{
			Italian:         "essere",
			English:         "to be",
			PartsOfSpeech:   "verb",
			DifficultyLevel: 1,
			VerbConjugation: verbConjugation,
		},
		{
			Italian:         "ciao",
			English:         "hello/goodbye",
			PartsOfSpeech:   "interjection",
			DifficultyLevel: 1,
			Notes:           sql.NullString{String: "Used for both greeting and farewell", Valid: true},
		},
	}

	// Convert and seed words
	for _, dbWord := range words {
		req := &models.CreateWordRequest{
			Italian:         dbWord.Italian,
			English:         dbWord.English,
			PartsOfSpeech:   dbWord.PartsOfSpeech,
			DifficultyLevel: dbWord.DifficultyLevel,
		}
		if dbWord.Notes.Valid {
			notes := dbWord.Notes.String
			req.Notes = &notes
		}
		if len(dbWord.VerbConjugation) > 0 {
			vc := string(dbWord.VerbConjugation)
			req.VerbConjugation = &vc
		}

		_, err := wordService.CreateWord(ctx, req)
		if err != nil {
			return err
		}
	}

	return nil
}


