<template>
  <c-layout column>
    <c-combobox
      :model-value="rule"
      :items="rules"
      :disabled="disabled"
      :label="label || $t('common.rule')"
      :name="ruleName"
      :loading="pending"
      :clearable="!required"
      item-title="name"
      item-value="_id"
      return-object
      @update:model-value="rule = $event"
    />
    <c-combobox
      :model-value="value.dictionary"
      :items="infosName"
      :disabled="disabled"
      :label="label || $t('common.dictionary')"
      :name="dictionaryName"
      :loading="pending"
      :clearable="!required"
      @update:model-value="updateField('dictionary', $event)"
    />
  </c-layout>
</template>

<script>
import { keyBy, map } from 'lodash';

import { formMixin } from '@/mixins/form';

export default {
  mixins: [formMixin],

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
      default: false,
    },
    required: {
      type: Boolean,
      default: false,
    },
  },

  emits: ['update:value'],

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

    infosName() {
      return map(this.rule?.infos ?? [], 'name');
    },
  },
};
</script>
