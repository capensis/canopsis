import { ACTION_TYPES } from '@/constants';

export default {
  emitTrigger: 'Émettre un déclencheur',
  withAuth: 'Avez-vous besoin de champs d\'authentification ?',
  emptyResponse: 'Réponse vide',
  isRegexp: 'La valeur peut être une expression régulière',
  headerKey: 'Clé d\'en-tête',
  headerValue: 'Valeur d\'en-tête',
  key: 'Clé',
  skipVerify: 'Ne pas vérifier les certificats HTTPS',
  headers: 'En-têtes',
  declareTicket: 'Déclarer un ticket',
  workflow: 'Comportement si cette action ne correspond pas :',
  remainingAction: 'Continuer avec les actions restantes',
  addAction: 'Ajouter une action',
  emptyActions: 'Aucune action ajoutée pour le moment',
  output: 'Format d\'action de sortie',
  forwardAuthor: 'Transmettre l\'auteur à l\'étape suivante',
  urlHelp: '<p>Les variables accessibles sont : <strong>.Alarm</strong>, <strong>.Entity</strong> et <strong>.Children</strong></p>'
    + '<i>Quelques exemples :</i>'
    + '<pre>"https://exampleurl.com?resource={{ .Alarm.Value.Resource }}"</pre>'
    + '<pre>"https://exampleurl.com?entity_id={{ .Entity.ID }}"</pre>'
    + '<pre>"https://exampleurl.com?children_count={{ len .Children }}"</pre>'
    + '<pre>"https://exampleurl.com?children={{ range .Children }}{{ .ID }}{{ end }}"</pre>',
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
    priorityExist: 'La priorité du scénario actuel est déjà utilisée. Voulez-vous changer la priorité actuelle du scénario en {priority} ?',
  },
};
