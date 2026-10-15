<template>
  <c-layout gap="2" column>
    <c-duration-field
      v-model:duration="duration"
      required
    />
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
  },

  setup() {
    const { t } = useI18n();
    const { duration, comment } = useFormModels({
      duration: 'duration',
      comment: 'comment',
    });

    return {
      t,
      FORM_MAX_STRING_LENGTH,
      duration,
      comment,
    };
  },
};
</script>
