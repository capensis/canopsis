<template lang="pug">
  v-layout(column)
    v-combobox(
      v-model="rule",
      v-validate="ruleRules",
      :items="rules",
      :disabled="disabled",
      :label="label || $tc('common.rule')",
      :name="ruleName",
      :error-messages="errors.collect(ruleName)",
      :loading="pending",
      :clearable="!required",
      item-text="name",
      item-value="_id",
      return-object
    )
    v-combobox(
      v-field="value.dictionary",
      v-validate="dictionaryRules",
      :items="infosName",
      :disabled="disabled",
      :label="label || $t('common.dictionary')",
      :return-object="false",
      :name="dictionaryName",
      :error-messages="errors.collect(dictionaryName)",
      :loading="pending",
      :clearable="!required",
      item-text="value",
      item-value="value"
    )
</template>

<script>
import { keyBy, map } from 'lodash';

import { formMixin } from '@/mixins/form';

export default {
  inject: ['$validator'],
  mixins: [formMixin],
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: Object,
      required: true,
    },
    rules: {
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
    disabled: {
      type: Boolean,
      default: false,
    },
    pending: {
      type: Boolean,
      required: false,
    },
    required: {
      type: Boolean,
      required: false,
    },
  },
  computed: {
    rulesById() {
      return keyBy(this.rules, '_id');
    },

    rule: {
      get() {
        return this.rulesById[this.value.rule] ?? null;
      },

      set(rule) {
        this.updateField('rule', rule?._id);
      },
    },

    ruleName() {
      return `${this.name}.ruleId`;
    },

    dictionaryName() {
      return `${this.name}.dictionary`;
    },

    ruleRules() {
      return {
        required: this.required || Boolean(this.value.dictionary),
      };
    },

    dictionaryRules() {
      return {
        required: this.required || Boolean(this.rule),
      };
    },

    infosName() {
      return map(this.rule?.infos ?? [], 'name');
    },
  },
};
</script>
