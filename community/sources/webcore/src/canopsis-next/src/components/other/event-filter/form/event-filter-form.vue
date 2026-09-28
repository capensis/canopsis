<template>
  <v-layout class="gap-3" column>
    <c-enabled-field v-field="form.enabled" with-background />

    <c-form-general-patterns-tabs
      v-field="form"
      :rule-id="ruleId"
      :type="type"
      reverse
    >
      <template #general="{ setRef, templateVars, copyVars }">
        <event-filter-general-form
          v-field="form"
          :ref="setRef"
          :template-vars="templateVars"
          :copy-vars="copyVars"
          :is-disabled-id-field="isDisabledIdField"
        />
      </template>
      <template #patterns="{ setRef }">
        <event-filter-patterns-form
          v-field="form"
          :ref="setRef"
          :event-attributes="eventAttributes"
          :attributes-pending="attributesPending"
        />
      </template>
    </c-form-general-patterns-tabs>
  </v-layout>
</template>

<script>
import { TEMPLATE_TESTING_TEST_TYPES } from '@/constants';

import EventFilterGeneralForm from './event-filter-general-form.vue';
import EventFilterPatternsForm from './event-filter-patterns-form.vue';

export default {
  components: {
    EventFilterGeneralForm,
    EventFilterPatternsForm,
  },
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      default: () => ({}),
    },
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
    const type = TEMPLATE_TESTING_TEST_TYPES.eventFilter;

    return { type };
  },
};
</script>
