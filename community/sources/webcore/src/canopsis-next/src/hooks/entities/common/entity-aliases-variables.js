import { computed } from 'vue';

import { ALARM_PAYLOADS_VARIABLES, ENTITY_PAYLOADS_VARIABLES, INFOS_NAME_VARIABLE } from '@/constants';

import { useEntityInfoProperty } from '@/hooks/store/modules/entity-info-property';

/**
 * Hook for managing entity aliases variables used in templates and forms.
 * Provides computed properties for different types of entity alias variables.
 *
 * @returns {Object} Object containing reactive variables for entity aliases
 */
export const useEntityAliasesVariables = () => {
  const { entityInfoPropertiesWithAlias } = useEntityInfoProperty();

  /**
   * Base variables computed from entity info properties that have aliases.
   * Maps each property to a variable object with alias flag, name as value, and alias as text.
   *
   * @type {import('vue').ComputedRef<Array<{alias: boolean, value: string, text: string}>>}
   */
  const variables = computed(() => entityInfoPropertiesWithAlias.value.map(property => ({
    alias: true,
    value: property.name,
    text: property.alias,
  })));

  /**
   * Variables formatted for alarm payload contexts.
   * Uses ALARM_PAYLOADS_VARIABLES.entityInfosValue pattern with property names substituted.
   *
   * @type {import('vue').ComputedRef<Array<{alias: boolean, value: string, text: string}>>}
   */
  const alarmAliasesVariables = computed(() => variables.value.map(variable => ({
    ...variable,

    value: ALARM_PAYLOADS_VARIABLES.entityInfosValue.replace(INFOS_NAME_VARIABLE, variable.value),
  })));

  /**
   * Variables formatted for entity payload contexts.
   * Uses ENTITY_PAYLOADS_VARIABLES.infosValue pattern with property names substituted.
   *
   * @type {import('vue').ComputedRef<Array<{alias: boolean, value: string, text: string}>>}
   */
  const entityAliasesVariables = computed(() => variables.value.map(variable => ({
    ...variable,

    value: ENTITY_PAYLOADS_VARIABLES.infosValue.replace(INFOS_NAME_VARIABLE, variable.value),
  })));

  return {
    variables,
    alarmAliasesVariables,
    entityAliasesVariables,
  };
};
