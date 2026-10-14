<template>
  <div>
    <v-tabs v-model="activeTab" color="primary" align-tabs="center">
      <v-tab value="summary">
        {{ $t('common.summary') }}
      </v-tab>
      <v-tab value="messages">
        {{ $t('testSuite.tabs.globalMessages') }}
      </v-tab>
      <v-tab value="gantt">
        {{ $t('testSuite.tabs.gantt') }}
      </v-tab>
      <v-tab value="details">
        {{ $t('testSuite.tabs.details') }}
      </v-tab>
      <v-tab v-if="hasScreenshots" value="screenshots">
        {{ $t('testSuite.tabs.screenshots') }}
      </v-tab>
      <v-tab v-if="hasVideos" value="videos">
        {{ $t('testSuite.tabs.videos') }}
      </v-tab>
    </v-tabs>

    <v-window v-model="activeTab">
      <v-window-item value="summary">
        <test-suite-summary-tab
          :test-suite="testSuite"
          class="pt-3"
        />
      </v-window-item>
      <v-window-item value="messages">
        <test-suite-messages-tab
          :test-suite="testSuite"
          class="pt-3"
        />
      </v-window-item>
      <v-window-item value="gantt">
        <test-suite-gantt-tab
          :test-suite="testSuite"
          class="pt-3"
        />
      </v-window-item>
      <v-window-item value="details">
        <test-suite-details-tab
          :test-suite="testSuite"
          class="pt-3"
        />
      </v-window-item>
      <v-window-item v-if="hasScreenshots" value="screenshots">
        <test-suite-screenshots-tab
          :test-suite="testSuite"
          class="pt-3"
        />
      </v-window-item>
      <v-window-item v-if="hasVideos" value="videos">
        <test-suite-videos-tab
          :test-suite="testSuite"
          class="pt-3"
        />
      </v-window-item>
    </v-window>
  </div>
</template>

<script>
import { entitiesTestSuitesMixin } from '@/mixins/entities/test-suite';

import TestSuiteDetailsTab from './test-suite-details-tab.vue';
import TestSuiteGanttTab from './test-suite-gantt-tab.vue';
import TestSuiteMessagesTab from './test-suite-messages-tab.vue';
import TestSuiteScreenshotsTab from './test-suite-screenshots-tab.vue';
import TestSuiteSummaryTab from './test-suite-summary-tab.vue';
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
      activeTab: 'summary',
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
      const response = await this.fetchTestSuiteItemSummaryWithoutStore({ id: this.testSuiteId });

      this.testSuite = response?.data ?? response ?? {};
    },
  },
};
</script>
