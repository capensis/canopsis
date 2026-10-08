<template>
  <c-form-general-patterns-tabs>
    <template #general="{ setRef }">
      <tag-general-form
        v-field="form"
        :ref="setRef"
        :is-imported="isImported"
        :is-new="isNew"
        :max-tag-name-length="maxTagNameLength"
      />
    </template>

    <template v-if="!isImported" #patterns="{ setRef }">
      <tag-patterns-form
        v-field="form.patterns"
        :ref="setRef"
      />
    </template>
  </c-form-general-patterns-tabs>
</template>

<script>
import { MAX_TAG_NAME_LENGTH } from '@/constants';

import TagGeneralForm from './tag-general-form.vue';
import TagPatternsForm from './tag-patterns-form.vue';

export default {
  components: { TagGeneralForm, TagPatternsForm },
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
