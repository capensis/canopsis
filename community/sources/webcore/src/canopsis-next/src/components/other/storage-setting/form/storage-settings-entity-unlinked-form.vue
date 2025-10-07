<template>
  <c-information-block :title="$t('storageSetting.entityUnlinked.title')">
    <template
      v-if="history"
      #subtitle=""
    >
      <storage-settings-history-message
        :history="history"
        archived-count-message-key="storageSetting.history.entity.archivedCount"
        deleted-count-message-key="storageSetting.history.entity.deletedCount"
        hide-deleted
      />
    </template>
    <c-enabled-duration-field
      v-field="form.archive_after"
      :label="$t('storageSetting.entityUnlinked.archiveAfter')"
      :suffix="$t('storageSetting.receivedFor')"
      name="entity_unlinked.archive_before"
      switcher
      hide-value-on-false
    />
  </c-information-block>
</template>

<script>
import { computed } from 'vue';

import { AVAILABLE_TIME_UNITS } from '@/constants';

import { useI18n } from '@/hooks/i18n';

import StorageSettingsHistoryMessage from '../partials/storage-settings-history-message.vue';

export default {
  inject: ['$validator'],
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
  setup(props) {
    const { tc } = useI18n();

    const timeUnits = computed(() => [
      AVAILABLE_TIME_UNITS.day,
      AVAILABLE_TIME_UNITS.week,
      AVAILABLE_TIME_UNITS.month,
      AVAILABLE_TIME_UNITS.year,
    ].map(({ value, text }) => ({
      value,
      text: tc(text, props.form.archive_before.value),
    })));

    return {
      timeUnits,
    };
  },
};
</script>
