<template>
  <c-layout gap="3" column>
    <c-name-field
      :model-value="form.value"
      :disabled="isImported || !isNew"
      :max-length="maxTagNameLength"
      name="value"
      autofocus
      required
      @update:model-value="updateField('value', $event)"
    />

    <c-form-block>
      <c-form-block-row :label="$t('common.color')" align-center>
        <c-color-picker-field
          :model-value="form.color"
          @update:model-value="updateField('color', $event)"
        />
      </c-form-block-row>
    </c-form-block>
  </c-layout>
</template>

<script>
import { MAX_TAG_NAME_LENGTH } from '@/constants';
import { useModelField } from '@/hooks/form/model-field';

export default {
  props: {
    form: {
      type: Object,
      default: () => ({}),
    },
    isImported: {
      type: Boolean,
      default: false,
    },
    isNew: {
      type: Boolean,
      default: false,
    },
    maxTagNameLength: {
      type: Number,
      default: MAX_TAG_NAME_LENGTH,
    },
  },

  emits: ['update:form'],

  setup(props, { emit }) {
    const { updateField } = useModelField(props, emit);

    return {
      updateField,
    };
  },
};
</script>
