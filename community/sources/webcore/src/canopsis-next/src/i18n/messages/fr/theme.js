import { THEME_FONT_SIZES } from '@/constants/theme';

export default {
  tabs: {
    mainColors: 'Couleurs principales',
    widgetTable: 'Couleurs des widgets et tableaux',
    severity: 'Couleurs de criticité',
    fontSize: 'Taille de police',
    extraDetails: 'Détails supplémentaires',
  },
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
    secondaryHelpText:
      'Couleur supplémentaire de l\'identité visuelle (pour les panneaux élargis, menus, etc.)',

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
  widget: {
    title: 'Widget',
    titleBackground: 'Arrière-plan du titre du widget',
    contentBackground: 'Arrière-plan du contenu du widget',
    border: 'Définir la bordure',
    borderColor: 'Couleur de la bordure du widget',
    borderRadius: 'Rayon (px)',
  },
  preview: {
    title: 'Aperçu',
    widgetTitle: 'Titre du widget',
    header: 'En-tête',
    tableRow: 'Ligne du tableau',
    secondTableRow: 'Deuxième ligne du tableau',
    hoveredRow: 'Ligne survolée',
  },
  table: {
    title: 'Tableau',

    background: 'En-tête du tableau',

    rowColor: 'Ligne des tableaux',

    shiftRowEnable: 'Alterner les couleurs d’arrière-plan des tableaux',
    shiftRowEnableHelpText:
      'Sélecteur pour activer/désactiver les changements de couleur pour les lignes des tableaux',

    shiftRowColor: 'Deuxième ligne du tableau',

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
  extraDetails: {
    ack: 'Acquittement / acquittement rapide',
    canceled: 'Annulée',
    children: 'Enfants',
    parent: 'Parent',
    comment: 'Commentaire',
    snooze: 'Mise en veille',
    tickets: 'Tickets',
    cancel_delay: 'Délai d’annulation',
  },
};
