<template>
  <div class="pa-2">
    <v-fade-transition v-if="testSuitesPending" mode="out-in">
      <v-progress-linear
        key="progress"
        class="progress-linear-absolute--top"
        height="2"
        indeterminate
      />
    </v-fade-transition>
    <v-row key="content" class="fill-height">
      <template v-if="testSuites.length">
        <v-col
          v-for="testSuite in testSuites"
          :key="testSuite._id"
          cols="6"
          md="4"
          lg="3"
        >
          <testing-weather-item
            :test-suite="testSuite"
            :widget="widget"
          />
        </v-col>
      </template>
      <v-alert
        v-else
        type="info"
      >
        {{ $t('common.noData') }}
      </v-alert>
    </v-row>
  </div>
</template>

<script>
import { MAX_LIMIT } from '@/constants';

import { entitiesTestSuitesMixin } from '@/mixins/entities/test-suite';
import { widgetFetchQueryMixin } from '@/mixins/widget/fetch-query';
import { widgetPeriodicRefreshMixin } from '@/mixins/widget/periodic-refresh';

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
