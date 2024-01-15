<template>
  <div class="text-editor">
    <v-label v-show="label">
      {{ label }}
    </v-label>
    <div
      class="text-editor"
      :class="{ 'error--text': hasError }"
      @blur="$emit('blur', $event)"
    >
      <div ref="textEditor" />
      <variables-menu
        v-if="variables"
        :variables="variables"
        :visible="variablesShown"
        :value="variablesMenuValue"
        :position-x="variablesMenuPosition.x"
        :position-y="variablesMenuPosition.y"
        @input="pasteVariable"
        @close="closeVariablesMenu"
      />
      <div class="text-editor__details">
        <v-messages
          :value="errorMessages"
          color="error"
        />
      </div>
    </div>
  </div>
</template>

<script>
import { isString, get } from 'lodash';
import { Jodit } from 'jodit';

import 'jodit/build/jodit.min.css';

import { BASE_URL, FILE_BASE_URL, LOCAL_STORAGE_ACCESS_TOKEN_KEY } from '@/config';

import localStorageService from '@/services/local-storage';

import { sanitizeHtml } from '@/helpers/html';
import { matchPayloadVariableBySelection } from '@/helpers/payload-json';

import VariablesMenu from './variables-menu.vue';

const {
  modules: {
    Dom,
    Ajax,
    Widget: { FileSelectorWidget },
  },
} = Jodit;

/**
 * We need to replace this method to avoid the problem with CORS and to validate files
 */
const originalSend = Ajax.prototype.send;

Ajax.prototype.send = function send(...args) {
  try {
    const fileValidator = get(this, 'jodit.options.uploader.fileValidator');

    if (fileValidator) {
      this.options.data.forEach(fileValidator);
    }

    if (this.options.headers?.['X-REQUESTED-WITH']) {
      delete this.options.headers['X-REQUESTED-WITH'];
    }

    return originalSend.call(this, ...args);
  } catch (err) {
    return Promise.reject(err);
  }
};

