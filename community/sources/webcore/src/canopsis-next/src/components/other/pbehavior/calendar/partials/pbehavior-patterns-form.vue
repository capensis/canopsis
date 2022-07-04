<template lang="pug">
  c-patterns-field(
    v-field="form",
    :entity-attributes="entityAttributes",
    :readonly="readonly",
    with-entity,
    required
  )
</template>

<script>
import { ENTITY_PATTERN_FIELDS } from '@/constants';

export default {
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      default: () => ({}),
    },
    readonly: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    entityAttributes() {
      return [
        {
          value: ENTITY_PATTERN_FIELDS.lastEventDate,
          options: { disabled: true },
        },
      ];
    },
  },
};
</script>
