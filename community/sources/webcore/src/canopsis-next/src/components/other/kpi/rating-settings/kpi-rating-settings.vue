<template lang="pug">
  v-layout(column)
    span.pl-4.subheading {{ $t('kpi.ratingSettings.helpInformation') }}
    kpi-rating-settings-list(
      :pagination.sync="pagination",
      :rating-settings="ratingSettings",
      :pending="ratingSettingsPending",
      :total-items="ratingSettingsMeta.total_count",
      :updatable="hasUpdateAnyKpiRatingSettingsAccess",
      @change-selected="changeSelectedRatingSettings"
    )
</template>

<script>
import { entitiesRatingSettingsMixin } from '@/mixins/entities/rating-settings';
import { localQueryMixin } from '@/mixins/query-local/query';
import { permissionsTechnicalKpiRatingSettingsMixin } from '@/mixins/permissions/technical/kpi-rating-settings';

import KpiRatingSettingsList from './kpi-rating-settings-list.vue';

export default {
  components: { KpiRatingSettingsList },
  mixins: [
    localQueryMixin,
    entitiesRatingSettingsMixin,
    permissionsTechnicalKpiRatingSettingsMixin,
  ],
  mounted() {
    this.fetchList();
  },
  methods: {
    fetchList() {
      this.fetchRatingSettingsList({ params: this.getQuery() });
    },

    async changeSelectedRatingSettings(changedRatingSettings) {
      await this.bulkUpdateRatingSettings({
        data: changedRatingSettings.map(ratingSetting => ({
          id: ratingSetting.id,
          enabled: ratingSetting.enabled,
        })),
      });

      this.fetchList();
    },
  },
};
</script>
