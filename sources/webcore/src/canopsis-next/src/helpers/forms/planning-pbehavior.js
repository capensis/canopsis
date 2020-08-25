import Vue from 'vue';
import moment from 'moment-timezone';
import { omit, isObject, isString, cloneDeep, isUndefined } from 'lodash';
import { CalendarEvent, DaySpan, Op, Schedule } from 'dayspan';

import uid from '@/helpers/uid';
import { convertDateToTimestampByTimezone, convertTimestampToMoment } from '@/helpers/date';
import { addKeyInEntity, getIdFromEntity, removeKeyFromEntity } from '@/helpers/entities';

/**
 * Clear exdate entity and convert to request.
 *
 * @param {Array} exdates
 * @return {{end: Number, type: String, begin: Number }[]}
 */
export const exdatesToRequest = exdates => exdates.map(({ type, begin, end }) => ({
  type: getIdFromEntity(type),
  begin: moment(begin).unix(),
  end: moment(end).unix(),
}));

/**
 * Convert exceptions to exceptions id array.
 *
 * @param {Array} exceptions
 * @return {String[]}
 */
export const exceptionsToRequest = exceptions => exceptions.map(exception => getIdFromEntity(exception));

/**
 * Convert pbehavior entity to form data.
 *
 * @param {Object} pbehavior
 * @param {String|Object} filter
 * @return {Object}
 */
export const pbehaviorToForm = (pbehavior = {}, filter = null) => {
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
    type: pbehavior.type, // TODO: add cloneDeep
    reason: pbehavior.reason, // TODO: add cloneDeep
    tstart: pbehavior.tstart ? convertTimestampToMoment(pbehavior.tstart).toDate() : null,
    tstop: pbehavior.tstop ? convertTimestampToMoment(pbehavior.tstop).toDate() : null,
    filter: isString(resultFilter) ? JSON.parse(resultFilter) : cloneDeep(resultFilter),
    exceptions: pbehavior.exceptions ? addKeyInEntity(cloneDeep(pbehavior.exceptions)) : [],
    comments: pbehavior.comments ? addKeyInEntity(cloneDeep(pbehavior.comments)) : [],
    exdates: pbehavior.exdates ? addKeyInEntity(cloneDeep(pbehavior.exdates)) : [], // TODO: convert timestamp to Date
  };
};

/**
 * Convert form to pbehavior entity.
 *
 * @param {Object} form
 * @param {String} timezone
 * @return {Object}
 */
export const formToPbehavior = (form, timezone) => ({
  ...form,

  reason: form.reason._id,
  type: form.type._id,
  comments: removeKeyFromEntity(form.comments),
  exdates: exdatesToRequest(form.exdates),
  exceptions: removeKeyFromEntity(form.exceptions).map(({ _id }) => _id),
  tstart: convertDateToTimestampByTimezone(form.tstart, timezone),
  tstop: form.tstop ? convertDateToTimestampByTimezone(form.tstop, timezone) : null,
});

/**
 * Convert calendar event to pbehavior form data
 *
 * @param {CalendarEvent} calendarEvent
 * @param {String|Object} filter
 * @return {Object}
 */
export const calendarEventToPbehaviorForm = (calendarEvent, filter) => {
  const { pbehavior, cachedForm = {} } = calendarEvent.data || {};

  const form = {
    ...pbehaviorToForm(pbehavior, filter),
    ...cachedForm,
  };

  form.tstart = calendarEvent.start.date.toDate();

  if (!pbehavior || pbehavior.tstop) {
    if (calendarEvent.schedule.durationUnit === 'days') {
      if (calendarEvent.end.date.diff(calendarEvent.start.date, 'days') <= 0) {
        form.tstop = calendarEvent.start.date.clone().endOf('day').toDate();
      } else {
        form.tstop = calendarEvent.end.date.clone().subtract(1, 'second').toDate();
      }
    } else {
      form.tstop = calendarEvent.end.date.toDate();
    }
  }

  return form;
};

/**
 * Convert form to calendar event.
 *
 * @param {Object} form
 * @param {CalendarEvent} calendarEvent
 * @param {String} timezone
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

  return result;
};
