import { DATETIME_FORMATS } from '@/constants';

import { convertDateToStringWithNewTimezone } from './date/date';
import { convertDurationToString } from './date/duration';

/**
 * @typedef {Object} AlarmExport
 * @property {string} display_name
 * @property {string} connector
 * @property {string} connector_name
 * @property {string} component
 * @property {string} resource
 * @property {string} initial_output
 * @property {string} output
 * @property {string} output
 * @property {number} events_count
 * @property {number} events_count
 * @property {string} duration
 * @property {string} current_date
 * @property {string} creation_date
 * @property {string} last_event_date
 * @property {string} last_update_date
 * @property {string} last_update_date
 * @property {string} acknowledge_date
 * @property {string} activation_date
 * @property {Object} infos
 * @property {Object} pbehavior_info
 * @property {Object} ticket_info
 * @property {AlarmEvent} last_comment
 * @property {string[]} tags
 * @property {AlarmLinks} links
 */

/**
 *
 * @param {Alarm} alarm
 * @param {string} timezone
 * @returns {AlarmExport}
 */
export const prepareAlarmForExport = (alarm, timezone) => ({
  display_name: alarm.v.display_name,
  state: alarm.v.state,
  status: alarm.v.status,
  connector: alarm.v.connector,
  connector_name: alarm.v.connector_name,
  component: alarm.v.component,
  resource: alarm.v.resource,
  initial_output: alarm.v.initial_output,
  output: alarm.v.output,
  events_count: alarm.v.events_count,
  duration: convertDurationToString(alarm.v.duration),
  current_date:
    convertDateToStringWithNewTimezone(new Date(), DATETIME_FORMATS.longWithTimezone, timezone),
  creation_date:
    convertDateToStringWithNewTimezone(alarm.t, DATETIME_FORMATS.longWithTimezone, timezone),
  last_event_date:
    convertDateToStringWithNewTimezone(alarm.v.last_event_date, DATETIME_FORMATS.longWithTimezone, timezone),
  last_update_date:
    convertDateToStringWithNewTimezone(alarm.v.last_update_date, DATETIME_FORMATS.longWithTimezone, timezone),
  acknowledge_date:
    convertDateToStringWithNewTimezone(alarm.v.ack?.t, DATETIME_FORMATS.longWithTimezone, timezone),
  activation_date:
    convertDateToStringWithNewTimezone(alarm.v.activation_date, DATETIME_FORMATS.longWithTimezone, timezone),
  infos: alarm.infos,
  pbehavior_info: alarm.v.pbehavior_info, // ????
  ticket_info: alarm.v.ticket_info, // ???
  last_comment: alarm.v.last_comment, // TODO
  tags: alarm.tags, // TODO
  links: alarm.links, // TODO
});
