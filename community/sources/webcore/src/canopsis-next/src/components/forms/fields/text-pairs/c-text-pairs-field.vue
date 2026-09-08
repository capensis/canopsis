<template>
  <v-layout
    :class="{ 'text-pairs__disabled': disabled }"
    class="text-pairs"
    column
  >
    <slot
      v-if="!items.length"
      name="no-data"
    />
    <c-form-block-array-field
      v-field="items"
      :add-button-label="addButtonLabel"
      :disabled="disabled"
      :form="items"
      :item-key="itemKey"
      :item-to-form="textPairToForm"
      :error-messages="errors.collect(name)"
      :required-error-message="requiredErrorMessage"
      :label="title || ''"
    >
      <template #item="{ index, remove }">
        <c-text-pair-field
          v-field="items[index]"
          :disabled="disabled"
          :value-required="valueRequired"
          :text-required="textRequired"
          :text-label="textLabel"
          :value-label="valueLabel"
          :item-text="itemText"
          :item-value="itemValue"
          :name="items[index][itemKey]"
          :variables="variables"
          :items="valueItems"
          @remove="remove"
        >
          <template #append-value="">
            <slot
              :item="items[index]"
              name="append-value"
            />
          </template>
        </c-text-pair-field>
      </template>
    </c-form-block-array-field>
  </v-layout>
</template>

<script>
import { textPairToForm } from '@/helpers/text-pairs';

export default {
  inject: ['$validator'],
  model: {
    prop: 'items',
    event: 'input',
  },
  props: {
    title: {
      type: String,
      default: null,
    },
    items: {
      type: Array,
      default: () => [],
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
      required: false,
    },
    itemValue: {
      type: String,
      required: false,
    },
    itemKey: {
      type: String,
      default: 'key',
    },
    name: {
      type: String,
      default: 'items',
    },
    textRequired: {
      type: Boolean,
      default: false,
    },
    valueRequired: {
      type: Boolean,
      default: false,
    },
    addButtonLabel: {
      type: String,
      required: false,
    },
    requiredErrorMessage: {
      type: String,
      default: '',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    valueItems: {
      type: Array,
      default: () => [],
    },
    variables: {
      type: Array,
      default: () => [],
    },
  },
  setup() {
    return {
      textPairToForm,
    };
  },
};
</script>
