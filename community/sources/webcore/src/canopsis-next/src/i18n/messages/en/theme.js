import { THEME_FONT_SIZES } from '@/constants/theme';

export default {
  colorThemes: 'Color themes',
  exampleText: 'Hello world!',
  defaultTheme: 'The theme is default, you cannot to edit theme!',
  errors: {
    notReadable: 'Text is not readable',
  },
  main: {
    title: 'Main UI elements',

    primary: 'Main brand color',
    primaryHelpText: 'Main brand color (Canopsis header)',

    secondary: 'Secondary brand color',
    secondaryHelpText: 'Additional brand color (for expanded panels, menus, etc)',

    accent: 'Neutral buttons color',
    accentHelpText: 'Color for neutral buttons (next/previous, etc)',

    error: 'Error color',
    errorHelpText: 'Color for error messages, negative action buttons, etc',

    info: 'Info color',
    infoHelpText: 'Color for neutral messages and notifications',

    success: 'Success/positive color',
    successHelpText: 'Color for positive/success messages and notifications',

    warning: 'Warning color',
    warningHelpText: 'Color for warning messages and notifications',

    background: 'Main background color',

    activeColor: 'Main active color',
    activeColorHelpText: 'Main color for texts and icons',
  },
  fontSize: {
    title: 'Font size settings',

    sizes: {
      [THEME_FONT_SIZES.small]: 'Small',
      [THEME_FONT_SIZES.medium]: 'Medium',
      [THEME_FONT_SIZES.large]: 'Large',
    },
  },
  table: {
    title: 'Table settings',

    background: 'Table background color',
    backgroundHelpText: 'BG color for the alarm list table',

    rowColor: 'Table row color',
    rowColorHelpText: 'BG color for the each table row',

    shiftRowEnable: 'Shift table background colors',
    shiftRowEnableHelpText: 'Switcher to enable/disable color shifts for table rows',

    shiftRowColor: 'Second table row background color',
    shiftRowColorHelpText: 'When enabled, rows colors are switching (every second row color is different)',

    hoverRowEnable: 'Change row color on hover',
    hoverRowEnableHelpText: 'Switcher to enable/disable table row color change on hover',

    hoverRowColor: 'Table row color on hover',
  },
  state: {
    title: 'Severity colors',

    ok: 'Ok',
    okHelpText: 'Color indication for the OK state',

    minor: 'Minor',
    minorHelpText: 'Color indication for the minor state',

    major: 'Major',
    majorHelpText: 'Color indication for the major state',

    critical: 'Critical',
    criticalHelpText: 'Color indication for the critical state',
  },
};
