<template>
  <c-information-block :title="$t('storageSetting.remediation.title')">
    <template
      v-if="history"
      #subtitle=""
    >
      <storage-settings-history-message :history="history" />
    </template>
    <c-enabled-duration-field
      v-field="form.delete_after"
      :label="$t('storageSetting.remediation.deleteAfter')"
      :suffix="$t('common.after')"
      name="remediation.delete_after"
      switcher
      hide-value-on-false
    />
    <c-enabled-duration-field
      v-field="form.delete_stats_after"
      :label="$t('storageSetting.remediation.deleteStatsAfter')"
      :after="form.delete_after"
      :suffix="$t('common.after')"
      name="remediation.delete_stats_after"
      switcher
      hide-value-on-false
    />
    <c-enabled-duration-field
      v-field="form.delete_mod_stats_after"
      :label="$t('storageSetting.remediation.deleteModStatsAfter')"
      :after="form.delete_stats_after"
      :suffix="$t('common.after')"
      name="remediation.delete_mod_stats_after"
      switcher
      hide-value-on-false
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
};
</script>
