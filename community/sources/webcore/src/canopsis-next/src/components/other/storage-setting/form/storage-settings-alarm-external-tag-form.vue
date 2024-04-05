<template>
  <c-information-block
    :title="$t('storageSetting.alarmExternalTag.title')"
    help-icon-color="info"
  >
    <template
      v-if="history"
      #subtitle=""
    >
      <storage-settings-history-message
        :history="history"
        deleted-count-message-key="storageSetting.history.alarmExternalTag.deletedCount"
        hide-archived
      />
    </template>
    <c-enabled-duration-field
      v-field="form.delete_after"
      :label="$t('storageSetting.alarmExternalTag.deleteAfter')"
      :name="alarmExternalTagDeleteAfterFieldName"
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
      type: Object,
      required: false,
    },
  },
  computed: {
    alarmExternalTagDeleteAfterFieldName() {
      return 'alarm_external_tag.delete_after';
    },
  },
};
</script>
