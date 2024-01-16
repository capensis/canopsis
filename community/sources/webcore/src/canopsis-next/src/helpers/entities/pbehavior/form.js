import { cloneDeep, isObject, omit } from 'lodash';

import { COLORS } from '@/config';
import {
  PATTERNS_FIELDS,
  PBEHAVIOR_ORIGINS,
  PBEHAVIOR_TYPE_TYPES,
  WEATHER_ENTITY_PBEHAVIOR_DEFAULT_TITLE,
} from '@/constants';

import { uid } from '@/helpers/uid';
import {
  convertDateToDateObjectByTimezone,
  convertDateToTimestampByTimezone,
  getLocaleTimezone,
  getNowTimestamp,
  isEndOfDay,
  isStartOfDay,
} from '@/helpers/date/date';
import { addKeyInEntities, getIdFromEntity, removeKeyFromEntities } from '@/helpers/array';
import { filterPatternsToForm, formFilterToPatterns } from '@/helpers/entities/filter/form';

/**
 * @typedef {Object} PbehaviorReason
 * @property {string} _id
 * @property {boolean} deletable
 * @property {string} description
 * @property {string} name
 */

/**
 * @typedef {Object} PbehaviorComment
 * @property {string} _id
 * @property {string} author
 * @property {string} message
 * @property {number} ts
 */

/**
 * @typedef {PbehaviorComment & ObjectKey} PbehaviorCommentForm
 */

/**
 * @typedef {Object} PbehaviorCommentDuplicate
 * @property {string} author
 * @property {string} message
 */

/**
 * @typedef {Object} PbehaviorExdate
 * @property {number} begin
 * @property {number} end
 * @property {PbehaviorType} type
 */

/**
 * @typedef {PbehaviorExdate & ObjectKey} PbehaviorExdateForm
 * @property {Date} begin
 * @property {Date} end
 */

/**
 * @typedef {Object} PbehaviorExdateRequest
 * @property {string} type
 */

/**
 * @typedef {FilterPatterns} Pbehavior
 * @property {string} _id
 * @property {string} author
 * @property {boolean} enabled
 * @property {string} name
 * @property {string} rrule
 * @property {boolean} start_on_trigger
 * @property {Duration} duration
 * @property {number} tstart
 * @property {number} tstop
 * @property {string} color
 * @property {PbehaviorType} type
 * @property {PbehaviorReason} reason
 * @property {PbehaviorComment[]} comments
 * @property {PbehaviorException[]} exceptions
 * @property {PbehaviorExdate[]} exdates
 */

/**
 * @typedef {Pbehavior} PbehaviorForm
 * @property {PbehaviorCommentForm[]} comments
 * @property {PbehaviorExceptionForm[]} exceptions
 * @property {PbehaviorExdateForm[]} exdates
 * @property {Duration} duration
 * @property {FilterPatternsForm} patterns
 */

/**
 * @typedef {Pbehavior} PbehaviorDuplicate
 * @property {PbehaviorCommentDuplicate[]} comments
 */

/**
 * @typedef {Pbehavior} PbehaviorRequest
 * @property {string} type
 * @property {string} reason
 * @property {string[]} exceptions
 * @property {PbehaviorExdateRequest[]} exdates
 */

/**
 * Check if pbehavior is paused
 *
 * @param {Pbehavior} pbehavior
 * @return {boolean}
 */
export const isPausedPbehavior = pbehavior => pbehavior.type.type === PBEHAVIOR_TYPE_TYPES.pause;

/**
 * Check is not active pbehavior type
 *
 * @param {string} type
 * @returns {boolean}
 */
export const isNotActivePbehaviorType = type => [
  PBEHAVIOR_TYPE_TYPES.pause,
  PBEHAVIOR_TYPE_TYPES.inactive,
  PBEHAVIOR_TYPE_TYPES.maintenance,
].includes(type);

/**
 * Check if pbehaviors have a paused type
 *
 * @param {Pbehavior[]} pbehaviors
 * @return {boolean}
 */
export const hasPausedPbehavior = pbehaviors => pbehaviors.some(isPausedPbehavior);

/**
 * Clear exdate entity and convert to request.
 *
 * @param {PbehaviorExdate[]} [exdates = []]
 * @return {PbehaviorExdateRequest[]}
 */
