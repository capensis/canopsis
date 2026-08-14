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

	hash := rulesHash(parsedRules)

	s.rulesMutex.Lock()
	defer s.rulesMutex.Unlock()
	s.rules = parsedRules
	s.rulesHash = hash
	s.logger.Debug().Strs("rules", ids).Msg("Loading event filter rules")

	return nil
}

func (s *ruleService) ProcessEvent(ctx context.Context, event *types.Event, partialResult ServiceResult) (ServiceResult, error) {
	rules, hash := s.getRules()

	// When resuming a parked event, the accumulated state only makes sense against the ruleset it was suspended on.
	// If the ruleset changed in the meantime,
	// signal a restart from the original event instead of resuming inconsistently.
	if partialResult.ExternalDataResponse != nil && partialResult.RulesHash != "" && partialResult.RulesHash != hash {
		return ServiceResult{}, ErrRulesetChanged
	}

	res := s.newResult(partialResult, hash)

	outcome := OutcomePass
	now := datetime.NewCpsTime()
	foundPrev := false
	for _, rule := range rules {
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

		if !s.isRuleActive(rule, now) {
			continue
		}

		matched, regexMatch := s.matchRule(rule, event)
		if !matched {
			continue
		}

		applicator, found := s.ruleApplicatorContainer.Get(rule.Type)
		if !found {
			s.logger.Warn().Str("rule_id", rule.ID).Str("rule_type", rule.Type).Msg("Event filter rule service: RuleApplicator doesn't exist")
			continue
		}

		externalData, ruleExternalRequestCount, externalDataRequest, err := s.resolveExternalData(ctx, rule, event, regexMatch, partialResult)
		if err != nil {
			outcome = rule.Config.OnFailure

			continue
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

func (s *ruleService) getRules() ([]ParsedRule, string) {
	s.rulesMutex.RLock()
	defer s.rulesMutex.RUnlock()

	return s.rules, s.rulesHash
}

// newResult creates the result which accumulates the state of the already processed rules.
func (s *ruleService) newResult(partialResult ServiceResult, rulesHash string) ServiceResult {
	res := ServiceResult{
		UpdatedEntityInfos:      partialResult.UpdatedEntityInfos,
		ExecutedEnrichRuleCount: partialResult.ExecutedEnrichRuleCount,
		ExternalRequestCount:    partialResult.ExternalRequestCount,
		RulesHash:               rulesHash,
	}
	if res.UpdatedEntityInfos == nil {
		res.UpdatedEntityInfos = make(map[string]UpdatedValue)
	}
	if res.ExternalRequestCount == nil {
		res.ExternalRequestCount = make(map[string]int64)
	}

	return res
}

// isRuleActive checks if now is in one of the resolved active periods of the rule.
func (s *ruleService) isRuleActive(rule ParsedRule, now datetime.CpsTime) bool {
	if rule.ResolvedStart == nil || rule.ResolvedStop == nil {
		return true
	}

	for _, exdate := range rule.ResolvedExdates {
		if now.After(exdate.Begin) && now.Before(exdate.End) {
			return false
		}
	}

	if now.Before(*rule.ResolvedStart) {
		return false
	}

	if now.After(*rule.ResolvedStop) {
		if rule.NextResolvedStart == nil || rule.NextResolvedStop == nil {
			return false
		}

		if now.Before(*rule.NextResolvedStart) || now.After(*rule.NextResolvedStop) {
			return false
		}
	}

	return true
}

// matchRule checks if the patterns of the rule match the event and returns the found regex matches.
func (s *ruleService) matchRule(rule ParsedRule, event *types.Event) (bool, RegexMatch) {
	if len(rule.EntityPattern) == 0 && len(rule.EventPattern) == 0 {
		s.logDebug(rule.ID, event, "rule is not matched")
		s.failureService.Add(rule.ID, rule.Description, rule.Updated, FailureTypeInvalidPattern, "missing entity and event patterns", nil)

		return false, RegexMatch{}
	}

	regexMatch := RegexMatch{}

	if len(rule.EventPattern) > 0 {
		matched, eventRegexMatches, err := match.MatchEventPatternWithRegexMatches(rule.EventPattern, event)
		if err != nil {
			s.logger.Err(err).Str("rule_id", rule.ID).Msg("Event filter rule service: invalid event pattern")
			s.failureService.Add(rule.ID, rule.Description, rule.Updated, FailureTypeInvalidPattern, "invalid event pattern: "+err.Error(), nil)

			return false, RegexMatch{}
		}

		if !matched {
			s.logDebug(rule.ID, event, "rule is not matched")

			return false, RegexMatch{}
		}

		regexMatch.EventRegexMatches = eventRegexMatches
	}

	if len(rule.EntityPattern) > 0 {
		if event.Entity == nil {
			s.logDebug(rule.ID, event, "entity is missing")

			return false, RegexMatch{}
		}

		matched, entityRegexMatches, err := match.MatchEntityPatternWithRegexMatches(rule.EntityPattern, event.Entity)
		if err != nil {
			s.logger.Err(err).Str("rule_id", rule.ID).Msg("Event filter rule service: invalid entity pattern")
			s.failureService.Add(rule.ID, rule.Description, rule.Updated, FailureTypeInvalidPattern, "invalid entity pattern: "+err.Error(), nil)

			return false, RegexMatch{}
		}

		if !matched {
			s.logDebug(rule.ID, event, "rule is not matched")

			return false, RegexMatch{}
		}

		regexMatch.Entity = entityRegexMatches
	}

	s.logDebug(rule.ID, event, "rule is matched")

	return true, regexMatch
}

func (s *ruleService) logDebug(ruleID string, event *types.Event, msg string) {
	if !event.Debug {
		return
	}

	s.logger.Info().
		Str("rule", ruleID).
		Str("event_type", event.EventType).
		Str("entity", event.GetEID()).
		Msg("event filter rule service: " + msg)
}

// resolveExternalData returns the external data of the rule either from the response of the request
// on which the event was parked or by fetching it.
func (s *ruleService) resolveExternalData(
	ctx context.Context,
	rule ParsedRule,
	event *types.Event,
	regexMatch RegexMatch,
	partialResult ServiceResult,
) (map[string]any, map[string]int64, *ExternalDataRequest, error) {
	resp := partialResult.ExternalDataResponse
	if resp == nil || resp.RuleID != rule.ID {
		externalData, requestCount, request, err := s.getExternalData(ctx, rule, event, regexMatch)
		if err != nil {
			s.logger.Err(err).Str("rule_id", rule.ID).Str("rule_type", rule.Type).Msg("Event filter rule service: failed to fetch external data")

			return nil, nil, nil, err
		}

		return externalData, requestCount, request, nil
	}

	if resp.Err != nil {
		failReason := "external data \"" + rule.ExternalData[resp.ErrIndex].Reference + "\" cannot be fetched: " + resp.Err.Error()
		s.failureService.Add(rule.ID, rule.Description, rule.Updated, FailureTypeExternalDataAPI, failReason, event)

		return nil, nil, nil, resp.Err
	}

	externalData := partialResult.ExternalData
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

	return externalData, nil, nil, nil
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

// rulesHash builds a stable digest of the ruleset from each rule's id and updated time.
func rulesHash(rules []ParsedRule) string {
	h := sha256.New()
	for _, r := range rules {
		h.Write([]byte(r.ID))
		h.Write([]byte{0})
		h.Write([]byte(strconv.FormatInt(r.Updated.Unix(), 10)))
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
