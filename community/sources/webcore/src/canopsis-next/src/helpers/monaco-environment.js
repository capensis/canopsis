import EditorWorker from 'monaco-editor/esm/vs/editor/editor.worker?worker';
import JsonWorker from 'monaco-editor/esm/vs/language/json/json.worker?worker';

/**
 * Configures Monaco web workers for Vite.
 * Monaco cannot resolve workers on its own without webpack's monaco plugin.
 *
 * @returns {void}
 */
export const setupMonacoEnvironment = () => {
  if (globalThis.MonacoEnvironment?.getWorker) {
    return;
  }

  globalThis.MonacoEnvironment = {
    getWorker: (_workerId, label) =>
      label === 'json' ? new JsonWorker() : new EditorWorker(),
  };
};
