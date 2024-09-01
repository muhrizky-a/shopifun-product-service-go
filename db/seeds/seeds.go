package seeds

import (
	"codebase-app/internal/adapter"
	"context"
	"os"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/jmoiron/sqlx"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
)

// Seed struct.
type Seed struct {
	db *sqlx.DB
}

// NewSeed return a Seed with a pool of connection to a dabase.
func newSeed(db *sqlx.DB) Seed {
	return Seed{
		db: db,
	}
}

func Execute(db *sqlx.DB, table string, total int) {
	seed := newSeed(db)
	seed.run(table, total)
}

// Run seeds.
func (s *Seed) run(table string, total int) {

	switch table {
	case "categories":
		s.categoriesSeed(total)
	case "all":
		s.categoriesSeed(total)
	case "delete-all":
		s.deleteAll()
	default:
		log.Warn().Msg("No seed to run")
	}

	if table != "" {
		log.Info().Msg("Seed ran successfully")
		log.Info().Msg("Exiting ...")
		if err := adapter.Adapters.Unsync(); err != nil {
			log.Fatal().Err(err).Msg("Error while closing database connection")
		}
		os.Exit(0)
	}
}

func (s *Seed) deleteAll() {
	tx, err := s.db.BeginTxx(context.Background(), nil)
	if err != nil {
		log.Error().Err(err).Msg("Error starting transaction")
		return
	}
	defer func() {
		if err != nil {
			err = tx.Rollback()
			log.Error().Err(err).Msg("Error rolling back transaction")
			return
		} else {
			err = tx.Commit()
			if err != nil {
				log.Error().Err(err).Msg("Error committing transaction")
			}
		}
	}()

	_, err = tx.Exec(`DELETE FROM categories`)
	if err != nil {
		log.Error().Err(err).Msg("Error deleting categories")
		return
	}
	log.Info().Msg("categories table deleted successfully")

	log.Info().Msg("=== All tables deleted successfully ===")
}

// categoriesSeed seeds the categories table.
func (s *Seed) categoriesSeed(total int) {
	tx, err := s.db.BeginTxx(context.Background(), nil)
	if err != nil {
		log.Error().Err(err).Msg("Error starting transaction")
		return
	}
	defer func() {
		if err != nil {
			err = tx.Rollback()
			log.Error().Err(err).Msg("Error rolling back transaction")
			return
		}

		err = tx.Commit()
		if err != nil {
			log.Error().Err(err).Msg("Error committing transaction")
		}
	}()

	userMaps := make([]map[string]any, 0)

	for i := 0; i < total; i++ {
		name := gofakeit.Name()
		dataUserToInsert := make(map[string]any)
		dataUserToInsert["id"] = ulid.Make().String()
		dataUserToInsert["name"] = name
		dataUserToInsert["description"] = "deskripsi kategori " + name

		userMaps = append(userMaps, dataUserToInsert)
	}

	_, err = tx.NamedExec(`
		INSERT INTO categories (name, description)
		VALUES (:name, :description)
	`, userMaps)
	if err != nil {
		log.Error().Err(err).Msg("Error creating categories")
		return
	}

	log.Info().Msg("categories table seeded successfully")
}
