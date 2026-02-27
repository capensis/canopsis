package pattern

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"regexp/syntax"
	"slices"
	"strings"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/patternfields"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/workers"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/v2/bson"
	mongodriver "go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"golang.org/x/sync/errgroup"
)

const (
	maxRetries = 10
)

type OptimizeWorker interface {
	CreateJob(ctx context.Context, r OptimizeRequest) (OptimizeJob, error)
	GetJob(ctx context.Context, id string) (OptimizeJob, error)
	UpdateJob(ctx context.Context, r OptimizeAcceptRequest) (OptimizeJob, error)
	DeleteJob(ctx context.Context, id string) (bool, error)
	ProcessJob(ctx context.Context, id string) error
	ProcessAbandonedJobs(ctx context.Context)
}

func NewOptimizeWorker(
	store Store,
	dbClient mongo.DbClient,
	jobPublisher workers.JobPublisher,
	transformer patternfields.Transformer,
	logger zerolog.Logger,
) OptimizeWorker {
	return &optimizeWorker{
		store:                   store,
		dbCollection:            dbClient.Collection(mongo.PatternOptimizeJobCollection),
		jobPublisher:            jobPublisher,
		transformer:             transformer,
		logger:                  logger,
		abandonedTickerInterval: time.Minute,
		pingInterval:            time.Second,
	}
}

type optimizeWorker struct {
	store                   Store
	dbCollection            mongo.DbCollection
	jobPublisher            workers.JobPublisher
	transformer             patternfields.Transformer
	abandonedTickerInterval time.Duration
	pingInterval            time.Duration
	logger                  zerolog.Logger
}

func (w *optimizeWorker) CreateJob(ctx context.Context, r OptimizeRequest) (OptimizeJob, error) {
	entityPattern, _, err := w.transformer.TransformAliases(ctx, r.EntityPattern, r)
	if err != nil {
		return OptimizeJob{}, err
	}

	job := OptimizeJob{
		ID:                    utils.NewID(),
		Status:                OptimizeStatusCreated,
		EntityPattern:         entityPattern,
		Created:               datetime.NewCpsTime(),
		Suggestions:           make([]Suggestion, 0),
		OptimizedFieldRegexps: make([]OptimizedFieldRegexp, 0),
	}

	_, err = w.dbCollection.InsertOne(ctx, job)
	if err != nil {
		return OptimizeJob{}, fmt.Errorf("failed to insert pattern optimize job: %w", err)
	}

	err = w.jobPublisher.Publish(ctx, job.ID)
	if err != nil {
		return OptimizeJob{}, fmt.Errorf("failed to publish pattern optimize job: %w", err)
	}

	return job, nil
}

func (w *optimizeWorker) ProcessJob(ctx context.Context, id string) (resErr error) {
	job := OptimizeJob{}

	err := w.dbCollection.FindOneAndUpdate(ctx,
		bson.M{
			"_id":    id,
			"status": bson.M{"$in": bson.A{OptimizeStatusCreated, OptimizeStatusRunning}},
		},
		bson.M{"$set": bson.M{
			"status":    OptimizeStatusRunning,
			"last_ping": datetime.NewCpsTime(),
		}},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&job)
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return nil
		}

		return err
	}

	g, gCtx := errgroup.WithContext(ctx)
	done := make(chan struct{})

	g.Go(func() (resErr error) {
		defer close(done)

		var update bson.M

		result, err := w.optimize(gCtx, job)
		if err != nil {
			if mongo.IsConnectionError(err) {
				return err
			}

			failReason := err.Error()
			w.logger.Err(err).Str("job", job.ID).Msg("failed to optimize pattern")

			update = bson.M{
				"$set": bson.M{
					"status":      OptimizeStatusFailed,
					"fail_reason": failReason,
				},
			}
		} else {
			update = bson.M{
				"$set": bson.M{
					"status":                  OptimizeStatusSucceeded,
					"suggestions":             result.Suggestions,
					"optimized_field_regexps": result.OptimizedFieldRegexps,
					"original_pattern_ms":     result.OriginalPatternMS,
					"original_pattern_count":  result.OriginalPatternCount,
				},
			}
		}

		updateRes, err := w.dbCollection.UpdateOne(gCtx,
			bson.M{
				"_id":    job.ID,
				"status": OptimizeStatusRunning,
			}, update,
		)
		if err != nil {
			return fmt.Errorf("failed to update pattern optimize job: %w", err)
		}

		if updateRes.ModifiedCount == 0 {
			return errors.New("pattern optimize job is processing by another worker or cancelled")
		}

		return nil
	})

	g.Go(func() error {
		ticker := time.NewTicker(w.pingInterval)
		lastPing := *job.LastPing
		for {
			select {
			case <-ctx.Done():
				return nil
			case <-done:
				return nil
			case <-ticker.C:
				newLastPing := datetime.NewCpsTime()
				updateRes, err := w.dbCollection.UpdateOne(gCtx,
					bson.M{
						"_id":       job.ID,
						"status":    OptimizeStatusRunning,
						"last_ping": lastPing,
					},
					bson.M{"$set": bson.M{
						"last_ping": newLastPing,
					}},
				)
				if err != nil {
					return fmt.Errorf("failed to update optimize job status: %w", err)
				}

				if updateRes.ModifiedCount == 0 {
					return errors.New("optimize job is processing by another worker or cancelled")
				}

				lastPing = newLastPing
			}
		}
	})

	return g.Wait()
}

