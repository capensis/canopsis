<template>
  <c-patterns-field
    v-field="patterns"
    :entity-attributes="entityAttributes"
    :entity-types="entityTypes"
    :entity-title="$t('stateSetting.addImpactingEntityPattern')"
    name="inherited_entity_pattern"
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

    entityAttributes() {
      return [
        {
          value: ENTITY_PATTERN_FIELDS.connector,
          options: { disabled: true },
        },
        {
          value: ENTITY_PATTERN_FIELDS.component,
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
