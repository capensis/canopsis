<template>
  <c-information-block
    :title="$t('storageSetting.eventsRecords.title')"
    :help-text="$t('storageSetting.eventsRecords.titleHelp')"
    help-icon-color="info"
  >
    <template
      v-if="history"
      #subtitle=""
    >
      <storage-settings-history-message :history="history" />
    </template>
    <storage-settings-duration-field
      v-field="form.delete_after"
      :label="$t('storageSetting.eventsRecords.deleteAfter') "
      :name="errorsDeleteAfterFieldName"
    />
  </c-information-block>
</template>

<script>
import StorageSettingsHistoryMessage from '../partials/storage-settings-history-message.vue';

import StorageSettingsDurationField from './partials/storage-settings-duration-field.vue';

export default {
  inject: ['$validator'],
  components: { StorageSettingsHistoryMessage, StorageSettingsDurationField },
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
      return 'eventsRecords.delete_after';
    },
  },
};
</script>
