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

    <v-card class="pa-3" color="secondary-lighten-2" flat>
      <v-window v-model="activeTab" class="pt-2">
        <v-window-item :value="0">
          <v-row justify="center">
            <v-col cols="8">
              <v-card class="pa-3">
                <kpi-filter-patterns-form :form="patterns" readonly />
              </v-card>
            </v-col>
          </v-row>
        </v-window-item>
      </v-window>
    </v-card>
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
