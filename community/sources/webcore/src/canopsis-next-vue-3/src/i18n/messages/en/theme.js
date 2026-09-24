import { THEME_FONT_SIZES } from '@/constants/theme';

export default {
  tabs: {
    mainColors: 'Main colors',
    widgetTable: 'Widget and table colors',
    severity: 'Severity colors',
    fontSize: 'Font size',
  },
  themes: 'Themes',
  exampleText: 'Hello world!',
  defaultTheme: 'The theme is default, you cannot to edit theme!',
  checkColors: 'Check colors',
  errors: {
    notReadable: 'Text is not readable',
  },
  main: {
    title: 'Main UI colors',

    primary: 'Main brand',
    primaryHelpText: 'Main brand color (Canopsis header)',

    secondary: 'Secondary brand',
    secondaryHelpText: 'Additional brand color (for expanded panels, menus, etc)',

    accent: 'Neutral buttons',

    error: 'Error icons / buttons',
    errorBackground: 'Error background',

    warning: 'Warning icons / buttons',
    warningBackground: 'Warning background',

    success: 'Success / positive icons',
    successBackground: 'Success / positive background',

    info: 'Info icons / buttons',
    infoBackground: 'Info background',

    background: 'Main background',

    activeColor: 'Main active',
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
  widget: {
    title: 'Widget',
    titleBackground: 'Widget title background',
    contentBackground: 'Widget content background',
    border: 'Set border',
    borderColor: 'Widget border color',
    borderRadius: 'Radius (px)',
  },
  preview: {
    title: 'Preview',
    widgetTitle: 'Widget title',
    header: 'Header',
    tableRow: 'Table row',
    secondTableRow: 'Second table row',
    hoveredRow: 'Hovered row',
  },
  table: {
    title: 'Table',

    background: 'Table header',

    rowColor: 'Table row',

    shiftRowEnable: 'Shift table background colors',
    shiftRowEnableHelpText: 'Switcher to enable/disable color shifts for table rows',

    shiftRowColor: 'Second table row',

    hoverRowEnable: 'Change row color on hover',

    hoverRowColor: 'Table row color on hover',
  },
  state: {
    title: 'Severity colors',

    ok: 'Ok',

    minor: 'Minor',

    major: 'Major',

    critical: 'Critical',
  },
};
