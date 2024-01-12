<template>
  <c-information-block
    :title="$t('storageSetting.remediation.title')"
    help-icon-color="info"
  >
    <template
      v-if="history"
      #subtitle=""
    >
      <storage-settings-history-message :history="history" />
    </template>
    <c-enabled-duration-field
      v-field="form.delete_after"
      :label="$t('storageSetting.remediation.deleteAfter')"
      :help-text="$t('storageSetting.remediation.deleteAfterHelpText')"
      :name="remediationDeleteAfterFieldName"
    />
    <c-enabled-duration-field
      v-field="form.delete_stats_after"
      :label="$t('storageSetting.remediation.deleteStatsAfter')"
      :help-text="$t('storageSetting.remediation.deleteStatsAfterHelpText')"
      :name="remediationDeleteStatsAfterFieldName"
      :after="form.delete_after"
    />
    <c-enabled-duration-field
      v-field="form.delete_mod_stats_after"
      :label="$t('storageSetting.remediation.deleteModStatsAfter')"
      :help-text="$t('storageSetting.remediation.deleteModStatsAfterHelpText')"
      :name="remediationDeleteModStatsAfterFieldName"
      :after="form.delete_stats_after"
    />
  </c-information-block>
</template>

<script>
import StorageSettingsHistoryMessage from '../partials/storage-settings-history-message.vue';

export default {
  inject: ['$validator'],
  components: { StorageSettingsHistoryMessage },
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
      type: Number,
      required: false,
    },
  },
  computed: {
    remediationDeleteAfterFieldName() {
      return 'remediation.delete_after';
    },

    remediationDeleteStatsAfterFieldName() {
      return 'remediation.delete_stats_after';
    },

    remediationDeleteModStatsAfterFieldName() {
      return 'remediation.delete_mod_stats_after';
    },
  },
  watch: {
    form() {
      this.$validator.validateAll([
        this.remediationDeleteAfterFieldName,
      ]);
    },
  },
};
</script>
