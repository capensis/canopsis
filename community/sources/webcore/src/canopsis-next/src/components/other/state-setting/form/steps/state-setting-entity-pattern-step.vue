<template>
  <c-patterns-field
    v-field="patterns"
    :entity-attributes="entityAttributes"
    :entity-types="entityTypes"
    name="rule_patterns"
    required
    with-entity
    entity-counters-type
  />
</template>

<script>
import { ENTITY_PATTERN_FIELDS, STATE_SETTING_ENTITY_TYPES } from '@/constants';

import { formValidationHeaderMixin } from '@/mixins/form';

export default {
  inject: ['$validator'],
  mixins: [formValidationHeaderMixin],
  model: {
    prop: 'patterns',
    event: 'input',
  },
  props: {
    patterns: {
      type: Object,
      default: () => ({}),
    },
  },
  computed: {
    entityTypes() {
      return [...STATE_SETTING_ENTITY_TYPES];
    },

    typeOptions() {
      return {
        valueField: {
          props: {
            types: this.entityTypes,
          },
        },
      };
    },

    entityAttributes() {
      return [
        {
          value: ENTITY_PATTERN_FIELDS.type,
          options: this.typeOptions,
        },
        {
          value: ENTITY_PATTERN_FIELDS.component,
          options: { disabled: true },
        },
        {
          value: ENTITY_PATTERN_FIELDS.connector,
          options: { disabled: true },
        },
        {
          value: ENTITY_PATTERN_FIELDS.type,
          options: { disabled: true },
        },
        {
          value: ENTITY_PATTERN_FIELDS.componentInfos,
          options: { disabled: true },
        },
        {
          value: ENTITY_PATTERN_FIELDS.lastEventDate,
          options: { disabled: true },
        },
      ];
    },
  },
};
</script>
