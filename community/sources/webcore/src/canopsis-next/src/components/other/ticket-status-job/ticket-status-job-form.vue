<template>
  <v-layout class="gap-3" column>
    <v-layout class="gap-3" wrap>
      <v-text-field
        :value="ruleType"
        :label="$t('jobs.ruleType')"
        readonly
        disabled
      />
      <v-text-field
        :value="form.rule_name"
        :label="$t('jobs.ruleName')"
        readonly
        disabled
      >
        <template #append>
          <c-help-icon
            :text="$t('jobs.ruleNameTooltip')"
            icon="help"
            icon-class="grey--text"
            left
          />
        </template>
      </v-text-field>
      <v-text-field
        :value="form.ticket_system_name"
        :label="$t('jobs.ticketSystemName')"
        readonly
        disabled
      />
    </v-layout>
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
import { computed, onMounted } from 'vue';

import { TEMPLATE_TESTING_TEST_TYPES } from '@/constants';

import { useI18n } from '@/hooks/i18n';
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
  setup(props) {
    const type = TEMPLATE_TESTING_TEST_TYPES.declareTicketRule;

    const { t } = useI18n();

    const ruleType = computed(() => t(`jobs.types.${props.form.rule_type}`));

    const {
      vars: templateVars,
      fetchList: fetchTemplateVarsList,
    } = useTemplateVarsList({ type });

    onMounted(fetchTemplateVarsList);

    return {
      ruleType,
      templateVars,
    };
  },
};
</script>
