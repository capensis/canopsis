<template>
  <c-patterns-field
    v-field="form"
    :entity-attributes="entityAttributes"
    :readonly="readonly"
    required
    with-entity
  />
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
      required: true,
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
