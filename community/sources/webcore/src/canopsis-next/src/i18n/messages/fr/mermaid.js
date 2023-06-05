import { MERMAID_THEMES } from '@/constants';

export default {
  theme: 'Thème de couleur',
  panzoom: {
    helpText: 'Raccourcis utiles :\n'
      + 'Ctrl + molette de la souris - zoom avant/arrière\n'
      + 'Maj + molette de la souris - défilement horizontal\n'
      + 'Alt + molette de la souris - défilement vertical\n'
      + 'Ctrl + Clic gauche de la souris + glisser - déplacer la zone',
  },
  themes: {
    [MERMAID_THEMES.default]: 'Défaut',
    [MERMAID_THEMES.base]: 'Base',
    [MERMAID_THEMES.dark]: 'Sombre',
    [MERMAID_THEMES.forest]: 'Forêt',
    [MERMAID_THEMES.neutral]: 'Neutre',
    [MERMAID_THEMES.canopsis]: 'Canopsis',
  },
  errors: {
    emptyMermaid: 'Le diagramme et les points doivent être ajoutés',
  },
};
