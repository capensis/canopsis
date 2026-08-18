<template>
  <c-layout column>
    <rating-field
      :model-value="form.rating"
      class="justify-center"
      hover
      @update:model-value="updateField('rating', $event)"
    />
    <c-textarea
      :model-value="form.comment"
      :label="t('common.comment')"
      autofocus
      @update:model-value="updateField('comment', $event)"
    />
  </c-layout>
</template>

<script>
import { useModelField } from '@/hooks/form/model-field';
import { useI18n } from '@/hooks/i18n';

import RatingField from '@/components/forms/fields/rating-field.vue';

export default {
  components: { RatingField },

  props: {
    form: {
      type: Object,
      required: true,
    },
  },

  emits: ['update:form'],

  setup(props, { emit }) {
    const { t } = useI18n();
    const { updateField } = useModelField(props, emit);

    return {
      t,
      updateField,
    };
  },
};
</script>
