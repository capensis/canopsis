<template>
  <v-layout :column="column">
    <v-flex :xs6="row">
      <v-combobox
        v-if="combobox"
        v-field="value.dictionary"
        v-validate="'required'"
        :items="items"
        :disabled="disabled"
        :label="label || $t('common.dictionary')"
        :return-object="false"
        :name="dictionaryName"
        :error-messages="errors.collect(dictionaryName)"
        :loading="pending"
        :hide-details="row"
        item-text="value"
        item-value="value"
      />
      <v-text-field
        v-else
        v-field="value.dictionary"
        v-validate="'required'"
        :disabled="disabled"
        :label="label || $t('common.dictionary')"
        :error-messages="errors.collect(dictionaryName)"
        :name="dictionaryName"
        :hide-details="row"
      />
    </v-flex>
    <v-flex
      :class="{ 'pl-3': row }"
      :xs6="row"
    >
      <v-select
        v-field="value.field"
        v-validate="'required'"
        :items="fieldItems"
        :disabled="disabled || !value.dictionary"
        :label="label || $t('common.field')"
        :name="fieldName"
        :error-messages="errors.collect(fieldName)"
        :hide-details="row"
      />
    </v-flex>
  </v-layout>
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
    combobox: {
      type: Boolean,
      default: false,
    },
    row: {
      type: Boolean,
      required: false,
    },
    column: {
      type: Boolean,
      required: false,
    },
    pending: {
      type: Boolean,
      required: false,
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
