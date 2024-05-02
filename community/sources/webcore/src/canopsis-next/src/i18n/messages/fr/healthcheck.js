import { HEALTHCHECK_ENGINES_NAMES, HEALTHCHECK_SERVICES_NAMES } from '@/constants';

export default {
  metricsUnavailable: 'Les métriques ne sont pas collectées',
  notRunning: '{name} n\'est pas disponible',
  queueOverflow: 'Débordement de file d\'attente',
  lackOfInstances: 'Nombre d\'instances insuffisant',
  diffInstancesConfig: 'Configuration des instances non valide',
  queueLength: 'Longueur de la file d\'attente {queueLength}/{maxQueueLength}',
  instancesCount: 'Instances {instances}/{minInstances}',
  activeInstances: 'Seules {instances} sont actives sur {minInstances}. Le nombre optimal d\'instances est de {optimalInstances}.',
  queueOverflowed: 'La file d\'attente est saturée : {queueLength} messages sur {maxQueueLength}.\nVeuillez vérifier les instances.',
  engineDown: '{name} est en panne, le système n\'est pas opérationnel.\nVeuillez vérifier le journal ou redémarrer le service.',
  engineDownOrSlow: '{name} est en panne ou répond trop lentement, le système n\'est pas opérationnel.\nVeuillez vérifier le journal ou redémarrer l\'instance.',
  timescaleDown: '{name} est en panne, les métriques et les KPI ne sont pas collectés.\nVeuillez vérifier le journal ou redémarrer l\'instance.',
  invalidEnginesOrder: 'Configuration des moteurs non valide',
  invalidInstancesConfiguration: 'Configuration des instances non valide : les instances du moteur consomment ou publient dans différentes files d\'attente.\nVeuillez vérifier les instances.',
  chainConfigurationInvalid: 'La configuration de la chaîne des moteurs n\'est pas valide.\nReportez-vous ci-dessous pour la séquence correcte des moteurs :',
  queueLimit: 'Limite de longueur de file d\'attente',
  defineQueueLimit: 'Définir la limite de longueur de file d\'attente des moteurs',
  notifyUsersQueueLimit: 'Les utilisateurs peuvent être avertis lorsque la limite de longueur de file d\'attente est dépassée',
  messagesLimit: 'Limite de traitement des messages',
  defineMessageLimit: 'Définir la limite de traitement des messages FIFO (par minute)',
  notifyUsersMessagesLimit: 'Les utilisateurs peuvent être avertis lorsque la limite de traitement des messages FIFO est dépassée',
  numberOfInstances: 'Nombre d\'instances',
  notifyUsersNumberOfInstances: 'Les utilisateurs peuvent être avertis lorsque le nombre d\'instances actives est inférieur à la valeur minimale. Le nombre optimal d\'instances est affiché lorsque l\'état du moteur n\'est pas disponible.',
  messagesHistory: 'Historique de traitement des messages FIFO',
  messagesLastHour: 'Traitement des messages FIFO pour la dernière heure',
  messagesPerHour: 'messages/heure',
  messagesPerMinute: 'messages/minute',
  unknown: 'Cet état du système n\'est pas disponible',
  systemStatusChipError: 'Le système n\'est pas opérationnel',
  systemStatusServerError: 'La configuration du système n\'est pas valide, veuillez contacter l\'administrateur',
  systemsOperational: 'Tous les systèmes sont opérationnels',
  validation: {
    max_value: 'Le champ doit être égal ou inférieur au nombre optimal d\'instances',
    min_value: 'Le champ doit être égal ou supérieur au nombre minimal d\'instances',
  },
  nodes: {
    [HEALTHCHECK_SERVICES_NAMES.mongo]: {
      name: 'MongoDB',
      edgeLabel: 'Vérification de l\'état',
    },

    [HEALTHCHECK_SERVICES_NAMES.rabbit]: {
      name: 'RabbitMQ',
      edgeLabel: 'Vérification de l\'état',
    },

    [HEALTHCHECK_SERVICES_NAMES.redis]: {
      name: 'Redis',
      edgeLabel: 'Données FIFO\nÉtat de Redis',
    },

    [HEALTHCHECK_SERVICES_NAMES.events]: {
      name: 'Événements',
    },

    [HEALTHCHECK_SERVICES_NAMES.api]: {
      name: 'Canopsis API',
    },

    [HEALTHCHECK_SERVICES_NAMES.enginesChain]: {
      name: 'Chaîne des moteurs',
    },

    [HEALTHCHECK_SERVICES_NAMES.healthcheck]: {
      name: 'Healthcheck',
    },

    [HEALTHCHECK_ENGINES_NAMES.snmp]: {
      name: 'SNMP',
    },

    [HEALTHCHECK_ENGINES_NAMES.webhook]: {
      name: 'Webhook',
      description: 'Gère les webhooks',
    },

    [HEALTHCHECK_ENGINES_NAMES.fifo]: {
      name: 'FIFO',
      edgeLabel: 'État de RabbitMQ\nFlux entrant des KPIs',
      description: 'Gère la file d\'attente des événements et des alarmes',
    },

    [HEALTHCHECK_ENGINES_NAMES.axe]: {
      name: 'AXE',
      description: 'Crée des alarmes et exécute des actions sur elles',
    },

    [HEALTHCHECK_ENGINES_NAMES.che]: {
      name: 'CHE',
      description: 'Applique les filtres d\'événements et crée des entités',
    },

    [HEALTHCHECK_ENGINES_NAMES.pbehavior]: {
      name: 'Pbehavior',
      description: 'Vérifie si l\'alarme est sous PBehavior',
    },

    [HEALTHCHECK_ENGINES_NAMES.action]: {
      name: 'Action',
      description: 'Déclenche le lancement des actions',
    },

    [HEALTHCHECK_ENGINES_NAMES.dynamicInfos]: {
      name: 'Infos dynamiques',
      description: 'Ajoute des informations dynamiques à l\'alarme',
    },

    [HEALTHCHECK_ENGINES_NAMES.correlation]: {
      name: 'Corrélation',
      description: 'Gère la corrélation',
    },

    [HEALTHCHECK_ENGINES_NAMES.remediation]: {
      name: 'Remédiation',
      description: 'Déclenche les consignes',
    },
  },
};
