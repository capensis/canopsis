<template>
  <c-information-block
    :title="$t('storageSetting.pbehavior.title')"
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
      :label="$t('storageSetting.pbehavior.deleteAfter')"
      :help-text="$t('storageSetting.pbehavior.deleteAfterHelpText')"
      :name="pbehaviorDeleteAfterFieldName"
    />
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
      type: Number,
      required: false,
    },
  },
  computed: {
    pbehaviorDeleteAfterFieldName() {
      return 'pbehavior.delete_after';
    },
  },
};
</script>
