<template>
  <c-patterns-field
    v-field="patterns"
    :entity-attributes="entityAttributes"
    :entity-types="entityTypes"
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
    entityTypes: {
      type: Array,
      default: () => [...STATE_SETTING_ENTITY_TYPES],
    },
  },
  computed: {
    entityAttributes() {
      return [
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
