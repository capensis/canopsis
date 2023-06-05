<template lang="pug">
  c-advanced-data-table(
    :headers="headers",
    :items="jobsItems",
    :loading="pending",
    :total-items="totalItems",
    :is-disabled-item="isSelectedJob",
    :pagination.sync="pagination",
    select-all,
    advanced-pagination
  )
    template(#toolbar="{ clearSearch, updateSearch }")
      v-layout(row)
        c-search-field(@submit="updateSearch", @clear="clearSearch")
    template(#actions="{ disabled, item }")
      v-btn(:disabled="disabled", icon, small, @click="$emit('select', [item])")
        v-icon add
    template(#mass-actions="{ selected, item, count }")
      v-btn.ma-2(
        outline,
        color="primary",
        @click="$emit('select', selected)"
      ) {{ $tc('remediation.job.addJobs', count, { count: count }) }}
</template>

<script>
import { localQueryMixin } from '@/mixins/query-local/query';
import { entitiesRemediationJobMixin } from '@/mixins/entities/remediation/job';

export default {
  mixins: [
    localQueryMixin,
    entitiesRemediationJobMixin,
  ],
  props: {
    jobs: {
      type: Array,
      default: () => [],
    },
  },
  data() {
    return {
      pending: false,
      jobsItems: [],
      totalItems: 0,
      query: {
        rowsPerPage: 5,
      },
    };
  },
  computed: {
    selectedIds() {
      return this.jobs.map(({ _id }) => _id);
    },

    headers() {
      return [
        { text: this.$t('common.name'), value: 'name' },
        { text: this.$t('common.type'), value: 'config.type' },
        { text: this.$t('common.actionsLabel'), value: 'actions', sortable: false },
      ];
    },
  },
  mounted() {
    this.fetchList();
  },
  methods: {
    isSelectedJob({ _id }) {
      return this.selectedIds.includes(_id);
    },

    async fetchList() {
      this.pending = true;

      const { data: jobs, meta } = await this.fetchRemediationJobsListWithoutStore({ params: this.getQuery() });

      this.jobsItems = jobs;
      this.totalItems = meta.total_count;
      this.pending = false;
    },
  },
};
</script>
