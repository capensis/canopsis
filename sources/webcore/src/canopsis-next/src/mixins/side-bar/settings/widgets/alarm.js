import { isString } from 'lodash';
import { DEFAULT_PERIODIC_REFRESH } from '@/constants';

export default {
  methods: {
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
     * @param {Array} [keysForPreparation=['widgetColumns', 'infoPopups', 'sort']]
     * @param {boolean} [isInitialization=false]
     * @returns {Object}
     */
    prepareAlarmWidgetParametersSettings(
      parameters,
      keysForPreparation = ['widgetColumns', 'widgetGroupColumns', 'infoPopups', 'sort'],
      isInitialization = false,
    ) {
      return {
        ...parameters,

        /**
         * widgetColumns preparation
         */
        widgetColumns: keysForPreparation.includes('widgetColumns') ? parameters.widgetColumns.map(v => ({
          ...v,
          value: this.prefixFormatter(v.value, isInitialization),
        })) : parameters.widgetColumns,

        /**
         * widgetGroupColumns preparation
         */
        widgetGroupColumns: keysForPreparation.includes('widgetGroupColumns') ? parameters.widgetGroupColumns.map(v => ({
          ...v,
          value: this.prefixFormatter(v.value, isInitialization),
        })) : parameters.widgetGroupColumns,

        /**
         * infoPopups preparation
         */
        infoPopups: keysForPreparation.includes('infoPopups') ? parameters.infoPopups.map(v => ({
          ...v,
          column: this.prefixFormatter(v.column, isInitialization),
        })) : parameters.infoPopups,

        /**
         * sort preparation
         */
        sort: keysForPreparation.includes('sort') ? {
          order: parameters.sort.order,
          column: this.prefixFormatter(parameters.sort.column, isInitialization),
        } : parameters.sort,

        /**
         * If there isn't periodic refresh then we are adding it
         */
        periodicRefresh: {
          unit: DEFAULT_PERIODIC_REFRESH.unit,
          ...parameters.periodicRefresh,
        },
      };
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
          ['widgetColumns', 'widgetGroupColumns', 'infoPopups', 'sort'],
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
            ['widgetColumns', 'widgetGroupColumns', 'infoPopups'],
            isInitialization,
          ),
        },
      };
    },
  },
};
