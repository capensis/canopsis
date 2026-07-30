<template>
  <c-select
    :model-value="modelValue"
    :items="items"
    :label="$t('llm.model')"
    :name="name"
    :disabled="disabled"
    :loading="loading"
    :required="required"
    item-title="name"
    item-value="name"
    @update:model-value="$emit('update:modelValue', $event)"
  >
    <template #item="{ props: itemProps, item }">
      <v-list-item v-bind="itemProps">
        <template #title>
          {{ item.raw.name }}
        </template>

        <template v-if="item.raw.recommended" #append>
          <span class="ml-4 text-medium-emphasis">
            {{ $t('llm.recommendedBadge') }}
          </span>
        </template>
      </v-list-item>
    </template>
  </c-select>
</template>

<script>
export default {
  props: {
    modelValue: {
      type: String,
      default: '',
    },
    items: {
      type: Array,
      default: () => [],
    },
    loading: {
      type: Boolean,
      default: false,
    },
    name: {
      type: String,
      default: 'model',
    },
    required: {
      type: Boolean,
      default: true,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
  },

  emits: ['update:modelValue'],
};
</script>
