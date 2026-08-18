import { useFormContext } from '@/hooks/form';

/**
 * Creates a writable model for a form value without registering a field.
 *
 * @param {string|import('vue').Ref<string>|import('vue').ComputedRef<string>} path
 * @return {import('vue').WritableComputedRef<*>}
 */
export const useFormModel = path => useFormContext().model(path);
