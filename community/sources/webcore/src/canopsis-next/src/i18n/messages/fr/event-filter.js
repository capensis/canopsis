import { EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES, EVENT_FILTER_FAILURE_TYPES, EVENT_FILTER_TYPES } from '@/constants';

export default {
  externalData: 'Données externes',
  actionsRequired: 'Veuillez ajouter au moins une action',
  configRequired: 'Aucune configuration définie. Veuillez ajouter au moins un paramètre de configuration',
  idHelp: 'Si ce champ n\'est pas renseigné, un identifiant unique sera généré automatiquement à la création de la règle',
  editPattern: 'Éditer le modèle',
  advanced: 'Avancée',
  addAField: 'Ajouter un champ',
  simpleEditor: 'Éditeur simple',
  field: 'Champ',
  value: 'Valeur',
  advancedEditor: 'Éditeur avancé',
  comparisonRules: 'Règles de comparaison',
  editActions: 'Actions',
  addAction: 'Ajouter une action',
  editAction: 'Éditer une action',
  actions: 'Actions',
  onSuccess: 'En cas de succès',
  onFailure: 'En cas d\'échec',
  configuration: 'Configuration',
  resource: 'ID de ressource ou modèle',
  component: 'ID de composant ou modèle',
  connector: 'ID ou modèle de connecteur',
  connectorName: 'Nom ou modèle de connecteur',
  duringPeriod: 'Appliqué pendant cette période uniquement',
  enrichmentOptions: 'Options d\'enrichissement',
  changeEntityOptions: 'Modifier les options d\'entité',
  eventsFilteredSinceLastUpdate: 'Evénements filtrés depuis la dernière mise à jour',
  errorsSinceLastUpdate: 'Erreurs depuis la dernière mise à jour',
  markAsRead: 'Marquer comme lu',
  filterByType: 'Filtrer par type',
  copyEventToClipboard: 'Copier l\'événement dans le presse papier',
  event: 'Evénement',
  eventCopied: 'Evénement copié dans le presse papier',
  syntaxIsValid: 'La syntaxe est valide',
  types: {
    [EVENT_FILTER_TYPES.drop]: 'Suppression',
    [EVENT_FILTER_TYPES.break]: 'Break',
    [EVENT_FILTER_TYPES.enrichment]: 'Enrichissement',
    [EVENT_FILTER_TYPES.changeEntity]: 'Changement d\'entité',
  },
  failureTypes: {
    [EVENT_FILTER_FAILURE_TYPES.invalidPattern]: 'Pattern invalide',
    [EVENT_FILTER_FAILURE_TYPES.invalidTemplate]: 'Template invalide',
    [EVENT_FILTER_FAILURE_TYPES.externalDataMongo]: 'Mongo',
    [EVENT_FILTER_FAILURE_TYPES.externalDataApi]: 'API externe',
    [EVENT_FILTER_FAILURE_TYPES.other]: 'Autre',
  },
  tooltips: {
    addValueRuleField: 'Ajouter une règle',
    editValueRuleField: 'Éditer la règle',
    addObjectRuleField: 'Ajouter un groupe de règles',
    editObjectRuleField: 'Éditer le groupe de règles',
    removeRuleField: 'Supprimer le groupe/la règle',
  },
  validation: {
    incorrectRegexOnSetTagsValue: 'Valeur non valide : la valeur de l\'action set_tags doit contenir une expression régulière pour extraire les groupes <name> et <value>',
  },
  actionsTypes: {
    [EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES.copy]: {
      text: 'Copier une valeur d\'un champ d\'événement à un autre',
      message: 'Cette action permet de copier la valeur ou une paire clé+valeur d\'un contrôle dans un événement.',
      description: '<h3 class="text-subtitle-1 font-weight-bold">Paramètres de l\'action</h3>'
      + '<ul>'
        + '<li><strong>value</strong> : nom du contrôle dont la valeur doit être copiée. Cela peut être un champ d\'événement, un sous-groupe d\'une expression régulière ou une donnée externe</li>'
        + '<li><strong>description</strong> (optionnel) : la description</li>'
        + '<li><strong>name</strong> : nom du champ d\'événement dans lequel la valeur doit être copiée</li>'
      + '</ul>',
    },
    [EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES.copyToEntityInfo]: {
      text: 'Copier une valeur d\'un champ d\'un événement vers une information d\'une entité',
      message: 'Cette action est utilisée pour copier la valeur du champ d\'un événement dans le champ d\'une entité.',
      description: '<h3 class="text-subtitle-1 font-weight-bold">Paramètres de l\'action</h3>'
      + '<ul>'
        + '<li><strong>description</strong> (optionnel) : la description</li>'
        + '<li><strong>name</strong> : nom du champ d\'une entité</li>'
        + '<li><strong>value</strong> : nom du contrôle dont la valeur doit être copiée. Cela peut être un champ d\'événement, un sous-groupe d\'une expression régulière ou une donnée externe</li>'
      + '</ul>',
    },
    [EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES.setEntityInfo]: {
      text: 'Définir une information d\'une entité sur une constante',
      message: 'Cette action permet de définir les informations dynamiques d\'une entité correspondant à l\'événement.',
      description: '<h3 class="text-subtitle-1 font-weight-bold">Paramètres de l\'action</h3>'
      + '<ul>'
        + '<li><strong>description</strong> (optionnel) : la description</li>'
        + '<li><strong>name</strong> : le nom du champ</li>'
        + '<li><strong>value</strong> : la valeur d\'un champ</li>'
      + '</ul>',
    },
    [EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES.setEntityInfoFromTemplate]: {
      text: 'Définir une chaîne d\'informations sur une entité à l\'aide d\'un modèle',
      message: 'Cette action permet de modifier les informations dynamiques d\'une entité correspondant à l\'événement.',
      description: '<h3 class="text-subtitle-1 font-weight-bold">Paramètres de l\'action</h3>'
      + '<ul>'
        + '<li><strong>description</strong> (optionnel) : la description</li>'
        + '<li><strong>name</strong> : le nom du champ</li>'
        + '<li><strong>value</strong> : le modèle utilisé pour déterminer la valeur de la donnée. '
          + 'Les modèles <code>{{.Event.NomDuChamp}}</code>, les expressions régulières ou les données externes peuvent être utilisés</li>'
      + '</ul>',
    },
    [EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES.setField]: {
      text: 'Définir un champ d\'un événement sur une constante',
      message: 'Cette action peut être utilisée pour modifier un champ de l\'événement.',
      description: '<h3 class="text-subtitle-1 font-weight-bold">Paramètres de l\'action</h3>'
      + '<ul>'
        + '<li><strong>description</strong> (optionnel) : la description</li>'
        + '<li><strong>name</strong> : le nom du champ</li>'
        + '<li><strong>value</strong> : la nouvelle valeur du champ</li>'
      + '</ul>',
    },
    [EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES.setFieldFromTemplate]: {
      text: 'Définir un champ de chaîne d\'un événement à l\'aide d\'un modèle',
      message: 'Cette action vous permet de modifier un champ d\'événement à partir d\'un modèle.',
      description: '<h3 class="text-subtitle-1 font-weight-bold">Paramètres de l\'action</h3>'
      + '<ul>'
        + '<li><strong>description</strong> (optionnel) : la description</li>'
        + '<li><strong>name</strong> : le nom du champ</li>'
        + '<li><strong>value</strong> : le modèle utilisé pour déterminer la valeur du champ. '
          + 'Les modèles <code>{{.Event.NomDuChamp}}</code>, les expressions régulières ou les données externes peuvent être utilisés</li>'
      + '</ul>',
    },
    [EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES.setEntityInfoFromDictionary]: {
      text: 'Définir plusieurs chaînes d\'informations sur une entité à partir d\'un dictionnaire',
      message: 'Cette action peut être utilisée pour définir plusieurs informations d\'entité à partir d\'un dictionnaire.',
      description: '<h3 class="text-subtitle-1 font-weight-bold">Paramètres de l\'action</h3>'
      + '<ul>'
        + '<li><strong>description</strong> (optionnel) : la description utilisée pour les informations de l\'entité. '
          + 'Si elle n\'est pas définie, les informations de l\'entité seront laissées vides</li>'
        + '<li><strong>value</strong> : le champ de l\'événement à partir duquel les informations sont récupérées. '
          + 'La valeur doit contenir un tableau de paires nom: valeur</li>'
      + '</ul>',
    },
    [EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES.setTags]: {
      text: 'Définir les balises d\'un champ à l\'aide d\'une correspondance d\'expression rationnelle',
      message: 'Cette action peut être utilisée pour définir des balises provenant d\'autres événements filtrés à l\'aide d\'une correspondance d\'expression rationnelle.',
      description: '<p>'
        + 'L\'action <strong>set_tags</strong> permet de créer dynamiquement des tags au format '
        + '<strong>« Nom: Valeur »</strong> à partir d\'un champ de <strong>l\'événement en cours de traitement</strong>, '
        + 'en utilisant des <strong>groupes de capture</strong> définis dans une expression régulière.'
        + '</p>'

        + '<p>'
        + 'Cette expression régulière doit être appliquée en amont, dans un <strong>filtre d\'événement</strong>, '
        + 'sur un champ textuel contenant les informations à transformer en tags.'
        + '</p>'

        + '<p>'
        + 'L\'expression régulière doit inclure deux groupes nommés :'
        + '</p>'
        + '<ul>'
        + '<li><code>(?P&lt;name&gt;...)</code> pour extraire le <strong>nom</strong> du tag</li>'
        + '<li><code>(?P&lt;value&gt;...)</code> pour extraire la <strong>valeur</strong> du tag</li>'
        + '</ul>'

        + '<p>'
        + 'Une fois le filtre appliqué et les groupes détectés, l\'action <strong>set_tags</strong> utilise ces valeurs '
        + 'pour générer automatiquement les tags correspondants.'
        + '</p>'

        + '<h3 class="text-subtitle-1 font-weight-bold mt-4 mb-2">Exemples d\'expressions régulières</h3>'

        + '<table class="striped">'
        + '<thead class="grey lighten-2">'
        + '<tr>'
        + '<th class="pa-2">Format attendu dans le champ source</th>'
        + '<th class="pa-2">Expression régulière</th>'
        + '</tr>'
        + '</thead>'
        + '<tbody>'
        + '<tr>'
        + '<td class="pa-2"><code>valeur nom;</code></td>'
        + '<td class="pa-2"><code>(?P&lt;value&gt;[a-zA-Z]+)\\s+(?P&lt;name&gt;[a-zA-Z]+);</code></td>'
        + '</tr>'
        + '<tr>'
        + '<td class="pa-2"><code>nom valeur;</code></td>'
        + '<td class="pa-2"><code>(?P&lt;name&gt;[a-zA-Z]+)\\s+(?P&lt;value&gt;[a-zA-Z]+);</code></td>'
        + '</tr>'
        + '<tr>'
        + '<td class="pa-2"><code>nom: valeur;</code></td>'
        + '<td class="pa-2"><code>(?P&lt;name&gt;[a-zA-Z]+):\\s+(?P&lt;value&gt;[a-zA-Z]+);</code></td>'
        + '</tr>'
        + '</tbody>'
        + '</table>'

        + '<h3 class="text-subtitle-1 font-weight-bold mt-4 mb-2">Paramètres de l\'action</h3>'
        + '<ul>'
        + '<li><strong>description</strong> (optionnel) : commentaire ou description libre de l\'action.</li>'
        + '<li><strong>value</strong> (obligatoire) : nom du champ de l\'événement sur lequel ont été appliqués '
        + 'les groupes de capture name et value.</li>'
        + '</ul>',
    },
    [EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES.setTagsFromTemplate]: {
      text: 'Définir les balises d\'un champ à l\'aide d\'un modèle',
      message: 'Cette action peut être utilisée pour définir des balises provenant d\'autres champs d\'événement à l\'aide d\'un modèle.',
      description: '<p>'
        + 'Cette action peut être utilisée pour définir des tags à partir d\'autres champs d\'événement en utilisant un modèle.'
        + '</p>'

        + '<p>'
        + 'L\'action <strong>set_tags_from_template</strong> permet d\'ajouter un <strong>tag unique</strong> au format '
        + '<strong>"Nom: Valeur"</strong>, construit dynamiquement à partir d\'un modèle basé sur les champs de l\'événement.'
        + '</p>'

        + '<p>'
        + 'Cette action est utile lorsque vous souhaitez définir un tag dont la <strong>valeur</strong> est calculée à partir '
        + 'du contenu d\'un ou plusieurs champs de l\'événement, en utilisant la syntaxe de templating Go (<code>{{.Event.Field}}</code>).'
        + '</p>'

        + '<p>'
        + 'Cette action permet uniquement la création d\'<strong>un seul tag à la fois</strong>.'
        + '</p>'

        + '<h3 class="text-subtitle-1 font-weight-bold mt-4 mb-2">Paramètres de l\'action</h3>'
        + '<ul>'
        + '<li><strong>description</strong> (optionnel) : commentaire ou description libre de l\'action.</li>'
        + '<li><strong>name</strong> (obligatoire) : nom du tag à créer.</li>'
        + '<li><strong>value</strong> (obligatoire) : modèle utilisé pour générer la valeur du tag.'
        + '<br>Il peut contenir :'
        + '<ul>'
        + '<li>des références aux champs de l\'événement (<code>{{.Event.field}}</code>)</li>'
        + '<li>des expressions régulières si le champ a été filtré auparavant</li>'
        + '<li>ou des données provenant d\'une source externe si elles ont été injectées dans le contexte.</li>'
        + '</ul>'
        + '</li>'
        + '</ul>',
    },
  },
};
