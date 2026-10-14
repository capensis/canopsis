export default {
  title: 'User interface',
  appTitle: 'App title',
  language: 'Default user interface language',
  footer: 'Login footer',
  description: 'Login page description',
  loginPage: 'Login page',
  matchedItems: 'Matched items',
  animations: 'Animations',
  instructionApprovement: 'Instruction approvement',
  advanced: 'Advanced',
  logo: 'Logo',
  infoPopupTimeout: 'Info popup timeout',
  errorPopupTimeout: 'Error popup timeout',
  allowChangeSeverityToInfo: 'Allow change severity to info',
  showHeaderOnKioskMode: 'Show header on kiosk mode',
  showHeaderOnKioskModeTooltip: 'Header won\'t be available in Full screen kiosk mode',
  maxMatchedItems: 'Max matched items',
  checkCountRequestTimeout: 'Check max matched items request timeout (seconds)',
  requiredInstructionApprove: 'Required instruction approvement',
  disabledTransitions: 'Disable some animations',
  disabledTransitionsTooltip: 'It will help to improve application performance',
  autoSuggestPbehaviorName: 'Automatically suggest a pbehavior name',
  defaultTheme: 'Default theme',
  dateFormat: 'Datetime format',
  dateFormats: {
    dateTime: '31/12/2000 00:00:00',
    dateTimeWithWeekday: '31/12/2000 (Mon) 00:00:00',
    weekdayWithDateTime: 'Mon 31/12/2000 00:00:00',
  },
  versionDescriptionTooltip: 'Tooltip popup on version banner',
  defaultVersionDescription: `<div>
  <div>Edition:<strong> {'{{ edition }}'}</strong></div>
  <div>Serial name:<strong> {'{{ serialName }}'}</strong></div>
  <div>Version updated:<strong> {'{{ timestamp versionUpdated }}'}</strong></div>
</div>`,
  tooltips: {
    dateFormat: 'This will be applied to all datetime/timestamp columns and inputs in Canopsis',
    maxMatchedItems:
      'it need to warn user when number of items that match patterns is above this value',
    checkCountRequestTimeout:
      'it need to define request timeout value for max matched items checking',
  },
};
