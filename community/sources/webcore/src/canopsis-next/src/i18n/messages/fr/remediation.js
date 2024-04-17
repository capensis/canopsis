import {
  REMEDIATION_INSTRUCTION_EXECUTION_STATUSES,
  REMEDIATION_INSTRUCTION_EXECUTION_STEP_STATUSES,
  REMEDIATION_INSTRUCTION_EXECUTION_STEP_TYPES,
  REMEDIATION_INSTRUCTION_TYPES,
} from '@/constants';

export default {
  tabs: {
    configurations: 'Configurations',
    jobs: 'Tâches',
    statistics: 'Statistiques de remédiation',
  },

  instruction: {
    name: 'Nom de l\'instruction',
    usingInstruction: 'Ne peut pas être supprimée, car en cours d\'utilisation',
    addStep: 'Ajouter une étape',
    addOperation: 'Ajouter une opération',
    addJob: 'Ajouter une tâche',
    endpoint: 'Point de terminaison',
    job: 'Tâche | Tâches',
    listJobs: 'Liste des tâches',
    endpointAvatar: 'EP',
    workflow: 'Si cette étape échoue :',
    jobWorkflow: 'Comportement si cette tâche échoue:',
    remainingStep: 'Continuer avec les étapes restantes',
    remainingJob: 'Continuer avec la tâche restante',
    timeToComplete: 'Temps d\'exécution (estimation)',
    emptySteps: 'Aucune étape ajoutée pour le moment',
    emptyOperations: 'Aucune opération ajoutée pour le moment',
    emptyJobs: 'Aucune tâche ajoutée pour le moment',
    timeoutAfterExecution: 'Délai d\'attente après l\'exécution de la consigne',
    requestApproval: 'Demande d\'approbation',
    type: 'Type de consigne',
    approvalPending: 'En attente d\'approbation',
    approvalDismissed: 'L\'instruction est rejetée',
    needApprove: 'Une approbation est nécessaire',
    executeInstruction: 'Exécuter la consigne "{instructionName}"',
    resumeInstruction: 'Reprendre la consigne "{instructionName}"',
    inProgressInstruction: '{instructionName} en cours...',
    types: {
      [REMEDIATION_INSTRUCTION_TYPES.simpleManual]: 'Manuel simplifié',
      [REMEDIATION_INSTRUCTION_TYPES.manual]: 'Manuel',
      [REMEDIATION_INSTRUCTION_TYPES.auto]: 'Automatique',
    },
    tooltips: {
      endpoint: 'Le point de terminaison doit être une question qui appelle une réponse Oui/Non',
    },
    table: {
      rating: 'Évaluation',
      monthExecutions: '№ d\'exécutions\nce mois-ci',
      lastExecutedOn: 'Dernière exécution le',
    },
    errors: {
      operationRequired: 'Veuillez ajouter au moins une opération',
      stepRequired: 'Veuillez ajouter au moins une étape',
      jobRequired: 'Veuillez ajouter au moins une tâche',
    },
  },

  configuration: {
    usingConfiguration: 'Ne peut pas être supprimée, car en cours d\'utilisation',
    host: 'Hôte',
  },

  instructionExecute: {
    timeToComplete: '{duration} pour terminer',
    completedAt: 'Terminé à {time}',
    failedAt: 'Échec à {time}',
    startedAt: 'Commencé à {time}\n(Date de lancement Canopsis)',
    closeConfirmationText: 'Souhaitez-vous reprendre cette consigne plus tard ?',
    queueNumber: '{number} {name} travaux sont dans la file d\'attente',
    runJobs: 'Exécuter des tâches',
    popups: {
      success: '{instructionName} a été exécutée avec succès',
      failed: '{instructionName} a échoué. Veuillez faire remonter ce problème',
      connectionError: 'Il y a un problème de connexion. Veuillez cliquer sur le bouton d\'actualisation ou recharger la page.',
      wasAborted: '{instructionName} a été abandonnée',
      wasPaused: 'La consigne {instructionName} sur l\'alarme {alarmName} a été interrompue à {date}. Vous pouvez la reprendre manuellement.',
      wasRemovedOrDisabled: 'La consigne {instructionName} a été supprimée ou désactivée.',
    },
    jobs: {
      title: 'Tâches attribuées :',
      startedAt: 'Date de déclenchement\n(par Canopsis)',
      launchedAt: 'Date de lancement\n(par l\'ordonnanceur)',
      launchedBy: 'Lancé par',
      completedAt: 'Fin de traitement\n(par l\'ordonnanceur)',
      waitAlert: 'L\'ordonnanceur ne répond pas, veuillez contacter votre administrateur',
      skip: 'Ignorer la tâche',
      await: 'Attendre',
      failedReason: 'Raison de l\'échec',
      output: 'Retour',
      instructionFailed: 'Échec de d\'une consigne',
      instructionComplete: 'Exécution des consignes terminée',
      stopped: 'Arrêté',
    },
    status: {
      tooltips: {
        [REMEDIATION_INSTRUCTION_EXECUTION_STATUSES.running]: 'En cours',
        [REMEDIATION_INSTRUCTION_EXECUTION_STATUSES.paused]: 'En pause',
        [REMEDIATION_INSTRUCTION_EXECUTION_STATUSES.completed]: 'Terminer avec succès',
        [REMEDIATION_INSTRUCTION_EXECUTION_STATUSES.aborted]: 'Annulé',
        [REMEDIATION_INSTRUCTION_EXECUTION_STATUSES.failed]: 'Échoué',
        [REMEDIATION_INSTRUCTION_EXECUTION_STATUSES.waitingResult]: 'En cours',
      },
    },
    stepsTitles: {
      [REMEDIATION_INSTRUCTION_EXECUTION_STEP_TYPES.manual]: {
        [REMEDIATION_INSTRUCTION_EXECUTION_STEP_STATUSES.completed]: 'L\'étape {name} est terminée',
        [REMEDIATION_INSTRUCTION_EXECUTION_STEP_STATUSES.failed]: 'L\'étape {name} a échoué',
        [REMEDIATION_INSTRUCTION_EXECUTION_STEP_STATUSES.aborted]: 'L\'étape {name} est abandonnée',
        [REMEDIATION_INSTRUCTION_EXECUTION_STEP_STATUSES.skipped]: 'L\'étape {name} est ignorée',
      },
      [REMEDIATION_INSTRUCTION_EXECUTION_STEP_TYPES.job]: {
        [REMEDIATION_INSTRUCTION_EXECUTION_STEP_STATUSES.completed]: '{name} est terminée',
        [REMEDIATION_INSTRUCTION_EXECUTION_STEP_STATUSES.failed]: '{name} a échoué',
        [REMEDIATION_INSTRUCTION_EXECUTION_STEP_STATUSES.aborted]: '{name} est abandonnée',
        [REMEDIATION_INSTRUCTION_EXECUTION_STEP_STATUSES.skipped]: '{name} est ignorée',
      },
    },
  },

  instructionsFilter: {
    button: 'Créer un filtre sur les consignes de remédiation',
    filterByInstructions: 'Pour les alarmes par instructions',
    with: 'Afficher les alarmes avec les instructions sélectionnées',
    without: 'Afficher les alarmes sans instructions sélectionnées',
    selectAll: 'Tout sélectionner',
    alarmsListDisplay: 'Affichage de la liste des alarmes',
    allAlarms: 'Afficher toutes les alarmes filtrées',
    showWithInProgress: 'Afficher les alarmes filtrées avec les instructions en cours',
    showWithoutInProgress: 'Afficher les alarmes filtrées sans instructions en cours',
    hideWithInProgress: 'Masquer les alarmes filtrées avec les instructions en cours',
    hideWithoutInProgress: 'Masquer les alarmes filtrées sans instructions en cours',
    selectedInstructions: 'Consignes sélectionnées',
    selectedInstructionsHelp: 'Les consignes du type sélectionné sont exclues de la liste',
    inProgress: 'En cours',
    chip: {
      with: 'AVEC',
      without: 'SANS',
      all: 'TOUT',
    },
  },

  instructionStat: {
    alarmsTimeline: 'Chronologie des alarmes',
    executedAt: 'Fin execution à',
    lastExecutedOn: 'Dernière exécution le',
    modifiedOn: 'Dernière modification le',
    averageCompletionTime: 'Temps moyen\nd\'achèvement',
    executionCount: 'Nombre\nd\'exécutions',
    totalExecutions: 'Total des exécutions',
    successfulExecutions: 'Exécutions réussies',
    alarmStates: 'Alarmes affectées par état',
    okAlarmStates: 'Nombre de résultats\nÉtats OK',
    instructionChanged: 'La consigne a été modifiée',
    alarmResolvedDate: 'Date de résolution de l\'alarme',
    showFailedExecutions: 'Afficher les exécutions d\'instructions ayant échoué',
    remediationDuration: 'Durée de la remédiation',
    timeoutAfterExecution: 'Timeout après exécution',
    actions: {
      needRate: 'Evaluez-les!',
      rate: 'Évaluer',
    },
  },

  pattern: {
    tabs: {
      pbehaviorTypes: {
        title: 'Types de comportements périodiques',
        fields: {
          activeOnTypes: 'Actif sur les types',
          disabledOnTypes: 'Désactivé sur les types',
        },
      },
    },
  },

  job: {
    configuration: 'Configuration',
    jobId: 'Identifiant de la tâche',
    addJobs: 'Ajouter {count} tâche | Ajouter {count} tâches',
    usingJob: 'La tâche ne peut être supprimée, car elle est en cours d\'utilisation',
    query: 'Requête',
    multipleExecutions: 'Autoriser l\'exécution parallèle',
    jobWaitInterval: 'Intervalle d\'attente des tâches',
    addPayload: 'Ajouter un payload',
    deletePayload: 'Supprimer le payload',
    payloadHelp: '<p>Les variables accessibles sont: <strong>.Alarm</strong> et <strong>.Entity</strong></p>'
      + '<i>Quelques exemples:</i>'
      + '<pre>{\n  resource: "{{ .Alarm.Value.Resource }}",\n  entity: "{{ .Entity.ID }}"\n}</pre>',
    errors: {
      invalidJSON: 'JSON non valide',
    },
  },

  statistic: {
    remediation: 'Remédiation',
    allInstructions: 'Toutes les instructions',
    manualInstructions: 'Instructions manuelles',
    autoInstructions: 'Instructions automatique',
    labels: {
      remediated: 'Remédié',
      withAssignedRemediations: 'Remédiables (Avec instructions assignées)',
    },
    tooltips: {
      remediated: '{value} alarmes remédiées',
      assigned: '{value} alarmes avec consigne',
    },
  },
};
