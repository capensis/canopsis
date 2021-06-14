package pbehavior_legacy

import (
	"time"

	"git.canopsis.net/canopsis/go-engines/lib/canopsis/types"
	"github.com/globalsign/mgo/bson"
	"github.com/rs/zerolog"
)

type service struct {
	adapter Adapter
	logger  zerolog.Logger
}

// NewService gives the correct pbehavior adapter.
func NewService(pbehaviorAdapter Adapter, logger zerolog.Logger) Service {
	service := service{
		adapter: pbehaviorAdapter,
		logger:  logger,
	}
	return &service
}

// Insert a new pbehavior.
func (s *service) Insert(pbehavior types.PBehaviorLegacy) error {
	err := s.adapter.Insert(pbehavior)
	return err
}

// Update updates an existing pbehavior or creates a new one in db.
func (s *service) Update(pbehavior types.PBehaviorLegacy) error {
	err := s.adapter.Update(pbehavior)
	return err
}

// Remove removes an existing pbehavior in db.
func (s *service) Remove(pbehavior types.PBehaviorLegacy) error {
	err := s.adapter.RemoveId(pbehavior.ID)
	return err
}

// Get retreive pbehaviors from db.
func (s *service) Get(filter bson.M) ([]types.PBehaviorLegacy, error) {
	pbhs, err := s.adapter.Get(filter)
	return pbhs, err
}

// GetByEntityIds retrieves every pbehaviors on a list of entities
func (s *service) GetByEntityIds(eids []string, enabled bool) ([]types.PBehaviorLegacy, error) {
	pbhs, err := s.adapter.GetByEntityIds(eids, enabled)
	return pbhs, err
}

// AlarmHasPBehavior checks if an alarm has an enabled pbehavior on it. Be
// careful, this method does not check if a pbehavior is active.
func (s *service) AlarmHasPBehavior(alarm types.Alarm) bool {

	pbs, err := s.Get(bson.M{"enabled": true, "eids": alarm.EntityID})
	if err != nil {
		s.logger.Warn().Err(err).Msg("failed to get enabled pbehaviors")
		return false
	}

	for _, pb := range pbs {
		if pb.IsImpacted(alarm.EntityID) && pb.Enabled {
			// FunFact: Alarm.EntityID == Entity.ID
			return true
		}
	}
	return false
}

// HasActivePBehavior returns true if the entity has an active pbehavior.
func (s *service) HasActivePBehavior(entityID string) bool {
	pbehaviors, err := s.GetByEntityIds([]string{entityID}, true)
	if err != nil {
		s.logger.Warn().Err(err).Msg("failed to get enabled pbehaviors")
		return false
	}

	for _, pbehavior := range pbehaviors {
		active, err := pbehavior.IsActive(time.Now())
		if err != nil {
			s.logger.Warn().Err(err).Str("pbehavior_id", pbehavior.ID).Msg("failed to check if pbehavior is active")
			continue
		}
		if active {
			return true
		}
	}

	return false
}
