import Vue from 'vue';
import { omit, isObject, cloneDeep } from 'lodash';
import {
  CalendarEvent,
  DaySpan,
  Op,
  Schedule,
} from 'dayspan';

import { OLD_PATTERNS_FIELDS, PATTERNS_FIELDS } from '@/constants';

import { filterPatternsToForm, formFilterToPatterns } from '@/helpers/forms/filter';

import uid from '@/helpers/uid';
import {
  convertDateToTimestampByTimezone,
  convertDateToDateObjectByTimezone,
  getLocaleTimezone,
} from '@/helpers/date/date';
import { addKeyInEntities, getIdFromEntity, removeKeyFromEntities } from '@/helpers/entities';

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
 * @typedef {PbehaviorComment} PbehaviorCommentForm
 * @property {string} key
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
 * @typedef {PbehaviorExdate} PbehaviorExdateForm
 * @property {string} key
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
  end: convertDateToDateObjectByTimezone(exdate.end, timezone),
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
    [OLD_PATTERNS_FIELDS.mongoQuery],
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
 * @param {CalendarEvent} calendarEvent
 * @param {Array} entityPattern
 * @param {string} [timezone = getLocaleTimezone()]
 * @return {PbehaviorForm}
 */
export const calendarEventToPbehaviorForm = (
  calendarEvent,
  entityPattern,
  timezone = getLocaleTimezone(),
) => {
  const {
    start,
    end,
    schedule,
    data: { pbehavior, cachedForm = {} },
  } = calendarEvent;

  const form = {
    ...pbehaviorToForm(pbehavior, entityPattern, timezone),
    ...cachedForm,
  };

  form.tstart = start.date.toDate();

  if (!pbehavior || pbehavior.tstop) {
    if (schedule.durationUnit === 'days') {
      if (end.date.diff(start.date, 'hours') <= 24) {
        form.tstop = start.date.clone().endOf('day').toDate();
      } else {
        form.tstop = end.date.clone().subtract(1, 'millisecond').toDate();
      }
    } else {
      form.tstop = end.date.toDate();
    }
  }

  return form;
};

/**
 * Convert form to calendar event.
 *
 * @param {PbehaviorForm} form
 * @param {CalendarEvent} calendarEvent
 * @param {string} timezone
 * @return {CalendarEvent}
 */
export const formToCalendarEvent = (form, calendarEvent, timezone) => {
  const span = new DaySpan(calendarEvent.start, calendarEvent.end);

  const schedule = calendarEvent.fullDay
    ? Schedule.forDay(span.start, span.days(Op.UP))
    : Schedule.forSpan(span);

  const details = {
    ...calendarEvent.data,

    pbehavior: formToPbehavior(form, timezone),
  };

  const event = Vue.$dayspan.createEvent(details, schedule);

  event.id = calendarEvent.event.id;

  return new CalendarEvent(
    calendarEvent.id,
    event,
    schedule.getSingleEventSpan(),
    calendarEvent.day,
  );
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
