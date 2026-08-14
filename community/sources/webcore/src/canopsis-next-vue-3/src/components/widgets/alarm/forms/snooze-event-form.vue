<template>
  <c-layout gap="2" column>
    <c-duration-field
      v-model:duration="duration"
      autofocus
      required
    />
    <c-description-field
      v-model="comment"
      :label="t('common.note')"
      :required="isNoteRequired"
      :max-length="255"
      name="comment"
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
  },

  setup() {
    const { t } = useI18n();
    const { duration, comment } = useFormModels({
      duration: 'duration',
      comment: 'comment',
    });

    return {
      t,
      duration,
      comment,
    };
  },
};
</script>