func (w *optimizeWorker) optimize(ctx context.Context, job OptimizeJob) (OptimizeResult, error) {
	initEntityIDs, duration, err := w.store.GetEntityIDs(ctx, job.EntityPattern)
	if err != nil {
		return OptimizeResult{}, err
	}

	res := OptimizeResult{
		OriginalPatternMS:     duration,
		OriginalPatternCount:  len(initEntityIDs),
		Suggestions:           make([]Suggestion, 0),
		OptimizedFieldRegexps: make([]OptimizedFieldRegexp, 0),
	}

	// if an initial pattern doesn't return any entity, it makes impossible to make any suggestion.
	if len(initEntityIDs) == 0 {
		return res, nil
	}

	allOptimizedFieldRegexpsMap := make(map[string]map[string]bool)
	allOptimizedFieldRegexps := make([]OptimizedFieldRegexp, 0)

	fieldSuggestions := make([][][]pattern.FieldCondition, len(job.EntityPattern))
	for idx, conditions := range job.EntityPattern {
		var optimizedFieldRegexps []OptimizedFieldRegexp
		fieldSuggestions[idx], optimizedFieldRegexps, err = w.suggestFieldConditions(ctx, conditions)
		if err != nil {
			return OptimizeResult{}, fmt.Errorf("failed to suggest field conditions: %w", err)
		}

		for _, optimizedFieldRegexp := range optimizedFieldRegexps {
			regexps, ok := allOptimizedFieldRegexpsMap[optimizedFieldRegexp.Field]
			if !ok {
				allOptimizedFieldRegexpsMap[optimizedFieldRegexp.Field] = map[string]bool{optimizedFieldRegexp.Regexp: true}
				allOptimizedFieldRegexps = append(allOptimizedFieldRegexps, optimizedFieldRegexp)
				continue
			}

			if _, ok := regexps[optimizedFieldRegexp.Regexp]; !ok {
				regexps[optimizedFieldRegexp.Regexp] = true
				allOptimizedFieldRegexps = append(allOptimizedFieldRegexps, optimizedFieldRegexp)
			}
		}
	}

	suggestedPatterns := getPatternsCombinations(fieldSuggestions)

	// if no optimization fields, then there is no need to suggest anything.
	if len(allOptimizedFieldRegexps) == 0 {
		return res, nil
	}

	g, gCtx := errgroup.WithContext(ctx)
	results := make([]Suggestion, len(suggestedPatterns))
	for i := 0; i < len(suggestedPatterns); i++ {
		g.Go(func() error {
			candidateEntityIDs, _, err := w.store.GetEntityIDs(gCtx, suggestedPatterns[i])
			if err != nil {
				return fmt.Errorf("failed to get entity IDs by a suggested pattern: %w", err)
			}

			results[i] = Suggestion{
				EntityPattern: suggestedPatterns[i],
				FoundEntities: len(candidateEntityIDs),
				Difference:    symmetricDiffCount(initEntityIDs, candidateEntityIDs),
			}

			return nil
		})
	}

	err = g.Wait()
	if err != nil {
		return OptimizeResult{}, err
	}

	slices.SortFunc(results, func(l, r Suggestion) int {
		return l.Difference - r.Difference
	})

	res.Suggestions = results[:min(len(results), maxReturnedSuggestions)]
	res.OptimizedFieldRegexps = allOptimizedFieldRegexps

	return res, nil
}

