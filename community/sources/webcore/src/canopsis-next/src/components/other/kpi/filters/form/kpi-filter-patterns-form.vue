<template>
  <c-patterns-field
    v-field="form"
    :entity-attributes="entityAttributes"
    :readonly="readonly"
    :pending="pending"
    required
    with-entity
  />
</template>

<script>
import { usePatternsFields, usePatternsFieldsFetching } from '@/hooks/store/modules/patterns-fields';

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
  setup() {
    const { fetchKpiFilterPatternFields } = usePatternsFields();

    const {
      pending,
      entityAttributes,
    } = usePatternsFieldsFetching(fetchKpiFilterPatternFields);

    return {
      pending,
      entityAttributes,
    };
  },
};
</script>
