import {
  ENTITIES_STATES_KEYS,
  STATE_SETTING_METHODS,
  STATE_SETTING_THRESHOLDS_METHODS,
  STATE_SETTING_THRESHOLDS_CONDITIONS,
  JUNIT_STATE_SETTING_METHODS,
} from '@/constants';

export default {
  title: 'Paramétrage de l\'état',
  targetEntityState: 'État de l\'entité cible',
  entitiesStates: 'États des entités impactant',
  computeMethod: 'Méthode de calcul d\'état',
  addImpactingEntityPattern: 'Modèle d\'entité impactant',
  conditionsError: 'Veuillez ajouter au moins une condition',
  entityThresholdSummary: 'Un état d\'entité ciblé est {state} lorsque le {method} d\'entités impactantes de l\'état {impactingEntitiesState} est {condition} {value}.',
  appliedFor: 'Appliqué pour',
  appliedForEntityType: 'Appliqué pour le type d\'entité',
  steps: {
    basics: 'Les bases',
    rulePatterns: 'Définir les entités cibles',
    conditions: 'Ajouter des conditions',
  },
  methods: {
    [STATE_SETTING_METHODS.inherited]: {
      label: 'L\'état est hérité des entités impactantes',
      tooltip: 'L\'état est défini par la ou les ressources impactantes.\n'
          + 'Lorsque plusieurs ressources sont définies, le pire état d\'entre elles est retenu.',
      stepTitle: 'L\'état de l\'entité cible est hérité d\'une ou plusieurs ressources impactantes. Lorsque plusieurs ressources correspondent au modèle, le pire état est retenu.',
    },
    [STATE_SETTING_METHODS.dependencies]: {
      label: 'L\'état est défini par une part ou un nombre d\'entités impactantes d\'un État spécifique',
      tooltip: 'Les états d\'entité peuvent être remplacés par une règle personnalisée définie par le nombre ou la part d\'entités impactantes d\'états spécifiques.',
      stepTitle: 'Les états des entités cibles peuvent être remplacés par des conditions basées sur un nombre ou une part d\'entités impactantes d\'un état spécifique.',
    },
  },
  thresholdMethods: {
    [STATE_SETTING_THRESHOLDS_METHODS.share]: 'Partager',
    [STATE_SETTING_THRESHOLDS_METHODS.number]: 'Nombre',
  },
  thresholdConditions: {
    [STATE_SETTING_THRESHOLDS_CONDITIONS.greater]: 'Plus grand que',
    [STATE_SETTING_THRESHOLDS_CONDITIONS.less]: 'Moins que',
  },
  states: {
    [ENTITIES_STATES_KEYS.ok]: 'OK',
    [ENTITIES_STATES_KEYS.minor]: 'Mineur',
    [ENTITIES_STATES_KEYS.major]: 'Majeur',
    [ENTITIES_STATES_KEYS.critical]: 'Critique',
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
