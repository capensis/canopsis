<template lang="pug">
  v-layout
    v-flex(xs6)
      v-combobox(
        v-field="value.dictionary",
        v-validate="'required'",
        :items="items",
        :disabled="disabled",
        :label="label || $t('common.dictionary')",
        :return-object="false",
        :name="dictionaryName",
        :error-messages="errors.collect(dictionaryName)",
        item-text="value",
        item-value="value",
        hide-details
      )
    v-flex.pl-3(xs6)
      v-select(
        v-field="value.field",
        v-validate="'required'",
        :items="fieldItems",
        :disabled="disabled || !value.dictionary",
        :label="label || $t('common.field')",
        :name="fieldName",
        :error-messages="errors.collect(fieldName)",
        hide-details
      )
</template>

<script>
import { PATTERN_RULE_INFOS_FIELDS } from '@/constants';

export default {
  inject: ['$validator'],
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: Object,
      required: true,
    },
    items: {
      type: Array,
      default: () => [],
    },
    label: {
      type: String,
      default: '',
    },
    name: {
      type: String,
      default: 'infos',
    },
    divider: {
      type: String,
      default: '.',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    dictionaryName() {
      return `${this.name}.dictionary`;
    },

    fieldName() {
      return `${this.name}.field`;
    },

    fieldItems() {
      return [
        {
          text: this.$t('common.name'),
          value: PATTERN_RULE_INFOS_FIELDS.name,
        },
        {
          text: this.$t('common.value'),
          value: PATTERN_RULE_INFOS_FIELDS.value,
        },
      ];
    },
  },
};
</script>
