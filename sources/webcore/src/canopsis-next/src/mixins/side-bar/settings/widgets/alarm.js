import isString from 'lodash/isString';

export default {
  methods: {
    prefixFormatter(value, isInitialization = false) {
      if (isString(value)) {
        if (isInitialization) {
          return value.replace(/^v\.entity\./, 'entity.').replace(/^v\./, 'alarm.');
        }

        return value.replace(/^entity\./, 'v.entity.').replace(/^alarm\./, 'v.');
      }

      return value;
    },

    prepareAlarmWidgetParametersSettings(parameters, isInitialization) {
      return {
        ...parameters,

        widgetColumns: parameters.widgetColumns.map(v => ({
          ...v,
          value: this.prefixFormatter(v.value, isInitialization),
        })),

        infoPopups: parameters.infoPopups.map(v => ({
          ...v,
          column: this.prefixFormatter(v.column, isInitialization),
        })),

        sort: this.prefixFormatter(parameters.sort, isInitialization),
      };
    },
  },
};
