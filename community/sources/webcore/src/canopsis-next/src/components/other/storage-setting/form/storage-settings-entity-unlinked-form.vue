<template lang="pug">
  c-information-block(
    :title="$t('storageSetting.entityUnlinked.title')",
    :help-text="$t('storageSetting.entityUnlinked.titleHelp')",
    help-icon-color="info"
  )
    template(v-if="history", #subtitle="")
      storage-settings-history-message(
        :history="history",
        archived-count-message-key="storageSetting.history.entity.archivedCount",
        deleted-count-message-key="storageSetting.history.entity.deletedCount",
        hide-deleted
      )
    v-layout(row, align-center)
      v-flex(xs5)
        span.v-label.text--secondary {{ $t('storageSetting.entityUnlinked.archiveBefore') }}
      v-flex(xs4)
        c-duration-field(
          v-field="form.archive_before",
          :units-label="$t('common.unit')",
          :units="timeUnits",
          :name="alarmArchiveAfterFieldName"
        )
    v-flex
      v-btn.ma-0.mb-4(
        :disabled="hasChildrenError",
        color="primary",
        @click="$emit('archive')"
      ) {{ $t('storageSetting.entityUnlinked.archiveUnlinked') }}
</template>

<script>
import { AVAILABLE_TIME_UNITS } from '@/constants';

import { validationChildrenMixin } from '@/mixins/form';

import StorageSettingsHistoryMessage from '../partials/storage-settings-history-message.vue';

export default {
  inject: ['$validator'],
  components: { StorageSettingsHistoryMessage },
  mixins: [validationChildrenMixin],
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
  computed: {
    timeUnits() {
      return [
        AVAILABLE_TIME_UNITS.day,
        AVAILABLE_TIME_UNITS.week,
        AVAILABLE_TIME_UNITS.month,
        AVAILABLE_TIME_UNITS.year,
      ].map(({ value, text }) => ({
        value,
        text: this.$tc(text, this.form.archive_before.value),
      }));
    },

    alarmArchiveAfterFieldName() {
      return 'entity_unlinked.archive_before';
    },
  },
};
</script>
