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

import { ALARM_PATTERN_FIELDS, ENTITY_PATTERN_FIELDS, MODALS } from '@/constants';

import { entitiesFilterMixin } from '@/mixins/entities/filter';
import { localQueryMixin } from '@/mixins/query-local/query';
import { permissionsTechnicalKpiFiltersMixin } from '@/mixins/permissions/technical/kpi-filters';

import KpiFiltersList from './kpi-filters-list.vue';

export default {
  components: { KpiFiltersList },
  mixins: [localQueryMixin, entitiesFilterMixin, permissionsTechnicalKpiFiltersMixin],
  computed: {
    alarmExcludedAttributes() {
      return [
        ALARM_PATTERN_FIELDS.lastUpdateDate,
        ALARM_PATTERN_FIELDS.lastEventDate,
        ALARM_PATTERN_FIELDS.resolvedAt,
        ALARM_PATTERN_FIELDS.ackAt,
        ALARM_PATTERN_FIELDS.creationDate,
      ];
    },

    entityExcludedItems() {
      return [
        ENTITY_PATTERN_FIELDS.lastEventDate,
      ];
    },
  },
  mounted() {
    this.fetchList();
  },
  methods: {
    fetchList() {
      this.fetchFiltersList({ params: this.getQuery() });
    },

    showEditFilterModal(filter) {
      this.$modals.show({
        name: MODALS.createFilter,
        config: {
          filter,
          title: this.$t('modals.createFilter.edit.title'),
          withTitle: true,
          withEntity: true,
          withAlarm: true,
          entityExcludedItems: this.entityExcludedItems,
          alarmExcludedAttributes: this.alarmExcludedAttributes,
          action: async (data) => {
            await this.updateFilter({ id: filter._id, data });

            this.fetchFiltersListWithPreviousParams();
          },
        },
      });
    },

    showDuplicateFilterModal(filter) {
      this.$modals.show({
        name: MODALS.createFilter,
        config: {
          filter: omit(filter, ['_id']),
          title: this.$t('modals.createFilter.duplicate.title'),
          withTitle: true,
          withEntity: true,
          withAlarm: true,
          entityExcludedItems: this.entityExcludedItems,
          alarmExcludedAttributes: this.alarmExcludedAttributes,
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
