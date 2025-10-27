<template>
  <v-layout column>
    <c-information-block :title="$t('storageSetting.healthCheck.title')">
      <template
        v-if="history"
        #subtitle=""
      >
        <storage-settings-history-message :history="history" />
      </template>
      <c-enabled-duration-field
        v-field="form.delete_after"
        :label="$t('storageSetting.healthCheck.deleteAfter')"
        :suffix="$t('common.after')"
        name="health_check.delete_after"
        switcher
        hide-value-on-false
      />
    </c-information-block>
  </v-layout>
</template>

<script>
import StorageSettingsHistoryMessage from '../partials/storage-settings-history-message.vue';

export default {
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
