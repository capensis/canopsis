<template>
  <c-information-block :title="$t('storageSetting.remediation.title')">
    <template
      v-if="history"
      #subtitle=""
    >
      <storage-settings-history-message :history="history" />
    </template>
    <c-enabled-duration-field
      :duration="form.delete_after"
      :label="$t('storageSetting.remediation.deleteAfter')"
      :suffix="$t('common.after')"
      name="remediation.delete_after"
      switcher
      hide-value-on-false
      @update:duration="updateField('delete_after', $event)"
    />
    <c-enabled-duration-field
      :duration="form.delete_stats_after"
      :label="$t('storageSetting.remediation.deleteStatsAfter')"
      :after="form.delete_after"
      :suffix="$t('common.after')"
      name="remediation.delete_stats_after"
      switcher
      hide-value-on-false
      @update:duration="updateField('delete_stats_after', $event)"
    />
    <c-enabled-duration-field
      :duration="form.delete_mod_stats_after"
      :label="$t('storageSetting.remediation.deleteModStatsAfter')"
      :after="form.delete_stats_after"
      :suffix="$t('common.after')"
      name="remediation.delete_mod_stats_after"
      switcher
      hide-value-on-false
      @update:duration="updateField('delete_mod_stats_after', $event)"
    />
  </c-information-block>
</template>

<script>
import { useModelField } from '@/hooks/form';

import StorageSettingsHistoryMessage from '../partials/storage-settings-history-message.vue';

export default {
  components: { StorageSettingsHistoryMessage },

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

  emits: ['update:form'],

  setup(props, { emit }) {
    const { updateField } = useModelField(props, emit);

    return { updateField };
  },
};
</script>
