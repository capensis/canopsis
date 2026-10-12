import { unref, watch, onMounted, onBeforeUnmount } from 'vue';

import { registerJavaScriptCompletion } from '@/helpers/monaco';

/**
 * A Vue composition function that registers JavaScript completions for a Monaco editor instance.
 * This function sets up the completion provider when the component is mounted and disposes of it before unmounting.
 *
 * @param {Object} params - The parameters for the function.
 * @param {Ref<Object>} params.codeEditor - A Vue ref object containing the Monaco editor instance.
 * @param {Ref<Object>} params.completions - A Vue ref object containing the completions to be registered.
 *
 * The function uses Vue's lifecycle hooks to manage the registration and disposal of the completion provider.
 * It unwraps the `codeEditor` and `completions` refs to access the underlying objects.
 * If the `codeEditor` has a `$monaco` property and `completions` are provided, it registers the completions
 * using the `registerJavaScriptCompletion` helper function.
 * The registered completions are disposed of when the component is about to be unmounted.
 */
export const useJavaScriptCompletions = ({ codeEditor, completions }) => {
  let registeredCompletions;

  /**
   * Registers JavaScript completions for the Monaco editor instance.
   *
   * This function unwraps the `codeEditor` and `completions` Vue refs to access the underlying objects.
   * If the `codeEditor` has a `$monaco` property and `completions` are provided, it registers the completions
   * using the `registerJavaScriptCompletion` helper function.
   *
   * The registered completions are stored in the `registeredCompletions` variable for later disposal.
   */
  const register = () => {
    const unwrappedCodeEditor = unref(codeEditor);
    const unwrappedCompletions = unref(completions);

    if (unwrappedCodeEditor.$monaco && unwrappedCompletions) {
      registeredCompletions = registerJavaScriptCompletion(unwrappedCodeEditor.$monaco, unwrappedCompletions);
    }
  };

  /**
   * Unregisters the JavaScript completions from the Monaco editor instance.
   *
   * This function checks if `registeredCompletions` is defined and has a `dispose` method.
   * If so, it calls the `dispose` method to clean up and remove the registered completions.
   */
  const unregister = () => registeredCompletions?.dispose?.();

  watch(completions, () => {
    unregister();
    register();
  });

  onMounted(register);

  onBeforeUnmount(unregister);
};
