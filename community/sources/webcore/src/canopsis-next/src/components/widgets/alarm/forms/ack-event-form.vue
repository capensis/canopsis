<template>
  <c-layout gap="2" column>
    <c-tooltip
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
    </c-tooltip>
    <c-description-field
      v-model="comment"
      :label="t('common.note')"
      :required="isNoteRequired"
      :max-length="FORM_MAX_STRING_LENGTH"
      name="comment"
      autofocus
    />
  </c-layout>
</template>

<script>
import { FORM_MAX_STRING_LENGTH } from '@/constants';

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
      FORM_MAX_STRING_LENGTH,
      comment,
      ackResources,
    };
  },
};
</script>
