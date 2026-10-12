<template>
  <c-layout column>
    <c-name-field
      v-model="name"
      :label="t('common.name')"
      :max-length="FORM_MAX_STRING_LENGTH"
      name="name"
      required
    />

    <c-layout class="mt-3" column>
      <comment-template-fields-form />
    </c-layout>
  </c-layout>
</template>

<script>
import { FORM_MAX_STRING_LENGTH } from '@/constants';

import { useFormModels } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';

import CommentTemplateFieldsForm from './comment-template-fields-form.vue';

export default {
  components: { CommentTemplateFieldsForm },

  setup() {
    const { t } = useI18n();
    const { name } = useFormModels(['name']);

    return {
      t,
      FORM_MAX_STRING_LENGTH,
      name,
    };
  },
};
</script>
