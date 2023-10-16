<template>
  <v-layout column>
    <c-information-block
      :title="$t('storageSetting.healthCheck.title')"
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
        :label="$t('storageSetting.healthCheck.deleteAfter')"
        :name="healthCheckDeleteAfterFieldName"
      />
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
    healthCheckDeleteAfterFieldName() {
      return 'health_check.delete_after';
    },
  },
};
</script>
