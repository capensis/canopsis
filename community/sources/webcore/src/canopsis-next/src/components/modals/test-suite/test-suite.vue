<template>
  <c-modal-wrapper
    :title-color="color"
    close
  >
    <template #title>
      <span>{{ testSuite.name }}</span>
    </template>
    <template #text>
      <c-layout v-if="testSuiteHistoryPending" justify="center">
        <v-progress-circular
          color="primary"
          indeterminate
        />
      </c-layout>
      <test-suite-history
        v-else
        :test-suite-history="testSuiteHistory"
      />
    </template>
    <template #actions>
      <v-btn
        variant="text"
        @click="$modals.hide({ id: modal.id })"
      >
        {{ $t('common.cancel') }}
      </v-btn>
    </template>
  </c-modal-wrapper>
</template>

<script>
import { MAX_LIMIT, MODALS } from '@/constants';

import { entitiesTestSuiteHistoryMixin } from '@/mixins/entities/test-suite-history';

import TestSuiteHistory from '@/components/other/test-suite/test-suite-history.vue';

export default {
  name: MODALS.testSuite,
  components: { TestSuiteHistory },
  mixins: [entitiesTestSuiteHistoryMixin],
  props: {
    modal: {
      type: Object,
      required: true,
    },
  },
  computed: {
    config() {
      return this.modal.config ?? {};
    },

    testSuite() {
      return this.config.testSuite;
    },

    color() {
      return this.config.color;
    },
  },
  mounted() {
    this.fetchTestSuiteHistorySummaryList({
      id: this.testSuite.test_suite_id,
      params: { limit: MAX_LIMIT },
    });
  },
};
</script>
