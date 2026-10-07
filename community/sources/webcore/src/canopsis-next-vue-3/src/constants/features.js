export const FEATURES_SERVICE_KEYS = {
  i18n: 'i18n',
  constants: {
    widgetTypes: 'constants.WIDGET_TYPES',
    componentsByWidgetTypes: 'constants.COMPONENTS_BY_WIDGET_TYPES',
    widgetIcons: 'constants.WIDGET_ICONS',
    sideBars: 'constants.SIDE_BARS',
    sideBarsByWidgetTypes: 'constants.SIDE_BARS_BY_WIDGET_TYPES',
    widgetTypesRules: 'constants.WIDGET_TYPES_RULES',
    userPermissions: {
      business: {
        alarmsList: {
          actions: 'constants.USER_PERMISSIONS.business.alarmsList.actions',
        },
      },
      api: {
        general: 'constants.USER_PERMISSIONS.api.general',
      },
    },
  },
  components: {
    modals: {
      components: 'components.modals.components',
      dialogPropsMap: 'components.modals.dialogPropsMap',
    },
    sidebars: {
      components: 'components.sidebars.components',
    },
    widgetWrapper: {
      components: 'components.widgetWrapper.components',
      computed: {
        preparedWidget: 'components.widgetWrapper.computed.preparedWidget',
      },
    },
    alarmListActionPanel: {
      computed: {
        actions: 'components.alarmListActionPanel.computed.actions',
      },
    },
    alarmListMassActionsPanel: {
      computed: {
        actions: 'components.alarmListMassActionsPanel.computed.actions',
      },
      methods: 'components.alarmListMassActionsPanel.methods',
    },
  },
  helpers: {
    forms: {
      widgets: {
        widgetParametersToForm: {
          widgetsMap: 'helpers.forms.widgets.widgetParametersToForm.widgetsMap',
        },
        formToWidgetParameters: {
          widgetsMap: 'helpers.forms.widgets.formToWidgetParameters.widgetsMap',
        },
      },
    },
    query: {
      convertUserPreferenceToQuery: {
        convertersMap: 'helpers.query.convertUserPreferenceToQuery.convertersMap',
      },
      convertWidgetToQuery: {
        convertersMap: 'helpers.query.convertWidgetToQuery.convertersMap',
      },
    },
  },
};
