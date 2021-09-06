<template lang="pug">
  v-layout.text-fields(:class="{ 'text-fields__disabled': disabled }", row, wrap)
    v-flex(v-show="title", xs12)
      h4 {{ title }}
    v-flex(xs12)
      slot(v-if="!items.length", name="no-data")
      v-layout.text-field(
        v-for="(item, index) in items",
        :key="item[itemKey]",
        justify-space-between,
        align-center
      )
        v-flex(xs12)
          v-text-field(
            v-if="!mixed",
            v-validate="validationRules",
            :value="item[itemValue]",
            :label="label",
            :disabled="disabled",
            :name="getFieldName(item[itemKey])",
            :error-messages="errors.collect(getFieldName(item[itemKey]))",
            @input="updateFieldInArrayItem(index, itemValue, $event)"
          )
          c-mixed-field(
            v-else,
            v-validate="validationRules",
            :value="item[itemValue]",
            :label="label",
            :name="getFieldName(item[itemKey])",
            :disabled="disabled",
            :error-messages="errors.collect(getFieldName(item[itemKey]))",
            @input="updateFieldInArrayItem(index, itemValue, $event)"
          )
        .text-fields__delete-button
          v-btn(v-if="!disabled", icon, @click="removeItemFromArray(index)")
            v-icon(color="error") delete
    v-flex(v-if="!disabled", xs12)
      v-layout
        v-btn.ml-0(color="primary", @click="addNewItem") {{ addButtonLabel || $t('common.add') }}
</template>

<script>
import { defaultPrimitiveArrayItemCreator } from '@/helpers/forms/shared/common';

import { formArrayMixin } from '@/mixins/form';

export default {
  inject: ['$validator'],
  mixins: [formArrayMixin],
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
    label: {
      type: String,
      default: '',
    },
    itemValue: {
      type: String,
      default: 'value',
    },
    itemKey: {
      type: String,
      default: 'key',
    },
    name: {
      type: String,
      default: 'item',
    },
    validationRules: {
      type: String,
      default: null,
    },
    addButtonLabel: {
      type: String,
      default: null,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    mixed: {
      type: Boolean,
      default: false,
    },
    itemCreator: {
      type: Function,
      default: defaultPrimitiveArrayItemCreator,
    },
  },
  methods: {
    getNamePrefix(key) {
      return `${this.name}[${key}]`;
    },

    getTextFieldName(key) {
      return `${this.getNamePrefix(key)}.${this.itemText}`;
    },

    getFieldName(key) {
      return `${this.getNamePrefix(key)}.${this.itemValue}`;
    },

    addNewItem() {
      this.addItemIntoArray(this.itemCreator());
    },
  },
};
</script>

<style lang="scss" scoped>
.text-fields {
  &:not(.text-fields__disabled) .text-pair {
    position: relative;
    padding-right: 50px;

    &__delete-button {
      position: absolute;
      right: 0;
      top: 50%;
      transform: translateY(-50%);
    }
  }
}
</style>
