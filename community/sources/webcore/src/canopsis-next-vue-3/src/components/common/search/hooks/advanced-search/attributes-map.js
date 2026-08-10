import { computed } from 'vue';

import {
  ADVANCED_SEARCH_STRING_WITH_ONE_OF_OPERATORS,
  ALARM_PATTERN_FIELDS,
  ENTITY_PATTERN_FIELDS,
  PATTERN_DATE_OPERATORS,
  PATTERN_NUMBER_OPERATORS,
  PATTERN_OPERATORS,
  PBEHAVIOR_FIELDS,
  DYNAMIC_INFO_FIELDS,
} from '@/constants';

const stringOptions = {
  operators: ADVANCED_SEARCH_STRING_WITH_ONE_OF_OPERATORS,
};

const numberOptions = {
  operators: PATTERN_NUMBER_OPERATORS,
};

const dateOptions = {
  operators: PATTERN_DATE_OPERATORS,
};

const booleanOptions = {
  operators: [PATTERN_OPERATORS.equal, PATTERN_OPERATORS.notEqual],
  values: [
    { value: true, text: 'true' },
    { value: false, text: 'false' },
  ],
};

export const useAdvancedSearchAlarmAttributes = () => ({
  attributesMap: computed(() => ({
    [ALARM_PATTERN_FIELDS.displayName]: stringOptions,
    [ALARM_PATTERN_FIELDS.connector]: stringOptions,
    [ALARM_PATTERN_FIELDS.connectorName]: stringOptions,
    [ALARM_PATTERN_FIELDS.component]: stringOptions,
    [ALARM_PATTERN_FIELDS.resource]: stringOptions,
    [ALARM_PATTERN_FIELDS.state]: numberOptions,
    [ALARM_PATTERN_FIELDS.status]: numberOptions,
    [ALARM_PATTERN_FIELDS.tags]: stringOptions,
    [ALARM_PATTERN_FIELDS.creationDate]: dateOptions,
    [ALARM_PATTERN_FIELDS.lastUpdateDate]: dateOptions,
    [ALARM_PATTERN_FIELDS.lastEventDate]: dateOptions,
  })),
});

export const useAdvancedSearchEntityAttributes = () => ({
  attributesMap: computed(() => ({
    [ENTITY_PATTERN_FIELDS.id]: stringOptions,
    [ENTITY_PATTERN_FIELDS.name]: stringOptions,
    [ENTITY_PATTERN_FIELDS.type]: stringOptions,
    [ENTITY_PATTERN_FIELDS.category]: stringOptions,
    [ENTITY_PATTERN_FIELDS.component]: stringOptions,
    [ENTITY_PATTERN_FIELDS.connector]: stringOptions,
    [ENTITY_PATTERN_FIELDS.impactLevel]: numberOptions,
    [ENTITY_PATTERN_FIELDS.enabled]: booleanOptions,
    [ENTITY_PATTERN_FIELDS.idleSince]: dateOptions,
    [ENTITY_PATTERN_FIELDS.imported]: dateOptions,
  })),
});

export const useAdvancedSearchPbehaviorAttributes = () => ({
  attributesMap: computed(() => ({
    [PBEHAVIOR_FIELDS.name]: stringOptions,
    [PBEHAVIOR_FIELDS.author]: stringOptions,
    [PBEHAVIOR_FIELDS.enabled]: booleanOptions,
    [PBEHAVIOR_FIELDS.reason]: stringOptions,
    [PBEHAVIOR_FIELDS.type]: stringOptions,
    [PBEHAVIOR_FIELDS.tstart]: dateOptions,
    [PBEHAVIOR_FIELDS.tstop]: dateOptions,
  })),
});

export const useAdvancedSearchDynamicInfoAttributes = () => ({
  attributesMap: computed(() => ({
    [DYNAMIC_INFO_FIELDS.id]: stringOptions,
    [DYNAMIC_INFO_FIELDS.name]: stringOptions,
    [DYNAMIC_INFO_FIELDS.author]: stringOptions,
    [DYNAMIC_INFO_FIELDS.enabled]: booleanOptions,
    [DYNAMIC_INFO_FIELDS.created]: dateOptions,
    [DYNAMIC_INFO_FIELDS.updated]: dateOptions,
  })),
});
