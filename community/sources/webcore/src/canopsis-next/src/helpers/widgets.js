import {
  ALARM_FIELDS,
  ALARM_FIELDS_TO_LABELS_KEYS,
  ALARM_UNSORTABLE_FIELDS, COLOR_INDICATOR_TYPES,
  DEFAULT_ALARMS_WIDGET_COLUMNS,
  DEFAULT_ALARMS_WIDGET_GROUP_COLUMNS,
  DEFAULT_CONTEXT_WIDGET_COLUMNS,
  DEFAULT_SERVICE_DEPENDENCIES_COLUMNS,
  ENTITY_FIELDS_TO_LABELS_KEYS,
  ENTITY_UNSORTABLE_FIELDS,
} from '@/constants';

import i18n from '@/i18n';

import { convertDateToStringWithFormatForToday } from './date/date';
import { convertDurationToString } from './date/duration';
import { setSeveralFields } from './immutable';
import { getInfosWidgetColumn, isLinksWidgetColumn } from './forms/shared/widget-column';

/**
 * Get translated label for widget column
 *
 * @param {WidgetColumn} [column = {}]
 * @param {Object<string, string>} [labelsMap = {}]
 * @returns {string}
 */
export const getColumnLabel = (column = {}, labelsMap = {}) => {
  if (column.label) {
    return column.label;
  }

  const infosColumn = getInfosWidgetColumn(column.value);

  if (infosColumn) {
    return i18n.tc(labelsMap[infosColumn], 2);
  }

  if (isLinksWidgetColumn(column.value)) {
    return i18n.tc(labelsMap[ALARM_FIELDS.links], 2);
  }

  return i18n.tc(labelsMap[column.value], 2);
};

/**
 * Get sortable property for widget column for table
 *
 * @param {WidgetColumn} [column = {}]
 * @param {string[]} [unsortableFields = []]
 * @returns {boolean}
 */
export const getSortable = (column = {}, unsortableFields = []) => (
  !unsortableFields.some(field => (column.value ?? '').startsWith(field))
);

/**
 * Prepared alarms list widget for displaying
 *
 * @param {Object} widget
 * @returns {Object}
 */
export const prepareAlarmListWidget = (widget = {}) => setSeveralFields(widget, {
  'parameters.widgetColumns': (columns = []) => (
    columns.map(column => ({
      ...column,

      sortable: getSortable(column, ALARM_UNSORTABLE_FIELDS),
      text: getColumnLabel(column, ALARM_FIELDS_TO_LABELS_KEYS),
    }))
  ),

  'parameters.widgetGroupColumns': (columns = DEFAULT_ALARMS_WIDGET_GROUP_COLUMNS) => (
    columns.map(column => ({
      ...column,

      sortable: getSortable(column, ALARM_UNSORTABLE_FIELDS),
      text: getColumnLabel(column, ALARM_FIELDS_TO_LABELS_KEYS),
    }))
  ),

  'parameters.serviceDependenciesColumns': (columns = DEFAULT_SERVICE_DEPENDENCIES_COLUMNS) => (
    columns.map(column => ({
      ...column,

      sortable: false,
      text: getColumnLabel(column, ENTITY_FIELDS_TO_LABELS_KEYS),
      value: column.value.startsWith('entity.') ? column.value : `entity.${column.value}`,
    }))
  ),

  'parameters.widgetExportColumns': (columns = []) => (
    columns.map(column => ({
      ...column,

      text: getColumnLabel(column, ALARM_FIELDS_TO_LABELS_KEYS),
    }))
  ),
});

/**
 * Prepared context widget for displaying
 *
 * @param {Object} widget
 * @returns {Object}
 */
export const prepareContextWidget = (widget = {}) => setSeveralFields(widget, {
  'parameters.widgetColumns': (columns = []) => (
    columns.map(column => ({
      ...column,

      sortable: getSortable(column, ENTITY_UNSORTABLE_FIELDS),
      text: getColumnLabel(column, ENTITY_FIELDS_TO_LABELS_KEYS),
    }))
  ),

  'parameters.serviceDependenciesColumns': (columns = DEFAULT_SERVICE_DEPENDENCIES_COLUMNS) => (
    columns.map(column => ({
      ...column,

      sortable: false,
      text: getColumnLabel(column, ENTITY_FIELDS_TO_LABELS_KEYS),
      value: column.value.startsWith('entity.') ? column.value : `entity.${column.value}`,
    }))
  ),

  'parameters.activeAlarmsColumns': (columns = DEFAULT_ALARMS_WIDGET_COLUMNS) => (
    columns.map(column => ({
      ...column,

      sortable: getSortable(column, ALARM_UNSORTABLE_FIELDS),
      text: getColumnLabel(column, ALARM_FIELDS_TO_LABELS_KEYS),
    }))
  ),

  'parameters.resolvedAlarmsColumns': (columns = DEFAULT_ALARMS_WIDGET_COLUMNS) => (
    columns.map(column => ({
      ...column,

      sortable: getSortable(column, ALARM_UNSORTABLE_FIELDS),
      text: getColumnLabel(column, ALARM_FIELDS_TO_LABELS_KEYS),
    }))
  ),

  'parameters.widgetExportColumns': (columns = []) => (
    columns.map(column => ({
      ...column,

      text: getColumnLabel(column, ENTITY_FIELDS_TO_LABELS_KEYS),
    }))
  ),
});

/**
 * Prepared service weather widget for displaying
 *
 * @param {Object} widget
 * @returns {Object}
 */
