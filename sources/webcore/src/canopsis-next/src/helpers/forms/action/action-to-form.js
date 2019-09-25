import moment from 'moment';
import { omit } from 'lodash';
import { ACTION_TYPES, DURATION_UNITS } from '@/constants';

import uuid from '@/helpers/uuid';
import pbehaviorToForm from '@/helpers/forms/pbehavior/pbehavior-to-form';
import pbehaviorToComments from '@/helpers/forms/pbehavior/pbehavior-to-comments';
import pbehaviorToExdates from '@/helpers/forms/pbehavior/pbehavior-to-exdates';

export default function (action = {}) {
  const defaultHook = {
    event_patterns: [],
    alarm_patterns: [],
    entity_patterns: [],
    triggers: [],
  };

  // Default 'snooze' action parameters
  const snoozeParameters = {
    message: '',
    duration: {
      duration: 1,
      durationType: DURATION_UNITS.minute.value,
    },
  };

  // Default 'pbehavior' action parameters
  const pbehaviorParameters = {
    general: {
      name: '',
      tstart: new Date(),
      tstop: new Date(),
      rrule: null,
      reason: '',
      type_: '',
    },
    comments: [],
    exdate: [],
  };

  // Get basic action parameters
  const generalParameters = {
    _id: action._id || uuid('action'),
    type: action.type || ACTION_TYPES.snooze,
    hook: action.hook || defaultHook,
  };

  // If action's type is "snooze", get snooze parameters
  if (action.type === ACTION_TYPES.snooze) {
    let duration = {
      duration: 1,
      durationType: DURATION_UNITS.minute.value,
    };

    if (action.parameters && action.parameters.duration) {
      const durationUnits = Object.values(DURATION_UNITS).map(unit => unit.value);

      // Check for the lowest possible unit to convert the duration in.
      const foundUnit = durationUnits.find(unit => moment.duration(action.parameters.duration, 'seconds').as(unit) % 1 === 0);

      duration = {
        duration: moment.duration(action.parameters.duration, 'seconds').as(foundUnit),
        durationType: foundUnit,
      };

      snoozeParameters.duration = duration;
    }

    if (action.parameters && action.parameters.message) {
      snoozeParameters.message = action.parameters.message;
    }
  }

  // If action's type is "pbehavior", get pbehavior parameters
  if (action.type === ACTION_TYPES.pbehavior) {
    if (action.parameters) {
      pbehaviorParameters.general = omit(pbehaviorToForm(action.parameters), ['filter']);

      if (action.parameters.comments) {
        pbehaviorParameters.comments = pbehaviorToComments(action.parameters.comments);
      }

      if (action.parameters.exdate) {
        pbehaviorParameters.exdate = pbehaviorToExdates(action.parameters.exdate);
      }
    }
  }

  return {
    generalParameters,
    snoozeParameters,
    pbehaviorParameters,
  };
}

