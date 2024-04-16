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
    <v-tab>{{ $t('common.testQuery') }}</v-tab>

    <v-tab-item eager>
      <declare-ticket-rule-general-form
        v-field="form"
        ref="general"
        class="mt-2"
      />
    </v-tab-item>
    <v-tab-item eager>
      <declare-ticket-rule-patterns-form
        v-field="form.patterns"
        ref="patterns"
        class="mt-2"
      />
    </v-tab-item>
    <v-tab-item>
      <v-layout>
        <v-flex
          offset-xs1
          xs10
        >
          <declare-ticket-rule-test-query
            :form="form"
            class="mt-2"
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
