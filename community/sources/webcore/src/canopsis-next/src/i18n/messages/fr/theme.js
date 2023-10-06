import { THEME_FONT_SIZES } from '@/constants/theme';

export default {
  themes: 'Thèmes',
  exampleText: 'Bonjour le monde!',
  defaultTheme: 'Le thème est par défaut, vous ne pouvez pas modifier le thème !',
  errors: {
    notReadable: 'Le texte n\'est pas lisible',
  },
  main: {
    title: 'Principaux éléments de l\'interface utilisateur',

    primary: 'Couleur principale de la marque',
    primaryHelpText: 'Couleur principale de la marque (en-tête Canopsis)',

    secondary: 'Couleur de marque secondaire',
    secondaryHelpText: 'Couleur de marque supplémentaire (pour les panneaux développés, les menus, etc.)',

    accent: 'Couleur neutre des boutons',
    accentHelpText: 'Couleur des boutons neutres (suivant/précédent, etc.)',

    error: 'Couleur d\'erreur',
    errorHelpText: 'Couleur des messages d\'erreur, des boutons d\'action négative, etc.',

    info: 'Couleur des informations',
    infoHelpText: 'Couleur pour les messages et notifications neutres',

    success: 'Succès/couleur positive',
    successHelpText: 'Couleur pour les messages et notifications positifs/succès',

    warning: 'Couleur d\'avertissement',
    warningHelpText: 'Couleur des messages d\'avertissement et des notifications',

    background: 'Couleur de fond principale',

    activeColor: 'Couleur active principale',
    activeColorHelpText: 'Couleur principale des textes et des icônes',
  },
  fontSize: {
    title: 'Paramètres de taille de police',

    sizes: {
      [THEME_FONT_SIZES.small]: 'Petite',
      [THEME_FONT_SIZES.medium]: 'Moyen',
      [THEME_FONT_SIZES.large]: 'Grande',
    },
  },
  table: {
    title: 'Paramètres du tableau',

    background: 'Couleur d’arrière-plan du tableau',
    backgroundHelpText: 'Couleur BG pour le tableau de la liste des alarmes',

    rowColor: 'Couleur des lignes du tableau',
    rowColorHelpText: 'Couleur BG pour chaque ligne du tableau',

    shiftRowEnable: 'Décaler les couleurs d’arrière-plan du tableau',
    shiftRowEnableHelpText: 'Sélecteur pour activer/désactiver les changements de couleur pour les lignes du tableau',

    shiftRowColor: 'Couleur d’arrière-plan de la deuxième ligne du tableau',
    shiftRowColorHelpText: 'Lorsqu\'elle est activée, les couleurs des lignes changent (une couleur de ligne sur deux est différente)',

    hoverRowEnable: 'Changer la couleur de la ligne au survol',
    hoverRowEnableHelpText: 'Sélecteur pour activer/désactiver le changement de couleur des lignes du tableau en survol',

    hoverRowColor: 'Couleur des lignes du tableau au survol',
  },
  state: {
    title: 'Couleurs de gravité',

    ok: 'Ok',
    okHelpText: 'Indication de couleur pour l\'état OK',

    minor: 'Mineur',
    minorHelpText: 'Indication de couleur pour l\'état mineur',

    major: 'Majeur',
    majorHelpText: 'Indication de couleur pour l\'état majeur',

    critical: 'Critique',
    criticalHelpText: 'Indication de couleur pour l\'état critique',
  },
};
