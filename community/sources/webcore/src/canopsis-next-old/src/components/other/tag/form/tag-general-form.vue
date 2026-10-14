<template>
  <v-layout class="gap-3" column>
    <c-name-field
      v-field="form.value"
      :disabled="isImported || !isNew"
      :max-length="maxTagNameLength"
      name="value"
      autofocus
      required
    />
    <c-form-block>
      <c-form-block-row :label="$t('common.color')" align-center>
        <c-color-picker-field v-field="form.color" />
      </c-form-block-row>
    </c-form-block>
  </v-layout>
</template>

<script>
import { MAX_TAG_NAME_LENGTH } from '@/constants';

export default {
  model: {
    prop: 'form',
    event: 'input',
  },
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
};
</script>
