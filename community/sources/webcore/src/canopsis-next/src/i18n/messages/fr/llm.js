import { LLM_AI_CHAT_SUGGESTION_TYPES, LLM_THINKING_LEVELS } from '@/constants';

export default {
  expandTabs: {
    details: 'Détails',
    promptsHistory: 'Historique des invites',
  },
  promptsHistory: {
    tabs: {
      allPrompts: 'Toutes les invites',
      byUser: 'Par utilisateur',
    },
    searchPlaceholder: 'Rechercher par utilisateur, modal ou nom',
    notRelatedToCanopsis: 'Sans lien avec Canopsis',
    groupByChat: 'Regrouper par conversation',
    seeChat: 'Voir la conversation',
    columns: {
      userName: 'Utilisateur',
      datetime: 'Date et heure',
      tokensUsed: 'Jetons utilisés',
      modal: 'Modal',
      name: 'Nom',
      usage: 'Usage',
      canopsisRelated: 'Lié à Canopsis',
      prompt: 'Invite',
      seeChat: 'Voir la conversation',
      promptsCount: 'Invites',
      lastUsed: 'Dernière utilisation',
    },
  },
  modelType: 'Type de modèle LLM',
  modelTypes: {
    gemini: 'Gemini',
  },
  modelName: 'Nom du modèle LLM',
  model: 'Modèle LLM',
  recommendedBadge: 'recommandé',
  apiKey: 'Clé API',
  apiKeyPlaceholder: 'Nouvelle clé API',
  thinkingLevel: 'Niveau de réflexion',
  isDefaultModel: 'Modèle par défaut',
  currentDefaultModelLine: 'Modèle par défaut actuel : {name}',
  lastUsedDate: 'Dernière utilisation',
  expandNoDetails: 'Aucun détail supplémentaire pour cette configuration.',
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
    modelPlaceholder: 'Modèle',
    ask: 'Demander',
    expandMessage: 'Afficher le message en entier',
    collapseMessage: 'Afficher moins',
    thinking: 'Réflexion en cours…',
    tryLabel: 'Essayez :',
    suggestions: {
      createPattern: 'Créer un modèle',
      editPattern: 'Éditer le modèle',
      validatePattern: 'Valider le modèle',
    },
    patternsMessage: '{patterns} modèle | {patterns} modèles',
    patternsEditedMessage: 'Vous avez modifié {patterns} modèle | Vous avez modifié {patterns} modèles',
    patternCreatedMessage: 'Modèles créés',
    patternUpdatedMessage: 'Modèles mis à jour',
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
  },
};
