import isString from 'lodash/isString';

export default {
  methods: {
    prefixFormatter(value) {
      if (isString(value)) {
        return value.replace('alarm.', 'v.');
      }

      return value;
    },

    prepareAlarmWidgetParametersSettings(parameters) {
      return {
        ...parameters,

        widgetColumns: parameters.widgetColumns.map(v => ({
          ...v,
          value: this.prefixFormatter(v.value),
        })),

        infoPopups: parameters.infoPopups.map(v => ({
          ...v,
          column: this.prefixFormatter(v.column),
        })),

        sort: this.prefixFormatter(parameters.sort),
      };
    },
  },
};
