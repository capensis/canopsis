<template lang="pug">
  v-layout(column)
    v-text-field(
      v-field="form.value",
      v-validate="valueRules",
      :disabled="isImported",
      :max-length="maxTagNameLength",
      :error-messages="errors.collect('value')",
      name="value"
    )
    c-color-picker-field(v-field="form.color")
    tag-patterns-form(v-if="!isImported", v-field="form.patterns")
</template>

<script>
import { MAX_TAG_NAME_LENGTH } from '@/constants';

import TagPatternsForm from './tag-patterns-form.vue';

export default {
  inject: ['$validator'],
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
  computed: {
    valueRules() {
      return {
        required: true,
        max: this.maxTagNameLength,
      };
    },
  },
};
</script>
