<template>
  <c-layout gap="2" column>
    <c-simple-tooltip
      v-if="!hideAckResources"
      :content="t('alarm.ackResourcesQuestion')"
      top
    >
      <template #activator="{ props }">
        <c-checkbox
          v-bind="props"
          v-model="ackResources"
          :label="t('alarm.ackResources')"
          color="primary"
        />
      </template>
    </c-simple-tooltip>
    <c-description-field
      v-model="comment"
      :label="t('common.note')"
      :required="isNoteRequired"
      :max-length="255"
      name="comment"
      autofocus
    />
  </c-layout>
</template>

<script>
import { useFormModels } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';

export default {
  props: {
    isNoteRequired: {
      type: Boolean,
      default: false,
    },
    hideAckResources: {
      type: Boolean,
      default: false,
    },
  },

  setup() {
    const { t } = useI18n();
    const { comment, ackResources } = useFormModels({
      comment: 'comment',
      ackResources: 'ack_resources',
    });

    return {
      t,
      comment,
      ackResources,
    };
  },
};
</script>
