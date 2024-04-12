package cli

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/tabwriter"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewStatusCmd(
	path string,
	client mongo.DbClient,
) Cmd {
	return &statusCmd{
		path:       path,
		collection: client.Collection(collectionName),
	}
}

type statusCmd struct {
	path       string
	collection mongo.DbCollection
}

func (c *statusCmd) Exec(ctx context.Context) error {
	ids, hasUp, hasDown, err := c.findMigrations()
	if err != nil {
		return err
	}
	executedIDs, executed, err := c.findExecutedMigrations(ctx)
	if err != nil {
		return err
	}

	var prev, current, next, latest string
	if len(executedIDs) > 0 {
		current = executedIDs[len(executedIDs)-1]
		if len(executedIDs) > 1 {
			prev = executedIDs[len(executedIDs)-2]
		}
	}

	available := len(ids)
	if len(ids) > 0 {
		latest = ids[len(ids)-1]
	}

	executedUnavailable := 0
	for _, id := range executedIDs {
		if !hasUp[id] && !hasDown[id] {
			ids = append(ids, id)
			executedUnavailable++
		}
	}
	sort.Strings(ids)

	if current == "" {
		if len(ids) > 0 {
			next = ids[0]
		}
	} else {
		for i, id := range ids {
			if current == id && i < len(ids)-1 {
				next = ids[i+1]
			}
		}
	}

	notMigrated := 0
	for _, id := range ids {
		if !executed[id] {
			notMigrated++
		}
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
	fmt.Fprintln(w, "== Configuration")
	fmt.Fprintf(w, "Directory:\t%s\n", c.path)
	fmt.Fprintf(w, "Migration collection:\t%s\n", collectionName)
	fmt.Fprintf(w, "Previous version:\t%s\n", prev)
	fmt.Fprintf(w, "Current version:\t%s\n", current)
	fmt.Fprintf(w, "Next version:\t%s\n", next)
	fmt.Fprintf(w, "Latest version:\t%s\n", latest)
	fmt.Fprintf(w, "Executed migrations:\t%d\n", len(executedIDs))
	fmt.Fprintf(w, "Executed unavailable migrations:\t%d\n", executedUnavailable)
	fmt.Fprintf(w, "Available migrations:\t%d\n", available)
	fmt.Fprintf(w, "New migrations:\t%d\n", notMigrated)

	fmt.Fprintf(w, "\n== Migration versions\n")
	for _, id := range ids {
		if executed[id] {
			fmt.Fprintf(w, "%s\tmigrated\n", id)
		} else {
			fmt.Fprintf(w, "%s\tnot migrated\n", id)
		}

		if hasUp[id] && !hasDown[id] {
			fmt.Fprintf(w, "%s\tmissing down migration script\n", id)
		} else if !hasUp[id] && hasDown[id] {
			fmt.Fprintf(w, "%s\tmissing up migration script\n", id)
		} else if !hasUp[id] && !hasDown[id] {
			fmt.Fprintf(w, "%s\tmissing up and down migration scripts\n", id)
		}
	}

	return w.Flush()
}

func (c *statusCmd) findMigrations() ([]string, map[string]bool, map[string]bool, error) {
	upFiles, err := filepath.Glob(filepath.Join(c.path, "*"+fileNameSuffixUp))
	if err != nil {
		return nil, nil, nil, fmt.Errorf("cannot read directory %q: %w", c.path, err)
	}
	downFiles, err := filepath.Glob(filepath.Join(c.path, "*"+fileNameSuffixDown))
	if err != nil {
		return nil, nil, nil, fmt.Errorf("cannot read directory %q: %w", c.path, err)
	}

	ids := make([]string, 0)
	hasUp := make(map[string]bool)
	hasDown := make(map[string]bool)

	for _, file := range upFiles {
		id := strings.TrimSuffix(filepath.Base(file), fileNameSuffixUp)
		ids = append(ids, id)
		hasUp[id] = true
	}

	for _, file := range downFiles {
		id := strings.TrimSuffix(filepath.Base(file), fileNameSuffixDown)
		hasDown[id] = true
		if !hasUp[id] {
			ids = append(ids, id)
		}
	}

	return ids, hasUp, hasDown, nil
}

func (c *statusCmd) findExecutedMigrations(ctx context.Context) ([]string, map[string]bool, error) {
	data := struct {
		ID string `bson:"_id"`
	}{}
	cursor, err := c.collection.Find(ctx, bson.M{}, options.Find().SetSort(bson.M{"_id": 1}))
	if err != nil {
		return nil, nil, fmt.Errorf("cannot fetch migrations: %w", err)
	}

	ids := make([]string, 0)
	executed := make(map[string]bool)
	for cursor.Next(ctx) {
		err := cursor.Decode(&data)
		if err != nil {
			return nil, nil, fmt.Errorf("cannot decode migration: %w", err)
		}

		ids = append(ids, data.ID)
		executed[data.ID] = true
	}

	return ids, executed, nil
}
