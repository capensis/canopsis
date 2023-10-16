<template>
  <v-tabs
    centered
    slider-color="primary"
  >
    <v-tab>{{ $t('common.summary') }}</v-tab>
    <v-tab-item>
      <test-suite-summary-tab
        class="pt-3"
        :test-suite="testSuite"
      />
    </v-tab-item>
    <v-tab>{{ $t('testSuite.tabs.globalMessages') }}</v-tab>
    <v-tab-item>
      <test-suite-messages-tab
        class="pt-3"
        :test-suite="testSuite"
      />
    </v-tab-item>
    <v-tab>{{ $t('testSuite.tabs.gantt') }}</v-tab>
    <v-tab-item>
      <test-suite-gantt-tab
        class="pt-3"
        :test-suite="testSuite"
      />
    </v-tab-item>
    <v-tab>{{ $t('testSuite.tabs.details') }}</v-tab>
    <v-tab-item>
      <test-suite-details-tab
        class="pt-3"
        :test-suite="testSuite"
      />
    </v-tab-item>
    <template v-if="hasScreenshots">
      <v-tab>{{ $t('testSuite.tabs.screenshots') }}</v-tab>
      <v-tab-item>
        <test-suite-screenshots-tab
          class="pt-3"
          :test-suite="testSuite"
        />
      </v-tab-item>
    </template>
    <template v-if="hasVideos">
      <v-tab>{{ $t('testSuite.tabs.videos') }}</v-tab>
      <v-tab-item>
        <test-suite-videos-tab
          class="pt-3"
          :test-suite="testSuite"
        />
      </v-tab-item>
    </template>
  </v-tabs>
</template>

<script>
import { entitiesTestSuitesMixin } from '@/mixins/entities/test-suite';

import TestSuiteSummaryTab from './test-suite-summary-tab.vue';
import TestSuiteMessagesTab from './test-suite-messages-tab.vue';
import TestSuiteGanttTab from './test-suite-gantt-tab.vue';
import TestSuiteDetailsTab from './test-suite-details-tab.vue';
import TestSuiteScreenshotsTab from './test-suite-screenshots-tab.vue';
import TestSuiteVideosTab from './test-suite-videos-tab.vue';

export default {
  components: {
    TestSuiteSummaryTab,
    TestSuiteMessagesTab,
    TestSuiteGanttTab,
    TestSuiteDetailsTab,
    TestSuiteScreenshotsTab,
    TestSuiteVideosTab,
  },
  mixins: [entitiesTestSuitesMixin],
  props: {
    testSuiteId: {
      type: String,
      required: true,
    },
  },
  data() {
    return {
      testSuite: {},
    };
  },
  computed: {
    hasScreenshots() {
      return this.testSuite.screenshots && this.testSuite.screenshots.length;
    },

    hasVideos() {
      return this.testSuite.videos && this.testSuite.videos.length;
    },
  },
  watch: {
    testSuiteId: 'fetchTestSuiteSummary',
  },
  mounted() {
    this.fetchTestSuiteSummary();
  },
  methods: {
    async fetchTestSuiteSummary() {
      this.testSuite = await this.fetchTestSuiteItemSummaryWithoutStore({ id: this.testSuiteId });
    },
  },
};
</script>
