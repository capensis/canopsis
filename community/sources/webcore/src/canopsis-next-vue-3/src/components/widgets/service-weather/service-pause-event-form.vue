<template>
  <c-layout column>
    <c-description-field
      :model-value="modelValue.comment"
      :label="$tc('common.comment')"
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
export default {
  props: {
    modelValue: {
      type: Object,
      default: () => ({}),
    },
  },
  emits: ['update:modelValue'],
  methods: {
    updateField(field, value) {
      this.$emit('update:modelValue', {
        ...this.modelValue,
        [field]: value,
      });
    },
  },
};
</script>
