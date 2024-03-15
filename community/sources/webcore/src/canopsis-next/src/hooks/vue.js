import { getCurrentInstance } from 'vue';

/**
 * Hook for get current component instance
 *
 * @return {Vue}
 */
export const useComponentInstance = () => {
  const instance = getCurrentInstance();

  return instance?.proxy || instance;
};

/**
 * Hook for get current component options
 *
 * @return {Object}
 */
export const useComponentOptions = () => {
  const instance = useComponentInstance();

  return instance.$options;
};

/**
 * Hook for get current component model
 *
 * @return {Object}
 */
export const useComponentModel = () => {
  const { model } = useComponentOptions();

  return {
    event: model?.event ?? 'input',
    prop: model?.event ?? 'value',
  };
};
