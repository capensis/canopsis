import { setField } from '@/helpers/immutable';

import { useComponentModel } from '../vue';

/**
 * A composable that provides functionality to update a model field and emit events to the parent component.
 * This hook is designed to work with v-field directive and Vue's component model pattern.
 * It handles both full model updates and individual field updates within the model.
 *
 * @param {Object} props - The component's props object containing the model data
 * @param {Function} emit - Vue's emit function to trigger events to the parent component
 * @returns {Object} An object containing methods and properties for model manipulation
 * @property {string} modelEvent - The event name to be emitted (usually 'update:modelValue' or custom v-model
 *                                event)
 * @property {string} modelProp - The prop name that contains the model data (usually 'modelValue' or custom
 *                               v-model prop)
 * @property {Function} updateModel - Function to update the entire model
 * @property {Function} updateField - Function to update a specific field within the model
 */
export const useModelField = (props, emit) => {
  const { event, prop } = useComponentModel();

  /**
   * Updates the entire model and emits the change event to the parent component
   *
   * @param {*} value - The new value for the model
   * @param {...*} rest - Additional arguments to pass to the emit function
   * @returns {*} The updated model value
   */
  const updateModel = (value, ...rest) => {
    emit(event, value, ...rest);

    return value;
  };

  /**
   * Updates a specific field within the model and emits the change event to the parent component
   *
   * @param {string|Array<string|number>} fieldName - The field name or path array to update
   * @param {*} value - The new value for the specified field
   * @param {...*} rest - Additional arguments to pass to the emit function
   * @returns {*} The updated model value
   */
  const updateField = (fieldName, value, ...rest) => (
    updateModel(setField(props[prop], fieldName, value), ...rest)
  );

  return {
    modelEvent: event,
    modelProp: prop,

    updateModel,
    updateField,
  };
};
