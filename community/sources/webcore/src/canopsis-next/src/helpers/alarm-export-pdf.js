import JsPDF from 'jspdf';

import { DATETIME_FORMATS, ENTITIES_STATES, ENTITIES_STATUSES } from '@/constants';

import ALARM_EXPORT_PDF_TEMPLATE from '@/assets/templates/alarm-export-pdf.html';

import { convertDateToStringWithNewTimezone } from './date/date';
import { convertDurationToString } from './date/duration';
import { createInstanceWithHelpers } from './handlebars/alarm-export-pdf-helpers';
import { compile } from './handlebars';

/**
 * @typedef {Object} AlarmExport
 * @property {string} display_name
 * @property {string} state
 * @property {string} status
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
 * Prepare alarm state fpr exporting
 *
 * @param {AlarmEvent} state
 * @returns {string}
 */
export const prepareAlarmStateForExport = state => ({
  [ENTITIES_STATES.ok]: '0 - OK',
  [ENTITIES_STATES.minor]: '0 - Minor',
  [ENTITIES_STATES.major]: '0 - Major',
  [ENTITIES_STATES.critical]: '0 - Critical',
}[state?.val] ?? `Invalid value (${state?.val})`);

/**
 * Prepare alarm status fpr exporting
 *
 * @param {AlarmEvent} status
 * @returns {string}
 */
export const prepareAlarmStatusForExport = status => ({
  [ENTITIES_STATUSES.closed]: 'Closed',
  [ENTITIES_STATUSES.ongoing]: 'Ongoing',
  [ENTITIES_STATUSES.flapping]: 'Flapping',
  [ENTITIES_STATUSES.stealthy]: 'Stealth',
  [ENTITIES_STATUSES.cancelled]: 'Canceled',
  [ENTITIES_STATUSES.noEvents]: 'No events',
}[status?.val] ?? `Invalid value (${status?.val})`);

/**
 * Prepare alarm for exporting
 *
 * @param {Alarm} [alarm = {}]
 * @param {string} [timezone]
 * @returns {AlarmExport}
 */
export const prepareAlarmForExport = (alarm = {}, timezone) => {
  const { v = {} } = alarm;

  return {
    display_name: v.display_name,
    state: prepareAlarmStateForExport(v.state),
    status: prepareAlarmStatusForExport(v.status),
    connector: v.connector,
    connector_name: v.connector_name,
    component: v.component,
    resource: v.resource,
    initial_output: v.initial_output,
    output: v.output,
    events_count: v.events_count,
    duration: convertDurationToString(v.duration),
    current_date:
      convertDateToStringWithNewTimezone(new Date(), DATETIME_FORMATS.longWithTimezone, timezone),
    creation_date:
      convertDateToStringWithNewTimezone(alarm.t, DATETIME_FORMATS.longWithTimezone, timezone),
    last_event_date:
      convertDateToStringWithNewTimezone(v.last_event_date, DATETIME_FORMATS.longWithTimezone, timezone),
    last_update_date:
      convertDateToStringWithNewTimezone(v.last_update_date, DATETIME_FORMATS.longWithTimezone, timezone),
    acknowledge_date:
      convertDateToStringWithNewTimezone(v.ack?.t, DATETIME_FORMATS.longWithTimezone, timezone),
    activation_date:
      convertDateToStringWithNewTimezone(v.activation_date, DATETIME_FORMATS.longWithTimezone, timezone),
    infos: alarm.infos,
    pbehavior_info: v.pbehavior_info
      ? {
        ...v.pbehavior_info,

        timestamp:
          convertDateToStringWithNewTimezone(v.pbehavior_info.timestamp, DATETIME_FORMATS.longWithTimezone, timezone),
      }
      : null,
    ticket_info: {
      _t: 'declareticket',
      a: 'root',
      m: 'Ticket declaration rule: test-declareticketrule-declareticket-execution-1-name. Ticket ID: test-ticket-declareticket-execution-1. Ticket URL: https://test/test-ticket-declareticket-execution-1. Ticket name: test-ticket-declareticket-execution-1 test-resource-declareticket-execution-1-1.',
      ticket: 'test-ticket-declareticket-execution-1',
      ticket_url: 'https://test/test-ticket-declareticket-execution-1',
      ticket_data: {
        name: 'test-ticket-declareticket-execution-1 test-resource-declareticket-execution-1-1',
      },
      ticket_comment: 'test-comment-declareticket-execution-1-1',
      ticket_rule_id: '{{ .ruleId }}',
      ticket_rule_name: 'Ticket declaration rule: test-declareticketrule-declareticket-execution-1-name',
      ticket_system_name: 'test-declareticketrule-declareticket-execution-1-system-name',
    },
    last_comment: v.last_comment
      ? {
        ...v.last_comment,

        t: convertDateToStringWithNewTimezone(v.last_comment.t, DATETIME_FORMATS.longWithTimezone, timezone),
      }
      : null,
    tags: alarm.tags,
    links: alarm.links,
  };
};

/**
 * Export alarm by special template in special timezone to PDF file
 *
 * @param {string} [template = ALARM_EXPORT_PDF_TEMPLATE]
 * @param {Alarm} [alarm = {}]
 * @param {string} [timezone]
 * @returns {Promise<unknown>}
 */
export const exportAlarmToPdf = async (template = ALARM_EXPORT_PDF_TEMPLATE, alarm = {}, timezone) => (
  new Promise((resolve, reject) => {
    const doc = new JsPDF();
    const handlebars = createInstanceWithHelpers();

    compile(
      template ?? ALARM_EXPORT_PDF_TEMPLATE,
      prepareAlarmForExport(alarm, timezone),
      handlebars,
    )
      .then(html => doc.html(html, {
        callback: () => {
          doc.save(`alarm-${alarm?._id}.pdf`);

          return resolve();
        },
        margin: [5, 5, 5, 5],
        autoPaging: 'text',
        x: 0,
        y: 0,
        width: 200,
        windowWidth: 1000,
      }))
      .catch((err) => {
        console.error(err);
        reject();
      });
  })
);
