<template>
  <v-tabs
    slider-color="primary"
    centered
  >
    <v-tab :class="{ 'error--text': hasGeneralError }">
      {{ $t('common.general') }}
    </v-tab>
    <v-tab :class="{ 'error--text': hasPatternsError }">
      {{ $tc('common.pattern') }}
    </v-tab>
    <v-tab>{{ $t('declareTicket.testQuery') }}</v-tab>

    <v-tab-item eager>
      <declare-ticket-rule-general-form
        class="mt-2"
        ref="general"
        v-field="form"
      />
    </v-tab-item>
    <v-tab-item eager>
      <declare-ticket-rule-patterns-form
        class="mt-2"
        ref="patterns"
        v-field="form.patterns"
      />
    </v-tab-item>
    <v-tab-item>
      <v-layout>
        <v-flex
          offset-xs1
          xs10
        >
          <declare-ticket-rule-test-query
            class="mt-2"
            :form="form"
          />
        </v-flex>
      </v-layout>
    </v-tab-item>
  </v-tabs>
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
