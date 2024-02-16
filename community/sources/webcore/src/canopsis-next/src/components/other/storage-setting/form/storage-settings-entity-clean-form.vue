<template>
  <c-information-block
    :title="$t('storageSetting.entityArchived.title')"
    :help-text="$t('storageSetting.entityArchived.titleHelp')"
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
        hide-archived
      />
    </template>
    <v-flex>
      <v-btn
        class="my-4"
        color="primary"
        @click="$emit('clean')"
      >
        {{ $t('storageSetting.entityArchived.cleanArchive') }}
      </v-btn>
    </v-flex>
  </c-information-block>
</template>

<script>
import StorageSettingsHistoryMessage from '../partials/storage-settings-history-message.vue';

export default {
  components: { StorageSettingsHistoryMessage },
  props: {
    history: {
      type: Object,
      required: false,
    },
  },
};
</script>
