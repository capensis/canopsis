<template>
  <v-layout column>
    <c-information-block
      :title="$t('storageSetting.webhook.title')"
      :help-text="$t('storageSetting.webhook.titleHelp')"
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
        :label="$t('storageSetting.webhook.deleteAfter')"
        :help-text="$t('storageSetting.webhook.deleteAfterHelpText')"
        :name="webhookDeleteAfterFieldName"
      />
      <c-enabled-field
        v-field="form.log_credentials"
        :name="webhookLogCredentialsFieldName"
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
  computed: {
    webhookDeleteAfterFieldName() {
      return 'webhook.delete_after';
    },

    webhookLogCredentialsFieldName() {
      return 'webhook.log_credentials';
    },
  },
};
</script>
