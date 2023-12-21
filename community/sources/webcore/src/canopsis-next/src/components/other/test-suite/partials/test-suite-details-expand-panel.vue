<template>
  <v-tabs
    v-model="activeTab"
    background-color="secondary lighten-1"
    slider-color="primary"
    dark
    centered
  >
    <v-tab>{{ $tc('common.information') }}</v-tab>
    <v-tab v-if="hasScreenshots">
      {{ $t('testSuite.tabs.screenshots') }}
    </v-tab>
    <v-tab v-if="hasVideos">
      {{ $t('testSuite.tabs.videos') }}
    </v-tab>

    <v-tabs-items
      v-model="activeTab"
      mandatory
    >
      <v-tab-item>
        <v-layout class="pa-3">
          <v-flex xs12>
            <v-card class="pa-3">
              <v-layout column>
                <v-flex offset-xs2>
                  <c-information-block-row
                    :label="$t('common.name')"
                    :value="testSuiteDetail.name"
                  />
                  <c-information-block-row
                    v-if="testSuiteDetail.description"
                    :label="$t('common.description')"
                    :value="testSuiteDetail.description"
                  />
                  <c-information-block-row
                    v-if="testSuiteDetail.classname"
                    :label="$t('testSuite.className')"
                    :value="testSuiteDetail.classname"
                  />
                  <c-information-block-row
                    v-if="testSuiteDetail.file"
                    :label="$t('common.file')"
                    :value="testSuiteDetail.file"
                  />
                  <c-information-block-row
                    v-if="hasLine"
                    :label="$t('testSuite.line')"
                    :value="testSuiteDetail.line"
                  />
                  <c-information-block-row :label="$t('common.timeTaken')">
                    {{ testSuiteDetail.time | fixed }}{{ $constants.TIME_UNITS.second }}
                  </c-information-block-row>
                  <system-message
                    v-if="testSuiteDetail.message"
                    :value="testSuiteDetail.message"
                    :file-name="testSuiteDetail.file || testSuiteDetail.name"
                  >
                    <template #label="">
                      <span class="font-weight-bold text-subtitle-1">{{ $t('testSuite.failureMessage') }}</span>
                    </template>
                  </system-message>
                </v-flex>
              </v-layout>
            </v-card>
          </v-flex>
        </v-layout>
      </v-tab-item>
      <v-tab-item v-if="hasScreenshots">
        <v-layout class="pa-3">
          <v-flex xs12>
            <v-card class="pa-3">
              <test-suite-screenshots :screenshots="testSuiteDetail.screenshots" />
            </v-card>
          </v-flex>
        </v-layout>
      </v-tab-item>
      <v-tab-item v-if="hasVideos">
        <v-layout class="pa-3">
          <v-flex xs12>
            <v-card class="pa-3">
              <test-suite-videos :videos="testSuiteDetail.videos" />
            </v-card>
          </v-flex>
        </v-layout>
      </v-tab-item>
    </v-tabs-items>
  </v-tabs>
</template>

<script>
import { isNumber } from 'lodash';

import SystemMessage from './system-message.vue';
import TestSuiteScreenshots from './test-suite-screenshots.vue';
import TestSuiteVideos from './test-suite-videos.vue';

export default {
  components: {
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
  data() {
    return {
      activeTab: 0,
    };
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
