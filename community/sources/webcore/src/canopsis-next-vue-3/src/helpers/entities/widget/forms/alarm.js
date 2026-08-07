import { PAGINATION_LIMIT } from '@/config';
import {
  COLOR_INDICATOR_TYPES,
  DEFAULT_ALARM_MORE_INFO_TEMPLATE,
  DEFAULT_ALARMS_WIDGET_COLUMNS,
  RESIZING_CELLS_CONTENTS_BEHAVIORS,
} from '@/constants';

import ALARM_EXPORT_PDF_TEMPLATE from '@/assets/templates/alarm-export-pdf.html?raw';

import { formToWidgetColumns, widgetColumnsToForm } from '../column/form';
import { formToWidgetSortColumns, widgetSortColumnsToForm } from '../sort-column/form';
import { formToWidgetTemplateValue, widgetTemplateValueToForm } from '../template/form';

/**
 * @typedef { 'wrap' | 'truncate' } AlarmsResizingBehaviors
 */

/**
 * @typedef {Object} WidgetColumnsParameters
 * @property {boolean} draggable
 * @property {boolean} resizable
 * @property {AlarmsResizingBehaviors} cells_content_behavior
 */

/**
 * @typedef {Object} WidgetInfoPopup
 * @property {string} column
 * @property {string} template
 */

/**
 * @typedef {Object} AlarmListBaseParameters
 * @property {number} itemsPerPage
 * @property {WidgetSortColumn} sort
 * @property {string} sortTemplate
 * @property {string} moreInfoTemplate
 * @property {string} moreInfoTemplateTemplate
 * @property {WidgetInfoPopup[]} infoPopups
 * @property {string} widgetColumnsTemplate
 * @property {WidgetColumn[]} widgetColumns
 * @property {string} exportPdfTemplate
 * @property {string} exportPdfTemplateTemplate
 * @property {boolean} showRootCauseByStateClick
 * @property {ColorIndicator} rootCauseColorIndicator
 * @property {WidgetColumnsParameters} columns
 */

/**
 * @typedef {AlarmListBaseParameters} AlarmListBaseParametersForm
 * @property {string | Symbol} widgetColumnsTemplate
 * @property {WidgetColumnForm[]} widgetColumns
 */

/**
 * Convert opened field widget.
 *
 * @param {boolean | null} [opened]
 * @returns {boolean | null}
 */
export const openedToForm = opened =>
  typeof opened === 'boolean' || opened === null ? opened : true;

/**
 * Convert columns parameters field widget.
 *
 * @param {WidgetColumnsParameters} [columns={}]
 * @returns {WidgetColumnsParameters}
 */
export const columnsParametersToForm = (columns = {}) => ({
  draggable: columns.draggable ?? false,
  resizable: columns.resizable ?? false,
  cells_content_behavior: columns.cells_content_behavior ?? RESIZING_CELLS_CONTENTS_BEHAVIORS.wrap,
});

/**
 * Convert alarm list infoPopups parameters to form.
 *
 * @param {WidgetInfoPopup[]} [infoPopups=[]]
 * @return {WidgetInfoPopup[]}
 */
export const infoPopupsToForm = (infoPopups = []) =>
  infoPopups.map(infoPopup => ({ ...infoPopup }));

/**
 * Convert alarm list base parameters to form.
 *
 * @param {AlarmListBaseParameters} [alarmListParameters={}]
 * @return {AlarmListBaseParametersForm}
 */
export const alarmListBaseParametersToForm = (alarmListParameters = {}) => ({
  sort: widgetSortColumnsToForm(alarmListParameters.sort),
  sortTemplate: widgetTemplateValueToForm(alarmListParameters.sortTemplate),
  itemsPerPage: alarmListParameters.itemsPerPage ?? PAGINATION_LIMIT,
  moreInfoTemplate: alarmListParameters.moreInfoTemplate ?? DEFAULT_ALARM_MORE_INFO_TEMPLATE,
  moreInfoTemplateTemplate: widgetTemplateValueToForm(alarmListParameters.moreInfoTemplateTemplate),
  infoPopups: infoPopupsToForm(alarmListParameters.infoPopups),
  widgetColumnsTemplate: widgetTemplateValueToForm(alarmListParameters.widgetColumnsTemplate),
  widgetColumns: widgetColumnsToForm(
    alarmListParameters.widgetColumns ?? DEFAULT_ALARMS_WIDGET_COLUMNS,
  ),
  columns: columnsParametersToForm(alarmListParameters.columns),
  exportPdfTemplate: alarmListParameters.exportPdfTemplate ?? ALARM_EXPORT_PDF_TEMPLATE,
  exportPdfTemplateTemplate: widgetTemplateValueToForm(
    alarmListParameters.exportPdfTemplateTemplate,
  ),
  showRootCauseByStateClick: alarmListParameters.showRootCauseByStateClick ?? true,
  rootCauseColorIndicator:
    alarmListParameters.rootCauseColorIndicator ?? COLOR_INDICATOR_TYPES.state,
});

/**
 * Convert form to alarm list base parameters.
 *
 * @param {AlarmListBaseParametersForm} [form={}]
 * @return {AlarmListBaseParameters}
 */
export const formToAlarmListBaseParameters = (form = {}) => ({
  ...form,
  moreInfoTemplateTemplate: formToWidgetTemplateValue(form.moreInfoTemplateTemplate),
  exportPdfTemplateTemplate: formToWidgetTemplateValue(form.exportPdfTemplateTemplate),
  widgetColumnsTemplate: formToWidgetTemplateValue(form.widgetColumnsTemplate),
  widgetColumns: formToWidgetColumns(form.widgetColumns),
  sort: formToWidgetSortColumns(form.sort),
  sortTemplate: formToWidgetTemplateValue(form.sortTemplate),
});
