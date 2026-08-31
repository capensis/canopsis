import { computed, onBeforeUnmount, ref, unref } from 'vue';

import { FILE_BASE_URL, LOCAL_STORAGE_ACCESS_TOKEN_KEY } from '@/config';

import localStorageService from '@/services/local-storage';

import { matchPayloadVariableBySelection } from '@/helpers/payload-json';
import { hasAtLeastOneVariable } from '@/helpers/variables';

import { useI18n } from '@/hooks/i18n';

/**
 * Hook for managing text editor upload functionality.
 *
 * @param {Object} params - The parameters object.
 * @param {Ref<Object>} params.editor - Reference to the Jodit editor instance.
 * @param {boolean} params.isPublic - Whether uploads should be public.
 * @param {number} [params.maxFileSize] - Maximum file size in bytes.
 * @returns {Object} Upload control and uploader options.
 */
export const useTextEditorUpload = ({ editor, isPublic, maxFileSize }) => {
  const { t } = useI18n();

  const controlsOptions = computed(() => ({
    file: {
      tags: ['a'],
    },
    image: {
      tags: ['img'],
    },
  }));

  /**
   * Custom upload function for handling file uploads with progress tracking.
   *
   * @param {FormData|File[]} request - Files to upload.
   * @param {Function} showProgress - Callback function that updates upload progress.
   * @returns {Promise<Object>}
   */
  const customUploadFunction = (request, showProgress) => new Promise((resolve) => {
      const { uploader } = unref(editor);
      const { fileValidator: uploaderFileValidator } = uploader.options ?? {};
      const formData = request instanceof FormData ? request : new FormData();

      if (!(request instanceof FormData)) {
        request.forEach(file => formData.append('files[]', file));
      }

      if (uploaderFileValidator) {
        try {
          Array.from(formData.values()).forEach(uploaderFileValidator);
        } catch (err) {
          resolve({ error: err.message });
          return;
        }
      }

      const xhr = new XMLHttpRequest();

      xhr.upload.addEventListener('progress', (event) => {
        if (event.lengthComputable) {
          showProgress((event.loaded / event.total) * 100);
        }
      });

      xhr.addEventListener('load', () => {
        try {
          resolve(JSON.parse(xhr.responseText));
        } catch (err) {
          resolve({ error: err.message });
        }
      });

      xhr.addEventListener('error', () => resolve({ error: xhr.statusText }));
      xhr.open('POST', `${FILE_BASE_URL}?public=${unref(isPublic)}`);
      xhr.setRequestHeader(
        'Authorization',
        `Bearer ${localStorageService.get(LOCAL_STORAGE_ACCESS_TOKEN_KEY)}`,
      );
      xhr.send(formData);
    });

  /**
   * Checks if the upload response indicates success.
   *
   * @param {Array} response - Upload response array.
   * @returns {boolean}
   */
  const isSuccess = response => response?.length;

  /**
   * Extracts error message from upload response.
   *
   * @param {Object|Array} response - Upload response object or array.
   * @returns {string}
   */
  const getMessage = response => response?.error || response?.filter?.(file => file.error).join(' ');

  /**
   * Processes upload response and formats it for the editor.
   *
   * @param {Object|Array} response - Upload response from server.
   * @returns {Object}
   */
  const process = (response) => {
    const files = response.filter(file => !file.error);

    return {
      files,
      baseurl: `${FILE_BASE_URL}/`,
      error: response.error,
      msg: response.msg,
    };
  };

  /**
   * Validates file size against maximum allowed size.
   *
   * @param {File} file - File object to validate.
   * @throws {Error} Throws error if file exceeds maximum size.
   */
  const fileValidator = (file) => {
    const maxSize = unref(maxFileSize);
    if (!maxSize) {
      return;
    }

    if (file instanceof File && file.size > maxSize) {
      throw new Error(t('validation.messages.size', [null, maxSize / 1024]));
    }
  };

  /**
   * Default success handler for file uploads.
   *
   * @param {Object} response - Upload response containing files array and baseurl.
   */
  const defaultHandlerSuccess = (response) => {
    if (response.files && response.files.length) {
      response.files.forEach((file) => {
        const [tagName, attr] =
          file.mediatype && file.mediatype.startsWith('image') ? ['img', 'src'] : ['a', 'href'];

        const attrValue = response.baseurl + file._id;
        const unwrappedEditor = unref(editor);
        const elm = unwrappedEditor.createInside.element(tagName);

        elm.setAttribute(attr, attrValue);

        if (tagName === 'a' && file.filename) {
          elm.setAttribute('target', '_blank');
          elm.innerText = file.filename;
        }

        if (tagName === 'img') {
          unwrappedEditor.selection.insertImage(
            elm,
            null,
            unwrappedEditor.options.imageDefaultWidth,
          );
        } else {
          unwrappedEditor.selection.insertNode(elm);
        }
      });
    }
  };

  const uploaderOptions = computed(() => ({
    enableDragAndDropFileToEditor: true,
    insertImageAsBase64URI: false,
    format: 'json',
    url: `${FILE_BASE_URL}?public=${unref(isPublic)}`,
    headers: { Authorization: `Bearer ${localStorageService.get(LOCAL_STORAGE_ACCESS_TOKEN_KEY)}` },
    customUploadFunction,
    getMessage,
    isSuccess,
    process,
    defaultHandlerSuccess,
    fileValidator,
  }));

  return {
    controlsOptions,
    uploaderOptions,
  };
};

