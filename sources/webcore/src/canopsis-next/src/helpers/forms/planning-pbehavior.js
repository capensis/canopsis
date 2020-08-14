import Vue from 'vue';
import moment from 'moment-timezone';
import { isObject, isString, cloneDeep, isUndefined } from 'lodash';
import { CalendarEvent, DaySpan, Op, Schedule } from 'dayspan';

import uid from '@/helpers/uid';
import { convertDateToTimestampByTimezone } from '@/helpers/date';

export function pbehaviorToForm(pbehavior = {}, timezone, filter) {
  let rrule = pbehavior.rrule || null;

  if (pbehavior.rrule && isObject(pbehavior.rrule)) {
    ({ rrule } = pbehavior.rrule);
  }

  const resultFilter = filter || pbehavior.filter || {};

  return {
    rrule,
    timezone,
    _id: pbehavior._id || uid('pbehavior'),
    enabled: isUndefined(pbehavior.enabled) ? true : pbehavior.enabled,
    author: pbehavior.author || '',
    name: pbehavior.name || '',
    type: pbehavior.type || '',
    reason: pbehavior.reason || '',
    filter: isString(resultFilter) ? JSON.parse(resultFilter) : cloneDeep(resultFilter),
    comments: cloneDeep(pbehavior.comments || []), // TODO: add key
    exdates: cloneDeep(pbehavior.exdates || []), // TODO: add key
  };
}

export function formToPbehavior(form, timezone) {
  return {
    ...form,

    comments: [],
    tstart: convertDateToTimestampByTimezone(form.tstart, timezone),
    tstop: convertDateToTimestampByTimezone(form.tstop, timezone),
  };
}

export function calendarEventToPbehaviorForm(calendarEvent, timezone, filter) {
  const { pbehavior, cachedForm = {} } = calendarEvent.data || {};

  const form = {
    ...pbehaviorToForm(pbehavior, timezone, filter),
    ...cachedForm,
  };

  form.tstart = calendarEvent.start.date.toDate();
  form.tstop = calendarEvent.schedule.durationUnit === 'days'
    ? moment(calendarEvent.end.date).subtract(1, 'second').toDate()
    : calendarEvent.end.date.toDate();

  return form;
}

export function formToCalendarEvent(form, calendarEvent, timezone) {
  const span = new DaySpan(calendarEvent.start, calendarEvent.end);

  const schedule = calendarEvent.fullDay
    ? Schedule.forDay(span.start, span.days(Op.UP))
    : Schedule.forSpan(span);

  const details = { ...calendarEvent.data, pbehavior: formToPbehavior(form, timezone) };
  const event = Vue.$dayspan.createEvent(details, schedule);

  event.id = calendarEvent.event.id;

  return new CalendarEvent(calendarEvent.id, event, span, calendarEvent.day);
}
