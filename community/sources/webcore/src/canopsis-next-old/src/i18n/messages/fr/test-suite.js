import { TEST_SUITE_STATUSES } from '@/constants';

export default {
  xmlFeed: 'Flux XML',
  hostname: 'Nom d\'hôte',
  lastUpdate: 'Dernière mise à jour',
  totalTests: 'Total des tests',
  disabledTests: 'Tests désactivés',
  copyMessage: 'Copier le message système',
  systemError: 'Erreur système',
  systemErrorMessage: 'Message d\'erreur système',
  systemOut: 'Retour système',
  systemOutMessage: 'Message de retour du système',
  compareWithHistorical: 'Comparer avec les données historiques',
  className: 'Nom du test',
  line: 'Ligne',
  failureMessage: 'Message d\'échec',
  noData: 'Aucun message système trouvé dans le formulaire XML',
  tabs: {
    globalMessages: 'Messages globaux',
    gantt: 'Gantt',
    details: 'Détails',
    screenshots: 'Captures d\'écran',
    videos: 'Vidéos',
  },
  statuses: {
    [TEST_SUITE_STATUSES.passed]: 'Passé',
    [TEST_SUITE_STATUSES.skipped]: 'Ignoré',
    [TEST_SUITE_STATUSES.error]: 'En erreur',
    [TEST_SUITE_STATUSES.failed]: 'Échoué',
    [TEST_SUITE_STATUSES.total]: 'Temps total pris',
  },
  popups: {
    systemMessageCopied: 'Message système copié dans le presse-papier',
  },
};
