<template>
  <c-patterns-field
    v-field="patterns"
    :entity-attributes="entityAttributes"
    :pending="pending"
    some-required
    with-entity
    entity-counters-type
  />
</template>

<script>
import { usePatternsFields, usePatternsFieldsFetching } from '@/hooks/store/modules/patterns-fields';

export default {
  inject: ['$validator'],
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
  setup() {
    const { fetchServicePatternFields } = usePatternsFields();

    const {
      pending,
      entityAttributes,
    } = usePatternsFieldsFetching(fetchServicePatternFields);

    return {
      pending,
      entityAttributes,
    };
  },
};
</script>
