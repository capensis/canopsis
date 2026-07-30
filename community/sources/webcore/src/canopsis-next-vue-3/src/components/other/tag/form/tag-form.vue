<template>
  <c-form-general-patterns-tabs :form="form" :with-template-testing="false">
    <template #general="{ setRef }">
      <tag-general-form
        :ref="setRef"
        :form="form"
        :is-imported="isImported"
        :is-new="isNew"
        :max-tag-name-length="maxTagNameLength"
        @update:form="updateModel"
      />
    </template>

    <template v-if="!isImported" #patterns="{ setRef }">
      <tag-patterns-form
        :ref="setRef"
        :form="form.patterns"
        @update:form="updateField('patterns', $event)"
      />
    </template>
  </c-form-general-patterns-tabs>
</template>

<script>
import { MAX_TAG_NAME_LENGTH } from '@/constants';

import { useModelField } from '@/hooks/form/model-field';

import TagGeneralForm from './tag-general-form.vue';
import TagPatternsForm from './tag-patterns-form.vue';

export default {
  components: {
    TagGeneralForm,
    TagPatternsForm,
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

  emits: ['update:form'],

  setup(props, { emit }) {
    const { updateModel, updateField } = useModelField(props, emit);

    return {
      updateModel,
      updateField,
    };
  },
};
</script>
