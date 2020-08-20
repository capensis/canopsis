import Vue from 'vue';
import moment from 'moment-timezone';
import { omit, isObject, isString, cloneDeep, isUndefined } from 'lodash';
import { CalendarEvent, DaySpan, Op, Schedule } from 'dayspan';

import uid from '@/helpers/uid';
import { convertDateToTimestampByTimezone, convertTimestampToMoment } from '@/helpers/date';
import { addKeyInEntity, removeKeyFromEntity } from '@/helpers/entities';

const preparePbehaviorType = type => (isObject(type) ? type._id : type);

export const exdatesToRequest = exdates => removeKeyFromEntity(exdates).map(({ type, begin, end }) => ({
  type: preparePbehaviorType(type),
  begin: moment(begin).unix(),
  end: moment(end).unix(),
}));

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
    type: pbehavior.type,
    reason: pbehavior.reason,
    tstart: pbehavior.tstart ? convertTimestampToMoment(pbehavior.tstart).toDate() : new Date(),
    tstop: pbehavior.tstop ? convertTimestampToMoment(pbehavior.tstop).toDate() : new Date(),
    filter: isString(resultFilter) ? JSON.parse(resultFilter) : cloneDeep(resultFilter),
    exceptions: pbehavior.exceptions ? addKeyInEntity(cloneDeep(pbehavior.exceptions)) : [],
    comments: pbehavior.comments ? addKeyInEntity(cloneDeep(pbehavior.comments)) : [],
    exdates: pbehavior.exdates ? addKeyInEntity(cloneDeep(pbehavior.exdates)) : [], // TODO: convert timestamp to Date
  };
};

export const formToPbehavior = (form, timezone) => ({
  ...form,

  reason: form.reason._id,
  type: form.type._id,
  comments: removeKeyFromEntity(form.comments),
  exdates: exdatesToRequest(form.exdates),
  exceptions: removeKeyFromEntity(form.exceptions).map(({ _id }) => _id),
  tstart: convertDateToTimestampByTimezone(form.tstart, timezone),
  tstop: convertDateToTimestampByTimezone(form.tstop, timezone),
});

export const calendarEventToPbehaviorForm = (calendarEvent, filter) => {
  const { pbehavior, cachedForm = {} } = calendarEvent.data || {};

  return {
    ...pbehaviorToForm(pbehavior, filter),
    ...cachedForm,
    tstart: calendarEvent.start.date.toDate(),
    tstop: calendarEvent.schedule.durationUnit === 'days'
      ? moment(calendarEvent.end.date).subtract(1, 'second').toDate()
      : calendarEvent.end.date.toDate(),
  };
};

export const formToCalendarEvent = (form, calendarEvent, timezone) => {
  const span = new DaySpan(calendarEvent.start, calendarEvent.end);

  const schedule = calendarEvent.fullDay
    ? Schedule.forDay(span.start, span.days(Op.UP))
    : Schedule.forSpan(span);

  const details = { ...calendarEvent.data, pbehavior: formToPbehavior(form, timezone) };
  const event = Vue.$dayspan.createEvent(details, schedule);

  event.id = calendarEvent.event.id;

  return new CalendarEvent(calendarEvent.id, event, span, calendarEvent.day);
};

export const pbehaviorToRequest = (pbehavior) => {
  const result = omit(pbehavior, ['type', 'reason', 'exdates']);

  result.type = preparePbehaviorType(pbehavior.type);
  result.reason = isObject(pbehavior.reason) ? pbehavior.reason._id : pbehavior.reason;

  if (pbehavior.exdates) {
    result.exdates = exdatesToRequest(pbehavior.exdates);
  }

  return result;
};
