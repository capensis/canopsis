<template lang="pug">
  advanced-data-table.white(
    :headers="headers",
    :items="remediationJobs",
    :loading="remediationJobsPending",
    :total-items="remediationJobsMeta.total_count",
    :is-disabled-item="isSelectedJob",
    :pagination.sync="pagination",
    select-all,
    advanced-pagination
  )
    template(slot="toolbar", slot-scope="props")
      v-layout(row, justify-space-between)
        v-flex(xs9)
          search-field(@submit="props.updateSearch", @clear="props.clearSearch")
        v-btn(
          v-if="props.selected.length",
          color="primary",
          @click="$emit('select', props.selected)"
        ) {{ $t('common.add') }}
    template(slot="actions", slot-scope="props")
      v-btn(:disabled="props.disabled", icon, small, @click="$emit('select', [props.item])")
        v-icon add
</template>

<script>
import entitiesRemediationJobsMixin from '@/mixins/entities/remediation/jobs';
import localQueryMixin from '@/mixins/query-local/query';

import SearchField from '@/components/forms/fields/search-field.vue';

export default {
  components: { SearchField },
  mixins: [
    entitiesRemediationJobsMixin,
    localQueryMixin,
  ],
  props: {
    jobs: {
      type: Array,
      default: () => [],
    },
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
  methods: {
    isSelectedJob({ _id }) {
      return this.selectedIds.includes(_id);
    },

    async fetchList() {
      this.fetchRemediationJobsList({ params: this.getQuery() });
    },
  },
};
</script>
