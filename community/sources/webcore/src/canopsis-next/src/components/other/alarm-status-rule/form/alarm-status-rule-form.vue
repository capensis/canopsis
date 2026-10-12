<template>
  <v-layout class="gap-2" column>
    <c-enabled-field
      v-field="form.enabled"
      :disabled="disablable || defaultRule"
      hide-details
      with-background
    />
    <c-form-general-patterns-tabs>
      <template #general="{ setRef }">
        <alarm-status-rule-general-form
          v-field="form"
          :ref="setRef"
          :flapping="flapping"
          :default-rule="defaultRule"
        />
      </template>
      <template v-if="!defaultRule" #patterns="{ setRef }">
        <alarm-status-rule-patterns-form
          v-field="form.patterns"
          :ref="setRef"
          :flapping="flapping"
          class="mt-2"
        />
      </template>
    </c-form-general-patterns-tabs>
  </v-layout>
</template>

<script>
import AlarmStatusRuleGeneralForm from './alarm-status-rule-general-form.vue';
import AlarmStatusRulePatternsForm from './alarm-status-rule-patterns-form.vue';

export default {
  inject: ['$validator'],
  components: {
    AlarmStatusRuleGeneralForm,
    AlarmStatusRulePatternsForm,
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
};
</script>
