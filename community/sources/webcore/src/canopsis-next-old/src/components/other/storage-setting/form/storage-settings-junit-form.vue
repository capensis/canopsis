<template>
  <c-information-block :title="$t('storageSetting.junit.title')">
    <template
      v-if="history"
      #subtitle=""
    >
      <storage-settings-history-message :history="history" />
    </template>
    <c-enabled-duration-field
      v-field="form.delete_after"
      :label="$t('storageSetting.junit.deleteAfter')"
      :help-text="$t('storageSetting.junit.deleteAfterHelpText')"
      :suffix="$t('storageSetting.olderThan')"
      name="junit.delete_after"
      switcher
      hide-value-on-false
    />
  </c-information-block>
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
