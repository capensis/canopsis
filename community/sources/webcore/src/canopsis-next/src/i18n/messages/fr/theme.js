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

    primary: 'Couleur principale du produit',
    primaryHelpText: 'Couleur principale du produit (en-tête Canopsis)',

    secondary: 'Couleur secondaire du produit',
    secondaryHelpText: 'Couleur secondaire (pour les panneaux développés, les menus, etc.)',

    accent: 'Couleur neutre des boutons',
    accentHelpText: 'Couleur des boutons neutres (suivant/précédent, etc.)',

    error: 'Couleur relative aux erreurs',
    errorHelpText: 'Couleur des messages d\'erreur, des boutons d\'action en échec, etc.',

    info: 'Couleur relatives aux informations',
    infoHelpText: 'Couleur pour les messages et notifications informatifs',

    success: 'Couleur relative aux Succès',
    successHelpText: 'Couleur pour les messages et notifications en succès',

    warning: 'Couleur relative aux avertissements',
    warningHelpText: 'Couleur des messages d\'avertissement et des notifications',

    background: 'Couleur de fond principale',

    activeColor: 'Couleur active principale',
    activeColorHelpText: 'Couleur principale des textes et des icônes',
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
    title: 'Paramètres du bac à alarmes',

    background: 'Couleur d\’arrière-plan du bac',
    backgroundHelpText: 'Couleur d\'arrière plan pour le bac à alarmes',

    rowColor: 'Couleur des lignes du bac',
    rowColorHelpText: 'Couleur des lignes du tableau',

    shiftRowEnable: 'Alterner les couleurs d’arrière-plan du bac',
    shiftRowEnableHelpText: 'Sélecteur pour activer/désactiver les changements de couleur pour les lignes du bac',

    shiftRowColor: 'Couleur de la deuxième ligne du bac',
    shiftRowColorHelpText: 'Lorsqu\'elle est activée, les couleurs des lignes changent (une couleur de ligne sur deux est différente)',

    hoverRowEnable: 'Changer la couleur de la ligne au survol',
    hoverRowEnableHelpText: 'Sélecteur pour activer/désactiver le changement de couleur des lignes du bac au survol',

    hoverRowColor: 'Couleur des lignes du bac au survol',
  },
  state: {
    title: 'Couleurs de criticités',

    ok: 'Ok',
    okHelpText: 'Indication de couleur pour l\'état OK',

    minor: 'Mineure',
    minorHelpText: 'Indication de couleur pour l\'état mineur',

    major: 'Majeure',
    majorHelpText: 'Indication de couleur pour l\'état majeur',

    critical: 'Critique',
    criticalHelpText: 'Indication de couleur pour l\'état critique',
  },
};
