<template>
  <c-description-field
    v-model="commentField"
    :label="t('common.comment', 1)"
    :max-length="FORM_MAX_STRING_LENGTH"
    name="comment"
    autofocus
  />
</template>

<script>
import { FORM_MAX_STRING_LENGTH } from '@/constants';

import { useFormModels } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';

export default {
  setup() {
    const { t } = useI18n();
    const { commentField } = useFormModels({ commentField: 'comment' });

    return {
      t,
      FORM_MAX_STRING_LENGTH,
      commentField,
    };
  },
};
</script>