func (w *optimizeWorker) suggestFieldConditions(ctx context.Context, initConditions []pattern.FieldCondition) ([][]pattern.FieldCondition, []OptimizedFieldRegexp, error) {
	// start with an empty suggested field condition.
	currentConditions := [][]pattern.FieldCondition{{}}
	optimizedFieldRegexps := make([]OptimizedFieldRegexp, 0)

	// takenFields shows that the entity field is taken by another condition and can't be used for a suggestion.
	takenFields := make(map[string]bool, len(initConditions))

	for _, initCondition := range initConditions {
		infoField, found := strings.CutPrefix(initCondition.Field, EntityInfosPrefix)
		if !found {
			w.appendCurrentConditions(initCondition, initCondition.Field, takenFields, &currentConditions)
			continue
		}

		if initCondition.Condition.Type != pattern.ConditionRegexp {
			w.appendCurrentConditions(initCondition, infoField, takenFields, &currentConditions)
			continue
		}

		regexpString, ok := initCondition.Condition.Value.(string)
		if !ok {
			w.appendCurrentConditions(initCondition, infoField, takenFields, &currentConditions)
			continue
		}

		parsedTree, err := syntax.Parse(regexpString, syntax.Perl)
		if err != nil {
			w.appendCurrentConditions(initCondition, infoField, takenFields, &currentConditions)
			continue
		}

		literalGroups, err := parseLiterals(parsedTree)
		if err != nil {
			w.appendCurrentConditions(initCondition, infoField, takenFields, &currentConditions)
			continue
		}

		uniqueLiteralsMap := make(map[string]bool)
		uniqueLiterals := make([]string, 0)

		for i := range literalGroups {
			for j := range literalGroups[i] {
				if !uniqueLiteralsMap[literalGroups[i][j]] {
					uniqueLiteralsMap[literalGroups[i][j]] = true
					uniqueLiterals = append(uniqueLiterals, literalGroups[i][j])
				}
			}
		}

		literalStats, err := w.store.GetLiteralsFieldStats(ctx, uniqueLiterals)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to count literals field stats: %w", err)
		}

		suggestions, err := suggestConditions(literalStats, literalGroups, takenFields)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to suggest conditions: %w", err)
		}

		if len(suggestions) == 0 {
			w.appendCurrentConditions(initCondition, infoField, takenFields, &currentConditions)
			continue
		}

		currentConditions = w.expandCurrentConditions(currentConditions, suggestions)
		optimizedFieldRegexps = append(optimizedFieldRegexps, OptimizedFieldRegexp{Field: infoField, Regexp: regexpString})
	}

	return currentConditions, optimizedFieldRegexps, nil
}

func (w *optimizeWorker) appendCurrentConditions(cond pattern.FieldCondition, field string, takenFields map[string]bool, conditions *[][]pattern.FieldCondition) {
	takenFields[field] = true
	for idx := range *conditions {
		(*conditions)[idx] = append((*conditions)[idx], cond)
	}
}

func (w *optimizeWorker) expandCurrentConditions(currentConditions [][]pattern.FieldCondition, suggestions [][]pattern.FieldCondition) [][]pattern.FieldCondition {
	maxSize := min(len(currentConditions)*len(suggestions), maxCalculatedSuggestions)
	expanded := make([][]pattern.FieldCondition, maxSize)

	for i, currentCond := range currentConditions {
		for j, alternative := range suggestions {
			idx := i*len(suggestions) + j
			if idx >= maxSize {
				return expanded
			}

			combined := make([]pattern.FieldCondition, len(currentCond)+len(alternative))
			copy(combined, currentCond)
			copy(combined[len(currentCond):], alternative)
			expanded[idx] = combined
		}
	}

	return expanded
}

