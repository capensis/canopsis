export const FEATURES_SERVICE_KEYS = {
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
