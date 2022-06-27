<template lang="pug">
  div
    v-tabs(v-model="activeTab", color="secondary lighten-1", slider-color="primary", dark, centered)
      v-tab {{ $tc('common.pattern', 2) }}
    v-layout.pa-3
      v-flex(xs12)
        v-card.pa-3
          v-tabs-items.pt-2(v-model="activeTab")
            v-tab-item(lazy)
              v-flex(xs12, lg10, offset-lg1)
                kpi-filter-patterns-form(:form="patterns", readonly)
</template>

<script>
import { OLD_PATTERNS_FIELDS, PATTERNS_FIELDS } from '@/constants';

import { filterPatternsToForm } from '@/helpers/forms/filter';

import KpiFilterPatternsForm from '../form/partials/kpi-filter-patterns-form.vue';

export default {
  components: { KpiFilterPatternsForm },
  props: {
    filter: {
      type: Object,
      required: true,
    },
  },
  data() {
    return {
      activeTab: 0,
    };
  },
  computed: {
    patterns() {
      return filterPatternsToForm(this.filter, [PATTERNS_FIELDS.entity], [OLD_PATTERNS_FIELDS.entity]);
    },
  },
};
</script>
