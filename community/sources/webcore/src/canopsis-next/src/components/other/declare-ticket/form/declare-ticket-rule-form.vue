<template lang="pug">
  v-tabs(slider-color="primary", color="transparent", fixed-tabs, centered)
    v-tab(:class="{ 'error--text': hasGeneralError }") {{ $t('common.general') }}
    v-tab-item
      declare-ticket-rule-general-form.mt-2(ref="general", v-field="form")
    v-tab(:class="{ 'error--text': hasPatternsError }") {{ $tc('common.pattern') }}
    v-tab-item
      declare-ticket-rule-patterns-form.mt-2(ref="patterns", v-field="form.patterns")
    v-tab {{ $t('declareTicket.testQuery') }}
    v-tab-item(lazy)
      v-layout(row)
        v-flex(offset-xs1, xs10)
          declare-ticket-rule-test-query.mt-2(:form="form")
</template>

<script>
import DeclareTicketRuleTestQuery from '../partials/declare-ticket-rule-test-query.vue';

import DeclareTicketRuleGeneralForm from './declare-ticket-rule-general-form.vue';
import DeclareTicketRulePatternsForm from './declare-ticket-rule-patterns-form.vue';

export default {
  components: { DeclareTicketRuleTestQuery, DeclareTicketRulePatternsForm, DeclareTicketRuleGeneralForm },
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      default: () => ({}),
    },
  },
  data() {
    return {
      hasGeneralError: false,
      hasPatternsError: false,
    };
  },
  mounted() {
    this.$watch(() => this.$refs.general.hasAnyError, (value) => {
      this.hasGeneralError = value;
    });

    this.$watch(() => this.$refs.patterns.hasAnyError, (value) => {
      this.hasPatternsError = value;
    });
  },
};
</script>
