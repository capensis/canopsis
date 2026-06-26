import { META_ALARMS_RULE_TYPES } from '@/constants';

export default {
  outputTemplate: 'Modèle de message',
  threshold: 'Seuil',
  thresholdType: 'Type de seuil',
  thresholdRate: 'Taux de déclenchement',
  thresholdHelpText: 'Après avoir atteint ce taux seuil, les alarmes qui correspondent aux modèles et créées pendant l\'intervalle de temps défini sont regroupées\n\n'
    + 'a - alarmes correspondant aux modèles d\'alarme et d\'entité\n'
    + 'b - entités du modèle d\'entité\n'
    + 'c - entités du modèle d\'entité total\n'
    + 'd - toutes les entités dans Canopsis\n'
    + '------------------------------------\n'
    + 'Cas 1 - Si seul le modèle d\'alarme est renseigné\n'
    + 'Taux seuil = a * 100 / d\n\n'
    + 'Cas 2 - Si le modèle d\'entité est renseigné et que le modèle d\'entité total ne l\'est pas\n'
    + 'Taux seuil = a * 100 / b\n\n'
    + 'Cas 3 - Si le modèle d\'entité total est renseigné\n'
    + 'Taux seuil = a * 100 / c',
  thresholdCount: 'Seuil de déclenchement',
  thresholdCountHelpText: 'Après avoir atteint ce seuil, les alarmes qui correspondent aux modèles et créées pendant l\'intervalle de temps défini sont regroupées',
  timeInterval: 'Intervalle de temps',
  timeIntervalHelpText: 'Les alarmes créées pendant cet intervalle de temps sont regroupées',
  childInactiveDelay: 'Délai d\'inactivité de l\'enfant',
  childInactiveDelayHelpText: 'L\'alarme correspondant à cette règle n\'est activée qu\'après le délai d\'inactivité',
  valuePath: 'Chemin de valeur | Chemins de valeur',
  autoResolve: 'Résolution automatique',
  idHelp: 'Si ce champ n\'est pas renseigné, un identifiant unique sera généré automatiquement à la création de la règle',
  corelId: 'Identifiant de corrélation',
  corelIdHelpText: 'Les alarmes avec le même attribut sélectionné sont regroupées',
  corelStatus: 'Statut de corrélation',
  corelStatusHelpText: 'Par ce paramètre, les alarmes sont divisées en parents et enfants',
  corelParent: 'Corrélation parent',
  corelParentHelpText: 'Les alarmes avec cette valeur du champ Corel Status sont définies comme parents',
  corelChild: 'Corrélation enfant',
  corelChildHelpText: 'Les alarmes avec cette valeur du champ Corel Status sont définies comme enfants',
  corelGroupingSummary: 'Les alarmes qui ont la même valeur dans le champ {corelID} et\n'
    + '{corelChild} = {corel}\n'
    + 'seront regroupées sous l\'alarme qui a\n'
    + '{corelParent} = {corel}',
  outputTemplateHelp: '<p>Les variables accessibles sont:</p>\n'
    + '<p><strong>.Count</strong>: Le nombre d\'alarmes conséquences attachées à la méta-alarme.</p>'
    + '<p><strong>.Children</strong>: L\'ensemble des variables de la dernière alarme conséquence attachée à la méta-alarme.</p>'
    + '<p><strong>.Rule</strong>: Les informations administratives de la méta-alarme en elle-même.</p>'
    + '<p>Quelques exemples:</p>'
    + '<p><strong>{{ .Count }} conséquences;</strong> Message de la dernière alarme conséquence : <strong>{{ .Children.Alarm.Value.State.Message }};</strong> Règle : <strong>{{ .Rule.Name }};</strong></p>'
    + '<p>Un message informatif statique</p>'
    + '<p>Corrélé par la règle <strong>{{ .Rule.Name }}</strong></p>',
  removeConfirmationText: 'Lors de la suppression d\'une règle de méta-alarme, toutes les méta-alarmes correspondantes seront également supprimées.\n'
    + 'Êtes-vous sûr de continuer?',
  selectType: 'Sélectionnez le type de règle de méta-alarme',
  grouping: 'Regroupement',
  groupingTabs: {
    existing: 'Regrouper sous une alarme existante',
    createNew: 'Regrouper et créer une nouvelle méta-alarme',
  },
  groupingLabels: {
    groupUnder: 'Regrouper sous',
    groupBy: 'Regrouper par',
  },
  valuePathHelpText: 'Les alarmes ayant la même combinaison de valeurs dans ces champs seront regroupées sous la même méta-alarme',
  componentTemplate: 'Modèle de composant',
  resourceTemplate: 'Modèle de ressource',
  copyTagsFromChildren: 'Copier les balises des alarmes pour enfants',
  filterByLabelEnabled: 'Filtrer par étiquette',
  filterByLabelEnabledTooltip: 'Certaines balises peuvent être définies au format Tag:Value, par exemple Env:Prod.\nAvec le filtre par libellé, seules les balises avec le libellé défini seront copiées des alarmes enfants vers la méta-alarme.',
  copyFromLastChild: 'Copie du dernier enfant',
  copyFromLastChildTooltip: 'Lorsque cette option est activée, la valeur infos est copiée à partir d\'une infos d\'une dernière alarme enfant. Le nom infos de l\'alarme enfant doit être défini dans ce cas.',
  massRemove: 'Supprimer les règles de méta-alarme',
  massEnable: 'Activer les règles de méta-alarme',
  massDisable: 'Désactiver les règles de méta-alarme',
  steps: {
    basics: 'Les bases',
    defineType: 'Définir le type',
    addParameters: 'Ajouter des paramètres',
  },
  types: {
    [META_ALARMS_RULE_TYPES.relation]: {
      label: 'Composant parent',
      text: 'Regrouper sous le composant parent',
      helpText: 'Définir les modèles pour lesquels toutes les alarmes déclenchées sur ses dépendances doivent être regroupées',
    },
    [META_ALARMS_RULE_TYPES.timebased]: {
      label: 'Intervalle de temps (lorsque la gravité change)',
      text: 'Regrouper par intervalle de dates de création',
      helpText: 'Toutes les alarmes qui correspondent aux modèles et déclenchées pendant un intervalle de temps défini sont regroupées',
    },
    [META_ALARMS_RULE_TYPES.attribute]: {
      label: 'Parent uniquement',
      text: 'Regrouper par modèle uniquement',
      helpText: 'Toutes les alarmes dont les attributs sont définis par les modèles de filtre sont regroupées',
    },
    [META_ALARMS_RULE_TYPES.complex]: {
      label: 'Avec seuil ou taux de déclenchement',
      text: 'Regrouper avec seuil ou taux de déclenchement',
      helpText: 'Toutes les alarmes dont les attributs sont définis par les modèles de filtre, l\'intervalle de temps et un seuil ou un taux de déclenchement sont regroupées',
    },
    [META_ALARMS_RULE_TYPES.valuegroup]: {
      label: 'Avec seuil ou taux de déclenchement et valeurs d\'attributs personnalisés',
      text: 'Regrouper par groupe de valeurs',
      helpText: 'Toutes les alarmes dont les attributs sont définis par les modèles de filtre, l\'intervalle de temps, le seuil ou le taux de déclenchement et le chemin de valeur sont regroupées',
    },
    [META_ALARMS_RULE_TYPES.corel]: {
      label: 'Parent défini personnalisé',
      text: 'Regrouper sous un parent défini personnalisé',
      helpText: 'Toutes les alarmes dont les attributs sont définis par les modèles de filtre, l\'intervalle de temps, le nombre de seuils et les identifiants de corrélation sont regroupées',
    },
  },
  patternsTabLabel: 'Modèles et paramètres',
  errors: {
    noValuePaths: 'Vous devez ajouter au moins un chemin de valeur',
  },
  field: {
    title: 'Règle de méta-alarme',
    noData: 'Aucune règle de méta-alarme n\'est trouvée selon les modèles définis',
  },
};
