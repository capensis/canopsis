<template>
  <c-patterns-field
    :value="patterns"
    :entity-attributes="entityAttributes"
    :pending="pending"
    some-required
    with-entity
    entity-counters-type
    @update:value="updateModel"
  />
</template>

<script>
import { useModelField } from '@/hooks/form/model-field';
import { usePatternsFields, usePatternsFieldsFetching } from '@/hooks/store/modules/patterns-fields';

export default {
  props: {
    patterns: {
      type: Object,
      default: () => ({}),
    },
  },

  emits: ['update:patterns'],

  setup(props, { emit }) {
    const { updateModel } = useModelField(props, emit);
    const { fetchServicePatternFields } = usePatternsFields();

    const {
      pending,
      entityAttributes,
    } = usePatternsFieldsFetching(fetchServicePatternFields);

    return {
      pending,
      entityAttributes,
      updateModel,
    };
  },
};
</script>
