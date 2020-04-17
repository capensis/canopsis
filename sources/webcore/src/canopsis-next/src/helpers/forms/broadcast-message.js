import moment from 'moment';

import { DEFAULT_BROADCAST_MESSAGE_COLOR } from '@/constants';

import convertTimestampToMoment from '@/helpers/date';

export function messageToForm({
  message = '',
  color = DEFAULT_BROADCAST_MESSAGE_COLOR,
  start,
  end,
} = {}) {
  return {
    message,
    color,
    start: start ? convertTimestampToMoment(start).toDate() : new Date(),
    end: end ? convertTimestampToMoment(end).toDate() : new Date(),
  };
}

export function formToMessage({
  message = '',
  color = '',
  start,
  end,
} = {}) {
  return {
    message,
    color,
    start: moment(start).unix(),
    end: moment(end).unix(),
  };
}
