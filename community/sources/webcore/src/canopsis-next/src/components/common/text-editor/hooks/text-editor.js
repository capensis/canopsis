import { isString, isFunction } from 'lodash';
import { ref, unref, computed, onBeforeUnmount } from 'vue';
import { Ajax } from 'jodit/esm/core/request';
import { ajaxInstances } from 'jodit/esm/modules/uploader/helpers/send';

import { FILE_BASE_URL, LOCAL_STORAGE_ACCESS_TOKEN_KEY } from '@/config';

import localStorageService from '@/services/local-storage';

import { matchPayloadVariableBySelection } from '@/helpers/payload-json';
import { hasAtLeastOneVariable } from '@/helpers/variables';

import { useI18n } from '@/hooks/i18n';

// eslint-disable-next-line import/no-webpack-loader-syntax
import VariablesIcon from '!!svg-inline-loader?modules!@/assets/images/variables.svg';

/**
 * Hook for managing text editor upload functionality.
 *
 * Provides upload configuration and handlers for file/image uploads in Jodit editor.
 *
 * @param {Object} params - The parameters object.
 * @param {Ref<Object>} params.editor - Reference to the Jodit editor instance.
 * @param {boolean} params.isPublic - Whether uploads should be public.
 * @param {number} [params.maxFileSize] - Maximum file size in bytes.
 *
 * @returns {Object} An object containing upload-related configuration and options.
 * @returns {ComputedRef<Object>} return.controlsOptions - Control options for file/image tags.
 * @returns {ComputedRef<Object>} return.uploaderOptions - Uploader configuration options.
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
   * @param {Array} request - Array of files to upload.
   * @param {Function} showProgress - Callback function to update upload progress (0-100).
   * @returns {Promise<Object>} Promise that resolves with upload response or error.
   */
  const customUploadFunction = (request, showProgress) => {
    const { uploader } = unref(editor);
    const { fileValidator: uploaderFileValidator } = uploader.options ?? {};

    if (uploaderFileValidator) {
      try {
        request.forEach(uploaderFileValidator);
      } catch (err) {
        return Promise.resolve({
          error: err.message,
        });
      }
    }

    const ajax = new Ajax({
      xhr: () => {
        const xhr = new XMLHttpRequest();
        if (uploader.j.ow.FormData !== undefined && xhr.upload) {
          showProgress(10);
          xhr.upload.addEventListener('progress', (event) => {
            if (event.lengthComputable) {
              let percentComplete = event.loaded / event.total;
              percentComplete *= 100;
              showProgress(percentComplete);
            }
          }, false);
        } else {
          showProgress(100);
        }
        return xhr;
      },
      method: uploader.o.method || 'POST',
      data: request,
      url: isFunction(uploader.o.url)
        ? uploader.o.url(request)
        : uploader.o.url,
      headers: uploader.o.headers,
      queryBuild: uploader.o.queryBuild,
      contentType: uploader.o.contentType.call(uploader, request),
      withCredentials: uploader.o.withCredentials || false,
    });

    let instances = ajaxInstances.get(uploader);

    if (!instances) {
      instances = new Set();
      ajaxInstances.set(uploader, instances);
    }

    instances.add(ajax);
    uploader.j.e.one('beforeDestruct', ajax.destruct);

    return ajax
      .send()
      .then(resp => resp.json())
      .catch(error => ({
        error,
      }))
      .finally(() => {
        ajax.destruct();

        if (instances) {
          instances.delete(ajax);
        }
      });
  };

  /**
   * Checks if the upload response indicates success.
   *
   * @param {Array} response - Upload response array.
   * @returns {boolean} True if response has items, false otherwise.
   */
  const isSuccess = response => response?.length;

  /**
   * Extracts error message from upload response.
   *
   * @param {Object|Array} response - Upload response object or array.
   * @returns {string} Error message if present, empty string otherwise.
   */
  const getMessage = response => response?.error || response?.filter?.(file => file.error).join(' ');

  /**
   * Processes upload response and formats it for the editor.
   *
   * @param {Object|Array} response - Upload response from server.
   * @returns {Object} Processed response with files, baseurl, error, and msg.
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
   * Inserts uploaded files (images or links) into the editor.
   *
   * @param {Object} response - Upload response containing files array and baseurl.
   */
  const defaultHandlerSuccess = (response) => {
    if (response.files && response.files.length) {
      response.files.forEach((file) => {
        const [tagName, attr] = file.mediatype && file.mediatype.startsWith('image')
          ? ['img', 'src']
          : ['a', 'href'];

        const attrValue = isString(file) ? file : response.baseurl + file._id;
        const unwrappedEditor = unref(editor);
        const elm = unwrappedEditor.createInside.element(tagName);

        elm.setAttribute(attr, attrValue);

        if (tagName === 'a' && file.filename) {
          elm.setAttribute('target', '_blank');

          elm.innerText = file.filename;
        }

        if (tagName === 'img') {
          unwrappedEditor.selection.insertImage(elm, null, unwrappedEditor.options.imageDefaultWidth);
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
 * Provides functionality for inserting and managing variables in the text editor,
 * including variable selection, menu positioning, and variable insertion.
 *
 * @param {Object} params - The parameters object.
 * @param {Ref<Object>} params.editor - Reference to the Jodit editor instance.
 * @param {Array} params.variables - Array of available variables.
 *
 * @returns {Object} An object containing variables-related state and methods.
 * @returns {ComputedRef<boolean>} return.hasVariables - Whether variables are available.
 * @returns {Ref<boolean>} return.variablesShown - Whether the variables menu is visible.
 * @returns {Ref<string>} return.variablesMenuValue - Current selected variable value.
 * @returns {Ref<Object>} return.variablesMenuPosition - Position of the variables menu.
 * @returns {ComputedRef<Object>} return.variablesButton - Button configuration for variables.
 * @returns {Function} return.pasteVariable - Function to paste a variable into the editor.
 * @returns {Function} return.closeVariablesMenu - Function to close the variables menu.
 */
export const useTextEditorVariables = ({ editor, variables }) => {
  const variablesShown = ref(false);
  const variablesMenuValue = ref('');
  const variablesMenuPosition = ref({
    x: 0,
    y: 0,
  });

  const hasVariables = computed(() => hasAtLeastOneVariable(unref(variables)));

  /**
   * Extracts variable value from a variable group match result.
   *
   * @param {Array} variableGroup - Variable group array from matchPayloadVariableBySelection.
   * @returns {string} Extracted variable value.
   */
  const getVariableValueFromGroup = (variableGroup) => {
    const [,, content] = variableGroup;

    const parts = content.trim().split(' ');

    return parts.length > 1 ? parts[1] : parts[0];
  };

  /**
   * Selects variable value at the current cursor position.
   *
   * Updates variablesMenuValue and adjusts selection to match the complete variable.
   */
  const selectVariableValueByCursor = () => {
    const selection = unref(editor).selection?.sel ?? {};
    const { anchorNode, anchorOffset, focusOffset } = selection;

    if (!anchorNode) {
      return;
    }

    const [selectionStart, selectionEnd] = [anchorOffset, focusOffset].sort();

    const variableGroup = matchPayloadVariableBySelection(anchorNode.nodeValue, selectionStart, selectionEnd);

    if (!variableGroup) {
      variablesMenuValue.value = undefined;
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
  };

  /**
   * Shows the variables menu at the specified position.
   *
   * @param {Object} instance - Jodit editor instance.
   * @param {HTMLElement} target - Target element that triggered the menu.
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
   * Closes the variables menu and removes selection change listener.
   */
  const closeVariablesMenu = () => {
    variablesShown.value = false;

    document.removeEventListener('selectionchange', selectVariableValueByCursor);
  };

  /**
   * Extracts the variable value from a variable string format.
   *
   * @param {string} variableString - Variable string in format "{{ value }}".
   * @returns {string} Extracted variable value without braces and spaces.
   */
  const extractVariableValue = variableString => variableString.replace('{{ ', '').replace(' }}', '');

  /**
   * Determines the text to insert based on whether replacing an existing variable.
   *
   * @param {string} variable - New variable string in format "{{ value }}".
   * @param {Array|null} variableGroup - Existing variable group match, or null if no match.
   * @returns {string} Text to insert into the editor.
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
   * Inserts a text node into the specified range and positions cursor after it.
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
   * Pastes a variable into the editor at the current selection.
   *
   * Replaces existing variable if one is selected, otherwise inserts new variable.
   * Preserves formatting and spaces around the insertion point.
   *
   * @param {string} variable - Variable string in format "{{ value }}" to insert.
   */
  const pasteVariable = (variable) => {
    selectVariableValueByCursor();
    const selection = unref(editor).selection?.sel ?? {};
    const { focusNode } = selection;

    if (!focusNode) {
      return;
    }

    const { anchorOffset, focusOffset } = selection;
    const [selectionStart, selectionEnd] = [anchorOffset, focusOffset].sort();

    const variableGroup = matchPayloadVariableBySelection(focusNode.nodeValue, selectionStart, selectionEnd);
    const textToInsert = getTextToInsert(variable, variableGroup);

    let range;

    if (selection.rangeCount > 0) {
      range = selection.getRangeAt(0);
    } else {
      range = document.createRange();
      range.setStart(focusNode, selectionStart);
      range.setEnd(focusNode, selectionEnd);
    }

    insertTextNode(range, textToInsert);
    selection.removeAllRanges();
    selection.addRange(range);
    unref(editor).selection.selectRange(range);

    closeVariablesMenu();
  };

  const variablesButton = computed(() => ({
    name: 'variables',
    mode: 3,
    exec: showVariablesMenu,
  }));

  const variablesExtraIcon = computed(() => ({
    variables: `<i class="material-icons v-icon v-icon--small" style="width: 18px; height: 18px;">${VariablesIcon}</i>`,
  }));

  onBeforeUnmount(() => document.removeEventListener('selectionchange', selectVariableValueByCursor));

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
