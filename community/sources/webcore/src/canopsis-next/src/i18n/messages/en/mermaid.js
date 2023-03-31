import { MERMAID_THEMES } from '@/constants';

export default {
  theme: 'Color theme',
  panzoom: {
    helpText: 'Useful shortcuts:\n'
      + 'Ctrl + mouse wheel - zoom in/out\n'
      + 'Shift + mouse wheel - horizontal scroll\n'
      + 'Alt + mouse wheel - vertical scroll\n'
      + 'Ctrl + Left mouse click + drag - pan the area',
  },
  themes: {
    [MERMAID_THEMES.default]: 'Default',
    [MERMAID_THEMES.base]: 'Base',
    [MERMAID_THEMES.dark]: 'Dark',
    [MERMAID_THEMES.forest]: 'Forest',
    [MERMAID_THEMES.neutral]: 'Neutral',
    [MERMAID_THEMES.canopsis]: 'Canopsis',
  },
  errors: {
    emptyMermaid: 'The diagram and points must be added',
  },
};
