import { ALARM_PAYLOADS_VARIABLES, ENTITY_PAYLOADS_VARIABLES, INFOS_NAME_VARIABLE } from '@/constants';

import { entitiesEntityInfoPropertyMixin } from '../entities/entity-info-property';

/**
 * @TODO remove this mixin in the future. Use entity-aliases-variables hook instead
 */
export const payloadEntityAliasesMixin = {
  mixins: [entitiesEntityInfoPropertyMixin],

  computed: {
    aliasesVariables() {
      return this.entityInfoPropertiesWithAlias.map(property => ({
        alias: true,
        value: property.name,
        text: property.alias,
      }));
    },

    alarmAliasesVariables() {
      return this.aliasesVariables.map(variable => ({
        ...variable,

        value: ALARM_PAYLOADS_VARIABLES.entityInfosValue.replace(INFOS_NAME_VARIABLE, variable.value),
      }));
    },

    entityAliasesVariables() {
      return this.aliasesVariables.map(variable => ({
        ...variable,

        value: ENTITY_PAYLOADS_VARIABLES.infosValue.replace(INFOS_NAME_VARIABLE, variable.value),
      }));
    },
  },
};
