import moment from 'moment';
import { omit, pick, isEmpty } from 'lodash';
import { ACTION_TYPES, ACTION_AUTHOR, ACTION_FORM_FIELDS_MAP_BY_TYPE } from '@/constants';

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
import { convertDurationToIntervalObject } from '@/helpers/date';
import { getConditionsForRemovingEmptyPatterns } from '@/helpers/forms/shared/patterns';

/**
 * If action's type is "snooze", get snooze parameters
 *
 * @param {Object} [parameters={}]
 * @returns {Object}
 */
function actionSnoozeParametersToForm(parameters = {}) {
  const data = {};

  if (parameters.duration) {
    const { unit, interval } = convertDurationToIntervalObject(parameters.duration);

    data.duration = {
      duration: interval,
      durationType: unit,
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
  data.generalParameters.enabled = action.enabled !== false;

  if (action.delay) {
    const [, value, unit] = action.delay.match(/^(\d+)(\w)$/);

    data.generalParameters.delay = {
      value: +value,
      unit,
    };
  }

  const actionToFormPrepareMap = {
    [ACTION_TYPES.snooze]: actionSnoozeParametersToForm,
    [ACTION_TYPES.pbehavior]: actionPbehaviorParametersToForm,
  };
  const prepareHandler = actionToFormPrepareMap[action.type];

  const parameters = prepareHandler
    ? prepareHandler(action.parameters)
    : action.parameters;

  const fieldKey = ACTION_FORM_FIELDS_MAP_BY_TYPE[action.type];

  data[fieldKey] = {
    ...data[fieldKey],
    ...parameters,
  };

  return data;
}

/**
 * Prepare snooze parameters from form
 *
 * @param snoozeParameters
 * @returns {{duration: number}}
 */
export function prepareSnoozeParameters({ snoozeParameters = {} }) {
  return ({
    ...snoozeParameters,
    duration: moment.duration(
      parseInt(snoozeParameters.duration.duration, 10),
      snoozeParameters.duration.durationType,
    ).asSeconds(),
  });
}

/**
 * Prepare pbehavior parameters from form
 *
 * @param pbehaviorParameters
 * @returns {{ tstart: number, exdate: Array, comments: Array, tstop: number }}
 */
export function preparePbehaviorParameters({ pbehaviorParameters = {} }) {
  const pbehavior = formToPbehavior(pbehaviorParameters.general);

  return {
    ...pbehavior,
    comments: commentsToPbehaviorComments(pbehaviorParameters.comments),
    exdate: exdatesToPbehaviorExdates(pbehaviorParameters.exdate),
  };
}

/**
 * Prepare action object by form object
 *
 * @param [generalParameters]
 * @param {Object} form
 * @param [form.pbehaviorParameters]
 * @param [form.snoozeParameters]
 * @param [form.changeStateParameters]
 * @param [form.ackParameters]
 * @param [form.ackremoveParameters]
 * @param [form.assocticketParameters]
 * @param [form.declareticketParameters]
 * @param [form.cancelParameters]
 * @returns {Object}
 */
export function formToAction({
  generalParameters = {},
  ...form
}) {
  const hasValue = v => !v;

  const data = unsetSeveralFieldsWithConditions(generalParameters, {
    ...getConditionsForRemovingEmptyPatterns([
      'hook.entity_patterns',
      'hook.alarm_patterns',
      'hook.event_patterns',
    ]),

    'delay.unit': hasValue,
    'delay.value': hasValue,
  });

  if (!isEmpty(data.delay)) {
    data.delay = `${data.delay.value}${data.delay.unit}`;
  } else {
    delete data.delay;
  }

  const formToActionPrepareMap = {
    [ACTION_TYPES.snooze]: prepareSnoozeParameters,
    [ACTION_TYPES.pbehavior]: preparePbehaviorParameters,
  };

  const prepareField = formToActionPrepareMap[generalParameters.type];
  const parameters = prepareField
    ? prepareField(form)
    : form[ACTION_FORM_FIELDS_MAP_BY_TYPE[generalParameters.type]];

  data.parameters = {
    ...parameters,
    author: ACTION_AUTHOR,
  };

  return data;
}
