<template>
  <div class="pa-2">
    <v-fade-transition
      v-if="testSuitesPending"
      key="progress"
      mode="out-in"
    >
      <v-progress-linear
        class="progress-linear-absolute--top"
        height="2"
        indeterminate
      />
    </v-fade-transition>
    <v-layout
      key="content"
      class="fill-height"
      wrap
    >
      <template v-if="testSuites.length">
        <v-flex
          v-for="testSuite in testSuites"
          :key="testSuite._id"
          xs6
          md4
          lg3
        >
          <testing-weather-item
            :test-suite="testSuite"
            :widget="widget"
          />
        </v-flex>
      </template>
      <v-alert
        v-else
        type="info"
      >
        {{ $t('common.noData') }}
      </v-alert>
    </v-layout>
  </div>
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
