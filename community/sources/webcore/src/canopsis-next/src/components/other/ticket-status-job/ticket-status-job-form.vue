<template>
  <v-layout class="gap-3" column>
    <v-text-field
      :value="form.ticket_system_name"
      :label="$t('jobs.ticketSystemName')"
      readonly
      disabled
    />
    <v-text-field
      v-field="form.ticket_id"
      v-validate="'required'"
      :label="$tc('common.ticket')"
      :error-messages="errors.collect('ticket')"
      name="ticket"
    />
    <declare-ticket-rule-check-ticket-status-field
      v-field="form.check_ticket_status"
      :template-vars="templateVars"
      class="c-alternative-bg-panel pa-5"
    />
  </v-layout>
</template>

<script>
import { onMounted } from 'vue';

import { TEMPLATE_TESTING_TEST_TYPES } from '@/constants';

import { useTemplateVarsList } from '@/hooks/vars/template';

import DeclareTicketRuleCheckTicketStatusField from '@/components/other/declare-ticket/form/fields/declare-ticket-rule-check-ticket-status-field.vue';

export default {
  inject: ['$validator'],
  components: {
    DeclareTicketRuleCheckTicketStatusField,
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
  },
  setup() {
    const type = TEMPLATE_TESTING_TEST_TYPES.declareTicketRule;

    const {
      vars: templateVars,
      fetchList: fetchTemplateVarsList,
    } = useTemplateVarsList({ type });

    onMounted(fetchTemplateVarsList);

    return {
      templateVars,
    };
  },
};
</script>
