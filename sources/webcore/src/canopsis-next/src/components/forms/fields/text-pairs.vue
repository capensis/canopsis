<template lang="pug">
  v-layout.text-pairs(:class="{ 'text-pairs__disabled': disabled }", row, wrap)
    v-flex(v-show="title", xs12)
      h4.ml-1 {{ title }}
    v-flex(xs12)
      slot(v-if="!items.length", name="no-data")
      v-layout.text-pair(
      v-for="(item, index) in items",
      :key="item[itemKey]",
      justify-space-between,
      align-center
      )
        v-flex.pa-1(xs6)
          v-text-field(
          :value="item[itemText]",
          :label="textLabel",
          :disabled="disabled",
          :name="getTextFieldName(index)",
          :error-messages="getCollectedErrorMessages(getTextFieldName(index))",
          v-validate="textValidationRules",
          @input="updateFieldInArrayItem(index, itemText, $event)"
          )
        v-flex.pa-1(xs6)
          v-text-field(
          v-if="!mixed",
          :value="item[itemValue]",
          :label="valueLabel",
          :disabled="disabled",
          :name="getValueFieldName(index)",
          :error-messages="getCollectedErrorMessages(getValueFieldName(index))",
          v-validate="valueValidationRules",
          @input="updateFieldInArrayItem(index, itemValue, $event)"
          )
          mixed-field(
          v-else
          :value="item[itemValue]",
          :name="getValueFieldName(index)",
          :disabled="disabled",
          :validationRules="valueValidationRules",
          @input="updateFieldInArrayItem(index, itemValue, $event)"
          )
        .text-pair__delete-button
          v-btn(v-if="!disabled", icon, @click="removeItemFromArray(index)")
            v-icon close
    v-flex(v-if="!disabled", xs12)
      v-layout
        v-btn.ml-1(color="primary", @click="addNewItem") {{ addButtonLabel || $t('common.add') }}
</template>

<script>
import { defaultTextPairCreator } from '@/helpers/text-pairs';

import formArrayMixin from '@/mixins/form/array';

import MixedField from '@/components/forms/fields/mixed-field.vue';

export default {
  inject: ['$validator'],
  components: { MixedField },
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
    mixed: {
      type: Boolean,
      default: false,
    },
    itemCreator: {
      type: Function,
      default: defaultTextPairCreator,
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

    getCollectedErrorMessages() {
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

<style lang="scss" scoped>
  .text-pairs {
    &:not(.text-pairs-disabled) .text-pair {
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
