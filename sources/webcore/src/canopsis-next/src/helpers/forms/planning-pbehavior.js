import Vue from 'vue';
import moment from 'moment';
import { isObject, isString, cloneDeep, isUndefined } from 'lodash';
import { CalendarEvent, Day, DaySpan, Op, Schedule } from 'dayspan';

import uid from '@/helpers/uid';
import convertTimestampToMoment from '@/helpers/date';

export function pbehaviorToForm(pbehavior = {}) {
  let rrule = pbehavior.rrule || null;

  if (pbehavior.rrule && isObject(pbehavior.rrule)) {
    ({ rrule } = pbehavior.rrule);
  }

  return {
    rrule,

    _id: pbehavior._id || uid('pbehavior'),
    enabled: isUndefined(pbehavior.enabled) ? true : pbehavior.enabled,
    author: pbehavior.author || '',
    name: pbehavior.name || '',
    tstart: pbehavior.tstart ? convertTimestampToMoment(pbehavior.tstart).toDate() : new Date(),
    tstop: pbehavior.tstop ? convertTimestampToMoment(pbehavior.tstop).toDate() : new Date(),
    type: pbehavior.type || '',
    reason: pbehavior.reason || '',
    filter: isString(pbehavior.filter) ? JSON.parse(pbehavior.filter) : cloneDeep(pbehavior.filter || {}),
    comments: cloneDeep(pbehavior.comments || []),
  };
}

export function formToPbehavior(form) {
  return {
    ...form,

    comments: [],
    tstart: moment(form.tstart).unix(),
    tstop: moment(form.tstop).unix(),
  };
}

export function calendarEventToPbehaviorForm(calendarEvent) {
  const { pbehavior } = calendarEvent.data || {};

  const form = pbehaviorToForm(pbehavior);

  form.tstart = calendarEvent.start.date.toDate();
  form.tstop = calendarEvent.end.date.toDate();

  return form;
}

export function formToCalendarEvent(form, calendarEvent) {
  const span = new DaySpan(
    Day.fromMoment(moment(form.tstart)),
    Day.fromMoment(moment(form.tstop)),
  );

  const schedule = calendarEvent.fullDay
    ? Schedule.forDay(span.start, span.days(Op.UP))
    : Schedule.forSpan(span);

  const details = { ...calendarEvent.data, pbehavior: formToPbehavior(form) };
  const event = Vue.$dayspan.createEvent(details, schedule, true);

  event.id = calendarEvent.event.id;

  return new CalendarEvent(calendarEvent.id, event, span, calendarEvent.day);
}
