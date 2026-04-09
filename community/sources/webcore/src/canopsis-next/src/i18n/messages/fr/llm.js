import {
  LLM_AI_CHAT_ERROR_CODES,
  LLM_AI_CHAT_SUGGESTION_TYPES,
  LLM_SOCKET_CONTEXTS,
  LLM_THINKING_LEVELS,
} from '@/constants';

export default {
  expandTabs: {
    promptsHistory: 'Historique des invites',
  },
  promptsHistory: {
    tabs: {
      allPrompts: 'Toutes les invites',
    },
    expandUserHistory: {
      lastPromptDate: 'Date de la dernière invite',
    },
    filterByUser: 'Filtrer par utilisateur',
    groupByChat: 'Regrouper par conversation',
    notRelatedToCanopsis: 'Sans lien avec Canopsis',
    allUsers: 'Tous les utilisateurs',
    search: 'Rechercher par utilisateur, modal ou nom',
    searchByUserName: "Rechercher par nom d'utilisateur",
    searchByModalOrName: 'Rechercher par modal ou nom',
    userHistoryTitle: 'Historique utilisateur',
    userHistoryEmpty: 'Aucun utilisateur pour cette configuration.',
    ruleNotSaved: 'Non enregistré',
    columns: {
      userName: 'Utilisateur',
      datetime: 'Date et heure',
      tokensUsed: 'Jetons utilisés',
      modal: 'Modal',
      usage: 'Utilisation',
      canopsisRelated: 'Lié à Canopsis',
      prompt: 'Invite',
      lastUsed: 'Dernière utilisation',
    },
    contextTitles: {
      [LLM_SOCKET_CONTEXTS.idleRule]: 'Règle d\'inactivité',
      [LLM_SOCKET_CONTEXTS.scenario]: 'Scénario',
      [LLM_SOCKET_CONTEXTS.flappingRule]: 'Règle de flapping',
      [LLM_SOCKET_CONTEXTS.resolveRule]: 'Règle de résolution',
      [LLM_SOCKET_CONTEXTS.alarmTag]: 'Étiquette d\'alarme',
      [LLM_SOCKET_CONTEXTS.linkRule]: 'Générateur de liens',
      [LLM_SOCKET_CONTEXTS.instruction]: 'Instruction de remédiation',
      [LLM_SOCKET_CONTEXTS.dynamicInfos]: 'Information dynamique',
      [LLM_SOCKET_CONTEXTS.metaAlarmRule]: 'Règle d\'alarme méta',
      [LLM_SOCKET_CONTEXTS.declareTicketRule]: 'Règle de déclaration de ticket',
      [LLM_SOCKET_CONTEXTS.pbehavior]: 'Comportement périodique',
      [LLM_SOCKET_CONTEXTS.entityService]: 'Service',
      [LLM_SOCKET_CONTEXTS.stateSettings]: 'Méthode de calcul d\'état',
      [LLM_SOCKET_CONTEXTS.entity]: 'Entité',
      [LLM_SOCKET_CONTEXTS.kpiFilter]: 'Filtre KPI',
      [LLM_SOCKET_CONTEXTS.eventFilter]: 'Règle de filtre d\'événements',
      [LLM_SOCKET_CONTEXTS.eventRecord]: 'Enregistrement d\'événements',
      [LLM_SOCKET_CONTEXTS.serviceWeather]: 'Météo du service',
      [LLM_SOCKET_CONTEXTS.widgetFilter]: 'Filtre de widget',
      [LLM_SOCKET_CONTEXTS.corporateAlarmPattern]: 'Filtre partagé d\'alarme',
      [LLM_SOCKET_CONTEXTS.corporateEntityPattern]: 'Filtre partagé d\'entité',
      [LLM_SOCKET_CONTEXTS.corporatePbehaviorPattern]: 'Filtre partagé de comportement périodique',
      [LLM_SOCKET_CONTEXTS.corporateWeatherServicePattern]: 'Modèle partagé de météo des services',
    },
  },
  modelType: 'Type de modèle LLM',
  modelName: 'Nom du modèle LLM',
  model: 'Modèle LLM',
  recommendedBadge: 'recommandé',
  apiKey: 'Clé API',
  apiKeyPlaceholder: 'Nouvelle clé API',
  thinkingLevel: 'Niveau de réflexion',
  isDefaultModel: 'Modèle par défaut',
  currentDefaultModelLine: 'Modèle par défaut actuel : {name}',
  lastUsedDate: 'Dernière utilisation',
  importantNotesMessage:
    '<div>Notes importantes :</div>'
    + '<ul>'
    + '<li>Les modèles Gemini peuvent être activés ou désactivés via ce module ; toutefois, '
    + '<strong>la gestion des limites de jetons et des quotas</strong>'
    + ' se configure directement dans la '
    + '<strong><a href="{geminiConsoleUrl}" target="_blank">Console Gemini AI</a></strong>.</li>'
    + '<li><strong>Une surveillance régulière</strong> des limites Google Cloud est requise, car les <strong>politiques d\'API</strong> et les conditions d\'utilisation sont susceptibles d\'évoluer.</li>'
    + '</ul>',
  thinkingLevels: {
    [LLM_THINKING_LEVELS.minimal]: 'Minimal',
    [LLM_THINKING_LEVELS.low]: 'Faible',
    [LLM_THINKING_LEVELS.medium]: 'Moyen',
    [LLM_THINKING_LEVELS.high]: 'Élevé',
  },
  chat: {
    howCanIHelp: 'Comment puis-je vous aider ?',
    promptPlaceholder: 'Décrivez votre besoin…',
    ask: 'Demander',
    thinking: 'Réflexion en cours…',
    tryLabel: 'Essayez :',
    restartConfirmation: 'Redémarrer la conversation ? <strong>L\'historique</strong> et les <strong>versions</strong> seront effacés.',
    suggestions: {
      createPattern: 'Créer un modèle',
      editPattern: 'Éditer le modèle',
      validatePattern: 'Valider le modèle',
    },
    patternsMessage: '{patterns} modèle | {patterns} modèles',
    patternsEditedMessage: 'Vous avez modifié {patterns} modèle | Vous avez modifié {patterns} modèles',
    patternCreatedMessage: 'Modèles créés',
    patternUpdatedMessage: 'Modèles mis à jour',
    emptyPatternsMessage: 'Vous avez supprimé tous les modèles.\nSi vous devez effacer le contexte, redémarrez l\'assistant IA.',
    patternCannotBeCreatedReasons: 'Le modèle ne peut pas être créé pour les raisons suivantes :\n',
    fixPatternPrompt: 'Corriger le modèle :\n{jsonString}',
    suggestionPrompts: {
      [LLM_AI_CHAT_SUGGESTION_TYPES.createPattern]: 'Créer un modèle avec les conditions suivantes :\n',
      [LLM_AI_CHAT_SUGGESTION_TYPES.editPattern]: 'Les modifications suivantes sont nécessaires dans le modèle :\n',
      [LLM_AI_CHAT_SUGGESTION_TYPES.validatePattern]: 'Corrigez le modèle suivant si nécessaire :\n',
    },
    pattern: {
      restoreVersion: 'Restaurer la version',
      seePattern: 'Voir le modèle',
      hidePattern: 'Masquer le modèle',
      version: 'Version {version}',
      versionRestored: 'Version {version} restaurée',
    },
    warningAlert: '<div>Veuillez noter que :</div>'
      + '<ul>'
      + '<li>Les résultats générés par l\'IA <strong>peuvent contenir des erreurs</strong> — vérifiez et validez toujours les modèles avant utilisation.</li>'
      + '<li>Cette conversation est <strong>temporaire</strong> et sera <strong>effacée</strong> à la fermeture de la fenêtre modale.</li>'
      + '</ul>',
    errors: {
      [LLM_AI_CHAT_ERROR_CODES.gone]:
        'Le modèle d\'IA sélectionné est actuellement indisponible.<br />'
        + '<strong>Veuillez contacter votre administrateur. Vous pouvez redémarrer la conversation avec un autre modèle (l\'historique sera perdu).</strong>',
      [LLM_AI_CHAT_ERROR_CODES.tooManyRequests]:
        'Vous avez dépassé votre quota actuel.<br />'
        + '<strong>Veuillez réessayer à {retryAt} ou redémarrer la conversation avec un autre modèle (l\'historique sera perdu).</strong>',
      [LLM_AI_CHAT_ERROR_CODES.internalError]:
        'Une erreur s\'est produite.<br />'
        + '<strong>Veuillez contacter votre administrateur</strong>',
      [LLM_AI_CHAT_ERROR_CODES.timeout]:
        'La requête a expiré.<br />'
        + '<strong>Veuillez réessayer plus tard.</strong>',
      noModels:
        'Aucun modèle d\'IA n\'est actuellement disponible.<br />'
        + '<strong>Veuillez contacter votre administrateur.</strong>',
      serverError:
        'Problème de connexion au socket.<br />'
        + '<strong>Veuillez contacter votre administrateur.</strong>',
    },
  },
};
