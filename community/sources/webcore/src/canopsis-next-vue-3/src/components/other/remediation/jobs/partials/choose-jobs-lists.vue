<template>
  <c-advanced-data-table
    v-model:options="options"
    :headers="headers"
    :items="jobsItems"
    :loading="pending"
    :total-items="totalItems"
    :is-disabled-item="isSelectedJob"
    :pagination-total-visible="3"
    select-all
    advanced-pagination
  >
    <template #toolbar="{ updateSearch }">
      <c-search @submit="updateSearch" />
    </template>
    <template #actions="{ disabled, item }">
      <v-btn
        :disabled="disabled"
        icon="add"
        variant="text"
        size="small"
        @click="$emit('select', [item])"
      />
    </template>
    <template #mass-actions="{ selected, count }">
      <v-btn
        class="ma-2"
        color="primary"
        variant="outlined"
        @click="$emit('select', selected)"
      >
        {{ $t('remediation.job.addJobs', count, { count }) }}
      </v-btn>
    </template>
  </c-advanced-data-table>
</template>

<script>
import { convertQueryToRequest } from '@/helpers/query';

import { useRemediationJob } from '@/hooks/store/modules/remediation-job';

export default {
  props: {
    jobs: {
      type: Array,
      default: () => [],
    },
  },
  emits: ['select'],
  setup() {
    const { fetchRemediationJobsListWithoutStore } = useRemediationJob();

    return {
      fetchRemediationJobsListWithoutStore,
    };
  },
  data() {
    return {
      pending: false,
      jobsItems: [],
      totalItems: 0,
      options: {
        page: 1,
        itemsPerPage: 5,
        sortBy: [],
        sortDesc: [],
      },
    };
  },
  computed: {
    selectedIds() {
      return this.jobs.map(({ _id }) => _id);
    },

    headers() {
      return [
        { title: this.$t('common.name'), key: 'name' },
        { title: this.$t('common.type'), key: 'config.type' },
        { title: this.$t('common.actionsLabel'), key: 'actions', sortable: false },
      ];
    },
  },
  watch: {
    options: {
      deep: true,
      handler() {
        this.fetchList();
      },
    },
  },
  mounted() {
    this.fetchList();
  },
  methods: {
    isSelectedJob({ _id }) {
      return this.selectedIds.includes(_id);
    },

    getQuery() {
      return convertQueryToRequest(this.options);
    },

    async fetchList() {
      this.pending = true;

      const { data: jobs, meta } = await this.fetchRemediationJobsListWithoutStore({
        params: this.getQuery(),
      });

      this.jobsItems = jobs;
      this.totalItems = meta.total_count;
      this.pending = false;
    },
  },
};
</script>