/**
 * Hook for managing text editor variables menu functionality.
 *
 * @param {Object} params - The parameters object.
 * @param {Ref<Object>} params.editor - Reference to the Jodit editor instance.
 * @param {Array} params.variables - Array of available variables.
 * @returns {Object} Variables state and methods.
 */
export const useTextEditorVariables = ({ editor, variables }) => {
  const variablesShown = ref(false);
  const variablesMenuValue = ref('');
  const variablesMenuPosition = ref({
    x: 0,
    y: 0,
  });
  const lastRange = ref(null);

  const hasVariables = computed(() => hasAtLeastOneVariable(unref(variables)));

  const isEditorNode = node => node && unref(editor)?.editor?.contains(node);

  const saveSelectionRange = () => {
    const selection = unref(editor).selection?.sel;

    if (!selection?.rangeCount || !isEditorNode(selection.anchorNode)) {
      return;
    }

    lastRange.value = selection.getRangeAt(0).cloneRange();
  };

  /**
   * Extracts variable value from a variable group match result.
   *
   * @param {Array} variableGroup - Variable group array from matchPayloadVariableBySelection.
   * @returns {string}
   */
  const getVariableValueFromGroup = (variableGroup) => {
    const [, , content] = variableGroup;
    const parts = content.trim().split(' ');

    return parts.length > 1 ? parts[1] : parts[0];
  };

  /**
   * Selects variable value at the current cursor position.
   */
  const selectVariableValueByCursor = () => {
    const selection = unref(editor).selection?.sel ?? {};
    const { anchorNode, anchorOffset, focusOffset } = selection;

    if (!isEditorNode(anchorNode)) {
      return;
    }

    const [selectionStart, selectionEnd] = [anchorOffset, focusOffset].sort();
    const variableGroup = matchPayloadVariableBySelection(
      anchorNode.nodeValue,
      selectionStart,
      selectionEnd,
    );

    if (!variableGroup) {
      variablesMenuValue.value = undefined;
      saveSelectionRange();
      return;
    }

    const [variable] = variableGroup;
    variablesMenuValue.value = `{{ ${getVariableValueFromGroup(variableGroup)} }}`;

    const [currentStart, currentEnd] = [anchorOffset, focusOffset].sort();
    const start = variableGroup.index;
    const end = variableGroup.index + variable.length;

    if (currentStart !== start || currentEnd !== end) {
      const range = document.createRange();

      range.setStart(anchorNode, start);
      range.setEnd(anchorNode, end);

      unref(editor).selection.selectRange(range);
    }

    saveSelectionRange();
  };

  /**
   * Shows variables menu at button position.
   *
   * @param {Object} instance - Jodit editor instance.
   * @param {HTMLElement} target - Target element.
   * @param {Object} options - Options object.
   * @param {Event} options.originalEvent - Original DOM event.
   */
  const showVariablesMenu = (instance, target, { originalEvent: event }) => {
    selectVariableValueByCursor();

    const { left, top, height } = event.target.getBoundingClientRect();

    variablesMenuPosition.value = {
      x: left,
      y: top + height,
    };
    variablesShown.value = true;

    document.addEventListener('selectionchange', selectVariableValueByCursor);
  };

  /**
   * Closes variables menu and removes listeners.
   */
  const closeVariablesMenu = () => {
    variablesShown.value = false;

    document.removeEventListener('selectionchange', selectVariableValueByCursor);
  };

  /**
   * Extracts the variable value from a variable string format.
   *
   * @param {string} variableString - Variable string in format "{{ value }}".
   * @returns {string}
   */
  const extractVariableValue = variableString => variableString.replace('{{ ', '').replace(' }}', '');

  /**
   * Wraps bare variable values into template syntax.
   *
   * @param {string} variable - Variable string or bare variable name.
   * @returns {string}
   */
  const normalizeVariable = (variable) => {
    const trimmedVariable = String(variable ?? '').trim();

    return trimmedVariable.startsWith('{{') ? trimmedVariable : `{{ ${trimmedVariable} }}`;
  };

  /**
   * Determines the text to insert.
   *
   * @param {string} variable - New variable string.
   * @param {Array|null} variableGroup - Existing variable group match.
   * @returns {string}
   */
  const getTextToInsert = (variable, variableGroup) => {
    if (!variableGroup) {
      return variable;
    }

    const [oldVariable] = variableGroup;
    const oldValue = getVariableValueFromGroup(variableGroup);
    const newValue = extractVariableValue(variable);

    return oldVariable.replace(oldValue, newValue);
  };

  /**
   * Inserts a text node into range and positions cursor after it.
   *
   * @param {Range} range - DOM Range object where text should be inserted.
   * @param {string} text - Text content to insert.
   */
  const insertTextNode = (range, text) => {
    range.deleteContents();
    const textNode = document.createTextNode(text);
    range.insertNode(textNode);
    range.setStartAfter(textNode);
    range.collapse(true);
  };

  /**
   * Pastes a variable into editor selection.
   *
   * @param {string} variable - Variable string in format "{{ value }}".
   */
  const pasteVariable = (variable) => {
    selectVariableValueByCursor();
    const selection = unref(editor).selection?.sel ?? {};
    const hasEditorSelection = selection.rangeCount > 0 && isEditorNode(selection.focusNode);
    const range = hasEditorSelection ? selection.getRangeAt(0) : lastRange.value?.cloneRange();
    const focusNode = range?.startContainer;

    if (!range || !focusNode) {
      return;
    }

    const [selectionStart, selectionEnd] = [range.startOffset, range.endOffset].sort();
    const variableGroup = focusNode.nodeValue
      ? matchPayloadVariableBySelection(focusNode.nodeValue, selectionStart, selectionEnd)
      : null;
    const textToInsert = getTextToInsert(normalizeVariable(variable), variableGroup);

    insertTextNode(range, textToInsert);
    selection.removeAllRanges();
    selection.addRange(range);
    unref(editor).selection.selectRange(range);
    lastRange.value = range.cloneRange();

    closeVariablesMenu();
  };

  const variablesButton = computed(() => ({
    name: 'variables',
    mode: 3,
    exec: showVariablesMenu,
  }));

  const variablesExtraIcon = computed(() => ({
    variables:
      '<i class="material-icons v-icon v-icon--small" style="width: 14px; height: 14px;">data_object</i>',
  }));

  onBeforeUnmount(() => document.removeEventListener('selectionchange', selectVariableValueByCursor),
  );

  return {
    hasVariables,
    variablesShown,
    variablesMenuValue,
    variablesMenuPosition,
    variablesButton,
    pasteVariable,
    variablesExtraIcon,
    closeVariablesMenu,
  };
};
