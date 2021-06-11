<template lang="pug">
  v-layout(column)
    storage-setting-block(:title="$t('storageSetting.junit.title')")
      template(v-if="history.junit", slot="subtitle") {{ junitSubTitle }}
      storage-setting-duration-field(
        v-field="form.junit.delete_after",
        :label="$t('storageSetting.junit.deleteAfter')",
        :help-text="$t('storageSetting.junit.deleteAfterHelpText')",
        :name="junitDeleteAfterFieldName"
      )
    storage-setting-block(:title="$t('storageSetting.remediation.title')")
      template(v-if="history.remediation", slot="subtitle") {{ remediationSubTitle }}
      storage-setting-duration-field(
        v-field="form.remediation.accumulate_after",
        :label="$t('storageSetting.remediation.accumulateAfter')",
        :name="remediationAccumulateAfterFieldName"
      )
      storage-setting-duration-field(
        v-field="form.remediation.delete_after",
        :label="$t('storageSetting.remediation.deleteAfter')",
        :help-text="$t('storageSetting.remediation.deleteAfterHelpText')",
        :name="remediationDeleteAfterFieldName"
      )
</template>

<script>
import { DATETIME_FORMATS } from '@/constants';

import StorageSettingBlock from './partials/storage-setting-block.vue';
import StorageSettingDurationField from './partials/storage-setting-duration-field.vue';

export default {
  inject: ['$validator'],
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
    junitDeleteAfterFieldName() {
      return 'junit.delete_after';
    },

    remediationAccumulateAfterFieldName() {
      return 'remediation.accumulate_after';
    },

    remediationDeleteAfterFieldName() {
      return 'remediation.delete_after';
    },

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
  watch: {
    'form.remediation': function remediationWatcher() {
      this.$validator.validateAll([
        this.remediationAccumulateAfterFieldName,
        this.remediationDeleteAfterFieldName,
      ]);
    },
  },
};
</script>
