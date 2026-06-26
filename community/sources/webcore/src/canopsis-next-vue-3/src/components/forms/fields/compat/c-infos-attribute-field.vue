<template>
  <c-layout :column="!row" :gap="row ? 3 : 2">
    <c-select
      :model-value="value.dictionary"
      :items="items"
      :name="name"
      :disabled="disabled"
      :label="$t('common.info')"
      item-title="text"
      item-value="value"
      @update:model-value="updateField('dictionary', $event)"
    />
    <c-select
      :model-value="value.field"
      :items="fields"
      :name="name"
      :disabled="disabled"
      :label="$t('common.field')"
      item-title="text"
      item-value="value"
      @update:model-value="updateField('field', $event)"
    />
  </c-layout>
</template>

<script>
import { PATTERN_RULE_INFOS_FIELDS } from '@/constants';
import { formMixin } from '@/mixins/form';

export default {
  emits: ['update:value'],
  mixins: [formMixin],
  props: {
    value: {
      type: Object,
      required: true,
    },
    items: {
      type: Array,
      default: () => [],
    },
    name: {
      type: String,
      default: '',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    row: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    fields() {
      return Object.values(PATTERN_RULE_INFOS_FIELDS).map(value => ({
        value,
        text: this.$t(`common.${value}`),
      }));
    },
  },
};
</script>
