<template>
  <div>
    <v-tabs
      v-model="activeTab"
      bg-color="secondary-lighten-1"
      color="primary"
      align-tabs="center"
    >
      <v-tab :value="0">
        {{ $t('common.pattern', 2) }}
      </v-tab>
    </v-tabs>

    <v-window v-model="activeTab" class="bg-secondary-lighten-2">
      <c-window-expand-item :value="0" cols="8">
        <kpi-filter-patterns-form :form="patterns" readonly />
      </c-window-expand-item>
    </v-window>
  </div>
</template>

<script>
import { PATTERNS_FIELDS } from '@/constants';

import { filterPatternsToForm } from '@/helpers/entities/filter/form';

import KpiFilterPatternsForm from '../form/kpi-filter-patterns-form.vue';

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
      return filterPatternsToForm(this.filter, [PATTERNS_FIELDS.entity]);
    },
  },
};
</script>
