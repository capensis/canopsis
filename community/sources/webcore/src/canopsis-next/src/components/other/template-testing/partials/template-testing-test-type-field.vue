<template>
  <c-select-field
    v-field="value"
    :items="items"
    :label="label || $t('templateTesting.filterByRuleType')"
    :name="name"
    :clearable="clearable"
    :required="required"
  />
</template>

<script>
import { computed } from 'vue';

import { TEMPLATE_TESTING_TEST_TYPES, TEMPLATE_TESTING_TESTS_TYPES_TO_PERMISSIONS } from '@/constants';

import { useI18n } from '@/hooks/i18n';
import { useCurrentUserPermissions } from '@/hooks/auth';

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
    const { tc } = useI18n();
    const { checkReadAccess } = useCurrentUserPermissions();

    const items = computed(() => Object.values(TEMPLATE_TESTING_TEST_TYPES)
      .filter(type => checkReadAccess(TEMPLATE_TESTING_TESTS_TYPES_TO_PERMISSIONS[type]))
      .map(type => ({
        value: type,
        text: tc(`templateTesting.testTypes.${type}`),
      })));

    return {
      items,
    };
  },
};
</script>
