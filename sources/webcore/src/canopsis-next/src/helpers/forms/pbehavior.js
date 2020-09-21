import moment from 'moment';
import { omit, isObject, isString, cloneDeep } from 'lodash';

import uid from '@/helpers/uid';
import convertTimestampToMoment from '@/helpers/date';

export function pbehaviorToForm(pbehavior = {}) {
  let rrule = pbehavior.rrule || null;

  if (pbehavior.rrule && isObject(pbehavior.rrule)) {
    ({ rrule } = pbehavior.rrule);
  }

  return {
    rrule,

    enabled: typeof pbehavior.enabled === 'undefined' ? true : pbehavior.enabled,
    author: pbehavior.author || '',
    name: pbehavior.name || '',
    tstart: pbehavior.tstart ? convertTimestampToMoment(pbehavior.tstart).toDate() : new Date(),
    tstop: pbehavior.tstop ? convertTimestampToMoment(pbehavior.tstop).toDate() : new Date(),
    type_: pbehavior.type_ || '',
    reason: pbehavior.reason || '',
    timezone: pbehavior.timezone || 'Europe/Paris',
    filter: isString(pbehavior.filter) ? JSON.parse(pbehavior.filter) : cloneDeep(pbehavior.filter || {}),
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

export function pbehaviorToComments(pbehavior = {}) {
  const comments = pbehavior.comments || [];

  return comments.map(comment => ({
    ...comment,

    key: uid(),
  }));
}

export function commentsToPbehaviorComments(comments) {
  return comments.map(comment => omit(comment, ['key', 'ts']));
}

export function pbehaviorToExdates(pbehavior = {}) {
  const exdate = pbehavior.exdate || [];

  return exdate.map(unix => ({
    value: new Date(unix * 1000),
    key: uid(),
  }));
}

export function exdatesToPbehaviorExdates(exdate) {
  return exdate.filter(({ value }) => value).map(({ value }) => moment(value).unix());
}
