<template>
  <v-layout column>
    <c-name-field
      v-field="form.value"
      :disabled="isImported"
      :max-length="maxTagNameLength"
      name="value"
      required
    />
    <c-color-picker-field v-field="form.color" />
    <tag-patterns-form
      v-if="!isImported"
      v-field="form.patterns"
    />
  </v-layout>
</template>

<script>
import { MAX_TAG_NAME_LENGTH } from '@/constants';

import TagPatternsForm from './tag-patterns-form.vue';

export default {
  components: { TagPatternsForm },
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
    maxTagNameLength: {
      type: Number,
      default: MAX_TAG_NAME_LENGTH,
    },
  },
};
</script>
