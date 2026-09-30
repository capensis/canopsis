<template>
  <c-layout column>
    <c-description-field
      :model-value="modelValue.comment"
      :label="$t('common.comment', 1)"
      name="comment"
      autofocus
      required
      @update:model-value="updateField('comment', $event)"
    />
    <c-pbehavior-reason-field
      :model-value="modelValue.reason"
      required
      return-object
      @update:model-value="updateField('reason', $event)"
    />
  </c-layout>
</template>

<script>
import { useModelField } from '@/hooks/form/model-field';

export default {
  props: {
    modelValue: {
      type: Object,
      default: () => ({}),
    },
  },
  emits: ['update:modelValue'],

  setup(props, { emit }) {
    const { updateField } = useModelField(props, emit);

    return { updateField };
  },
};
</script>