func (w *optimizeWorker) GetJob(ctx context.Context, id string) (OptimizeJob, error) {
	job := OptimizeJob{}
	err := w.dbCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&job)
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return OptimizeJob{}, nil
		}

		return OptimizeJob{}, err
	}

	return job, nil
}

func (w *optimizeWorker) UpdateJob(ctx context.Context, r OptimizeAcceptRequest) (OptimizeJob, error) {
	update := bson.M{"status": OptimizeStatusRejected}
	if r.Accept != nil && *r.Accept {
		update = bson.M{"status": OptimizeStatusAccepted, "accepted_suggestion": r.Index}
	}

	job := OptimizeJob{}

	err := w.dbCollection.FindOneAndUpdate(
		ctx,
		bson.M{"_id": r.ID},
		bson.M{"$set": update},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&job)
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return OptimizeJob{}, nil
		}

		return OptimizeJob{}, err
	}

	return job, nil
}

func (w *optimizeWorker) DeleteJob(ctx context.Context, id string) (bool, error) {
	deleted, err := w.dbCollection.DeleteOne(ctx, bson.M{
		"_id":    id,
		"status": bson.M{"$in": bson.A{OptimizeStatusCreated, OptimizeStatusRunning}},
	})

	return deleted > 0, err
}

func (w *optimizeWorker) ProcessAbandonedJobs(ctx context.Context) {
	ticker := time.NewTicker(w.abandonedTickerInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			now := datetime.NewCpsTime()
			cursor, err := w.dbCollection.Find(ctx, bson.M{
				"status":    bson.M{"$in": bson.A{OptimizeStatusCreated, OptimizeStatusRunning}},
				"last_ping": bson.M{"$lt": now.Time.Add(-2 * w.pingInterval).Unix()},
			})
			if err != nil {
				w.logger.Err(err).Msg("failed to find abandoned optimize jobs")
				continue
			}

			for cursor.Next(ctx) {
				job := OptimizeJob{}
				if err = cursor.Decode(&job); err != nil {
					w.logger.Err(err).Msg("failed to decode abandoned optimize job")
					continue
				}

				if job.Retries == maxRetries {
					_, err = w.dbCollection.UpdateOne(ctx, bson.M{"_id": job.ID}, bson.M{"$set": bson.M{
						"status":      OptimizeStatusFailed,
						"fail_reason": "max retries exceeded",
					}})
					if err != nil {
						w.logger.Err(err).Msg("failed to update abandoned failed optimize job")
					}

					continue
				}

				job.Retries++
				res, err := w.dbCollection.UpdateOne(ctx, bson.M{"_id": job.ID}, bson.M{"$set": bson.M{
					"retries":   job.Retries,
					"last_ping": now,
				}})
				if err != nil {
					w.logger.Err(err).Msg("failed to update abandoned optimize job")
					continue
				}

				if res.ModifiedCount == 0 {
					w.logger.Warn().Msgf("abandoned optimize job was already processed: %s", job.ID)
				}

				err = w.jobPublisher.Publish(ctx, job.ID)
				if err != nil {
					w.logger.Err(err).Msg("failed to publish abandoned optimize job")
					continue
				}
			}

			if err = cursor.Err(); err != nil {
				w.logger.Err(err).Msg("failed to fetch abandoned optimize jobs")
			}

			err = cursor.Close(ctx)
			if err != nil {
				w.logger.Err(err).Msg("failed to close cursor")
			}
		}
	}
}

