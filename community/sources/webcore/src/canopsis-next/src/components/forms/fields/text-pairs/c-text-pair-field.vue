<template lang="pug">
  v-layout.text-pair(justify-space-between, align-center)
    v-flex.pa-1(xs6)
      v-text-field(
        v-field="item[itemText]",
        v-validate="textValidationRules",
        :label="textLabel",
        :disabled="disabled",
        :name="textFieldName",
        :error-messages="errors.collect(textFieldName)"
      )
    v-flex.pa-1(xs6)
      v-text-field(
        v-if="!mixed",
        v-field="item[itemValue]",
        v-validate="valueValidationRules",
        :label="valueLabel",
        :disabled="disabled",
        :name="valueFieldName",
        :error-messages="errors.collect(valueFieldName)"
      )
      c-mixed-field(
        v-else,
        v-field="item[itemValue]",
        v-validate="valueValidationRules",
        :name="valueFieldName",
        :disabled="disabled",
        :error-messages="errors.collect(valueFieldName)"
      )
    c-action-btn(v-if="!disabled", type="delete", @click="$emit('remove')")
</template>

<script>
export default {
  inject: ['$validator'],
  model: {
    prop: 'item',
    event: 'input',
  },
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
    textValidationRules: {
      type: String,
      default: 'required',
    },
    valueValidationRules: {
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
  },
  computed: {
    textFieldName() {
      return `${this.name}.${this.itemText}`;
    },

    valueFieldName() {
      return `${this.name}.${this.itemValue}`;
    },
  },
};
</script>
