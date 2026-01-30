import { ADVANCED_SEARCH_GROUPS } from '@/constants';

export default {
  title: 'Recherche Avancée',
  not: 'NOT',
  valueTypeListEmptyMessage: 'Appuyez sur au moins 1 symbole et appuyez sur <kbd>Entrée</kbd> pour appliquer',
  switchAdvancedSearchActiveToTrue: 'Passer à la recherche avancée',
  switchAdvancedSearchActiveToFalse: 'Passer à la recherche simple',
  noDataList: 'Il n\'y a aucun élément à saisir',
  inputPlaceholder: 'Rechercher ou filtrer les résultats...',
  definedDifferent: 'La valeur diffère de celle définie',

  groups: {
    [ADVANCED_SEARCH_GROUPS.basic]: 'Basique',
    [ADVANCED_SEARCH_GROUPS.messages]: 'Messages',
    [ADVANCED_SEARCH_GROUPS.ticket]: 'Ticket',
    [ADVANCED_SEARCH_GROUPS.dates]: 'Dates',
    [ADVANCED_SEARCH_GROUPS.actions]: 'Actions',
    [ADVANCED_SEARCH_GROUPS.entity]: 'Entité',
    [ADVANCED_SEARCH_GROUPS.alias]: 'Alias',
    [ADVANCED_SEARCH_GROUPS.pbehavior]: 'Comportements périodiques',
  },

  searchForThisText: 'Appuyez sur <kbd>Entrée</kbd> pour rechercher ce texte',
  searchByText: 'Rechercher par texte',
  listDisabledMessage: 'Il n\'est pas possible de combiner des modèles\n(alarme, entité et comportement) avec OU (uniquement ET)',
};
