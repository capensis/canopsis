<template>
  <c-patterns-field
    :value="form"
    :entity-attributes="entityAttributes"
    :readonly="readonly"
    :pending="pending"
    required
    with-entity
    @update:value="$emit('update:form', $event)"
  />
</template>

<script>
import {
  usePatternsFields,
  usePatternsFieldsFetching,
} from '@/hooks/store/modules/patterns-fields';

export default {
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

  emits: ['update:form'],

  setup(props) {
    const { fetchKpiFilterPatternFields } = usePatternsFields();

    const { pending, entityAttributes } = usePatternsFieldsFetching(
      fetchKpiFilterPatternFields,
      props.readonly,
    );

    return {
      pending,
      entityAttributes,
    };
  },
};
</script>