export const prepareServiceWeatherWidget = (widget = {}) => setSeveralFields(widget, {
  'parameters.serviceDependenciesColumns': (columns = DEFAULT_SERVICE_DEPENDENCIES_COLUMNS) => (
    columns.map(column => ({
      ...column,

      sortable: false,
      text: getColumnLabel(column, ENTITY_FIELDS_TO_LABELS_KEYS),
      value: column.value.startsWith('entity.') ? column.value : `entity.${column.value}`,
    }))
  ),

  'parameters.alarmsList.widgetColumns': (columns = DEFAULT_ALARMS_WIDGET_COLUMNS) => (
    columns.map(column => ({
      ...column,

      sortable: getSortable(column, ALARM_UNSORTABLE_FIELDS),
      text: getColumnLabel(column, ALARM_FIELDS_TO_LABELS_KEYS),
    }))
  ),
});

/**
 * Prepared stats calendar/counter widget for displaying
 *
 * @param {Object} widget
 * @returns {Object}
 */
export const prepareStatsCalendarAndCounterWidget = (widget = {}) => setSeveralFields(widget, {
  'parameters.alarmsList.widgetColumns': (columns = DEFAULT_ALARMS_WIDGET_COLUMNS) => (
    columns.map(column => ({
      ...column,

      sortable: getSortable(column, ALARM_UNSORTABLE_FIELDS),
      text: getColumnLabel(column, ALARM_FIELDS_TO_LABELS_KEYS),
    }))
  ),
});

/**
 * Prepared map widget for displaying
 *
 * @param {Object} widget
 * @returns {Object}
 */
export const prepareMapWidget = (widget = {}) => setSeveralFields(widget, {
  'parameters.alarmsColumns': (columns = DEFAULT_ALARMS_WIDGET_COLUMNS) => (
    columns.map(column => ({
      ...column,

      sortable: getSortable(column, ALARM_UNSORTABLE_FIELDS),
      text: getColumnLabel(column, ALARM_FIELDS_TO_LABELS_KEYS),
    }))
  ),

  'parameters.entitiesColumns': (columns = DEFAULT_CONTEXT_WIDGET_COLUMNS) => (
    columns.map(column => ({
      ...column,

      sortable: getSortable(column, ENTITY_UNSORTABLE_FIELDS),
      text: getColumnLabel(column, ENTITY_FIELDS_TO_LABELS_KEYS),
    }))
  ),
});

/**
 * Get filter for alarms list widget column value
 *
 * @param {string} value
 * @returns {Function | undefined}
 */
export const getAlarmsListWidgetColumnValueFilter = (value) => {
  switch (value) {
    case ALARM_FIELDS.lastUpdateDate:
    case ALARM_FIELDS.creationDate:
    case ALARM_FIELDS.lastEventDate:
    case ALARM_FIELDS.activationDate:
    case ALARM_FIELDS.ackAt:
    case ALARM_FIELDS.stateAt:
    case ALARM_FIELDS.statusAt:
    case ALARM_FIELDS.resolved:
    case ALARM_FIELDS.timestamp:
    case ALARM_FIELDS.entityLastPbehaviorDate:
      return convertDateToStringWithFormatForToday;

    case ALARM_FIELDS.duration:
    case ALARM_FIELDS.currentStateDuration:
    case ALARM_FIELDS.activeDuration:
    case ALARM_FIELDS.snoozeDuration:
    case ALARM_FIELDS.pbhInactiveDuration:
      return convertDurationToString;

    default:
      return undefined;
  }
};

/**
 * Get component getter for alarms list widget column
 *
 * @param {string} value
 * @param {boolean} [onlyIcon]
 * @param {Widget | {}} [widget = {}]
 * @returns {Function}
 */
export const getAlarmsListWidgetColumnComponentGetter = ({ value, onlyIcon }, widget = {}) => {
  switch (value) {
    case ALARM_FIELDS.state:
      return context => ({
        bind: {
          is: 'alarm-column-value-state',
          alarm: context.alarm,
          small: context.small,
        },
      });

    case ALARM_FIELDS.status:
      return context => ({
        bind: {
          is: 'alarm-column-value-status',
          alarm: context.alarm,
          small: context.small,
        },
      });

    case ALARM_FIELDS.impactState:
      return context => ({
        bind: {
          is: 'color-indicator-wrapper',
          type: COLOR_INDICATOR_TYPES.impactState,
          entity: context.alarm.entity,
          alarm: context.alarm,
        },
      });

    case ALARM_FIELDS.links:
      return context => ({
        bind: {
          onlyIcon,

          is: 'c-alarm-links-chips',
          alarm: context.alarm,
          small: context.small,
          inlineCount: widget.parameters?.inlineLinksCount,
        },
        on: {
          activate: context.$listeners.activate,
        },
      });

    case ALARM_FIELDS.extraDetails:
      return context => ({
        bind: {
          is: 'alarm-column-value-extra-details',
          alarm: context.alarm,
          small: context.small,
        },
      });

    case ALARM_FIELDS.tags:
      return context => ({
        bind: {
          is: 'c-alarm-tags-chips',
          alarm: context.alarm,
          selectedTag: context.selectedTag,
          small: context.small,
        },
        on: {
          select: context.$listeners['select:tag'],
        },
      });
  }

  if (value.startsWith('links.')) {
    const category = value.replace('links.', '');

    return context => ({
      bind: {
        category,
        onlyIcon,

        is: 'c-alarm-links-chips',
        alarm: context.alarm,
        small: context.small,
        inlineCount: widget.parameters?.inlineLinksCount,
      },
    });
  }

  return context => ({
    bind: {
      is: 'c-ellipsis',
      text: context.value,
    },
  });
};
