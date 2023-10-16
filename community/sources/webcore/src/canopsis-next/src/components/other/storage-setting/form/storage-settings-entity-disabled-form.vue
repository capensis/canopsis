<template>
  <c-information-block
    :title="$t('storageSetting.entity.title')"
    :help-text="$t('storageSetting.entity.titleHelp')"
    help-icon-color="info"
  >
    <template
      v-if="history"
      #subtitle=""
    >
      <storage-settings-history-message
        :history="history"
        archived-count-message-key="storageSetting.history.entity.archivedCount"
        deleted-count-message-key="storageSetting.history.entity.deletedCount"
        hide-deleted="hide-deleted"
      />
    </template>
    <v-checkbox
      v-field="form.with_dependencies"
      :label="$t('storageSetting.entity.archiveDependencies')"
      color="primary"
    >
      <template #append="">
        <c-help-icon
          :text="$t('storageSetting.entity.archiveDependenciesHelp')"
          color="info"
          max-width="300"
          top="top"
        />
      </template>
    </v-checkbox>
    <v-flex>
      <v-btn
        class="ma-0 mb-4"
        color="primary"
        @click="$emit('archive')"
      >
        {{ $t('storageSetting.entity.archiveDisabled') }}
      </v-btn>
    </v-flex>
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
};
</script>
