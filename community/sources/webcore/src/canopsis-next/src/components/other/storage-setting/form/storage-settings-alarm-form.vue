<template>
  <c-information-block
    :title="$t('storageSetting.alarm.title')"
    :help-text="$t('storageSetting.alarm.titleHelp')"
    help-icon-color="info"
  >
    <template
      v-if="history"
      #subtitle=""
    >
      <storage-settings-history-message
        :history="history"
        archived-count-message-key="storageSetting.history.alarm.archivedCount"
        deleted-count-message-key="storageSetting.history.alarm.deletedCount"
      />
    </template>
    <c-enabled-duration-field
      v-field="form.archive_after"
      :label="$t('storageSetting.alarm.archiveAfter')"
      :name="alarmArchiveAfterFieldName"
    />
    <c-enabled-duration-field
      v-field="form.delete_after"
      :label="$t('storageSetting.alarm.deleteAfter')"
      :name="alarmDeleteAfterFieldName"
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
      type: Object,
      required: false,
    },
  },
  computed: {
    alarmArchiveAfterFieldName() {
      return 'alarm.archive_after';
    },

    alarmDeleteAfterFieldName() {
      return 'alarm.delete_after';
    },
  },
  watch: {
    form() {
      this.$validator.validateAll([
        this.alarmArchiveAfterFieldName,
        this.alarmDeleteAfterFieldName,
      ]);
    },
  },
};
</script>
