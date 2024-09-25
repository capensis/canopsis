import { ACTION_TYPES } from '@/constants';

export default {
  withAuth: 'Avez-vous besoin de champs d\'authentification ?',
  key: 'Clé',
  declareTicket: 'Déclarer un ticket',
  workflow: 'Comportement si cette action ne correspond pas :',
  remainingAction: 'Continuer avec les actions restantes',
  addAction: 'Ajouter une action',
  emptyActions: 'Aucune action ajoutée pour le moment',
  output: 'Format d\'action de sortie',
  forwardAuthor: 'Transmettre l\'auteur à l\'étape suivante',
  skipForChild: 'Sauter pour les enfants de la méta-alarme',
  skipForInstruction: 'Ignorer si l\'événement a déclenché une instruction automatique',
  outputHelp: '<p>Les variables accessibles sont : <strong>.Alarm</strong> et <strong>.Entity</strong></p>'
    + '<i>Quelques exemples:</i>'
    + '<pre>Resource - {{ .Alarm.Value.Resource }}. Entity - {{ .Entity.ID }}.</pre>',
  payloadHelp: '<p>Les variables accessibles sont : <strong>.Alarm</strong>, <strong>.Entity</strong> et <strong>.Children</strong></p>'
    + '<i>Quelques exemples:</i>'
    + '<pre>{\n'
    + '  resource: "{{ .Alarm.Value.Resource }}",\n'
    + '  entity: "{{ .Entity.ID }}",\n'
    + '  children_count: "{{ len .Children }}",\n'
    + '  children: {{ range .Children }}{{ .ID }}{{ end }}\n'
    + '}</pre>',
  actions: {
    [ACTION_TYPES.snooze]: 'Mettre en veille',
    [ACTION_TYPES.pbehavior]: 'Définir un comportement périodique',
    [ACTION_TYPES.changeState]: 'Changer l\'état (Change et vérouille la criticité)',
    [ACTION_TYPES.ack]: 'Acquitter',
    [ACTION_TYPES.ackremove]: 'Supprimer l\'acquittement',
    [ACTION_TYPES.assocticket]: 'Associer un ticket',
    [ACTION_TYPES.cancel]: 'Annuler',
    [ACTION_TYPES.webhook]: 'Webhook',
  },
  tabs: {
    pattern: 'Modèle',
  },
  errors: {
    actionRequired: 'Veuillez ajouter au moins une action',
    deprecatedTriggerExist: 'Ce scénario n\'est plus pris en charge en raison de son ancien format et donc désactivé. \n'
      + 'Veuillez mettre à jour les déclencheurs de scénario ou créer une nouvelle règle de déclaration de ticket.',
    testQueryRequireSteps: 'La requête de test n\'est pas disponible : aucun webhook n\'a été ajouté au scénario',
  },
  tooltips: {
    pbehaviorActionsNamePrefix: 'Le nom va être `{{prefix}} {{entity_id}} {{start}}-{{stop}}`',
  },
};
