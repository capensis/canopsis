import moment from 'moment';
import { omit, pick } from 'lodash';
import { ACTION_TYPES, DURATION_UNITS, ACTION_AUTHOR } from '@/constants';

import { unsetSeveralFieldsWithConditions } from '@/helpers/immutable';
import { generateAction } from '@/helpers/entities';
import {
  pbehaviorToForm,
  pbehaviorToComments,
  pbehaviorToExdates,
  formToPbehavior,
  commentsToPbehaviorComments,
  exdatesToPbehaviorExdates,
} from '@/helpers/forms/pbehavior';

export function actionToForm(action) {
  const data = generateAction();

  if (!action) {
    return data;
  }

  data.generalParameters = pick(action, ['_id', 'type', 'hook']);

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

      data.snoozeParameters.duration = duration;
    }

    if (action.parameters && action.parameters.message) {
      data.snoozeParameters.message = action.parameters.message;
    }
  }

  // If action's type is "pbehavior", get pbehavior parameters
  if (action.type === ACTION_TYPES.pbehavior) {
    if (action.parameters) {
      data.pbehaviorParameters.general = omit(pbehaviorToForm(action.parameters), ['filter']);

      if (action.parameters.comments) {
        data.pbehaviorParameters.comments = pbehaviorToComments(action.parameters);
      }

      if (action.parameters.exdate) {
        data.pbehaviorParameters.exdate = pbehaviorToExdates(action.parameters);
      }
    }
  }

  return data;
}

export function formToAction({ generalParameters = {}, pbehaviorParameters = {}, snoozeParameters = {} }) {
  let data = { ...generalParameters };

  const patternsCondition = value => !value || !value.length;

  data = unsetSeveralFieldsWithConditions(data, {
    'hook.event_patterns': patternsCondition,
    'hook.alarm_patterns': patternsCondition,
    'hook.entity_patterns': patternsCondition,
  });

  if (generalParameters.type === ACTION_TYPES.snooze) {
    const duration = moment.duration(
      parseInt(snoozeParameters.duration.duration, 10),
      snoozeParameters.duration.durationType,
    ).asSeconds();

    data.parameters = {
      ...snoozeParameters,
      duration,
    };
  } else if (generalParameters.type === ACTION_TYPES.pbehavior) {
    const pbehavior = formToPbehavior(pbehaviorParameters.general);

    pbehavior.comments =
      commentsToPbehaviorComments(pbehaviorParameters.comments);
    pbehavior.exdate = exdatesToPbehaviorExdates(pbehaviorParameters.exdate);

    data.parameters = { ...pbehavior };
  }

  data.parameters.author = ACTION_AUTHOR;

  return data;
}
