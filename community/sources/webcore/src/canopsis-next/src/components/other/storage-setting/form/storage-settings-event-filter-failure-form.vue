<template>
  <c-information-block
    :title="$t('storageSetting.eventFilterFailure.title')"
    :help-text="$t('storageSetting.eventFilterFailure.titleHelp')"
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
      :label="$t('storageSetting.eventFilterFailure.deleteAfter')"
      :help-text="$t('storageSetting.eventFilterFailure.deleteAfterHelpText')"
      :name="errorsDeleteAfterFieldName"
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
  computed: {
    errorsDeleteAfterFieldName() {
      return 'errors.delete_after';
    },
  },
};
</script>
