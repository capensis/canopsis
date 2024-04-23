<template>
  <v-data-table
    :headers="headers"
    :items="alarmExecutions"
    hide-default-footer
  >
    <template #item="{ item }">
      <tr>
        <td class="text-left pre-wrap">
          {{ item.alarm.v.connector_name }}
        </td>
        <td class="text-left pre-wrap">
          {{ item.alarm.v.connector }}
        </td>
        <td class="text-left pre-wrap">
          {{ item.alarm.v.component }}
        </td>
        <td class="text-left pre-wrap">
          {{ item.alarm.v.resource }}
        </td>
        <template v-if="!isOneExecution">
          <td class="text-left pre-wrap">
            {{ item.ruleName }}
          </td>
          <td class="text-left">
            <alarm-webhook-execution-status
              :running="isExecutionRunning(item)"
              :success="isExecutionSucceeded(item)"
              :fail-reason="item.fail_reason"
            />
          </td>
        </template>
      </tr>
    </template>
  </v-data-table>
</template>

<script>
import { isWebhookExecutionRunning, isWebhookExecutionSucceeded } from '@/helpers/entities/webhook-execution/entity';

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
      return isWebhookExecutionRunning(execution);
    },

    isExecutionSucceeded(execution) {
      return isWebhookExecutionSucceeded(execution);
    },
  },
};
</script>
