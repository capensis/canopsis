<template>
  <c-layout column>
    <c-information-block :title="$t('storageSetting.healthCheck.title')">
      <template v-if="history" #subtitle="">
        <storage-settings-history-message :history="history" />
      </template>
      <c-enabled-duration-field
        :duration="form.delete_after"
        :label="$t('storageSetting.healthCheck.deleteAfter')"
        :suffix="$t('common.after')"
        name="health_check.delete_after"
        switcher
        hide-value-on-false
        @update:duration="updateField('delete_after', $event)"
      />
    </c-information-block>
  </c-layout>
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
