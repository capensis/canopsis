<template lang="pug">
  v-layout(row, wrap)
    v-flex(v-show="title", xs12)
      h4.ml-1 {{ title }}
    v-flex(xs12)
      slot(v-if="!items.length", name="no-data")
      v-layout(
      v-for="(item, index) in items",
      :key="item[itemKey]",
      justify-space-between,
      align-center
      )
        v-flex.pa-1
          v-text-field(
          :value="item[itemText]",
          :label="textLabel",
          :disabled="disabled",
          :name="getTextFieldName(index)",
          :error-messages="getErrorMessages(getTextFieldName(index))",
          v-validate="textValidationRules",
          @input="updateFieldInArrayItem(index, itemText, $event)"
          )
        v-flex.pa-1
          v-text-field(
          :value="item[itemValue]",
          :label="valueLabel",
          :disabled="disabled",
          :name="getValueFieldName(index)",
          :error-messages="getErrorMessages(getValueFieldName(index))",
          v-validate="valueValidationRules",
          @input="updateFieldInArrayItem(index, itemValue, $event)"
          )
        v-btn(v-if="!disabled", icon, @click="removeItemFromArray(index)")
          v-icon close
    v-flex(v-if="!disabled", xs12)
      v-layout
        v-btn.ml-1(color="primary", @click="addNewItem") {{ addButtonLabel || $t('common.add') }}
</template>

<script>
import { defaultItemCreator } from '@/helpers/text-pairs';

import formArrayMixin from '@/mixins/form/array';

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
    itemKey: {
      type: String,
      default: 'key',
    },
    name: {
      type: String,
      default: 'item',
    },
    textValidationRules: {
      type: String,
      default: 'required',
    },
    valueValidationRules: {
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
    itemCreator: {
      type: Function,
      default: defaultItemCreator,
    },
  },
  computed: {
    getNamePrefix() {
      return index => `${this.name}[${index}]`;
    },
    getTextFieldName() {
      return index => `${this.getNamePrefix(index)}.${this.itemText}`;
    },
    getValueFieldName() {
      return index => `${this.getNamePrefix(index)}.${this.itemValue}`;
    },
    getErrorMessages() {
      return (name) => {
        if (this.errors) {
          return this.errors.collect(name);
        }

        return [];
      };
    },
  },
  methods: {
    addNewItem() {
      this.addItemIntoArray(this.itemCreator());
    },
  },
};
</script>
