import Vue from 'vue';
import moment from 'moment-timezone';
import { omit, isObject, isString, cloneDeep, isUndefined } from 'lodash';
import { CalendarEvent, DaySpan, Op, Schedule } from 'dayspan';

import uid from '@/helpers/uid';
import { convertDateToTimestampByTimezone, convertTimestampToMoment } from '@/helpers/date';
import { addKeyInEntity, removeKeyFromEntity } from '@/helpers/entities';

export function pbehaviorToForm(pbehavior = {}, filter = null) {
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
    type: pbehavior.type,
    reason: pbehavior.reason,
    tstart: pbehavior.tstart ? convertTimestampToMoment(pbehavior.tstart).toDate() : new Date(),
    tstop: pbehavior.tstop ? convertTimestampToMoment(pbehavior.tstop).toDate() : new Date(),
    filter: isString(resultFilter) ? JSON.parse(resultFilter) : cloneDeep(resultFilter),
    comments: addKeyInEntity(cloneDeep(pbehavior.comments || [])),
    exdates: addKeyInEntity(cloneDeep(pbehavior.exdates || [])), // TODO: convert timestamp to Date
  };
}

export function formToPbehavior(form, timezone) {
  return {
    ...form,

    reason: '8a48507a-7eba-463f-953f-41b93fce9745', // TODO should be replaced in version 6
    type: form.type._id,
    comments: removeKeyFromEntity(form.comments),
    exdates: removeKeyFromEntity(form.exdates),
    tstart: convertDateToTimestampByTimezone(form.tstart, timezone),
    tstop: convertDateToTimestampByTimezone(form.tstop, timezone),
  };
}

export function calendarEventToPbehaviorForm(calendarEvent, filter) {
  const { pbehavior, cachedForm = {} } = calendarEvent.data || {};

  const form = {
    ...pbehaviorToForm(pbehavior, filter),
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

export function pbehaviorToRequest(pbehavior) {
  const result = omit(pbehavior, ['_id', 'type', 'reason', 'exdates']);

  result.type = isObject(pbehavior.type) ? pbehavior.type._id : pbehavior.type;
  result.reason = isObject(pbehavior.reason) ? pbehavior.reason._id : pbehavior.reason;

  if (!pbehavior._id.includes('pbehavior')) { // TODO: fix that
    result._id = pbehavior._id;
  }

  if (pbehavior.exdates) {
    result.exdates = pbehavior.exdates
      .map(exdate => ({ ...exdate, type: isObject(exdate.type) ? exdate.type._id : exdate.type }));
  }

  return result;
}
