<template>
  <v-card class="ma-4 mt-0">
    <v-tabs
      v-model="activeTab"
      slider-color="primary"
      centered
    >
      <v-tab
        v-for="tab in tabs"
        :key="tab.id"
        :href="`#${tab.id}`"
      >
        {{ $t(tab.label) }}
      </v-tab>
    </v-tabs>
    <v-tabs-items v-model="activeTab">
      <v-tab-item
        v-for="tab in tabs"
        :key="tab.id"
        :value="tab.id"
      >
        <v-card-text>
          <jobs-table
            :columns="tab.columns"
            :items="items"
            :pending="pending"
            :total-items="totalItems"
            :options="options"
            @update:options="$emit('update:options', $event)"
            @action="$emit('action', $event)"
          />
        </v-card-text>
      </v-tab-item>
    </v-tabs-items>
  </v-card>
</template>

<script>
import { ref } from 'vue';

import JobsTable from '@/components/other/jobs/jobs-table.vue';

const JOBS_TABS = [
  {
    id: 'instructions',
    label: 'jobs.tabs.instructions',
    columns: [
      { value: 'ruleName', sortable: true },
      { value: 'status', sortable: true },
      { value: 'startDate', sortable: true },
      { value: 'finishDate', sortable: true, showWhen: 'statusNotInProgress' },
      { value: 'failReason', sortable: false, showWhen: 'statusFailed' },
      { value: 'actions', sortable: false },
    ],
  },
  {
    id: 'webhooks',
    label: 'jobs.tabs.webhooks',
    columns: [
      { value: 'ruleName', sortable: true },
      { value: 'ruleType', sortable: true },
      { value: 'status', sortable: true },
      { value: 'startDate', sortable: true },
      { value: 'finishDate', sortable: true, showWhen: 'statusNotInProgress' },
      { value: 'failReason', sortable: false, showWhen: 'statusFailed' },
      { value: 'actions', sortable: false },
    ],
  },
  {
    id: 'ticket-status',
    label: 'jobs.tabs.ticketStatus',
    columns: [
      { value: 'ticketSystemName', sortable: true },
      { value: 'active', sortable: true },
      { value: 'status', sortable: true },
      { value: 'startDate', sortable: true },
      { value: 'finishDate', sortable: true, showWhen: 'statusNotInProgress' },
      { value: 'failReason', sortable: false, showWhen: 'statusFailed' },
      { value: 'actions', sortable: false },
    ],
  },
  {
    id: 'auth-token',
    label: 'jobs.tabs.authToken',
    columns: [
      { value: 'authTokenName', sortable: true },
      { value: 'status', sortable: true },
      { value: 'startDate', sortable: true },
      { value: 'finishDate', sortable: true, showWhen: 'statusNotInProgress' },
      { value: 'failReason', sortable: false, showWhen: 'statusFailed' },
      { value: 'expirationDate', sortable: true, showWhen: 'statusSucceed' },
      { value: 'actions', sortable: false },
    ],
  },
];

export default {
  components: { JobsTable },
  props: {
    items: {
      type: Array,
      required: true,
    },
    pending: {
      type: Boolean,
      default: false,
    },
    totalItems: {
      type: Number,
      default: 0,
    },
    options: {
      type: Object,
      required: true,
    },
  },
  setup() {
    const tabs = JOBS_TABS;
    const activeTab = ref(tabs[2].id);

    return {
      activeTab,
      tabs,
    };
  },
};
</script>
