<template>
  <div>
    <v-tabs
      v-model="activeTab"
      bg-color="secondary-lighten-1"
      color="primary"
      theme="dark"
      align-tabs="center"
    >
      <v-tab value="information">
        {{ $tc('common.information') }}
      </v-tab>
      <v-tab v-if="hasScreenshots" value="screenshots">
        {{ $t('testSuite.tabs.screenshots') }}
      </v-tab>
      <v-tab v-if="hasVideos" value="videos">
        {{ $t('testSuite.tabs.videos') }}
      </v-tab>
    </v-tabs>

    <v-window v-model="activeTab">
      <v-window-item value="information">
        <v-row class="pa-3">
          <v-col cols="12">
            <v-card class="pa-3">
              <c-layout column>
                <v-col offset="2">
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
                    {{ convertNumberToFixedString(testSuiteDetail.time) }}{{ $constants.TIME_UNITS.second }}
                  </c-information-block-row>
                  <system-message
                    v-if="testSuiteDetail.message"
                    :value="testSuiteDetail.message"
                    :file-name="testSuiteDetail.file || testSuiteDetail.name"
                  >
                    <template #label>
                      <span class="font-weight-bold text-subtitle-1">
                        {{ $t('testSuite.failureMessage') }}
                      </span>
                    </template>
                  </system-message>
                </v-col>
              </c-layout>
            </v-card>
          </v-col>
        </v-row>
      </v-window-item>
      <v-window-item v-if="hasScreenshots" value="screenshots">
        <v-row class="pa-3">
          <v-col cols="12">
            <v-card class="pa-3">
              <test-suite-screenshots :screenshots="testSuiteDetail.screenshots" />
            </v-card>
          </v-col>
        </v-row>
      </v-window-item>
      <v-window-item v-if="hasVideos" value="videos">
        <v-row class="pa-3">
          <v-col cols="12">
            <v-card class="pa-3">
              <test-suite-videos :videos="testSuiteDetail.videos" />
            </v-card>
          </v-col>
        </v-row>
      </v-window-item>
    </v-window>
  </div>
</template>

<script>
import { isNumber } from 'lodash';

import { convertNumberToFixedString } from '@/helpers/string';

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
      activeTab: 'information',
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
  methods: {
    convertNumberToFixedString,
  },
};
</script>
