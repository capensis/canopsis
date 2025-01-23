export default {
  title: 'User interface',
  appTitle: 'App title',
  language: 'Default user interface language',
  footer: 'Login footer',
  description: 'Login page description',
  logo: 'Logo',
  infoPopupTimeout: 'Info popup timeout',
  errorPopupTimeout: 'Error popup timeout',
  allowChangeSeverityToInfo: 'Allow change severity to info',
  showHeaderOnKioskMode: 'Show header on kiosk mode',
  maxMatchedItems: 'Max matched items',
  checkCountRequestTimeout: 'Check max matched items request timeout (seconds)',
  requiredInstructionApprove: 'Required instruction approvement',
  disabledTransitions: 'Disable some animations',
  disabledTransitionsTooltip: 'It will help to improve application performance',
  autoSuggestPbehaviorName: 'Automatically suggest a pbehavior name',
  defaultTheme: 'Default theme',
  versionDescriptionTooltip: 'Tooltip popup on version banner',
  defaultVersionDescription: `<div>
  <div>@:common.edition:<strong> {{ edition }}</strong></div>
  <div>@:common.serialName:<strong> {{ serialName }}</strong></div>
  <div>@:common.versionUpdated:<strong> {{ timestamp versionUpdated }}</strong></div>
</div>`,
  tooltips: {
    maxMatchedItems: 'it need to warn user when number of items that match patterns is above this value',
    checkCountRequestTimeout: 'it need to define request timeout value for max matched items checking',
  },
};
