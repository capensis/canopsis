<template lang="pug">
  c-information-block(
    :title="$t('storageSetting.entity.title')",
    :help-text="$t('storageSetting.entity.titleHelp')",
    help-icon-color="info"
  )
    template(v-if="history", #subtitle="")
      storage-settings-history-message(
        :history="history",
        archived-count-message-key="storageSetting.history.entity.archivedCount",
        deleted-count-message-key="storageSetting.history.entity.deletedCount",
        hide-deleted
      )
    v-checkbox(
      v-field="form.with_dependencies",
      :label="$t('storageSetting.entity.archiveDependencies')",
      color="primary"
    )
      template(#append="")
        c-help-icon(:text="$t('storageSetting.entity.archiveDependenciesHelp')", color="info", max-width="300", top)
    v-flex
      v-btn.ma-0.mb-4(color="primary", @click="$emit('archive')") {{ $t('storageSetting.entity.archiveDisabled') }}
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
