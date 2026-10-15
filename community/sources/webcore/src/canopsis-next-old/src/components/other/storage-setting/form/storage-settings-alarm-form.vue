<template>
  <c-information-block :title="$tc('common.alarm', 2)">
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
      :suffix="$t('common.after')"
      name="alarm.archive_after"
      switcher
      hide-value-on-false
    />
    <c-enabled-duration-field
      v-field="form.delete_after"
      :label="$t('storageSetting.alarm.deleteAfter')"
      :suffix="$t('common.after')"
      :after="form.archive_after"
      name="alarm.delete_after"
      switcher
      hide-value-on-false
      @input="validateDeleteAfter"
    />
  </c-information-block>
</template>

<script>
import { nextTick } from 'vue';

import { useValidator } from '@/hooks/validator/validator';

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
  setup() {
    const validator = useValidator();

    const validateDeleteAfter = () => nextTick(() => validator.validate('alarm.delete_after'));

    return {
      validateDeleteAfter,
    };
  },
};
</script>
