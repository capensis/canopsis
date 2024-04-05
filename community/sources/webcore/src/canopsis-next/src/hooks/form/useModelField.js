import { setField } from '@/helpers/immutable';

import { useComponentModel } from '../vue';

/**
 * Hook for update model value
 *
 * @param {Object} props
 * @param {Function} emit
 * @return {{ updateModel: function }}
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
