import { isObject, cloneDeep } from 'lodash';

import convertTimestampToMoment from '@/helpers/date';

export default function (pbehavior = {}) {
  let rrule = pbehavior.rrule || null;

  if (pbehavior.rrule && isObject(pbehavior.rrule)) {
    ({ rrule } = pbehavior.rrule);
  }

  return {
    author: pbehavior.author || '',
    name: pbehavior.name || '',
    tstart: pbehavior.tstart ? convertTimestampToMoment(pbehavior.tstart).toDate() : new Date(),
    tstop: pbehavior.tstop ? convertTimestampToMoment(pbehavior.tstop).toDate() : new Date(),
    filter: cloneDeep(pbehavior.filter || {}),
    type_: pbehavior.type_ || '',
    reason: pbehavior.reason || '',
    rrule,
  };
}