export const exdatesToRequest = (exdates = []) => exdates.map(({ type, begin, end }) => ({
  type: getIdFromEntity(type),
  begin,
  end,
}));

/**
 * Convert exdate to form
 *
 * @param {PbehaviorExdate} exdate
 * @param {string} [timezone = getLocaleTimezone()]
 * @return {PbehaviorExdateForm}
 */
export const exdateToForm = (exdate, timezone = getLocaleTimezone()) => ({
  ...exdate,
  key: uid(),
  begin: convertDateToDateObjectByTimezone(exdate.begin, timezone),
  end: convertDateToDateObjectByTimezone(isStartOfDay(exdate.end) ? exdate.end - 1 : exdate.end, timezone),
});

/**
 * Convert exdates array to exdates form array
 *
 * @param {PbehaviorExdate[]} exdates
 * @param {string} timezone
 * @returns {PbehaviorExdateForm[]}
 */
export const exdatesToForm = (exdates, timezone) => (
  exdates ? exdates.map(exdate => exdateToForm(exdate, timezone)) : []
);

/**
 * Convert exdate form to exdate
 *
 * @param {PbehaviorExdateForm} formExdate
 * @param {string} [timezone = getLocaleTimezone()]
 * @return {PbehaviorExdate}
 */
export const formToExdate = (formExdate, timezone = getLocaleTimezone()) => ({
  type: formExdate.type,
  begin: convertDateToTimestampByTimezone(formExdate.begin, timezone),
  end: convertDateToTimestampByTimezone(formExdate.end, timezone),
});

/**
 * Convert exceptions to exceptions id array.
 *
 * @param {PbehaviorException[]} exceptions
 * @return {string[]}
 */
export const exceptionsToRequest = (exceptions = []) => exceptions.map(exception => getIdFromEntity(exception));

/**
 * Convert exceptions array to exceptions form array
 *
 * @param {PbehaviorException[]} exceptions
 * @returns {PbehaviorExceptionForm[]}
 */
export const exceptionsToForm = (exceptions = []) => addKeyInEntities(cloneDeep(exceptions));

/**
 * Convert pbehavior entity to form data.
 *
 * @param {Pbehavior} [pbehavior = {}]
 * @param {string|Object} [entityPattern]
 * @param {string} [timezone = getLocaleTimezone()]
 * @return {PbehaviorForm}
 */
export const pbehaviorToForm = (
  pbehavior = {},
  entityPattern,
  timezone = getLocaleTimezone(),
) => {
  let rrule = pbehavior.rrule ?? null;

  if (pbehavior.rrule && isObject(pbehavior.rrule)) {
    ({ rrule } = pbehavior.rrule);
  }

  const patterns = filterPatternsToForm(
    entityPattern
      ? { entity_pattern: entityPattern }
      : pbehavior,
    [PATTERNS_FIELDS.entity],
  );

  return {
    rrule,
    patterns,
    _id: pbehavior._id ?? uid('pbehavior'),
    color: pbehavior.color ?? '',
    enabled: pbehavior.enabled ?? true,
    name: pbehavior.name ?? '',
    type: cloneDeep(pbehavior.type),
    reason: cloneDeep(pbehavior.reason),
    tstart: pbehavior.tstart ? convertDateToDateObjectByTimezone(pbehavior.tstart, timezone) : null,
    tstop: pbehavior.tstop ? convertDateToDateObjectByTimezone(pbehavior.tstop, timezone) : null,
    comments: pbehavior.comments ? addKeyInEntities(cloneDeep(pbehavior.comments)) : [],
    exceptions: exceptionsToForm(pbehavior.exceptions),
    exdates: exdatesToForm(pbehavior.exdates, timezone),
  };
};

/**
 * @param {Pbehavior} pbehavior
 * @returns {PbehaviorDuplicate}
 */
export const pbehaviorToDuplicateForm = pbehavior => ({
  ...pbehavior,
  comments: pbehavior.comments.map(({ message, author }) => ({ message, author })),
});

/**
 * Convert form exceptions to exceptions object
 *
 * @param {PbehaviorExceptionForm[]} exceptions
 * @returns {PbehaviorException[]}
 */
export const formExceptionsToExceptions = exceptions => removeKeyFromEntities(exceptions);

/**
 * Convert form exdates to exdates object
 *
 * @param {PbehaviorExdateForm[]} exdates
 * @param {string} timezone
 * @returns {PbehaviorExdate[]}
 */
