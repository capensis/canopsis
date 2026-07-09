package eventfilter

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strconv"
	"sync"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/externaldata"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern/match"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/template"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"github.com/rs/zerolog"
)

type ruleService struct {
	rulesAdapter            RuleAdapter
	ruleApplicatorContainer RuleApplicatorContainer
	externalDataGetter      externaldata.Getter
	eventCounter            EventCounter
	failureService          FailureService
	templateExecutor        template.Executor
	rules                   []ParsedRule
	rulesHash               string
	rulesMutex              sync.RWMutex
	logger                  zerolog.Logger
}

func (s *ruleService) LoadRules(ctx context.Context, types []string) error {
	rules, err := s.rulesAdapter.GetByTypes(ctx, types)
	if err != nil {
		return err
	}

	parsedRules := make([]ParsedRule, len(rules))
	ids := make([]string, len(rules))
	for i := 0; i < len(rules); i++ {
		ids[i] = rules[i].ID
		parsedRules[i] = ParseRule(rules[i], s.templateExecutor)
	}

	hash := s.rulesFingerprint(parsedRules)

	s.rulesMutex.Lock()
	defer s.rulesMutex.Unlock()
	s.rules = parsedRules
	s.rulesHash = hash
	s.logger.Debug().Strs("rules", ids).Msg("Loading event filter rules")

	return nil
}

