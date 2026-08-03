<template>
  <c-information-block :title="$t('common.alarm', 2)">
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
      :duration="form.archive_after"
      :label="$t('storageSetting.alarm.archiveAfter')"
      :suffix="$t('common.after')"
      name="alarm.archive_after"
      switcher
      hide-value-on-false
      @update:duration="updateField('archive_after', $event)"
    />
    <c-enabled-duration-field
      :duration="form.delete_after"
      :label="$t('storageSetting.alarm.deleteAfter')"
      :suffix="$t('common.after')"
      :after="form.archive_after"
      name="alarm.delete_after"
      switcher
      hide-value-on-false
      @update:duration="updateDeleteAfter"
    />
  </c-information-block>
</template>

<script>
import { nextTick } from 'vue';

import { useModelField } from '@/hooks/form';
import { useValidator } from '@/hooks/validator/validator';

import StorageSettingsHistoryMessage from '../partials/storage-settings-history-message.vue';

export default {
  components: { StorageSettingsHistoryMessage },

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

  emits: ['update:form'],

  setup(props, { emit }) {
    const validator = useValidator();
    const { updateField } = useModelField(props, emit);

    const validateDeleteAfter = () => nextTick(() => validator.validate('alarm.delete_after'));
    const updateDeleteAfter = (duration) => {
      updateField('delete_after', duration);
      validateDeleteAfter();
    };

    return {
      updateField,
      updateDeleteAfter,
    };
  },
};
</script>
