<template>
  <modal-wrapper
    :title-color="color"
    close
  >
    <template #title="">
      <span>{{ testSuite.name }}</span>
    </template>
    <template #text="">
      <v-layout
        v-if="testSuiteHistoryPending"
        justify-center
      >
        <v-progress-circular
          color="primary"
          indeterminate
        />
      </v-layout>
      <test-suite-history
        v-else
        :test-suite-history="testSuiteHistory"
      />
    </template>
    <template #actions="">
      <v-btn
        depressed
        text
        @click="$modals.hide"
      >
        {{ $t('common.cancel') }}
      </v-btn>
    </template>
  </modal-wrapper>
</template>

<script>
import { MAX_LIMIT, MODALS } from '@/constants';

import { modalInnerMixin } from '@/mixins/modal/inner';
import { entitiesTestSuiteHistoryMixin } from '@/mixins/entities/test-suite-history';

import TestSuiteHistory from '@/components/other/test-suite/test-suite-history.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.testSuite,
  components: { ModalWrapper, TestSuiteHistory },
  mixins: [modalInnerMixin, entitiesTestSuiteHistoryMixin],
  computed: {
    testSuite() {
      return this.config.testSuite;
    },

    color() {
      return this.config.color;
    },
  },
  mounted() {
    this.fetchTestSuiteHistorySummaryList({ id: this.testSuite.test_suite_id, params: { limit: MAX_LIMIT } });
  },
};
</script>
