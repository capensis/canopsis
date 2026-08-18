<template>
  <c-layout gap="3" column>
    <c-enabled-field v-model="enabledField" with-background />

    <c-form-general-patterns-tabs
      :form="values"
      :rule-id="ruleId"
      :type="type"
      reverse
    >
      <template #general="{ setRef, templateVars, copyVars }">
        <event-filter-general-form
          :ref="setRef"
          :template-vars="templateVars"
          :copy-vars="copyVars"
          :is-disabled-id-field="isDisabledIdField"
        />
      </template>

      <template #patterns="{ setRef }">
        <event-filter-patterns-form
          :ref="setRef"
          :form="values"
          :event-attributes="eventAttributes"
          :attributes-pending="attributesPending"
        />
      </template>
    </c-form-general-patterns-tabs>
  </c-layout>
</template>

<script>
import { TEMPLATE_TESTING_TEST_TYPES } from '@/constants';

import { useFormContext } from '@/hooks/form';

import EventFilterGeneralForm from './event-filter-general-form.vue';
import EventFilterPatternsForm from './event-filter-patterns-form.vue';

export default {
  components: {
    EventFilterGeneralForm,
    EventFilterPatternsForm,
  },
  props: {
    ruleId: {
      type: String,
      default: undefined,
    },
    isDisabledIdField: {
      type: Boolean,
      default: false,
    },
    eventAttributes: {
      type: Array,
      default: () => [],
    },
    attributesPending: {
      type: Boolean,
      default: false,
    },
  },
  setup() {
    const form = useFormContext();

    const enabledField = form.model('enabled');

    const type = TEMPLATE_TESTING_TEST_TYPES.eventFilter;

    return {
      values: form.values,
      enabledField,
      type,
    };
  },
};
</script>
