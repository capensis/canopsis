<template>
  <c-layout column>
    <c-enabled-field
      :model-value="modelValue.forward_author"
      :label="t('scenario.forwardAuthor')"
      :name="`${name}.forward_author`"
      class="mt-0"
      @update:model-value="updateField('forward_author', $event)"
    />

    <v-expand-transition>
      <c-payload-text-field
        v-if="!modelValue.forward_author"
        :value="modelValue.author"
        :label="t('common.author')"
        :name="`${name}.author`"
        :variables="variables"
        @update:value="updateField('author', $event)"
      />
    </v-expand-transition>
  </c-layout>
</template>

<script>
import { useModelField } from '@/hooks/form/model-field';
import { useI18n } from '@/hooks/i18n';

export default {
  props: {
    modelValue: {
      type: Object,
      default: () => ({}),
    },
    name: {
      type: String,
      default: 'parameters',
    },
    variables: {
      type: Array,
      default: () => [],
    },
  },

  emits: ['update:modelValue'],

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
