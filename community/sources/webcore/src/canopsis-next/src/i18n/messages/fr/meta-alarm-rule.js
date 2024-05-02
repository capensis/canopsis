import { META_ALARMS_RULE_TYPES } from '@/constants';

export default {
  outputTemplate: 'Modèle de message',
  thresholdType: 'Type de seuil',
  thresholdRate: 'Taux de déclenchement',
  thresholdRateHelpText: 'Après avoir atteint ce taux seuil, les alarmes qui correspondent aux modèles et créées pendant l\'intervalle de temps défini sont regroupées',
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
  valuePathHelpText: 'Attribut personnalisé pour regrouper les alarmes définies par un chemin de valeur',
  steps: {
    basics: 'Les bases',
    defineType: 'Définir le type',
    addParameters: 'Ajouter des paramètres',
  },
  types: {
    [META_ALARMS_RULE_TYPES.relation]: {
      text: 'Relation parent-enfant',
      helpText: 'Toutes les alarmes déclenchées sur les entités dépendantes sont regroupées',
    },
    [META_ALARMS_RULE_TYPES.timebased]: {
      text: 'Regroupement par intervalle de temps',
      helpText: 'Toutes les alarmes déclenchées pendant un intervalle de temps défini sont regroupées',
    },
    [META_ALARMS_RULE_TYPES.attribute]: {
      text: 'Regroupement par attribut',
      helpText: 'Toutes les alarmes qui correspondent à un modèle avec des attributs définis sont regroupées',
    },
    [META_ALARMS_RULE_TYPES.complex]: {
      text: 'Regroupement complexe avec seuil ou taux de déclenchement',
      helpText: 'Toutes les alarmes qui correspondent à un modèle avec des attributs définis pendant l\'intervalle de temps défini sont regroupées',
    },
    [META_ALARMS_RULE_TYPES.valuegroup]: {
      text: 'Regroupement par groupe de valeurs',
      helpText: 'Il s\'agit d\'un regroupement complexe avec des chemins de valeurs comme paramètres supplémentaires pour le regroupement',
    },
    [META_ALARMS_RULE_TYPES.corel]: {
      text: 'Regroupement par identifiants de corrélation',
      helpText: 'Regroupement des alarmes déjà corrélées : toutes les alarmes d\'un même identifiant de corrélation sont regroupées',
    },
  },
  parametersTitle: {
    [META_ALARMS_RULE_TYPES.relation]: 'Relation parent-enfant',
    [META_ALARMS_RULE_TYPES.timebased]: 'Relation basée sur le temps',
    [META_ALARMS_RULE_TYPES.attribute]: 'Relation basée sur les attributs',
    [META_ALARMS_RULE_TYPES.complex]: 'Regroupement complexe avec un seuil ou un taux de déclenchement',
    [META_ALARMS_RULE_TYPES.valuegroup]: 'Regroupement par groupe de valeurs',
    [META_ALARMS_RULE_TYPES.corel]: 'Regroupement par identifiants de corrélation',
  },
  parametersDescription: {
    [META_ALARMS_RULE_TYPES.relation]: 'Définir les modèles de filtre pour les entités dont toutes les alarmes déclenchées sur ses dépendances doivent être regroupées',
    [META_ALARMS_RULE_TYPES.timebased]: 'Toutes les alarmes qui correspondent aux modèles et déclenchées pendant un intervalle de temps défini sont regroupées',
    [META_ALARMS_RULE_TYPES.attribute]: 'Toutes les alarmes dont les attributs sont définis par des modèles de filtre sont regroupées',
    [META_ALARMS_RULE_TYPES.complex]: 'Toutes les alarmes dont les attributs sont définis par des modèles de filtre, un intervalle de temps et un seuil ou un taux de déclenchement sont regroupées',
    [META_ALARMS_RULE_TYPES.valuegroup]: 'Toutes les alarmes dont les attributs sont définis par des modèles de filtre, un intervalle de temps, un seuil ou un taux et le chemin de valeur sont regroupées',
    [META_ALARMS_RULE_TYPES.corel]: 'Toutes les alarmes dont les attributs sont définis par des modèles de filtre, un intervalle de temps, un nombre de seuils et des identifiants de corrélation sont regroupées.',
  },
  errors: {
    noValuePaths: 'Vous devez ajouter au moins un chemin de valeur',
  },
};
