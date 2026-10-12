<template>
  <v-data-table :headers="headers" :items="alarmExecutions" hide-default-footer>
    <template #item="{ item }">
      <tr>
        <td class="text-left pre-wrap">
          {{ getExecutionItem(item).alarm.v.connector_name }}
        </td>
        <td class="text-left pre-wrap">
          {{ getExecutionItem(item).alarm.v.connector }}
        </td>
        <td class="text-left pre-wrap">
          {{ getExecutionItem(item).alarm.v.component }}
        </td>
        <td class="text-left pre-wrap">
          {{ getExecutionItem(item).alarm.v.resource }}
        </td>
        <template v-if="!isOneExecution">
          <td class="text-left pre-wrap">
            {{ getExecutionItem(item).ruleName }}
          </td>
          <td class="text-left">
            <alarm-webhook-execution-status
              :running="isExecutionRunning(getExecutionItem(item))"
              :success="isExecutionSucceeded(getExecutionItem(item))"
              :fail-reason="getExecutionItem(item).fail_reason"
            />
          </td>
        </template>
      </tr>
    </template>
  </v-data-table>
</template>

<script>
import {
  isWebhookExecutionRunning,
  isWebhookExecutionSucceeded,
} from '@/helpers/entities/webhook-execution/entity';

import AlarmWebhookExecutionStatus from '@/components/other/alarm/partials/alarm-webhook-execution-status.vue';

export default {
  components: { AlarmWebhookExecutionStatus },
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
          title: this.$t('common.connectorName'),
          sortable: false,
        },
        {
          title: this.$t('common.connector'),
          sortable: false,
        },
        {
          title: this.$t('common.component'),
          sortable: false,
        },
        {
          title: this.$t('common.resource'),
          sortable: false,
        },
      ];

      if (!this.isOneExecution) {
        headers.push(
          {
            title: this.$t('common.rule'),
            sortable: false,
          },
          {
            title: this.$t('common.status'),
            sortable: false,
          },
        );
      }

      return headers;
    },
  },
  methods: {
    getExecutionItem(item) {
      return item.raw ?? item;
    },

    isExecutionRunning(execution) {
      return isWebhookExecutionRunning(execution);
    },

    isExecutionSucceeded(execution) {
      return isWebhookExecutionSucceeded(execution);
    },
  },
};
</script>
