import { setField } from '@/helpers/immutable';

import { useComponentModel } from '../vue';

/**
 * Provides functionality to update a model field and emit an event to the parent component.
 * It offers two methods: `updateModel` for updating the entire model,
 * and `updateField` for updating a specific field within the model.
 *
 * @param {Object} props - The props object of the component, expected to contain the model property.
 * @param {Function} emit - The function to emit events to the parent component.
 * @returns {{
 *  updateModel: (value: *) => *,
 *  updateField: (fieldName: string, value: *) => *
 * }}
 */
export const useModelField = (props, emit) => {
  const { event, prop } = useComponentModel();

  /**
   * Update full model
   *
   * @param {*} value
   * @return {Array|Object}
   */
  const updateModel = (value) => {
    emit(event, value);

    return value;
  };

  /**
   * Emit event to parent with new object and with updated field
   *
   * @param {string|Array} fieldName
   * @param {*} value
   */
  const updateField = (fieldName, value) => updateModel(setField(props[prop], fieldName, value));

  return {
    updateModel,
    updateField,
  };
};
