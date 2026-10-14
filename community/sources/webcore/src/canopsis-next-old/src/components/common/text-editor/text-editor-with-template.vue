<template>
  <v-layout column>
    <c-widget-template-field
      :value="form.template"
      :templates="templates"
      @input="updateTemplate"
    />
    <text-editor-field
      v-validate="rules"
      :value="form.text"
      :label="label"
      :error-messages="errors.collect('text')"
      :variables="variables"
      :sanitize-options="sanitizeOptions"
      :dark="$system.dark"
      name="text"
      @input="updateText"
    />
  </v-layout>
</template>

<script>
import { CUSTOM_WIDGET_TEMPLATE } from '@/constants';

import { useModelField } from '@/hooks/form/model-field';

import TextEditorField from './text-editor.vue';

export default {
  inject: ['$validator', '$system'],

  components: { TextEditorField },

  model: {
    prop: 'form',
    event: 'input',
  },

  props: {
    form: {
      type: Object,
      default: () => ({}),
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
  },

  setup(props, { emit }) {
    const { updateModel, updateField } = useModelField(props, emit);

    /**
     * Handles template selection and updates the model with the selected template and its content.
     *
     * @param {Object} param - The template update payload.
     * @param {string} param.value - The selected template value.
     * @param {string} param.content - The content associated with the selected template.
     */
    const updateTemplate = ({ value, content }) => {
      if (value === props.form.template) {
        return;
      }

      updateModel({
        template: value,
        text: content,
      });
    };

    /**
     * Handles text changes. If a custom template is used, updates both text and template fields in the model.
     * Otherwise, updates only the text field.
     *
     * @param {string} text - The new text value.
     */
    const updateText = (text) => {
      if (props.form.template !== CUSTOM_WIDGET_TEMPLATE && text !== props.form.text) {
        updateModel({
          text,
          template: CUSTOM_WIDGET_TEMPLATE,
        });

        return;
      }

      updateField('text', text);
    };

    return {
      updateTemplate,
      updateText,
    };
  },
};
</script>
