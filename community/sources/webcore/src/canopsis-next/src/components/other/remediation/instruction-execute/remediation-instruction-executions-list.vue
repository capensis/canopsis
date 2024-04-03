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
    <template #launched_at="{ item }">
      {{ item.launched_at | date('long', '-') }}
    </template>
    <template #completed_at="{ item }">
      {{ item.completed_at | date('long', '-') }}
    </template>
    <template #status="{ item }">
      <c-enabled :value="item.enabled" />
    </template>
    <template #expand="{ item }">
      <v-layout class="pa-3 secondary lighten-1">
        <v-flex xs12>
          <v-card class="tab-item-card">
            <v-card-text>
              <remediation-instruction-executions-expand-panel :steps="item.steps" />
            </v-card-text>
          </v-card>
        </v-flex>
      </v-layout>
    </template>
  </c-advanced-data-table>
</template>

<script>
import RemediationInstructionExecutionsExpandPanel
  from './partials/remediation-instruction-executions-expand-panel.vue';

export default {
  components: { RemediationInstructionExecutionsExpandPanel },
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
  computed: {
    headers() {
      return [
        { text: this.$t('common.instructionName'), value: 'name' },
        { text: this.$t('common.type'), value: 'type' },
        { text: this.$t('remediation.instructionExecute.jobs.launchedAt'), value: 'launched_at' },
        { text: this.$t('remediation.instructionExecute.jobs.completedAt'), value: 'completed_at' },
        { text: this.$t('common.launchedBy'), value: 'author.display_name' },
        { text: this.$t('common.status'), value: 'status' },
      ];
    },
  },
};
</script>
