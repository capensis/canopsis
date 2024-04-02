import { ALARM_LINK_ICON_CHIP_COLUMN_GAP, ALARM_LINK_ICON_CHIP_WIDTH, ALARM_LINK_TD_PADDINGS } from '@/constants';

/**
 * @typedef {0 | 1 | 2} AlarmListWidgetDense
 */

/**
 * Generate alarm details id by widgetId
 *
 * @param {string} alarmId
 * @param {string} widgetId
 * @returns {string}
 */
export const generateAlarmDetailsId = (alarmId, widgetId) => `${alarmId}_${widgetId}`;

/**
 * Get dataPreparer for alarmDetails entity
 *
 * @param {string} widgetId
 * @returns {Function}
 */
export const getAlarmDetailsDataPreparer = widgetId => data => (
  data.map(item => ({
    ...item,

    /**
     * We are generating new id based on alarmId and widgetId to avoiding collision with two widgets
     * on the same view with opened expand panel on the same alarm
     */
    _id: generateAlarmDetailsId(item._id, widgetId),
  }))
);

/**
 * Calculate alarm links column width for table
 *
 * @param {AlarmListWidgetDense} dense
 * @param {number} linksInRowCount
 * @returns {`${number}px`}
 */
export const calculateAlarmLinksColumnWidth = (dense, linksInRowCount) => (
  `${(ALARM_LINK_ICON_CHIP_WIDTH * linksInRowCount)
  + (ALARM_LINK_ICON_CHIP_COLUMN_GAP[dense] * (linksInRowCount - 1))
  + (ALARM_LINK_TD_PADDINGS[dense] * 2)}px`
);
