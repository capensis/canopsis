import { isString } from 'lodash';
import { DEFAULT_ALARMS_WIDGET_GROUP_COLUMNS, DEFAULT_PERIODIC_REFRESH } from '@/constants';

const WIDGET_PARAMETERS_FIELDS = {
  widgetColumns: 'widgetColumns',
  widgetGroupColumns: 'widgetGroupColumns',
  periodicRefresh: 'periodicRefresh',
  infoPopups: 'infoPopups',
  sort: 'sort',
};

const defaultFieldsForAlarmsPreparation = [
  WIDGET_PARAMETERS_FIELDS.widgetColumns,
  WIDGET_PARAMETERS_FIELDS.widgetGroupColumns,
  WIDGET_PARAMETERS_FIELDS.infoPopups,
  WIDGET_PARAMETERS_FIELDS.sort,
];

export default {
  methods: {
    getAlarmWidgetPreparationMap() {
      return {
        [WIDGET_PARAMETERS_FIELDS.widgetColumns]: this.widgetColumnsPreparation,
        [WIDGET_PARAMETERS_FIELDS.widgetGroupColumns]: this.widgetGroupColumnsPreparation,
        [WIDGET_PARAMETERS_FIELDS.periodicRefresh]: this.periodicRefreshPreparation,
        [WIDGET_PARAMETERS_FIELDS.infoPopups]: this.infoPopupsPreparation,
        [WIDGET_PARAMETERS_FIELDS.sort]: this.sortPreparation,
      };
    },

    /**
     * widgetColumns preparation
     */
    widgetColumnsPreparation(widgetColumns, isInitialization) {
      return widgetColumns.map(v => ({
        ...v,
        value: this.prefixFormatter(v.value, isInitialization),
      }));
    },

    /**
     * widgetGroupColumns preparation
     */
    widgetGroupColumnsPreparation(widgetGroupColumns, isInitialization) {
      return widgetGroupColumns
        ? widgetGroupColumns.map(v => ({
          ...v,
          value: this.prefixFormatter(v.value, isInitialization),
        }))
        : DEFAULT_ALARMS_WIDGET_GROUP_COLUMNS.map(({ labelKey, value }) => ({
          label: this.$t(labelKey),
          value: this.prefixFormatter(value, isInitialization),
        }));
    },

    /**
     * sort preparation
     */
    sortPreparation(sort, isInitialization) {
      return {
        order: sort.order,
        column: this.prefixFormatter(sort.column, isInitialization),
      };
    },

    /**
     * infoPopups preparation
     */
    infoPopupsPreparation(infoPopups, isInitialization) {
      return infoPopups.map(v => ({
        ...v,
        column: this.prefixFormatter(v.column, isInitialization),
      }));
    },

    /**
     * If there isn't periodic refresh then we are adding it
     */
    periodicRefreshPreparation(periodicRefresh) {
      return { unit: DEFAULT_PERIODIC_REFRESH.unit, ...periodicRefresh };
    },

    /**
     * Prefix formatter for value
     *
     * @param {string|undefined} value
     * @param {boolean} [isInitialization=false]
     * @returns {*}
     */
    prefixFormatter(value, isInitialization = false) {
      if (isString(value) && value !== '') {
        if (isInitialization) {
          return value.replace(/^v\./, 'alarm.v.');
        }

        return value.replace(/^alarm\./, '');
      }

      return value;
    },

    /**
     * Preparation for only alarms list parameters
     *
     * @param {Object} parameters
     * @param {Array} [keysForPreparation=['widgetColumns', 'widgetGroupColumns', 'infoPopups', 'sort']]
     * @param {boolean} [isInitialization=false]
     * @returns {Object}
     */
    prepareAlarmWidgetParametersSettings(
      parameters,
      keysForPreparation = defaultFieldsForAlarmsPreparation,
      isInitialization = false,
    ) {
      const alarmWidgetPreparationMap = this.getAlarmWidgetPreparationMap();

      return keysForPreparation.reduce((acc, field) => {
        const preparationFunc = alarmWidgetPreparationMap[field];

        acc[field] = preparationFunc
          ? preparationFunc(parameters[field], isInitialization)
          : parameters[field];

        return acc;
      }, { ...parameters });
    },

    /**
     * Full preparation for alarms list widget
     *
     * @param {Object} widget
     * @param {boolean} [isInitialization=false]
     * @returns {Object}
     */
    prepareAlarmWidgetSettings(widget, isInitialization = false) {
      return {
        ...widget,

        parameters: this.prepareAlarmWidgetParametersSettings(
          widget.parameters,
          [
            WIDGET_PARAMETERS_FIELDS.widgetColumns,
            WIDGET_PARAMETERS_FIELDS.widgetGroupColumns,
            WIDGET_PARAMETERS_FIELDS.periodicRefresh,
            WIDGET_PARAMETERS_FIELDS.infoPopups,
            WIDGET_PARAMETERS_FIELDS.sort,
          ],
          isInitialization,
        ),
      };
    },

    /**
     * Full preparation for widgets which has alarmsList property in parameters
     *
     * @param {Object} widget
     * @param {boolean} [isInitialization=false]
     * @returns {Object}
     */
    prepareWidgetWithAlarmParametersSettings(widget, isInitialization = false) {
      return {
        ...widget,

        parameters: {
          ...widget.parameters,

          alarmsList: this.prepareAlarmWidgetParametersSettings(
            widget.parameters.alarmsList,
            [WIDGET_PARAMETERS_FIELDS.widgetColumns, WIDGET_PARAMETERS_FIELDS.infoPopups],
            isInitialization,
          ),
        },
      };
    },
  },
};