func (s *ruleService) ProcessEvent(ctx context.Context, event *types.Event, partialResult ServiceResult) (ServiceResult, error) {
	s.rulesMutex.RLock()
	defer s.rulesMutex.RUnlock()

	// When resuming a parked event, the accumulated state only makes sense against the ruleset it was suspended on.
	// If the ruleset changed in the meantime,
	// signal a restart from the original event instead of resuming inconsistently.
	if partialResult.ExternalDataResponse != nil && partialResult.RulesHash != "" && partialResult.RulesHash != s.rulesHash {
		return ServiceResult{}, ErrRulesetChanged
	}

	res := ServiceResult{
		UpdatedEntityInfos:      partialResult.UpdatedEntityInfos,
		ExecutedEnrichRuleCount: partialResult.ExecutedEnrichRuleCount,
		ExternalRequestCount:    partialResult.ExternalRequestCount,
		RulesHash:               s.rulesHash,
	}
	if res.UpdatedEntityInfos == nil {
		res.UpdatedEntityInfos = make(map[string]UpdatedValue)
	}
	if res.ExternalRequestCount == nil {
		res.ExternalRequestCount = make(map[string]int64)
	}

	outcome := OutcomePass
	now := datetime.NewCpsTime()
	foundPrev := false
	for _, rule := range s.rules {
		if outcome != OutcomePass {
			break
		}

		if partialResult.ExternalDataResponse != nil {
			if partialResult.ExternalDataResponse.RuleID == rule.ID {
				foundPrev = true
			} else if !foundPrev {
				continue
			}
		}

		if rule.ResolvedStart != nil && rule.ResolvedStop != nil {
			inExDate := false
			for _, exdate := range rule.ResolvedExdates {
				if now.After(exdate.Begin) && now.Before(exdate.End) {
					inExDate = true
					break
				}
			}

			if inExDate {
				continue
			}

			if now.Before(*rule.ResolvedStart) {
				continue
			}

			if now.After(*rule.ResolvedStop) {
				if rule.NextResolvedStart == nil || rule.NextResolvedStop == nil {
					continue
				}

				if now.Before(*rule.NextResolvedStart) || now.After(*rule.NextResolvedStop) {
					continue
				}
			}
		}

		var err error
		var matched bool
		var eventRegexMatches match.EventRegexMatches
		var entityRegexMatches match.EntityRegexMatches

		if len(rule.EntityPattern) == 0 && len(rule.EventPattern) == 0 {
			if event.Debug {
				s.logger.Info().Str("rule", rule.ID).Str("event_type", event.EventType).Str("entity", event.GetEID()).Msg("Event filter rule service: rule is not matched")
			}

			s.failureService.Add(rule.ID, rule.Description, rule.Updated, FailureTypeInvalidPattern, "missing entity and event patterns", nil)
			continue
		}

		if len(rule.EventPattern) > 0 {
			matched, eventRegexMatches, err = match.MatchEventPatternWithRegexMatches(rule.EventPattern, event)
			if err != nil {
				s.logger.Err(err).Str("rule_id", rule.ID).Msg("Event filter rule service: invalid event pattern")
				s.failureService.Add(rule.ID, rule.Description, rule.Updated, FailureTypeInvalidPattern, "invalid event pattern: "+err.Error(), nil)
				continue
			}

			if !matched {
				if event.Debug {
					s.logger.Info().Str("rule", rule.ID).Str("event_type", event.EventType).Str("entity", event.GetEID()).Msg("Event filter rule service: rule is not matched")
				}

				continue
			}
		}

		if len(rule.EntityPattern) > 0 {
			if event.Entity == nil {
				if event.Debug {
					s.logger.Info().Str("rule", rule.ID).Str("event_type", event.EventType).Str("entity", event.GetEID()).Msg("Event filter rule service: entity is missing")
				}

				continue
			}

			matched, entityRegexMatches, err = match.MatchEntityPatternWithRegexMatches(rule.EntityPattern, event.Entity)
			if err != nil {
				s.logger.Err(err).Str("rule_id", rule.ID).Msg("Event filter rule service: invalid entity pattern")
				s.failureService.Add(rule.ID, rule.Description, rule.Updated, FailureTypeInvalidPattern, "invalid entity pattern: "+err.Error(), nil)
				continue
			}

			if !matched {
				if event.Debug {
					s.logger.Info().Str("rule", rule.ID).Str("event_type", event.EventType).Str("entity", event.GetEID()).Msg("Event filter rule service: rule is not matched")
				}

				continue
			}
		}

		if event.Debug {
			s.logger.Info().Str("rule", rule.ID).Str("event_type", event.EventType).Str("entity", event.GetEID()).Msg("Event filter rule service: rule is matched")
		}

		applicator, found := s.ruleApplicatorContainer.Get(rule.Type)
		if !found {
			s.logger.Warn().Str("rule_id", rule.ID).Str("rule_type", rule.Type).Msg("Event filter rule service: RuleApplicator doesn't exist")
			continue
		}

		regexMatch := RegexMatch{
			EventRegexMatches: eventRegexMatches,
			Entity:            entityRegexMatches,
		}

		var externalData map[string]any
		var ruleExternalRequestCount map[string]int64
		var externalDataRequest *ExternalDataRequest
		if resp := partialResult.ExternalDataResponse; resp != nil && resp.RuleID == rule.ID {
			if resp.Err == nil {
				externalData = partialResult.ExternalData
				if externalData == nil {
					externalData = make(map[string]any, len(rule.ExternalData))
				}
				i := 0
				for _, d := range rule.ExternalData {
					if d.Type == externaldata.RefTypeAPI {
						externalData[d.Reference] = resp.Result[i]
						i++
					}
				}
			} else {
				failReason := "external data \"" + rule.ExternalData[resp.ErrIndex].Reference + "\" cannot be fetched: " + resp.Err.Error()
				s.failureService.Add(rule.ID, rule.Description, rule.Updated, FailureTypeExternalDataAPI, failReason, event)

				outcome = rule.Config.OnFailure
				continue
			}
		} else {
			externalData, ruleExternalRequestCount, externalDataRequest, err = s.getExternalData(ctx, rule, event, regexMatch)
			if err != nil {
				s.logger.Err(err).Str("rule_id", rule.ID).Str("rule_type", rule.Type).Msg("Event filter rule service: failed to fetch external data")
				outcome = rule.Config.OnFailure

				continue
			}
		}

		for k, v := range ruleExternalRequestCount {
			res.ExternalRequestCount[k] += v
		}

		if externalDataRequest != nil {
			res.ExternalData = externalData
			res.ExternalDataRequest = externalDataRequest

			return res, nil
		}

		ar, err := applicator.Apply(ctx, rule, event, res.UpdatedEntityInfos, regexMatch, externalData)
		outcome = ar.Outcome
		if err != nil {
			s.logger.Err(err).Str("rule_id", rule.ID).Str("rule_type", rule.Type).Msg("Event filter rule service: failed to apply")
			continue
		}

		res.UpdatedEntityInfos = ar.UpdatedEntityInfos
		s.eventCounter.Add(rule.ID, rule.Updated)

		if rule.Type == RuleTypeEnrichment {
			res.ExecutedEnrichRuleCount++
		}
	}

	if outcome == OutcomeDrop {
		return res, ErrDropOutcome
	}

	return res, nil
}

