<template lang="pug">
  v-layout(column)
    span.pl-4.subheading.grey--text.text--darken-2 {{ $t('kpiFilters.helpInformation') }}
    kpi-filters-list(
      :pagination.sync="pagination",
      :filters="filters",
      :pending="filtersPending",
      :total-items="filtersMeta.total_count",
      :duplicable="hasCreateAnyKpiFiltersAccess",
      :removable="hasDeleteAnyKpiFiltersAccess",
      :updatable="hasUpdateAnyKpiFiltersAccess",
      @edit="showEditFilterModal",
      @duplicate="showDuplicateFilterModal",
      @remove="showDeleteFilterModal"
    )
</template>

<script>
import { omit } from 'lodash';

import { MODALS } from '@/constants';

import { entitiesFilterMixin } from '@/mixins/entities/filter';
import { localQueryMixin } from '@/mixins/query-local/query';
import { permissionsTechnicalKpiFiltersMixin } from '@/mixins/permissions/technical/kpi-filters';

import KpiFiltersList from './kpi-filters-list.vue';

export default {
  components: { KpiFiltersList },
  mixins: [localQueryMixin, entitiesFilterMixin, permissionsTechnicalKpiFiltersMixin],
  mounted() {
    this.fetchList();
  },
  methods: {
    fetchList() {
      this.fetchFiltersList({ params: this.getQuery() });
    },

    showEditFilterModal(filter) {
      this.$modals.show({
        name: MODALS.patterns,
        config: {
          title: this.$t('modals.filter.edit.title'),
          name: true,
          entity: true,
          patterns: filter,
          action: async (data) => {
            await this.updateFilter({ id: filter._id, data });

            this.fetchFiltersListWithPreviousParams();
          },
        },
      });
    },

    showDuplicateFilterModal(filter) {
      this.$modals.show({
        name: MODALS.patterns,
        config: {
          title: this.$t('modals.filter.duplicate.title'),
          name: true,
          entity: true,
          patterns: omit(filter, ['_id']),
          action: async (data) => {
            await this.createFilter({ data });

            this.fetchFiltersListWithPreviousParams();
          },
        },
      });
    },

    showDeleteFilterModal(id) {
      this.$modals.show({
        name: MODALS.confirmation,
        config: {
          action: async () => {
            await this.removeFilter({ id });

            this.fetchFiltersListWithPreviousParams();
          },
        },
      });
    },
  },
};
</script>