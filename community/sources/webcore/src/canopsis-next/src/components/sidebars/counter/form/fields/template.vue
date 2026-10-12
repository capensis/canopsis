<template>
  <field-text-editor
    :model-value="modelValue"
    :title="title"
    :variables="variables"
    @update:model-value="$emit('update:modelValue', $event)"
  />
</template>

<script>
import FieldTextEditor from '@/components/sidebars/form/fields/text-editor.vue';

export default {
  components: { FieldTextEditor },

  props: {
    modelValue: {
      type: String,
      required: true,
    },
    title: {
      type: String,
      required: true,
    },
    variables: {
      type: Array,
      required: false,
    },
  },

  emits: ['update:modelValue'],
};
</script>
