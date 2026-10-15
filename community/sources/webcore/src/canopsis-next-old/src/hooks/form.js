import { useComponentModel } from './vue';

/**
 * Hook for update model value
 *
 * @param {Object} props
 * @param {Function} emit
 * @return {{ updateModel: function }}
 */
export const useModelField = (props, emit) => { // TODO: remove it
  const { event } = useComponentModel();

  const updateModel = (value) => {
    emit(event, value);

    return value;
  };

  return { updateModel };
};
