<template>
  <c-layout gap="2" column>
    <c-enabled-field
      v-model="enabledField"
      :disabled="disablable || defaultRule"
      hide-details
      with-background
    />

    <c-form-general-patterns-tabs>
      <template #general="{ setRef }">
        <alarm-status-rule-general-form
          :ref="setRef"
          :flapping="flapping"
          :default-rule="defaultRule"
        />
      </template>

      <template v-if="!defaultRule" #patterns="{ setRef }">
        <alarm-status-rule-patterns-form :ref="setRef" :flapping="flapping" />
      </template>
    </c-form-general-patterns-tabs>
  </c-layout>
</template>

<script>
import { useFormModel } from '@/hooks/form';

import AlarmStatusRuleGeneralForm from './alarm-status-rule-general-form.vue';
import AlarmStatusRulePatternsForm from './alarm-status-rule-patterns-form.vue';

export default {
  components: {
    AlarmStatusRuleGeneralForm,
    AlarmStatusRulePatternsForm,
  },

  props: {
    flapping: {
      type: Boolean,
      default: false,
    },
    disablable: {
      type: Boolean,
      default: false,
    },
    defaultRule: {
      type: Boolean,
      default: false,
    },
  },

  setup() {
    const enabledField = useFormModel('enabled');

    return {
      enabledField,
    };
  },
};
</script>