func suggestConditions(literalFieldStats map[string][]LiteralFieldStats, literalGroups [][]string, takenFields map[string]bool) ([][]pattern.FieldCondition, error) {
	var alternatives [][]pattern.FieldCondition
	uniqueConditions := make(map[string]bool)

	// Generate all possible field assignment combinations for the literals
	fieldCombinations := getLiteralToFieldCombinations(literalFieldStats, takenFields)

	for _, fieldCombination := range fieldCombinations {
		if len(fieldCombination) == 0 {
			continue
		}

		// Build groups of literals per field
		// Start with a single empty mapping
		literalsByField := []map[string][]string{{}}

		for _, literalGroup := range literalGroups {
			if len(literalGroup) == 1 {
				// Single literal: add to its assigned field
				field, ok := fieldCombination[literalGroup[0]]
				if !ok {
					continue
				}

				for i := range literalsByField {
					literalsByField[i][field] = append(literalsByField[i][field], literalGroup[0])
				}
			} else {
				// Multiple literals in a group (OR condition): determine which fields they map to
				targetFieldsMap := make(map[string]bool)
				targetFields := make([]string, 0)

				for _, literal := range literalGroup {
					field, ok := fieldCombination[literal]
					if !ok {
						continue
					}

					if !targetFieldsMap[field] {
						targetFields = append(targetFields, field)
						targetFieldsMap[field] = true
					}
				}

				if len(targetFields) == 1 {
					// All literals map to the same field - group them together
					for i := range literalsByField {
						literalsByField[i][targetFields[0]] = append(literalsByField[i][targetFields[0]], literalGroup...)
					}
				} else {
					// Literals map to different fields - create cartesian product
					newLiteralsByField := make([]map[string][]string, len(literalsByField)*len(targetFields))

					for i, field := range targetFields {
						for j, mapping := range literalsByField {
							idx := i*len(literalsByField) + j
							newLiteralsByField[idx] = make(map[string][]string)
							maps.Copy(newLiteralsByField[idx], mapping)
							newLiteralsByField[idx][field] = append(newLiteralsByField[idx][field], literalGroup...)
						}
					}

					literalsByField = newLiteralsByField
				}
			}
		}

		// Convert field->literals mappings to FieldConditions
		for _, fieldLiterals := range literalsByField {
			conditions := buildFieldConditions(fieldLiterals, takenFields)
			key, err := conditionsKey(conditions)
			if err != nil {
				return nil, fmt.Errorf("failed to build conditions key: %w", err)
			}

			if !uniqueConditions[key] {
				uniqueConditions[key] = true
				alternatives = append(alternatives, conditions)
			}
		}
	}

	return alternatives, nil
}

func conditionsKey(conditions []pattern.FieldCondition) (string, error) {
	var sb strings.Builder
	for i, c := range conditions {
		if i > 0 {
			sb.WriteByte('|')
		}
		sb.WriteString(c.Field)
		sb.WriteByte(':')
		sb.WriteString(c.Condition.Type)
		sb.WriteByte(':')
		_, err := fmt.Fprintf(&sb, "%v", c.Condition.Value)
		if err != nil {
			return "", fmt.Errorf("failed to format condition value: %w", err)
		}
	}

	return sb.String(), nil
}

func buildFieldConditions(fieldLiterals map[string][]string, takenFields map[string]bool) []pattern.FieldCondition {
	conditions := make([]pattern.FieldCondition, 0, len(fieldLiterals))

	fieldNames := make([]string, 0, len(fieldLiterals))
	for fieldName := range fieldLiterals {
		fieldNames = append(fieldNames, fieldName)
	}

	// ensure order
	slices.Sort(fieldNames)

	for _, fieldName := range fieldNames {
		fieldCondition := pattern.FieldCondition{}

		literals := fieldLiterals[fieldName]
		if len(literals) == 1 {
			fieldCondition.Condition = pattern.NewStringCondition(pattern.ConditionEqual, literals[0])
		} else {
			fieldCondition.Condition = pattern.NewStringArrayCondition(pattern.ConditionIsOneOf, literals)
		}

		switch fieldName {
		case EntityFieldComponent, EntityFieldName:
			fieldCondition.Field = fieldName
		default:
			fieldCondition.Field = EntityInfosPrefix + fieldName
			fieldCondition.FieldType = pattern.FieldTypeString
		}

		takenFields[fieldName] = true
		conditions = append(conditions, fieldCondition)
	}

	return conditions
}

func symmetricDiffCount(a, b []string) int {
	count := len(a)
	set := make(map[string]bool, len(a))

	for _, v := range a {
		set[v] = true
	}
	for _, v := range b {
		if set[v] {
			count--
		} else {
			count++
		}
	}

	return count
}
