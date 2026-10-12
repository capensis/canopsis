<template>
  <c-layout column>
    <span class="text-body-2 my-2">
      {{ $t('common.column', 2) }}
    </span>

    <v-alert v-if="required && !modelValue.length" color="info" variant="tonal">
      {{ $t('widgetTemplate.errors.sortColumnsRequired') }}
    </v-alert>

    <c-sort-columns-field
      :model-value="modelValue"
      :items="items"
      :name="name"
      @update:model-value="$emit('update:modelValue', $event)"
    />
  </c-layout>
</template>

<script>
export default {
  props: {
    modelValue: {
      type: Array,
      default: () => [],
    },
    required: {
      type: Boolean,
      default: false,
    },
    name: {
      type: String,
      default: 'sort_columns',
    },
    items: {
      type: Array,
      default: () => [],
    },
  },

  emits: ['update:modelValue'],
};
</script>