func (s *ruleService) getExternalData(
	ctx context.Context,
	rule ParsedRule,
	event *types.Event,
	regexMatch RegexMatch,
) (map[string]any, map[string]int64, *ExternalDataRequest, error) {
	if len(rule.ExternalData) == 0 {
		return nil, nil, nil, nil
	}

	if s.externalDataGetter == nil {
		s.failureService.Add(rule.ID, rule.Description, rule.Updated, FailureTypeOther, "External data is not supported", nil)

		return nil, nil, nil, errors.New("external data is not supported")
	}

	er, err := s.externalDataGetter.Get(ctx, externaldata.Rule{
		ID:            rule.ID,
		Name:          rule.Description,
		ExternalData:  rule.ExternalData,
		TriggerUserID: event.UserID,
		TriggerAuthor: event.Author,
	}, Template{
		Event:      event,
		RegexMatch: regexMatch,
	})
	if err != nil {
		var failureType int64
		failReason := ""
		isParamsInvalid := false
		if getterTplErr, ok := errors.AsType[*externaldata.GetterTplError](err); ok {
			failureType = FailureTypeInvalidTemplate
			failReason = getterTplErr.FailReason()
			isParamsInvalid = getterTplErr.IsParamsInvalid()
		} else if getterErr, ok := errors.AsType[*externaldata.GetterError](err); ok {
			failureType, ok = failureTypeMapping[getterErr.Type()]
			if !ok {
				failureType = FailureTypeOther
			}

			failReason = getterErr.FailReason()
			isParamsInvalid = getterErr.IsParamsInvalid()
		}

		if failReason != "" {
			if isParamsInvalid {
				s.failureService.Add(rule.ID, rule.Description, rule.Updated, failureType, failReason, nil)
			} else {
				s.failureService.Add(rule.ID, rule.Description, rule.Updated, failureType, failReason, event)
			}
		}

		return nil, nil, nil, err
	}

	var request *ExternalDataRequest
	if er.WebhookEvent != nil {
		request = &ExternalDataRequest{
			RuleID:        rule.ID,
			CorrelationID: er.CorrelationID,
			WebhookEvent:  *er.WebhookEvent,
		}
	}

	return er.ExternalData, er.RequestCount, request, nil
}

// rulesFingerprint builds a stable digest of the ruleset from each rule's id and updated time.
func (s *ruleService) rulesFingerprint(rules []ParsedRule) string {
	entries := make([]string, len(rules))
	for i, r := range rules {
		entries[i] = r.ID + ":" + strconv.FormatInt(r.Updated.Unix(), 10)
	}

	h := sha256.New()
	for _, e := range entries {
		h.Write([]byte(e))
		h.Write([]byte{0})
	}

	return hex.EncodeToString(h.Sum(nil))
}

func NewRuleService(
	ruleAdapter RuleAdapter,
	container RuleApplicatorContainer,
	externalDataGetter externaldata.Getter,
	eventCounter EventCounter,
	failureService FailureService,
	templateExecutor template.Executor,
	logger zerolog.Logger,
) Service {
	return &ruleService{
		rulesMutex:              sync.RWMutex{},
		rulesAdapter:            ruleAdapter,
		ruleApplicatorContainer: container,
		externalDataGetter:      externalDataGetter,
		eventCounter:            eventCounter,
		failureService:          failureService,
		templateExecutor:        templateExecutor,
		logger:                  logger,
	}
}
