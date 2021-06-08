<template lang="pug">
  v-layout(column)
    storage-setting-block(:title="$t('storageSetting.junit.title')")
      template(v-if="history.junit", slot="subtitle") {{ junitSubTitle }}
      storage-setting-duration-field(
        v-field="form.junit.delete_after",
        :label="$t('storageSetting.junit.deleteAfter')",
        :help-text="$t('storageSetting.junit.deleteAfterHelpText')",
        name="junit.delete_after"
      )
    storage-setting-block(:title="$t('storageSetting.remediation.title')")
      template(v-if="history.remediation", slot="subtitle") {{ remediationSubTitle }}
      storage-setting-duration-field(
        v-field="form.remediation.accumulate_after",
        :label="$t('storageSetting.remediation.accumulateAfter')",
        name="remediation.accumulate_after"
      )
      storage-setting-duration-field(
        v-field="form.remediation.delete_after",
        :label="$t('storageSetting.remediation.deleteAfter')",
        name="remediation.delete_after"
      )
</template>

<script>
import { DATETIME_FORMATS } from '@/constants';

import StorageSettingBlock from './partials/storage-setting-block.vue';
import StorageSettingDurationField from './partials/storage-setting-duration-field.vue';

export default {
  components: { StorageSettingDurationField, StorageSettingBlock },
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      required: true,
    },
    history: {
      type: Object,
      required: true,
    },
  },
  computed: {
    junitSubTitle() {
      return this.$t('storageSetting.history.junit', {
        launchedAt: this.$options.filters.date(this.history.junit, DATETIME_FORMATS.long, true),
      });
    },
    remediationSubTitle() {
      return this.$t('storageSetting.history.remediation', {
        launchedAt: this.$options.filters.date(this.history.remediation, DATETIME_FORMATS.long, true),
      });
    },
  },
};
</script>
