import moment from 'moment';
import { omit, pick } from 'lodash';
import { ACTION_TYPES, TIME_UNITS, ACTION_AUTHOR } from '@/constants';

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

/**
 * If action's type is "snooze", get snooze parameters
 *
 * @param {Object} [parameters={}]
 * @returns {Object}
 */
function actionSnoozeParametersToForm(parameters = {}) {
  const data = {};

  if (parameters.duration) {
    const durationUnits = [
      TIME_UNITS.year,
      TIME_UNITS.month,
      TIME_UNITS.week,
      TIME_UNITS.week,
      TIME_UNITS.day,
      TIME_UNITS.hour,
      TIME_UNITS.minute,
      TIME_UNITS.second,
    ];

    // Check for the lowest possible unit to convert the duration in.
    const durationType = durationUnits.find(unit => moment.duration(parameters.duration, 'seconds').as(unit) % 1 === 0);

    data.duration = {
      duration: moment.duration(parameters.duration, 'seconds').as(durationType),
      durationType,
    };
  }

  if (parameters && parameters.message) {
    data.message = parameters.message;
  }

  return data;
}

/**
 * If action's type is "pbehavior", get pbehavior parameters
 *
 * @param {Object} [parameters={}]
 * @returns {Object}
 */
function actionPbehaviorParametersToForm(parameters = {}) {
  const data = {};

  data.general = omit(pbehaviorToForm(parameters), ['filter']);

  if (parameters.comments) {
    data.comments = pbehaviorToComments(parameters);
  }

  if (parameters.exdate) {
    data.exdate = pbehaviorToExdates(parameters);
  }

  return data;
}

/**
 * Prepare form object from action object
 *
 * @param {Object} [action]
 * @returns {Object}
 */
export function actionToForm(action) {
  const data = generateAction();

  if (!action) {
    return data;
  }

  data.generalParameters = pick(action, ['_id', 'type', 'hook']);

  switch (action.type) {
    case ACTION_TYPES.snooze:
      data.snoozeParameters = {
        ...data.snoozeParameters,
        ...actionSnoozeParametersToForm(action.parameters),
      };
      break;
    case ACTION_TYPES.pbehavior:
      data.pbehaviorParameters = {
        ...data.pbehaviorParameters,
        ...actionPbehaviorParametersToForm(action.parameters),
      };
      break;
    case ACTION_TYPES.changeState:
      data.changeStateParameters = {
        ...data.changeStateParameters,
        ...action.parameters,
      };
      break;
  }

  return data;
}

/**
 * Prepare action object by form object
 *
 * @param [generalParameters={}]
 * @param [pbehaviorParameters={}]
 * @param [snoozeParameters={}]
 * @param [changeStateParameters={}]
 * @returns {Object}
 */
export function formToAction({
  generalParameters = {},
  pbehaviorParameters = {},
  snoozeParameters = {},
  changeStateParameters = {},
}) {
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
  } else if (generalParameters.type === ACTION_TYPES.changeState) {
    data.parameters = { ...changeStateParameters };
  }

  data.parameters.author = ACTION_AUTHOR;

  return data;
}
