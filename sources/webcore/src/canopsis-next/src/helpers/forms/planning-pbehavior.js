import Vue from 'vue';
import moment from 'moment-timezone';
import { omit, isObject, isString, cloneDeep, isUndefined } from 'lodash';
import { CalendarEvent, DaySpan, Op, Schedule } from 'dayspan';

import uid from '@/helpers/uid';
import {
  convertDateToTimestampByTimezone,
  convertTimestampToMomentByTimezone,
} from '@/helpers/date/date';
import { addKeyInEntity, getIdFromEntity, removeKeyFromEntity } from '@/helpers/entities';

/**
 * Clear exdate entity and convert to request.
 *
 * @param {Array} exdates
 * @return {{end: number, type: string, begin: number }[]}
 */
export const exdatesToRequest = (exdates = []) => exdates.map(({ type, begin, end }) => ({
  type: getIdFromEntity(type),
  begin,
  end,
}));

/**
 * Convert exdate to form
 *
 * @param {Object} exdate
 * @param {string} [timezone = moment.tz.guess()]
 * @return {{end: Date, type: Object, begin: Date }[]}
 */
export const exdateToForm = (exdate, timezone = moment.tz.guess()) => ({
  ...exdate,
  key: uid(),
  begin: convertTimestampToMomentByTimezone(exdate.begin, timezone).toDate(),
  end: convertTimestampToMomentByTimezone(exdate.end, timezone).toDate(),
});

/**
 * Convert exdate form to exdate
 *
 * @param {Object} formExdate
 * @param {string} [timezone = moment.tz.guess()]
 * @return {{type: string, begin: number, end: number}}
 */
export const formToExdate = (formExdate, timezone = moment.tz.guess()) => ({
  type: formExdate.type,
  begin: convertDateToTimestampByTimezone(formExdate.begin, timezone),
  end: convertDateToTimestampByTimezone(formExdate.end, timezone),
});

/**
 * Convert exceptions to exceptions id array.
 *
 * @param {Array} exceptions
 * @return {string[]}
 */
export const exceptionsToRequest = (exceptions = []) => exceptions.map(exception => getIdFromEntity(exception));

/**
 * Convert pbehavior entity to form data.
 *
 * @param {Object} [pbehavior = {}]
 * @param {string|Object} [filter = null]
 * @param {string} [timezone = moment.tz.guess()]
 * @return {Object}
 */
export const pbehaviorToForm = (
  pbehavior = {},
  filter = null,
  timezone = moment.tz.guess(),
) => {
  let rrule = pbehavior.rrule || null;

  if (pbehavior.rrule && isObject(pbehavior.rrule)) {
    ({ rrule } = pbehavior.rrule);
  }

  const resultFilter = filter || pbehavior.filter || {};

  return {
    rrule,
    _id: pbehavior._id || uid('pbehavior'),
    enabled: isUndefined(pbehavior.enabled) ? true : pbehavior.enabled,
    author: pbehavior.author || '',
    name: pbehavior.name || '',
    type: cloneDeep(pbehavior.type),
    reason: cloneDeep(pbehavior.reason),
    tstart: pbehavior.tstart ? convertTimestampToMomentByTimezone(pbehavior.tstart, timezone).toDate() : null,
    tstop: pbehavior.tstop ? convertTimestampToMomentByTimezone(pbehavior.tstop, timezone).toDate() : null,
    filter: isString(resultFilter) ? JSON.parse(resultFilter) : cloneDeep(resultFilter),
    exceptions: pbehavior.exceptions ? addKeyInEntity(cloneDeep(pbehavior.exceptions)) : [],
    comments: pbehavior.comments ? addKeyInEntity(cloneDeep(pbehavior.comments)) : [],
    exdates: pbehavior.exdates ? pbehavior.exdates.map(exdate => exdateToForm(exdate, timezone)) : [],
  };
};

/**
 * @param {Object} pbehavior
 * @returns {Object}
 */
export const pbehaviorToDuplicateForm = pbehavior => ({
  ...pbehavior,
  comments: pbehavior.comments.map(({ message, author }) => ({ message, author })),
});

/**
 * Convert form to pbehavior entity.
 *
 * @param {Object} form
 * @param {string} timezone
 * @return {Object}
 */
export const formToPbehavior = (form, timezone = moment.tz.guess()) => ({
  ...form,

  enabled: isUndefined(form.enabled) ? true : form.enabled,
  reason: form.reason,
  type: form.type,
  comments: removeKeyFromEntity(form.comments),
  exdates: form.exdates ? form.exdates.map(exdateForm => formToExdate(exdateForm, timezone)) : [],
  exceptions: removeKeyFromEntity(form.exceptions),
  tstart: form.tstart ? convertDateToTimestampByTimezone(form.tstart, timezone) : null,
  tstop: form.tstop ? convertDateToTimestampByTimezone(form.tstop, timezone) : null,
});

/**
 * Convert calendar event to pbehavior form data
 *
 * @param {CalendarEvent} calendarEvent
 * @param {string|Object} filter
 * @param {string} [timezone = moment.tz.guess()]
 * @return {Object}
 */
export const calendarEventToPbehaviorForm = (
  calendarEvent,
  filter,
  timezone = moment.tz.guess(),
) => {
  const {
    start,
    end,
    schedule,
    data: { pbehavior, cachedForm = {} },
  } = calendarEvent;

  const form = {
    ...pbehaviorToForm(pbehavior, filter, timezone),
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
 * @param {Object} form
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

  return new CalendarEvent(calendarEvent.id, event, span, calendarEvent.day);
};

/**
 * Convert pbehavior to request data.
 *
 * @param {Object} pbehavior
 * @return {Object}
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
