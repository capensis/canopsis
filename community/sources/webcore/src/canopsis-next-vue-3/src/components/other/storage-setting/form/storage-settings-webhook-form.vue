<template>
  <c-layout column>
    <c-information-block :title="$t('storageSetting.webhook.title')">
      <template
        v-if="history"
        #subtitle=""
      >
        <storage-settings-history-message :history="history" />
      </template>
      <c-enabled-duration-field
        :duration="form.delete_after"
        :label="$t('storageSetting.webhook.deleteAfter')"
        :suffix="$t('common.after')"
        name="webhook.delete_after"
        switcher
        hide-value-on-false
        @update:duration="updateField('delete_after', $event)"
      />
      <c-enabled-field
        :model-value="form.log_credentials"
        name="webhook.log_credentials"
        class="pb-2"
        @update:model-value="updateField('log_credentials', $event)"
      >
        <template #label="">
          {{ $t('storageSetting.webhook.logCredentials') }}
          <c-help-icon
            :text="$t('storageSetting.webhook.logCredentialsHelpText')"
            icon-class="ml-2"
            color="info"
            top
          />
        </template>
      </c-enabled-field>
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
