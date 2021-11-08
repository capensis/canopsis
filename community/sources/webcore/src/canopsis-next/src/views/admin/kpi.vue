<template lang="pug">
  div
    c-page-header
    v-card.ma-4.mt-0
      v-tabs(v-model="activeTab", slider-color="primary", centered)
        v-tab {{ $tc('common.graph', 2) }}
        v-tab {{ $t('kpi.tabs.dataSlices') }}
        v-tab {{ $t('kpi.tabs.detailedMetrics') }}
        v-tab {{ $t('common.rating') }}
        v-tab {{ $t('kpi.tabs.sliPatterns') }}

      v-tabs-items(v-model="activeTab")
        v-card-text
          v-tab-item
            kpi-graphs
          v-tab-item(lazy)
            kpi-data-slices
          v-tab-item(lazy)
            kpi-detailed-metrics
          v-tab-item(lazy)
            kpi-rating
          v-tab-item(lazy)
            kpi-sli-patterns
</template>

<script>
import KpiGraphs from '@/components/other/kpi/graphs/kpi-graphs.vue';
import KpiDataSlices from '@/components/other/kpi/data-slices/kpi-data-slices.vue';
import KpiDetailedMetrics from '@/components/other/kpi/detailed-metrics/kpi-detailed-metrics.vue';
import KpiRating from '@/components/other/kpi/rating/kpi-rating.vue';
import KpiSliPatterns from '@/components/other/kpi/sli-patterns/kpi-sli-patterns.vue';

export default {
  components: { KpiGraphs, KpiDataSlices, KpiDetailedMetrics, KpiRating, KpiSliPatterns },
  data() {
    return {
      activeTab: 0,
    };
  },
};
</script>
