<template>
  <v-layout column>
    <span class="pl-4">{{ $t('kpi.filters.helpInformation') }}</span>
    <kpi-filters-list
      :options.sync="options"
      :filters="filters"
      :pending="filtersPending"
      :total-items="filtersMeta.total_count"
      :duplicable="hasCreateAnyKpiFiltersAccess"
      :removable="hasDeleteAnyKpiFiltersAccess"
      :updatable="hasUpdateAnyKpiFiltersAccess"
      @edit="showEditFilterModal"
      @duplicate="showDuplicateFilterModal"
      @remove="showDeleteFilterModal"
    />
  </v-layout>
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
        name: MODALS.createKpiFilter,
        config: {
          filter,
          title: this.$t('modals.createFilter.edit.title'),
          action: async (data) => {
            await this.updateFilter({ id: filter._id, data });

            this.fetchFiltersListWithPreviousParams();
          },
        },
      });
    },

    showDuplicateFilterModal(filter) {
      this.$modals.show({
        name: MODALS.createKpiFilter,
        config: {
          filter: omit(filter, ['_id']),
          title: this.$t('modals.createFilter.duplicate.title'),
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
