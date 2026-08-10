import { unref } from 'vue';

import { CUSTOM_WIDGET_TEMPLATE } from '@/constants';
import { useModelField } from '@/hooks/form/model-field';

/**
 * Hook for managing widget template field functionality.
 *
 * @param {Object} props - Component props containing template and model data
 * @param {string|import('vue').Ref<string>} valueKey - Key to extract value from the template object
 * @param {Function} emit - Vue emit function for model updates
 * @returns {{
 *   updateTemplate: Function,
 *   updateValue: Function,
 * }}
 */
export const useWidgetTemplateField = (props, valueKey, emit) => {
  const { updateModel } = useModelField(props, emit);

  /**
   * Emits an 'update:template' event with the new template and value.
   *
   * @param {string|Symbol} newTemplate
   * @param {*} newValue
   */
  const emitUpdateTemplate = (newTemplate, newValue) =>
    emit('update:template', newTemplate, newValue);

  /**
   * Updates both template and value from selected widget template.
   *
   * @param {Object} newTemplate
   */
  const updateTemplate = newTemplate =>
    emitUpdateTemplate(newTemplate?.value, newTemplate?.[unref(valueKey)]);

  /**
   * Updates model value and switches to custom template when needed.
   *
   * @param {*} newValue
   */
  const updateValue = newValue =>
    props.template === CUSTOM_WIDGET_TEMPLATE
      ? updateModel(newValue)
      : emitUpdateTemplate(CUSTOM_WIDGET_TEMPLATE, newValue);

  return {
    updateTemplate,
    updateValue,
  };
};
