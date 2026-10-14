import { VIEW_SCREEN_MODES } from '@/constants';

export default {
  sharedViewUrl: 'URL de vue partagée',
  shareView: 'Partager la vue {name}',
  deleteRow: 'Supprimer la ligne',
  deleteWidget: 'Supprimer le widget',
  fullScreen: 'Plein écran',
  copyWidgetId: 'Copier l\'identifiant du widget',
  autoHeightButton: 'Si ce bouton est sélectionné, la hauteur sera calculée automatiquement.',
  selectAll: 'Sélectionnez tous les groupes et vues',
  duplicateAsPrivate: 'Dupliquer en vue privée',
  duplicateAsRegular: 'Dupliquer en vue normale',
  periodicRefresh: 'Rafraichissement périodique',
  groupIds: 'Choisissez une groupe, ou créez-en un nouveau',
  groupTags: 'Étiquettes de groupe',
  noGroupsFound: 'Aucun groupe correspondant. Appuyez sur <kbd>enter</kbd> pour en créer un nouveau.',
  errors: {
    emptyTabs: 'Merci de créer un onglet',
  },
  screenMode: {
    [VIEW_SCREEN_MODES.default]: {
      title: 'Mode par défaut',
      tooltip: 'Alt / Cmd + Shift + 1',
    },
    [VIEW_SCREEN_MODES.fullscreen]: {
      title: 'Plein écran',
      tooltip: 'Alt / Cmd + Shift + 2',
    },
    [VIEW_SCREEN_MODES.kiosk]: {
      title: 'Mode kiosque uniquement',
      tooltip: 'Alt / Cmd + Shift + 3',
    },
    [VIEW_SCREEN_MODES.kioskFullscreen]: {
      title: 'Kiosque + plein écran',
      tooltip: 'Alt / Cmd + Shift + 4',
    },
  },
};
