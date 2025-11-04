<template>
  <c-select-field
    v-field="value"
    :items="items"
    :label="label || $t('common.filterByType')"
    :name="name"
    :disabled="disabled"
    :required="required"
    :clearable="clearable"
  />
</template>

<script>
import { computed } from 'vue';

import { TEMPLATE_TESTING_DATA_TYPES } from '@/constants';

import { useI18n } from '@/hooks/i18n';

export default {
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: Number,
      required: false,
    },
    label: {
      type: String,
      default: '',
    },
    name: {
      type: String,
      default: 'type',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    required: {
      type: Boolean,
      default: false,
    },
    clearable: {
      type: Boolean,
      default: false,
    },
  },
  setup() {
    const { t } = useI18n();

    const items = computed(() => Object.values(TEMPLATE_TESTING_DATA_TYPES).map(value => ({
      value,
      text: t(`templateTesting.dataTypes.${value}`),
    })));

    return {
      items,
    };
  },
};
</script>
