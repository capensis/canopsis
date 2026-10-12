<template>
  <c-advanced-data-table
    :options="options"
    :items="executions"
    :loading="pending"
    :headers="headers"
    :total-items="totalItems"
    advanced-pagination
    hide-actions
    expand
    @update:options="$emit('update:options', $event)"
  >
    <template #type="{ item }">
      {{ t(`remediation.instruction.types.${item.type}`) }}
    </template>
    <template #started_at="{ item }">
      {{ dateToString(item.started_at) }}
    </template>
    <template #completed_at="{ item }">
      {{ dateToString(item.completed_at) }}
    </template>
    <template #status="{ item }">
      <remediation-instruction-execution-status-icon :status="item.status" />
    </template>
    <template #expand="{ item }">
      <v-card class="tab-item-card" variant="flat">
        <v-card-text>
          <remediation-instruction-executions-expand-panel :steps="item.steps" />
        </v-card-text>
      </v-card>
    </template>
  </c-advanced-data-table>
</template>

<script>
import { computed } from 'vue';

import { DATETIME_FORMATS } from '@/constants';

import { convertDateToString } from '@/helpers/date/date';

import { useI18n } from '@/hooks/i18n';

import RemediationInstructionExecutionStatusIcon from './partials/remediation-instruction-execution-status-icon.vue';
import RemediationInstructionExecutionsExpandPanel from './partials/remediation-instruction-executions-expand-panel.vue';

export default {
  components: {
    RemediationInstructionExecutionStatusIcon,
    RemediationInstructionExecutionsExpandPanel,
  },

  props: {
    executions: {
      type: Array,
      default: () => [],
    },
    options: {
      type: Object,
      default: () => ({}),
    },
    totalItems: {
      type: Number,
      default: 0,
    },
    pending: {
      type: Boolean,
      default: false,
    },
  },

  emits: ['update:options'],

  setup() {
    const { t } = useI18n();

    const headers = computed(() => [
      { text: t('remediation.instruction.name'), value: 'name' },
      { text: t('common.type'), value: 'type' },
      { text: t('remediation.instructionExecute.jobs.launchedAt'), value: 'started_at' },
      { text: t('remediation.instructionExecute.jobs.completedAt'), value: 'completed_at' },
      { text: t('common.author'), value: 'author.display_name' },
      { text: t('common.status'), value: 'status' },
    ]);
    const dateToString = date => convertDateToString(date, DATETIME_FORMATS.long, '-');

    return {
      dateToString,
      headers,
      t,
    };
  },
};
</script>
