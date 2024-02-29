import JsPDF from 'jspdf';
import { Visitor, parseWithoutProcessing } from 'handlebars';

import {
  ALARM_EXPORT_PDF_FIELDS,
  ALARM_EXPORT_PDF_FIELDS_TO_ORIGINAL_FIELDS,
  DATETIME_FORMATS,
  ALARM_STATES,
  ALARM_STATUSES,
} from '@/constants';

import ALARM_EXPORT_PDF_TEMPLATE from '@/assets/templates/alarm-export-pdf.html';

import { compile } from '../handlebars';
import { createInstanceWithHelpers } from '../handlebars/alarm-export-pdf-helpers';
import { convertDateToStringWithNewTimezone } from '../date/date';
import { convertDurationToString } from '../date/duration';

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
 * @property {string} resolved
 * @property {string} activation_date
 * @property {Object} infos
 * @property {Object} pbehavior_info
 * @property {Object} ticket_info
 * @property {AlarmStep[]} comments
 * @property {string[]} tags
 * @property {AlarmLinks} links
 */

/**
 * @typedef {Alarm} AlarmWithComments
 * @property {AlarmStep[]} comments
 */

class AlarmExportToPdfVisitor extends Visitor {
  constructor() {
    super();

    this.mutating = true;
  }

  /* eslint-disable-next-line class-methods-use-this */
  PathExpression(path) {
    const { parts = [] } = path;
    const lastItem = parts.at(-1);

    return {
      ...path,
      parts: ALARM_EXPORT_PDF_FIELDS_TO_ORIGINAL_FIELDS[lastItem] ? [lastItem] : parts,
    };
  }
}

/**
 * Prepare alarm state fpr exporting
 *
 * @param {AlarmStep} state
 * @returns {string}
 */
export const prepareAlarmStateForExport = state => ({
  [ALARM_STATES.ok]: '0 - OK',
  [ALARM_STATES.minor]: '1 - Minor',
  [ALARM_STATES.major]: '2 - Major',
  [ALARM_STATES.critical]: '3 - Critical',
}[state?.val] ?? `Invalid value (${state?.val})`);

/**
 * Prepare alarm status fpr exporting
 *
 * @param {AlarmStep} status
 * @returns {string}
 */
const prepareAlarmStatusForExport = status => ({
  [ALARM_STATUSES.closed]: 'Closed',
  [ALARM_STATUSES.ongoing]: 'Ongoing',
  [ALARM_STATUSES.flapping]: 'Flapping',
  [ALARM_STATUSES.stealthy]: 'Stealth',
  [ALARM_STATUSES.cancelled]: 'Canceled',
  [ALARM_STATUSES.noEvents]: 'No events',
}[status?.val] ?? `Invalid value (${status?.val})`);

/**
 * Prepare alarm for exporting
 *
 * @param {AlarmWithComments} [alarm = {}]
 * @param {string} [timezone]
 * @returns {AlarmExport}
 */
const prepareAlarmForExport = (alarm = {}, timezone) => {
  const { v = {} } = alarm;

  return {
    alarm,

    [ALARM_EXPORT_PDF_FIELDS.displayName]: v.display_name,
    [ALARM_EXPORT_PDF_FIELDS.state]: prepareAlarmStateForExport(v.state),
    [ALARM_EXPORT_PDF_FIELDS.status]: prepareAlarmStatusForExport(v.status),
    [ALARM_EXPORT_PDF_FIELDS.connector]: v.connector,
    [ALARM_EXPORT_PDF_FIELDS.connectorName]: v.connector_name,
    [ALARM_EXPORT_PDF_FIELDS.component]: v.component,
    [ALARM_EXPORT_PDF_FIELDS.resource]: v.resource,
    [ALARM_EXPORT_PDF_FIELDS.initialOutput]: v.initial_output,
    [ALARM_EXPORT_PDF_FIELDS.output]: v.output,
    [ALARM_EXPORT_PDF_FIELDS.eventsCount]: v.events_count,
    [ALARM_EXPORT_PDF_FIELDS.duration]: convertDurationToString(v.duration),
    [ALARM_EXPORT_PDF_FIELDS.currentDate]:
      convertDateToStringWithNewTimezone(new Date(), DATETIME_FORMATS.longWithTimezone, timezone),
    [ALARM_EXPORT_PDF_FIELDS.creationDate]:
      convertDateToStringWithNewTimezone(alarm.t, DATETIME_FORMATS.longWithTimezone, timezone),
    [ALARM_EXPORT_PDF_FIELDS.lastEventDate]:
      convertDateToStringWithNewTimezone(v.last_event_date, DATETIME_FORMATS.longWithTimezone, timezone),
    [ALARM_EXPORT_PDF_FIELDS.lastUpdateDate]:
      convertDateToStringWithNewTimezone(v.last_update_date, DATETIME_FORMATS.longWithTimezone, timezone),
    [ALARM_EXPORT_PDF_FIELDS.acknowledgeDate]:
      convertDateToStringWithNewTimezone(v.ack?.t, DATETIME_FORMATS.longWithTimezone, timezone),
    [ALARM_EXPORT_PDF_FIELDS.resolved]:
      convertDateToStringWithNewTimezone(v.resolved, DATETIME_FORMATS.longWithTimezone, timezone),
    [ALARM_EXPORT_PDF_FIELDS.activationDate]:
      convertDateToStringWithNewTimezone(v.activation_date, DATETIME_FORMATS.longWithTimezone, timezone),
    [ALARM_EXPORT_PDF_FIELDS.infos]: alarm.infos,
    [ALARM_EXPORT_PDF_FIELDS.ticket]: v.ticket,
    [ALARM_EXPORT_PDF_FIELDS.tags]: alarm.tags,
    [ALARM_EXPORT_PDF_FIELDS.links]: alarm.links,
    [ALARM_EXPORT_PDF_FIELDS.pbehaviorInfo]: v.pbehavior_info
      ? {
        ...v.pbehavior_info,

        timestamp:
          convertDateToStringWithNewTimezone(v.pbehavior_info.timestamp, DATETIME_FORMATS.longWithTimezone, timezone),
      }
      : null,
    [ALARM_EXPORT_PDF_FIELDS.comments]: alarm.comments?.map(comment => ({
      ...comment,

      t: convertDateToStringWithNewTimezone(comment.t, DATETIME_FORMATS.longWithTimezone, timezone),
    })),

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
export const exportAlarmToPdf = async (template = ALARM_EXPORT_PDF_TEMPLATE, alarm = {}, timezone) => {
  const doc = new JsPDF();
  const handlebars = createInstanceWithHelpers();
  const ast = parseWithoutProcessing(template ?? ALARM_EXPORT_PDF_TEMPLATE);
  const scanner = new AlarmExportToPdfVisitor();
  const html = await compile(
    scanner.accept(ast),
    prepareAlarmForExport(alarm, timezone),
    handlebars,
  );

  return doc.html(html, {
    margin: [5, 5, 5, 5],
    autoPaging: 'text',
    x: 0,
    y: 0,
    width: 200,
    windowWidth: 1000,
    html2canvas: {
      logging: false,
    },
  }).save(`alarm-${alarm?._id}.pdf`);
};
