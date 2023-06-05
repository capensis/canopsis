<template lang="pug">
  div.pa-2
    v-fade-transition(v-if="testSuitesPending", key="progress", mode="out-in")
      v-progress-linear.progress-linear-absolute--top(height="2", indeterminate)
    v-layout.fill-height(key="content", wrap)
      template(v-if="testSuites.length")
        v-flex(v-for="testSuite in testSuites", :key="testSuite._id", xs6, md4, lg3)
          testing-weather-item(:test-suite="testSuite", :widget="widget")
      v-alert(v-else, type="info", :value="true") {{ $t('common.noData') }}
</template>

<script>
import { MAX_LIMIT } from '@/constants';

import { widgetPeriodicRefreshMixin } from '@/mixins/widget/periodic-refresh';
import { widgetFetchQueryMixin } from '@/mixins/widget/fetch-query';
import { entitiesTestSuitesMixin } from '@/mixins/entities/test-suite';

import TestingWeatherItem from './testing-weather-item.vue';

export default {
  components: {
    TestingWeatherItem,
  },
  mixins: [
    widgetPeriodicRefreshMixin,
    widgetFetchQueryMixin,
    entitiesTestSuitesMixin,
  ],
  props: {
    widget: {
      type: Object,
      required: true,
    },
  },
  mounted() {
    this.fetchList();
  },
  methods: {
    fetchList() {
      this.fetchTestSuitesList({ widgetId: this.widget._id, params: { limit: MAX_LIMIT } });
    },
  },
};
</script>
