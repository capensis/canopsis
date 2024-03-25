import {
  STATE_SETTING_METHODS,
  STATE_SETTING_THRESHOLDS_METHODS,
  STATE_SETTING_THRESHOLDS_CONDITIONS,
  JUNIT_STATE_SETTING_METHODS,
} from '@/constants';

export default {
  title: 'Paramétrage de l\'état',
  targetEntityState: 'État de l\'entité cible',
  entitiesStates: 'État des dépendances',
  computeMethod: 'Méthode de calcul d\'état',
  dependenciesEntityPattern: 'Modèle pour les dépendances',
  conditionsError: 'Veuillez ajouter au moins une condition',
  targetEntityThresholdSummary: 'Un état d\'entité ciblé est {state} lorsque la {method} des dépendances de l\'état {dependenciesEntitiesState} est {condition} {value}.',
  entityThresholdSummary: 'L\'état {name} est {state} lorsque la {method} des dépendances de l\'état {dependenciesEntitiesState} est {condition} {value}',
  appliedFor: 'Appliqué pour',
  appliedForEntityType: 'Appliqué pour le type d\'entité',
  stateIsInheritFrom: 'L\'état de {name} est hérité de',
  seeFilterPattern: 'Voir le modèle de filtre',
  dependsCount: 'Nombre total de dépendances',
  stateDependsCount: 'Nombre de dépendances de l\'état {state}',
  steps: {
    basics: 'Les bases',
    rulePatterns: 'Définir les entités cibles',
    conditions: 'Ajouter des conditions',
  },
  methods: {
    [STATE_SETTING_METHODS.inherited]: {
      label: 'L\'État est hérité des dépendances',
      tooltip: 'L\'État est défini par une ou plusieurs dépendances.\n'
          + 'Lorsque plusieurs dépendances sont définies, le pire état d’entre elles est retenu.',
      stepTitle: 'L\'état de l\'entité cible est hérité d\'une ou plusieurs dépendances. Lorsque plusieurs dépendances correspondent au modèle, le pire état est retenu.',
    },
    [STATE_SETTING_METHODS.dependencies]: {
      label: 'L\'État est défini par un calcul (pourcentage ou nombre) appliqué sur les états des dépendances.',
      tooltip: 'Il est possible de définir chaque état cible à partir d\'un calcul sur les états des dépendances.',
      stepTitle: 'Les états de l\'entité cible peuvent être remplacés par des conditions basées sur les états des dépendances en nombre ou en pourcentage.',
    },
  },
  thresholdMethods: {
    [STATE_SETTING_THRESHOLDS_METHODS.share]: 'Pourcentage',
    [STATE_SETTING_THRESHOLDS_METHODS.number]: 'Nombre',
  },
  thresholdConditions: {
    [STATE_SETTING_THRESHOLDS_CONDITIONS.greater]: 'Plus grand que',
    [STATE_SETTING_THRESHOLDS_CONDITIONS.less]: 'Moins que',
  },
  junit: {
    worstLabel: 'Le pire critère :',
    worstHelpText: 'Canopsis compte l\'état pour chaque critère défini. L\'état final de la suite de tests JUnit est considéré comme le pire des états résultants.',
    criterion: 'Critère',
    methods: {
      [JUNIT_STATE_SETTING_METHODS.worst]: 'Pire',
      [JUNIT_STATE_SETTING_METHODS.worstOfShare]: 'Pire des états',
    },
  },
};
