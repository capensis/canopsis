export default {
  outputTemplate: 'Modèle de message',
  thresholdType: 'Type de seuil',
  thresholdRate: 'Taux de déclenchement',
  thresholdCount: 'Seuil de déclenchement',
  timeInterval: 'Intervalle de temps',
  childInactiveDelay: 'Délai d\'inactivité de l\'enfant',
  childInactiveDelayTooltip: 'L\'alarme correspondant à cette règle n\'est activée qu\'après le délai d\'inactivité',
  valuePath: 'Chemin de valeur | Chemins de valeur',
  autoResolve: 'Résolution automatique',
  idHelp: 'Si ce champ n\'est pas renseigné, un identifiant unique sera généré automatiquement à la création de la règle',
  corelId: 'Identifiant de corrélation',
  corelIdHelp: '<p>Les variables accessibles sont: <strong>.Alarm</strong> et <strong>.Entity</strong></p>'
    + '<i>Quelques exemples:</i> <span>"{{ .Alarm.Value.Connector }}", "{{ .Entity.Component }}"</span>',
  corelStatus: 'Statut de corrélation',
  corelStatusHelp: '<p>Les variables accessibles sont: <strong>.Alarm</strong> et <strong>.Entity</strong></p>'
    + '<i>Quelques exemples:</i> <span>"{{ .Alarm.Value.Connector }}", "{{ .Entity.Component }}"</span>',
  corelParent: 'Corrélation parent',
  corelChild: 'Corrélation enfant',
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
  errors: {
    noValuePaths: 'Vous devez ajouter au moins un chemin de valeur',
  },
};
