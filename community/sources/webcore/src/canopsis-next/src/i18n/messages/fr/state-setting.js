import {
  ENTITIES_STATES_KEYS,
  STATE_SETTING_METHODS,
  STATE_SETTING_CONDITIONS_METHODS,
  STATE_SETTING_CONDITIONS,
} from '@/constants';

export default {
  targetEntityState: 'État de l\'entité cible',
  entitiesStates: 'États des entités impactant',
  computeMethod: 'Méthode de calcul d\'état',
  addImpactingEntityPattern: 'Modèle d\'entité impactant',
  conditionsError: 'Veuillez ajouter au moins une condition',
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
    },
    [STATE_SETTING_METHODS.dependenciesState]: {
      label: 'L\'état est défini par une part ou un nombre d\'entités impactantes d\'un État spécifique',
      tooltip: 'Les états d\'entité peuvent être remplacés par une règle personnalisée définie par le nombre ou la part d\'entités impactantes d\'états spécifiques.',
    },
  },
  calculationMethods: {
    [STATE_SETTING_CONDITIONS_METHODS.share]: 'Partager',
    [STATE_SETTING_CONDITIONS_METHODS.number]: 'Nombre',
  },
  conditions: {
    [STATE_SETTING_CONDITIONS.greater]: 'Plus grand que',
    [STATE_SETTING_CONDITIONS.less]: 'Moins que',
  },
  states: {
    [ENTITIES_STATES_KEYS.ok]: 'OK',
    [ENTITIES_STATES_KEYS.minor]: 'Minor',
    [ENTITIES_STATES_KEYS.major]: 'Major',
    [ENTITIES_STATES_KEYS.critical]: 'Critical',
  },
};
