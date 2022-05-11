import {
  formWidgetColumnsToColumns,
  widgetColumnsToForm,
} from '@/helpers/forms/widgets/alarm';

export const sideBarSettingsWidgetAlarmMixin = {
  methods: {
    /**
     * Full preparation for widgets which has alarmsList property in parameters
     *
     * @param {Object} widget
     * @param {boolean} [isInitialization=false]
     * @returns {Object}
     */
    prepareWidgetWithAlarmParametersSettings(widget, isInitialization = false) {
      const { alarmsList } = widget.parameters;

      return {
        ...widget,

        parameters: {
          ...widget.parameters,
          alarmsList: {
            ...alarmsList,
            widgetColumns: isInitialization
              ? widgetColumnsToForm(alarmsList.widgetColumns)
              : formWidgetColumnsToColumns(alarmsList.widgetColumns),
            infoPopups: isInitialization
              ? widgetColumnsToForm(alarmsList.infoPopups)
              : formWidgetColumnsToColumns(alarmsList.infoPopups),
          },
        },
      };
    },
  },
};
