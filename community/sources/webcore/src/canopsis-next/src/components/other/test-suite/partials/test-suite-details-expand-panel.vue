<template lang="pug">
  v-tabs(color="secondary lighten-1", dark, centered, slider-color="primary")
    v-tab {{ $tc('common.information') }}
    v-tab-item
      v-layout.pa-3
        v-flex(xs12)
          v-card.pa-3
            v-layout(column)
              v-flex(offset-xs2)
                test-suite-summary-row(:label="$t('common.name')", :value="testSuiteDetail.name")
                test-suite-summary-row(
                  v-if="testSuiteDetail.description",
                  :label="$t('common.description')",
                  :value="testSuiteDetail.description"
                )
                test-suite-summary-row(
                  v-if="testSuiteDetail.classname",
                  :label="$t('testSuite.className')",
                  :value="testSuiteDetail.classname"
                )
                test-suite-summary-row(
                  v-if="testSuiteDetail.file",
                  :label="$t('common.file')",
                  :value="testSuiteDetail.file"
                )
                test-suite-summary-row(
                  v-if="hasLine",
                  :label="$t('testSuite.line')",
                  :value="testSuiteDetail.line"
                )
                test-suite-summary-row(
                  :label="$t('testSuite.timeTaken')"
                ) {{ testSuiteDetail.time | fixed }}{{ $constants.TIME_UNITS.second }}
                system-message(
                  v-if="testSuiteDetail.message",
                  :value="testSuiteDetail.message",
                  :file-name="testSuiteDetail.file || testSuiteDetail.name"
                )
                  span.font-weight-bold.subheading(slot="label") {{ $t('testSuite.failureMessage') }}
    template(v-if="hasScreenshots")
      v-tab {{ $t('testSuite.tabs.screenshots') }}
      v-tab-item
        v-layout.pa-3
          v-flex(xs12)
            v-card.pa-3
              test-suite-screenshots(:screenshots="testSuiteDetail.screenshots")
    template(v-if="hasVideos")
      v-tab {{ $t('testSuite.tabs.videos') }}
      v-tab-item
        v-layout.pa-3
          v-flex(xs12)
            v-card.pa-3
              test-suite-videos(:videos="testSuiteDetail.videos")
</template>

<script>
import { isNumber } from 'lodash';

import SystemMessage from './system-message.vue';
import TestSuiteSummaryRow from './test-suite-summary-row.vue';
import TestSuiteScreenshots from './test-suite-screenshots.vue';
import TestSuiteVideos from './test-suite-videos.vue';

export default {
  components: {
    TestSuiteSummaryRow,
    SystemMessage,
    TestSuiteScreenshots,
    TestSuiteVideos,
  },
  props: {
    testSuiteDetail: {
      type: Object,
      required: true,
    },
  },
  computed: {
    hasLine() {
      return isNumber(this.testSuiteDetail.line);
    },

    hasScreenshots() {
      return this.testSuiteDetail.screenshots && this.testSuiteDetail.screenshots.length;
    },

    hasVideos() {
      return this.testSuiteDetail.videos && this.testSuiteDetail.videos.length;
    },
  },
};
</script>
