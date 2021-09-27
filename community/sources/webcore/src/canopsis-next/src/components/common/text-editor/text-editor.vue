<template lang="pug">
  div
    span.theme--light.v-label.text-editor__label.mb-2(
      v-show="label"
    ) {{ label }}
    div.text-editor(:class="{ 'error--text': hasError }", @blur="$emit('blur', $event)")
      div(ref="textEditor")
      div.text-editor__details
        div.v-messages.theme--light.error--text
          div.v-messages__wrapper
            div.v-messages__message(v-for="errorMessage in errorMessages") {{ errorMessage }}
</template>

<script>
import { isString } from 'lodash';
import { Jodit } from 'jodit';

import 'jodit/build/jodit.min.css';

import { BASE_URL, FILE_BASE_URL, LOCAL_STORAGE_ACCESS_TOKEN_KEY } from '@/config';

import localStorageService from '@/services/local-storage';

const {
  modules: {
    Dom,
    Ajax,
    Widget: { FileSelectorWidget },
  },
} = Jodit;

/**
 * We need to replace this method to avoid the problem with CORS
 */
const originalSend = Ajax.prototype.send;

Ajax.prototype.send = function send(...args) {
  delete this.options.headers['X-REQUESTED-WITH'];
  this.options.withCredentials = false;

  return originalSend.call(this, ...args);
};

export default {
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
      default: () => [],
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
  },
  data() {
    return {
      editor: null,
    };
  },
  computed: {
    hasError() {
      return this.errorMessages.length;
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

      if (this.extraButtons.length) {
        config.extraButtons = this.extraButtons;
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
        defaultHandlerError: this.uploaderDefaultHandlerError,
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
    value(newValue) {
      if (this.editor.value !== newValue) {
        this.editor.setEditorValue(newValue);
      }
    },
  },
  mounted() {
    this.editor = new Jodit(this.$refs.textEditor, this.editorConfig);
    this.editor.setEditorValue(this.value);
    this.editor.events.on('change', this.onChange);
  },
  beforeDestroy() {
    this.editor.events.off('change', this.onChange);
    this.editor.destruct();

    delete this.editor;
  },
  methods: {
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
      return response.files.filter(file => file.error).join(' ');
    },

    /**
     * Uploader error handler
     *
     * @param {Object} err
     */
    uploaderError(err) {
      this.editor.events.fire('errorPopap', [err, 'error', 7000]);
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
          const elm = this.editor.create.inside.element(tagName);

          elm.setAttribute(attr, attrValue);

          if (tagName === 'a' && file.filename) {
            elm.setAttribute('target', '_blank');

            elm.innerText = file.filename;
          }

          if (tagName === 'img') {
            this.editor.selection.insertImage(elm, null, this.editor.options.imageDefaultWidth);
          } else {
            this.editor.selection.insertNode(elm);
          }
        });
      }
    },

    /**
     * Uploader default handler for error
     *
     * @param {Object} response
     */
    uploaderDefaultHandlerError(response) {
      this.editor.events.fire('errorPopap', [this.editor.options.getMessage(response)]);
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

      const imgElements = current instanceof HTMLDocument ? [...current.querySelectorAll('img')] : [];
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

<style>
.jodit_fullsize_box {
  z-index: 100000 !important;
}
</style>

<style lang="scss" scoped>
.text-editor {
  &__label {
    font-size: .85em;
    display: block;
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

  &.error--text /deep/ .jodit_container {
    margin-bottom: 8px;

    .jodit_workplace {
      border-color: currentColor;
    }
  }
}
</style>
