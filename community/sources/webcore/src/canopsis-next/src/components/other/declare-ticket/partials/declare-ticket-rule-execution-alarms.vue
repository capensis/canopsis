<template lang="pug">
  v-data-table(:headers="headers", :items="alarmExecutions", hide-actions)
    template(#items="{ item, index }")
      td.text-xs-left.pre-wrap {{ item.alarm.v.connector_name }}
      td.text-xs-left.pre-wrap {{ item.alarm.v.connector }}
      td.text-xs-left.pre-wrap {{ item.alarm.v.component }}
      td.text-xs-left.pre-wrap {{ item.alarm.v.resource }}
      template(v-if="!isOneExecution")
        td.text-xs-left.pre-wrap {{ item.ruleName }}
        td.text-xs-left
          declare-ticket-rule-execution-status(
            :running="isExecutionRunning(item)",
            :success="isExecutionSucceeded(item)",
            :fail-reason="item.fail_reason"
          )
</template>

<script>
import {
  isDeclareTicketExecutionRunning,
  isDeclareTicketExecutionSucceeded,
} from '@/helpers/forms/declare-ticket-rule';

import ProgressCell from '@/components/common/table/progress-cell.vue';

import DeclareTicketRuleExecutionStatus from './declare-ticket-rule-execution-status.vue';

export default {
  components: { DeclareTicketRuleExecutionStatus, ProgressCell },
  props: {
    alarmExecutions: {
      type: Array,
      default: () => [],
    },
    isOneExecution: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    headers() {
      const headers = [
        {
          text: this.$t('common.connectorName'),
          sortable: false,
        },
        {
          text: this.$t('common.connector'),
          sortable: false,
        },
        {
          text: this.$t('common.component'),
          sortable: false,
        },
        {
          text: this.$t('common.resource'),
          sortable: false,
        },
      ];

      if (!this.isOneExecution) {
        headers.push({
          text: this.$tc('common.rule'),
          sortable: false,
        }, {
          text: this.$t('common.status'),
          sortable: false,
        });
      }

      return headers;
    },
  },
  methods: {
    isExecutionRunning(execution) {
      return isDeclareTicketExecutionRunning(execution);
    },

    isExecutionSucceeded(execution) {
      return isDeclareTicketExecutionSucceeded(execution);
    },
  },
};
</script>
