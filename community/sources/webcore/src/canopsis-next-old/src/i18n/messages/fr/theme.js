import { THEME_FONT_SIZES } from '@/constants/theme';

export default {
  themes: 'Thèmes',
  exampleText: 'Bonjour le monde!',
  defaultTheme: 'Le thème est par défaut, vous ne pouvez pas modifier le thème !',
  checkColors: 'Vérifier les couleurs',
  errors: {
    notReadable: 'Le texte n\'est pas lisible',
  },
  main: {
    title: 'Couleurs principales de l\'interface',

    primary: 'Marque principale',
    primaryHelpText: 'Couleur principale de l\'identité visuelle (en-tête Canopsis)',

    secondary: 'Marque secondaire',
    secondaryHelpText: 'Couleur supplémentaire de l\'identité visuelle (pour les panneaux élargis, menus, etc.)',

    accent: 'Boutons neutres',

    error: 'Icônes / boutons d\'erreur',
    errorBackground: 'Arrière-plan d\'erreur',

    warning: 'Icônes / boutons d\'avertissement',
    warningBackground: 'Arrière-plan d\'avertissement',

    success: 'Icônes / boutons de réussite / positif',
    successBackground: 'Arrière-plan de réussite / positif',

    info: 'Icônes / boutons d\'information',
    infoBackground: 'Arrière-plan d\'information',

    background: 'Arrière-plan principal\n',

    activeColor: 'Principal actif',
    activeColorHelpText: 'Couleur principale pour les textes et les icônes',
  },
  fontSize: {
    title: 'Paramètres de taille de police',

    sizes: {
      [THEME_FONT_SIZES.small]: 'Petite',
      [THEME_FONT_SIZES.medium]: 'Moyenne',
      [THEME_FONT_SIZES.large]: 'Grande',
    },
  },
  table: {
    title: 'Paramètres des tableaux',

    background: 'Arrière-plan des tableaux',

    rowColor: 'Ligne des tableaux',

    shiftRowEnable: 'Alterner les couleurs d’arrière-plan des tableaux',
    shiftRowEnableHelpText: 'Sélecteur pour activer/désactiver les changements de couleur pour les lignes des tableaux',

    shiftRowColor: 'Arrière-plan de la deuxième ligne des tableaux',

    hoverRowEnable: 'Changer la couleur de la ligne au survol',

    hoverRowColor: 'Couleur des lignes des tableaux au survol',
  },
  state: {
    title: 'Couleurs de criticités',

    ok: 'Ok',

    minor: 'Mineure',

    major: 'Majeure',

    critical: 'Critique',
  },
};
