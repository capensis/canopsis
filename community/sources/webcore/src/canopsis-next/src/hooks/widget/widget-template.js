import { unref } from 'vue';

import { CUSTOM_WIDGET_TEMPLATE } from '@/constants';

import { useModelField } from '@/hooks/form/model-field';

/**
 * Hook for managing widget template field functionality
 *
 * @param {Object} props - Component props containing template and model data
 * @param {string} props.template - Current template value
 * @param {Function} emit - Vue emit function for model updates
 * @returns {Object} Object containing template update functions
 * @property {Function} updateTemplate - Updates both template and columns
 * @property {Function} updateValue - Updates model based on template type
 */
export const useWidgetTemplateField = (props, valueKey, emit) => {
  const { updateModel } = useModelField(props, emit);

  /**
   * Updates both template and columns values
   *
   * @param {WidgetTemplate} newTemplate - New template value to set
   */
  const updateTemplate = newTemplate => updateModel(newTemplate?.value, newTemplate?.[unref(valueKey)]);

  /**
   * Updates model value based on template type
   * If template is CUSTOM_WIDGET_TEMPLATE, only updates value
   * Otherwise updates both template and value
   *
   * @param {Object} newValue - New value
   */
  const updateValue = newValue => (props.template === CUSTOM_WIDGET_TEMPLATE
    ? updateModel(newValue)
    : updateModel(CUSTOM_WIDGET_TEMPLATE, newValue));

  return {
    updateTemplate,
    updateValue,
  };
};
