package link

import (
	"context"
	"errors"
	"sync"
	"text/template"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/externaldata"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/js"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
	libtemplate "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/template"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var ErrNoRule = errors.New("rule not found")

// ParsedRule is a link rule with compiled templates and source code.
type ParsedRule struct {
	ID            string
	Type          string
	AlarmPattern  pattern.Alarm
	EntityPattern pattern.Entity
	ExternalData  []externaldata.ParsedRefParameters
	Links         []ParsedParameters
	CodeExecutor  js.Executor
}

type ParsedParameters struct {
	Parameters
	UrlTpl   *template.Template
	LabelTpl *template.Template
}

type RuleProvider interface {
	Load(ctx context.Context) error
	Rules() []ParsedRule
	Rule(id string) (ParsedRule, error)
}

func NewRuleProvider(client mongo.DbClient, tplExecutor libtemplate.Executor, logger zerolog.Logger) RuleProvider {
	return &ruleProvider{
		collection: client.Collection(mongo.LinkRuleMongoCollection),
		parser:     newRuleParser(tplExecutor, logger),
	}
}

type ruleProvider struct {
	collection mongo.DbCollection
	parser     *ruleParser

	mx    sync.RWMutex
	rules []ParsedRule
}

func (p *ruleProvider) Load(ctx context.Context) error {
	cursor, err := p.collection.Find(ctx, bson.M{"enabled": true}, options.Find().SetSort(bson.M{"_id": 1}))
	if err != nil {
		return err
	}

	defer cursor.Close(ctx)

	parsedRules := make([]ParsedRule, 0)
	for cursor.Next(ctx) {
		rule := Rule{}
		err = cursor.Decode(&rule)
		if err != nil {
			return err
		}

		parsedRule, ok := p.parser.parse(rule)
		if !ok {
			continue
		}

		parsedRules = append(parsedRules, parsedRule)
	}

	if err = cursor.Err(); err != nil {
		return err
	}

	p.mx.Lock()
	defer p.mx.Unlock()
	p.rules = parsedRules

	return nil
}

func (p *ruleProvider) Rules() []ParsedRule {
	p.mx.RLock()
	defer p.mx.RUnlock()

	return p.rules
}

func (p *ruleProvider) Rule(id string) (ParsedRule, error) {
	p.mx.RLock()
	defer p.mx.RUnlock()

	for _, rule := range p.rules {
		if rule.ID == id {
			return rule, nil
		}
	}

	return ParsedRule{}, ErrNoRule
}

func newRuleParser(tplExecutor libtemplate.Executor, logger zerolog.Logger) *ruleParser {
	return &ruleParser{
		tplExecutor: tplExecutor,
		logger:      logger,
	}
}

type ruleParser struct {
	tplExecutor libtemplate.Executor
	logger      zerolog.Logger
}

func (p *ruleParser) parse(rule Rule) (ParsedRule, bool) {
	parsedRule := ParsedRule{
		ID:            rule.ID,
		Type:          rule.Type,
		AlarmPattern:  rule.AlarmPattern,
		EntityPattern: rule.EntityPattern,
		ExternalData:  externaldata.ParseRefParameters(rule.ExternalData, p.tplExecutor),
	}

	if rule.SourceCode != "" {
		var err error
		parsedRule.CodeExecutor, err = js.Compile(rule.ID, rule.SourceCode)
		if err != nil {
			p.logger.Err(err).Str("rule", rule.ID).Msg("invalid source code in link rule")

			return ParsedRule{}, false
		}

		return parsedRule, true
	}

	parsedRule.Links = make([]ParsedParameters, len(rule.Links))
	for i, params := range rule.Links {
		if params.URL == "" {
			p.logger.Error().Str("rule", rule.ID).Msg("empty url template in link rule")

			return ParsedRule{}, false
		}

		parsedUrl := p.tplExecutor.Parse(params.URL)
		if parsedUrl.Err != nil {
			p.logger.Err(parsedUrl.Err).Str("rule", rule.ID).Msg("invalid url template in link rule")

			return ParsedRule{}, false
		}

		if params.Label == "" {
			p.logger.Error().Str("rule", rule.ID).Msg("empty label template in link rule")

			return ParsedRule{}, false
		}

		parsedLabel := p.tplExecutor.Parse(params.Label)
		if parsedLabel.Err != nil {
			p.logger.Err(parsedLabel.Err).Str("rule", rule.ID).Msg("invalid label template in link rule")

			return ParsedRule{}, false
		}

		parsedRule.Links[i].Parameters = params
		parsedRule.Links[i].UrlTpl = parsedUrl.Tpl
		parsedRule.Links[i].LabelTpl = parsedLabel.Tpl
	}

	return parsedRule, true
}
