<template>
  <c-layout gap="3" align="center">
    <c-text-field
      :model-value="item[itemText]"
      :label="textLabel"
      :disabled="disabled"
      :name="textFieldName"
      :required="textRequired"
      @update:model-value="updateField(itemText, $event)"
    />

    <c-select
      v-if="items.length"
      :model-value="item[itemValue]"
      :label="valueLabel"
      :disabled="disabled"
      :items="items"
      :name="valueFieldName"
      :required="valueRequired"
      item-title="text"
      item-value="value"
      @update:model-value="updateField(itemValue, $event)"
    />
    <c-payload-text-field
      v-else
      :value="item[itemValue]"
      :label="valueLabel"
      :disabled="disabled"
      :name="valueFieldName"
      :variables="variables"
      :required="valueRequired"
      @update:value="updateField(itemValue, $event)"
    >
      <template #append="">
        <slot name="append-value" />
      </template>
    </c-payload-text-field>

    <c-action-btn
      v-if="!disabled"
      type="delete"
      @click="$emit('remove')"
    />
  </c-layout>
</template>

<script>
import { computed } from 'vue';

import { formMixin } from '@/mixins/form';

export default {
  emits: ['update:item'],
  inject: ['$validator'],
  mixins: [formMixin],
  props: {
    item: {
      type: Object,
      default: () => ({}),
    },
    textLabel: {
      type: String,
      default: '',
    },
    valueLabel: {
      type: String,
      default: '',
    },
    itemText: {
      type: String,
      default: 'text',
    },
    itemValue: {
      type: String,
      default: 'value',
    },
    name: {
      type: String,
      default: 'item',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    textRequired: {
      type: Boolean,
      default: false,
    },
    valueRequired: {
      type: Boolean,
      default: false,
    },
    items: {
      type: Array,
      default: () => [],
    },
    variables: {
      type: Array,
      default: () => [],
    },
  },
  setup(props) {
    const textFieldName = computed(() => `${props.name}.${props.itemText}`);
    const valueFieldName = computed(() => `${props.name}.${props.itemValue}`);

    return {
      textFieldName,
      valueFieldName,
    };
  },
};
</script>
