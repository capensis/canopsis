package cli

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
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
	ids, hasUp, hasDown, invalidFile, err := c.findMigrations()
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
			next = executedIDs[len(executedIDs)-2]
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

	if len(invalidFile) > 0 {
		fmt.Fprintf(w, "\n== Invalid migration files\n")

		for _, name := range invalidFile {
			fmt.Fprintln(w, name)
		}
	}

	return w.Flush()
}

func (c *statusCmd) findMigrations() ([]string, map[string]bool, map[string]bool, []string, error) {
	files, err := ioutil.ReadDir(c.path)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("cannot read directory %q: %w", c.path, err)
	}

	suffixUp := fmt.Sprintf("%s%s%s", fileNameDelimiter, fileNameSuffixUp, fileExtJs)
	suffixDown := fmt.Sprintf("%s%s%s", fileNameDelimiter, fileNameSuffixDown, fileExtJs)
	ids := make([]string, 0)
	hasUp := make(map[string]bool)
	hasDown := make(map[string]bool)
	invalidFile := make([]string, 0)

	for _, file := range files {
		if strings.HasSuffix(file.Name(), suffixUp) {
			id := strings.TrimSuffix(file.Name(), suffixUp)
			hasUp[id] = true
			if !hasDown[id] {
				ids = append(ids, id)
			}
		} else if strings.HasSuffix(file.Name(), suffixDown) {
			id := strings.TrimSuffix(file.Name(), suffixDown)
			hasDown[id] = true
			if !hasUp[id] {
				ids = append(ids, id)
			}
		} else {
			invalidFile = append(invalidFile, file.Name())
		}
	}

	return ids, hasUp, hasDown, invalidFile, nil
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