export const formExdatesToExdates = (exdates = [], timezone) => exdates.map(
  exdateForm => formToExdate(exdateForm, timezone),
);

/**
 * Convert form to pbehavior entity.
 *
 * @param {PbehaviorForm} form
 * @param {string} timezone
 * @return {Pbehavior}
 */
export const formToPbehavior = (form, timezone = getLocaleTimezone()) => ({
  ...omit(form, ['patterns']),

  enabled: form.enabled ?? true,
  reason: form.reason,
  type: form.type,
  comments: removeKeyFromEntities(form.comments),
  exdates: formExdatesToExdates(form.exdates, timezone),
  exceptions: formExceptionsToExceptions(form.exceptions),
  tstart: form.tstart ? convertDateToTimestampByTimezone(form.tstart, timezone) : null,
  tstop: form.tstop ? convertDateToTimestampByTimezone(form.tstop, timezone) : null,
  ...formFilterToPatterns(form.patterns),
});

/**
 * Convert calendar event to pbehavior form data
 *
 * @param {Object} event
 * @param {Array} entityPattern
 * @param {string} [timezone = getLocaleTimezone()]
 * @return {PbehaviorForm}
 */
export const calendarEventToPbehaviorForm = (
  event,
  entityPattern,
  timezone = getLocaleTimezone(),
) => {
  const {
    start,
    end,
    data: { cachedForm = {}, pbehavior },
  } = event;

  const form = {
    ...pbehaviorToForm(pbehavior, entityPattern, timezone),
    ...cachedForm,
  };

  form.tstart = start;

  if (!pbehavior || pbehavior.tstop) {
    if (event.durationUnit === 'days') {
      if (end.date.diff(start.date, 'hours') <= 24) {
        form.tstop = start.date.clone().endOf('day').toDate();
      } else {
        form.tstop = end.date.clone().subtract(1, 'millisecond').toDate();
      }
    } else {
      form.tstop = end;
    }
  }

  return form;
};

/**
 * Check pbehavior is full day
 *
 * @param {LocalDate} start
 * @param {LocalDate} stop
 * @return {boolean}
 */
export const isFullDayEvent = (start, stop) => {
  const noEnding = start && !stop;

  return isStartOfDay(start) && (noEnding || isEndOfDay(stop));
};

/**
 * Convert form to calendar event.
 *
 * @param {PbehaviorForm} form
 * @param {Object} event
 * @param {string} timezone
 * @return {Object}
 */
export const formToCalendarEvent = (form, event, timezone) => {
  const timed = !isFullDayEvent(form.tstart, form.tstop);

  return {
    ...event,
    timed,
    color: form.color,

    pbehavior: formToPbehavior(form, timezone),
  };
};

/**
 * Convert pbehavior to request data.
 *
 * @param {Pbehavior} pbehavior
 * @return {PbehaviorRequest}
 */
export const pbehaviorToRequest = (pbehavior) => {
  const result = omit(pbehavior, ['type', 'reason', 'exdates']);

  result.type = getIdFromEntity(pbehavior.type);
  result.reason = getIdFromEntity(pbehavior.reason);

  if (pbehavior.exdates) {
    result.exdates = exdatesToRequest(pbehavior.exdates);
  }

  if (pbehavior.exceptions) {
    result.exceptions = exceptionsToRequest(pbehavior.exceptions);
  }

  return result;
};
/**
 * Create downtime pbehavior, without stop time
 *
 * @param {Entity} entity
 * @param {PbehaviorReason} reason
 * @param {string} comment
 * @param {PbehaviorType} type
 * @return {Pbehavior}
 */
export const createDowntimePbehavior = ({ entity, reason, comment, type }) => pbehaviorToRequest(formToPbehavior({
  reason,
  type,
  origin: PBEHAVIOR_ORIGINS.serviceWeather,
  color: COLORS.secondary,
  name: `${WEATHER_ENTITY_PBEHAVIOR_DEFAULT_TITLE}-${entity.name}-${uid()}`,
  tstart: getNowTimestamp(),
  comment,
  entity: entity._id,
}));

/**
 * Get color for pbehavior
 *
 * @param pbehavior
 * @returns {string}
 */
export const getPbehaviorColor = (pbehavior = {}) => pbehavior.color || pbehavior.type?.color || COLORS.secondary;
