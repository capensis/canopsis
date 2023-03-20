export default {
  impacts: 'Impacts',
  dependencies: 'Dépendances',
  noEventsFilter: 'Aucun filtre d\'événements',
  impactChain: 'Chaîne d\'impact',
  resolvedAlarms: 'Alarmes résolues',
  activeAlarm: 'Alarme active',
  impactDepends: 'Impacts/Dépendances',
  treeOfDependencies: 'Arbre de dépendances',
  infosSearchLabel: 'Rechercher une info',
  eventStatisticsMessage: '{ok} OK événements\n{ko} KO événements',
  eventStatistics: 'Statistiques d\'événement',
  addSelection: 'Ajouter une sélection',
  advancedSearch: '<span>Aide sur la recherche avancée :</span>\n'
    + '<p>- [ NOT ] &lt;NomColonne&gt; &lt;Opérateur&gt; &lt;Valeur&gt;</p> [ AND|OR [ NOT ] &lt;NomColonne&gt; &lt;Opérateur&gt; &lt;Valeur&gt; ]\n'
    + '<p>Le "-" avant la recherche est obligatoire</p>\n'
    + '<p>Opérateurs:\n'
    + '    <=, <,=, !=,>=, >, LIKE (Pour les expressions régulières MongoDB)</p>\n'
    + '<p>Les types de valeurs : Chaîne de caractères entre guillemets doubles, Booléen ("TRUE", "FALSE"), Entier, Nombre flottant, "NULL"</p>\n'
    + '<dl><dt>Exemples :</dt><dt>- Name = "name_1"</dt>\n'
    + '    <dd>Entités dont le names est "name_1"</dd><dt>- Name="name_1" AND Type="service"</dt>\n'
    + '    <dd>Entités dont le names est "name_1" et la types est "service"</dd><dt>- infos.custom.value="Custom value" OR Type="resource"</dt>\n'
    + '    <dd>Entités dont le infos.custom.value est "Custom value" ou la type est "resource"</dd><dt>- infos.custom.value LIKE 1 OR infos.custom.value LIKE 2</dt>\n'
    + '    <dd>Entités dont le infos.custom.value contient 1 or 2</dd><dt>- NOT Name = "name_1"</dt>\n'
    + '    <dd>Entités dont le name n\'est pas "name_1"</dd>\n'
    + '</dl>',
  actions: {
    titles: {
      editEntity: 'Éditer l\'entité',
      duplicateEntity: 'Dupliquer l\'entité',
      deleteEntity: 'Supprimer l\'entité',
      pbehavior: 'Comportement périodique',
      variablesHelp: 'Liste des variables disponibles',
      massEnable: 'Activer les entités',
      massDisable: 'Désactiver les entités',
    },
  },
  fab: {
    common: 'Ajouter une nouvelle entité',
    addService: 'Ajouter une nouvelle entité de service',
  },
  popups: {
    massDeleteWarning: 'La suppression en masse ne peut pas être appliquée pour certains des éléments sélectionnés, ils ne seront donc pas supprimés.',
  },
};
