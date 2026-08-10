<template>
  <c-layout column>
    <c-widget-template-field
      :value="currentForm.template"
      :templates="templates"
      @update:model-value="updateTemplate"
    />

    <text-editor-field
      :model-value="currentForm.text"
      :label="label"
      :variables="variables"
      :sanitize-options="sanitizeOptions"
      :required="required"
      name="text"
      @update:model-value="updateText"
    />
  </c-layout>
</template>

<script>
import { computed } from 'vue';

import TextEditorField from '@/components/forms/fields/text-editor-field.vue';
import { CUSTOM_WIDGET_TEMPLATE } from '@/constants';

export default {
  components: { TextEditorField },

  props: {
    modelValue: {
      type: Object,
      default: undefined,
    },
    form: {
      type: Object,
      default: undefined,
    },
    templates: {
      type: Array,
      default: () => [],
    },
    variables: {
      type: Array,
      default: () => [],
    },
    label: {
      type: String,
      required: false,
    },
    rules: {
      type: Object,
      required: false,
    },
    sanitizeOptions: {
      type: Object,
      required: false,
    },
    required: {
      type: Boolean,
      default: false,
    },
  },

  emits: ['update:modelValue'],

  setup(props, { emit }) {
    const currentForm = computed(() => props.modelValue ?? props.form ?? {});

    const updateModel = (value) => {
      emit('update:modelValue', value);
    };

    /**
     * Handles template selection and updates the model with the selected template and its content.
     *
     * @param {Object} template - The template update payload.
     */
    const updateTemplate = (template) => {
      if (template.value === currentForm.value.template) {
        return;
      }

      updateModel({
        template: template.value,
        text: template.content,
      });
    };

    /**
     * Handles text changes. If a custom template is used, updates both text and template fields in the model.
     *
     * @param {string} text - The new text value.
     */
    const updateText = (text) => {
      if (
        currentForm.value.template !== CUSTOM_WIDGET_TEMPLATE &&
        text !== currentForm.value.text
      ) {
        updateModel({
          text,
          template: CUSTOM_WIDGET_TEMPLATE,
        });

        return;
      }

      updateModel({
        ...currentForm.value,
        text,
      });
    };

    return {
      currentForm,
      updateTemplate,
      updateText,
    };
  },
};
</script>
