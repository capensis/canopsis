import { escape } from 'lodash';
import { unref } from 'vue';

import { MODALS } from '@/constants';

import { useModals } from '@/hooks/modals';

const PROCESS_HTML_EVENT = 'processHTML';
const TEXT_EDITOR_PASTE_ACTIONS = {
  html: 'html',
  text: 'text',
  cancel: 'cancel',
};

/**
 * Converts pasted HTML into plain text without executing markup.
 *
 * @param {string} html - HTML fragment from the clipboard.
 * @returns {string}
 */
const convertHtmlToPlainText = (html) => {
  const { body } = new DOMParser().parseFromString(html, 'text/html');

  return body?.innerText || body?.textContent || '';
};

/**
 * Checks whether a Jodit instance can still accept selection or insert operations.
 *
 * @param {Object} instance - Jodit editor instance.
 * @returns {boolean}
 */
const isEditorAvailable = instance => Boolean(instance && !instance.isInDestruct);

/**
 * Restores the selection saved with `editor.s.save()` and refocuses the editor.
 *
 * @param {Object} instance - Jodit editor instance.
 * @returns {boolean}
 */
const restoreEditorSelection = (instance) => {
  if (!isEditorAvailable(instance)) {
    return false;
  }

  instance.s.focus();
  instance.s.restore();

  return true;
};

/**
 * Inserts clipboard HTML either as markup or as escaped plain text.
 *
 * @param {Object} instance - Jodit editor instance.
 * @param {string} html - Original pasted HTML.
 * @param {boolean} asHtml - Whether to keep HTML or insert plain text.
 */
const insertPastedContent = (instance, html, asHtml) => {
  if (!restoreEditorSelection(instance)) {
    return;
  }

  if (asHtml) {
    instance.s.insertHTML(html);
    return;
  }

  const plainText = convertHtmlToPlainText(html);

  instance.s.insertHTML(escape(plainText).replace(/\r\n|\r|\n/g, '<br>'));
};

/**
 * Hook that replaces Jodit's native HTML paste dialog with the application modal.
 *
 * @param {Object} params - The parameters object.
 * @param {Ref<Object>} params.editor - Reference to the Jodit editor instance.
 * @returns {{ pasteOptions: Object, handleProcessHTML: Function }} Paste options and processHTML handler.
 */
export const useTextEditorPaste = ({ editor }) => {
  const modals = useModals();

  let isPasteModalOpened = false;

  const pasteOptions = {
    askBeforePasteHTML: false,
    askBeforePasteFromWord: false,
  };

  /**
   * Intercepts HTML clipboard content, owns the paste, and opens the custom modal.
   *
   * @param {ClipboardEvent|DragEvent} event - Original paste/drop event.
   * @param {string} value - HTML-like clipboard value that triggered the event.
   * @param {{ html?: string, plain?: string, rtf?: string }} texts - Clipboard MIME payloads.
   * @returns {true}
   */
  const handleProcessHTML = (_event, value, texts) => {
    const instance = unref(editor);

    if (!isEditorAvailable(instance)) {
      return true;
    }

    instance.e.stopPropagation(PROCESS_HTML_EVENT);

    if (isPasteModalOpened) {
      return true;
    }

    const html = texts?.html || value;

    if (!html) {
      return true;
    }

    instance.s.save();
    isPasteModalOpened = true;

    let settled = false;

    const finish = (action) => {
      if (settled) {
        return;
      }

      settled = true;
      isPasteModalOpened = false;

      if (action === TEXT_EDITOR_PASTE_ACTIONS.html) {
        insertPastedContent(instance, html, true);
        return;
      }

      if (action === TEXT_EDITOR_PASTE_ACTIONS.text) {
        insertPastedContent(instance, html, false);
        return;
      }

      restoreEditorSelection(instance);
    };

    modals.show({
      name: MODALS.textEditorPaste,
      config: {
        html,
        action: finish,
        cancel: () => finish(TEXT_EDITOR_PASTE_ACTIONS.cancel),
      },
    });

    return true;
  };

  return {
    pasteOptions,
    handleProcessHTML,
  };
};
