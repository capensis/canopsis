<template>
  <c-information-block :title="$t('storageSetting.entityUnlinked.title')">
    <template v-if="history" #subtitle="">
      <storage-settings-history-message
        :history="history"
        archived-count-message-key="storageSetting.history.entity.archivedCount"
        deleted-count-message-key="storageSetting.history.entity.deletedCount"
        hide-deleted
      />
    </template>
    <c-enabled-duration-field
      :duration="form.archive_after"
      :label="$t('storageSetting.entityUnlinked.archiveAfter')"
      :suffix="$t('storageSetting.receivedFor')"
      :units="timeUnits"
      name="entity.archive_before"
      switcher
      hide-value-on-false
      @update:duration="updateField('archive_after', $event)"
    />
  </c-information-block>
</template>

<script>
import { computed } from 'vue';

import { AVAILABLE_TIME_UNITS } from '@/constants';

import { useModelField } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';

import StorageSettingsHistoryMessage from '../partials/storage-settings-history-message.vue';

export default {
  components: { StorageSettingsHistoryMessage },

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

  emits: ['update:form', 'archive'],

  setup(props, { emit }) {
    const { t } = useI18n();
    const { updateField } = useModelField(props, emit);

    const timeUnits = computed(() =>
      [
        AVAILABLE_TIME_UNITS.day,
        AVAILABLE_TIME_UNITS.week,
        AVAILABLE_TIME_UNITS.month,
        AVAILABLE_TIME_UNITS.year,
      ].map(({ value, text }) => ({
        value,
        text: t(text, props.form.archive_after.value),
      })),
    );

    return {
      timeUnits,
      updateField,
    };
  },
};
</script>
