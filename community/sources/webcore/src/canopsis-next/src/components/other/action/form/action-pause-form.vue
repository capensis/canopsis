<template>
  <div>
    <action-note-form
      :model-value="modelValue"
      :name="name"
      :depth="depth"
      :template-vars="templateVars"
      @update:model-value="$emit('update:modelValue', $event)"
    />
    <c-form-block-row :label="t('common.type')" :depth="depth" indented>
      <c-text-field
        :model-value="modelValue.pause_type"
        :label="t('common.type')"
        :max-length="FORM_MAX_STRING_LENGTH"
        :name="`${name}.pause_type`"
        @update:model-value="updateField('pause_type', $event)"
      />
    </c-form-block-row>
    <c-form-block-row :label="t('settings.pause.propagate')" :depth="depth" indented>
      <v-switch
        :model-value="modelValue.inherited"
        :label="t('settings.pause.propagate')"
        color="primary"
        hide-details
        @update:model-value="updateField('inherited', $event)"
      />
    </c-form-block-row>
  </div>
</template>

<script>
import { FORM_MAX_STRING_LENGTH } from '@/constants';

import { useModelField } from '@/hooks/form/model-field';
import { useI18n } from '@/hooks/i18n';

import ActionNoteForm from './action-note-form.vue';

export default {
  components: { ActionNoteForm },

  props: {
    modelValue: {
      type: Object,
      default: () => ({}),
    },
    name: {
      type: String,
      default: 'parameters',
    },
    depth: {
      type: Number,
      default: 0,
    },
    templateVars: {
      type: Object,
      default: () => ({}),
    },
  },

  emits: ['update:modelValue'],

  setup(props, { emit }) {
    const { t } = useI18n();
    const { updateField } = useModelField(props, emit);

    return {
      t,
      updateField,
      FORM_MAX_STRING_LENGTH,
    };
  },
};
</script>
