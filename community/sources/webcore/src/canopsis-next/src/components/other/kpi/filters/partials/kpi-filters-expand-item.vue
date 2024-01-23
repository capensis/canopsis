<template>
  <div>
    <v-tabs
      v-model="activeTab"
      background-color="secondary lighten-1"
      slider-color="primary"
      dark
      centered
    >
      <v-tab>{{ $tc('common.pattern', 2) }}</v-tab>
    </v-tabs>
    <v-layout class="pa-3">
      <v-flex xs12>
        <v-card class="pa-3">
          <v-tabs-items
            v-model="activeTab"
            class="pt-2"
          >
            <v-tab-item>
              <v-flex
                xs12
                lg10
                offset-lg1
              >
                <kpi-filter-patterns-form
                  :form="patterns"
                  readonly
                />
              </v-flex>
            </v-tab-item>
          </v-tabs-items>
        </v-card>
      </v-flex>
    </v-layout>
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