export default {
  components: { VariablesMenu },
  props: {
    value: {
      type: String,
    },
    label: {
      type: String,
      default: '',
    },
    buttons: {
      type: Array,
      default: () => [
        'source', '|',
        'bold', 'italic', 'strikethrough', 'underline', '|',
        'ul', 'ol', '|',
        'font', 'fontsize', 'brush', 'paragraph', '|',
        'image', 'table', 'link', '|',
        'align', 'undo', 'redo', '|',
      ],
    },
    public: {
      type: Boolean,
      default: false,
    },
    extraButtons: {
      type: Array,
      default: () => [],
    },
    config: {
      type: Object,
      default: () => ({}),
    },
    errorMessages: {
      type: Array,
      default: () => [],
    },
    maxFileSize: {
      type: Number,
      required: false,
    },
    variables: {
      type: Array,
      required: false,
    },
    dark: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      sanitized: false,
      variablesShown: false,
      variablesMenuValue: '',
      variablesMenuPosition: {
        x: 0,
        y: 0,
      },
    };
  },
  computed: {
    sanitizedValue() {
      return sanitizeHtml(this.value);
    },

    hasError() {
      return this.errorMessages.length;
    },

    variablesButton() {
      return {
        name: 'variables',
        mode: 3,
        getContent: () => {
          const controlButton = document.createElement('span');

          controlButton.classList.add('text-editor__variables-button');
          controlButton.addEventListener('click', this.showVariablesMenu);

          return controlButton;
        },
      };
    },

    editorConfig() {
      const config = {
        language: this.$i18n.locale,
        toolbarSticky: false,
        controls: this.controlsConfig,
        uploader: this.uploaderConfig,
        sourceEditorCDNUrlsJS: [
          `${BASE_URL}scripts/libs/ace/1.3.3/ace.js`,
        ],
        beautifyHTMLCDNUrlsJS: [
          `${BASE_URL}scripts/libs/js-beautify/1.7.5/beautify.min.js`,
          `${BASE_URL}scripts/libs/js-beautify/1.7.5/beautify-html.min.js`,
        ],

        ...this.config,
      };

      if (this.buttons.length) {
        config.buttons = this.buttons;
        config.buttonsMD = this.buttons;
        config.buttonsSM = this.buttons;
        config.buttonsXS = this.buttons;
      }

      if (this.dark && !this.config.theme) {
        config.theme = 'dark';
      }

      config.extraButtons = [];

      if (this.variables) {
        config.extraButtons.push(this.variablesButton);
      }

      if (this.extraButtons.length) {
        config.extraButtons.push(...this.extraButtons);
      }

      return config;
    },

    uploaderConfig() {
      return {
        enableDragAndDropFileToEditor: true,
        insertImageAsBase64URI: false,
        format: 'json',
        filesVariableName: 'files',
        url: `${FILE_BASE_URL}?public=${this.public}`,
        headers: { Authorization: `Bearer ${localStorageService.get(LOCAL_STORAGE_ACCESS_TOKEN_KEY)}` },
        prepareData: this.uploaderPrepareData,
        isSuccess: this.uploaderIsSuccess,
        process: this.uploaderProcess,
        getMessage: this.uploaderGetMessage,
        error: this.uploaderError,
        defaultHandlerSuccess: this.uploaderDefaultHandlerSuccess,
        fileValidator: this.uploaderFileValidator,
      };
    },

    controlsConfig() {
      return {
        file: {
          popup: this.controlsFilePopup,
          tags: ['a'],
          tooltip: 'Insert file',
        },
        image: {
          popup: this.controlsImagePopup,
          tags: ['img'],
          tooltip: 'Insert image',
        },
      };
    },
  },
  watch: {
    editorConfig() {
      this.destroyJodit();
      this.createJodit();
    },
    value(newValue) {
      if (this.value && !this.sanitized) {
        this.$editor.setEditorValue(this.sanitizedValue);
        this.sanitized = true;
      } else if (this.$editor.value !== newValue) {
        this.$editor.setEditorValue(newValue);
      }
    },
  },
  mounted() {
    this.createJodit();
  },
  beforeDestroy() {
    this.destroyJodit();
  },
  methods: {
    createJodit() {
      this.$editor = new Jodit(this.$refs.textEditor, this.editorConfig);
      this.$editor.setEditorValue(this.sanitizedValue);
      this.$editor.events.on('change', this.onChange);
    },

    destroyJodit() {
      this.$editor.events.off('change', this.onChange);
      this.$editor.destruct();

      delete this.$editor;
    },

    selectVariableValueByCursor() {
      const selection = this.$editor.selection.sel;
      const { anchorNode, anchorOffset, focusOffset } = selection;

      if (!anchorNode) {
        return;
      }

      const [selectionStart, selectionEnd] = [anchorOffset, focusOffset].sort();
      const variableGroup = matchPayloadVariableBySelection(anchorNode.nodeValue, selectionStart, selectionEnd);

      if (!variableGroup) {
        this.variablesMenuValue = undefined;
        return;
      }

      const [variable] = variableGroup;
      this.variablesMenuValue = this.getVariableValueFromGroup(variableGroup);

      const [currentStart, currentEnd] = [anchorOffset, focusOffset].sort();
      const start = variableGroup.index;
      const end = variableGroup.index + variable.length;

      if (currentStart !== start || currentEnd !== end) {
        const range = document.createRange();

        range.setStart(anchorNode, start);
        range.setEnd(anchorNode, end);

        selection.removeAllRanges();
        selection.addRange(range);
      }
    },

    showVariablesMenu(event) {
      this.selectVariableValueByCursor();

      const { left, top, height } = event.target.getBoundingClientRect();

      this.variablesMenuPosition = {
        x: left,
        y: top + height,
      };
      this.variablesShown = true;

      document.addEventListener('selectionchange', this.selectVariableValueByCursor);
    },

    closeVariablesMenu() {
      this.variablesShown = false;

      document.removeEventListener('selectionchange', this.selectVariableValueByCursor);
    },

    getVariableValueFromGroup(group) {
      const [,, content] = group;

      const parts = content.trim().split(' ');

      return parts.length > 1 ? parts[1] : parts[0];
    },

    pasteVariable(variable) {
      this.selectVariableValueByCursor();

      const selection = this.$editor.selection.sel;
      const { anchorNode } = selection;

      const { anchorOffset, focusOffset } = selection;
      const [selectionStart, selectionEnd] = [anchorOffset, focusOffset].sort();

      const variableGroup = matchPayloadVariableBySelection(anchorNode.nodeValue, selectionStart, selectionEnd);

      if (variableGroup) {
        const [oldVariable] = variableGroup;
        const oldValue = this.getVariableValueFromGroup(variableGroup);

        this.$editor.selection.insertHTML(oldVariable.replace(oldValue, variable));
      } else {
        this.$editor.selection.insertHTML(`{{ ${variable} }}`);
      }

      this.closeVariablesMenu();
    },

    /**
     * Editor change event handler
     *
     * @param {string} value
     */
    onChange(value) {
      this.$emit('input', value);
    },

    /**
     * Prepare data for file upload
     *
     * @param {FormData} data
     * @returns {FormData}
     */
    uploaderPrepareData(data) {
      data.delete('path');
      data.delete('source');

      /**
       * There is fix for Jodit problem with doubling of extension
       */

      // eslint-disable-next-line no-restricted-syntax
      for (const [key, file] of data.entries()) {
        if (file instanceof File) {
          const mime = file.type.match(/\/([a-z0-9]+)/i);
          const extension = mime && mime[1] ? mime[1].toLowerCase() : '';

          if (extension) {
            const pattern = `\\.${extension}`;
            const matches = file.name.match(new RegExp(pattern, 'g'));

            if (matches.length > 1) {
              data.set(key, file, file.name.replace(new RegExp(`${pattern}$`), ''));
            }
          }
        }
      }

      return data;
    },

    /**
     * Is success checker for uploader
     *
     * @param response
     * @returns {boolean}
     */
    uploaderIsSuccess(response) {
      return response.length;
    },

    /**
     * Process handler for uploader
     *
     * @param {Object} response
     * @returns {{msg: *, baseurl: string, files: (*|*[]), error: *}}
     */
    uploaderProcess(response) {
      const files = response.filter(file => !file.error);

      return {
        files,
        baseurl: `${FILE_BASE_URL}/`,
        error: response.error,
        msg: response.msg,
      };
    },

    /**
     * Uploader get message handler
     *
     * @param {Object} response
     * @return {string}
     */
    uploaderGetMessage(response) {
      return response.filter(file => file.error).join(' ');
    },

    /**
     * Uploader error handler
     *
     * @param {Object} err
     */
    uploaderError(err) {
      this.$editor.events.fire('errorMessage', err.message, 'error', 7000);
    },

    /**
     * Uploader default handler for success
     *
     * @param {Object} response
     */
    uploaderDefaultHandlerSuccess(response) {
      if (response.files && response.files.length) {
        response.files.forEach((file) => {
          const [tagName, attr] = file.mediatype && file.mediatype.startsWith('image')
            ? ['img', 'src']
            : ['a', 'href'];

          const attrValue = isString(file) ? file : response.baseurl + file._id;
          const elm = this.$editor.create.inside.element(tagName);

          elm.setAttribute(attr, attrValue);

          if (tagName === 'a' && file.filename) {
            elm.setAttribute('target', '_blank');

            elm.innerText = file.filename;
          }

          if (tagName === 'img') {
            this.$editor.selection.insertImage(elm, null, this.$editor.options.imageDefaultWidth);
          } else {
            this.$editor.selection.insertNode(elm);
          }
        });
      }
    },

    /**
     * File size validator
     *
     * @param {File} file
     */
    uploaderFileValidator(file) {
      if (!this.maxFileSize) {
        return;
      }

      if (file instanceof File && file.size > this.maxFileSize) {
        throw new Error(this.$t('validation.messages.size', [null, this.maxFileSize / 1024]));
      }
    },

    /**
     * File control popup
     *
     * @param {Object} editor
     * @param {HTMLDocument|HTMLElement} current
     * @param {Object} self
     * @param {Function} close
     * @returns {Object}
     */
    controlsFilePopup(editor, current, self, close) {
      /**
       * Insert link into editor selection
       *
       * @param {string} url
       * @param {string} [title = '']
       */
      const insertLink = (url, title = '') => {
        const linkElement = `<a href="${url}" title="${title}" target="_blank">${title || url}</a>`;

        editor.selection.insertNode(editor.create.inside.fromHTML(linkElement));
      };

      /**
       * filebrowser and upload handler for file popup control
       *
       * @param {Object} data
       */
      const uploadHandler = ({ baseurl, files = [] } = {}) => {
        for (let i = 0; i < files.length; i += 1) {
          const file = files[i];
          const url = baseurl + file._id;

          insertLink(url, file.filename);
        }

        close();
      };

      const isLink = current.nodeName === 'A';
      let sourceAnchor = null;

      if (current && (isLink || Dom.closest(current, 'A', editor.editor))) {
        sourceAnchor = isLink ? current : Dom.closest(current, 'A', editor.editor);
      }

      return FileSelectorWidget(editor, {
        filebrowser: uploadHandler,
        upload: uploadHandler,
        url: (url, text) => {
          if (sourceAnchor) {
            sourceAnchor.setAttribute('target', '_blank');
            sourceAnchor.setAttribute('href', url);
            sourceAnchor.setAttribute('title', text);
          } else {
            insertLink(url, text);
          }
          close();
        },
      }, sourceAnchor, close, false);
    },

    /**
     * Image control popup
     *
     * @param {Object} editor
     * @param {HTMLDocument|HTMLElement} current
     * @param {Object} self
     * @param {Function} close
     * @returns {Object}
     */
    controlsImagePopup(editor, current, self, close) {
      /**
       * filebrowser and upload handler for image popup control
       *
       * @param {Object} data
       */
      const uploadHandler = async ({ baseurl, files = [] } = {}) => {
        for (let i = 0; i < files.length; i += 1) {
          const file = files[i];
          const url = baseurl + file._id;

          // eslint-disable-next-line no-await-in-loop
          await editor.selection.insertImage(url, null, editor.options.imageDefaultWidth);
        }

        close();
      };

      const imgElements = current instanceof HTMLDocument
        ? [...current.querySelectorAll('img')]
        : [];
      const isImage = current.tagName === 'IMG';
      let sourceImage = null;

      if (current && current.nodeType !== Node.TEXT_NODE && (isImage || imgElements.length)) {
        sourceImage = isImage ? current : imgElements[0];
      }

      return FileSelectorWidget(editor, {
        filebrowser: uploadHandler,
        upload: uploadHandler,
        url: async (url, text) => {
          const image = sourceImage || editor.create.inside.element('img');

          image.setAttribute('src', url);
          image.setAttribute('alt', text);

          if (!sourceImage) {
            await editor.selection.insertImage(image, null, editor.options.imageDefaultWidth);
          }

          close();
        },
      }, sourceImage, close);
    },
  },
};
</script>

<style lang="scss">
.jodit_fullsize_box {
  z-index: 100000 !important;
}

.text-editor {
  &__variables-button {
    display: flex;
    width: 100%;
    height: 100%;

    &:after {
      content: '(x)';
      display: block;
      width: 100%;
      height: 100%;
      color: black;
      font-size: 15px;

      .jodit_dark_theme & {
        color: silver;
      }
    }
  }
}
</style>

<style lang="scss" scoped>
.text-editor {
  .v-label {
    font-size: .85em;
    display: block;
    margin-bottom: 8px;
  }

  &__details {
    display: -webkit-box;
    display: -ms-flexbox;
    display: flex;
    -webkit-box-flex: 1;
    -ms-flex: 1 0 auto;
    flex: 1 0 auto;
    max-width: 100%;
    overflow: hidden;
  }

  &.error--text ::v-deep .jodit_container {
    margin-bottom: 8px;

    .jodit_workplace {
      border-color: currentColor;
    }
  }

  & ::v-deep .jodit_toolbar_popup {
    z-index: 25;
  }

  & ::v-deep .jodit_error_box_for_messages .error {
     color: white;
  }
}
</style>
