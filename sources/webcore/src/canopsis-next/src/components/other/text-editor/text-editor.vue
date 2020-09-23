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
import { Jodit } from 'jodit';

import 'jodit/build/jodit.min.css';

import { BASE_URL, API_BASE_URL, API_ROUTES } from '@/config';

const { modules: { Dom, Widget: { FileSelectorWidget } } } = Jodit;

const isJoditObject = jodit =>
  jodit && jodit instanceof Object && typeof jodit.constructor === 'function' && jodit.constructor.name === 'Jodit';

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
      uploadedFiles: [],
    };
  },
  computed: {
    hasError() {
      return this.errorMessages.length;
    },

    editorConfig() {
      const FILE_BASE_URL = `${API_BASE_URL}${API_ROUTES.file}`;

      const config = {
        language: this.$i18n.locale,
        toolbarSticky: false,
        controls: {
          file: {
            popup(editor, current, self, close) {
              const insert = (url, title = '') => {
                const linkElement = `<a href="${url}" title="${title}">${title || url}</a>`;

                editor.selection.insertNode(editor.create.inside.fromHTML(linkElement));
              };
              let sourceAnchor = null;
              if (current &&
                (current.nodeName === 'A' ||
                  Dom.closest(current, 'A', editor.editor))) {
                sourceAnchor =
                  current.nodeName === 'A'
                    ? current
                    : Dom.closest(current, 'A', editor.editor);
              }
              return FileSelectorWidget(editor, {
                filebrowser: (data) => {
                  if (data.files && data.files.length) {
                    let i;
                    for (i = 0; i < data.files.length; i += 1) {
                      insert(data.baseurl + data.files[i].id, data.files[i].fileName);
                    }
                  }
                  close();
                },
                upload: (data) => {
                  let i;
                  if (data.files && data.files.length) {
                    for (i = 0; i < data.files.length; i += 1) {
                      insert(data.baseurl + data.files[i].id, data.files[i].fileName);
                    }
                  }
                  close();
                },
                url: (url, text) => {
                  if (sourceAnchor) {
                    sourceAnchor.setAttribute('href', url);
                    sourceAnchor.setAttribute('title', text);
                  } else {
                    insert(url, text);
                  }
                  close();
                },
              }, sourceAnchor, close, false);
            },
            tags: ['a'],
            tooltip: 'Insert file',
          },
        },
        uploader: {
          enableDragAndDropFileToEditor: true,
          insertImageAsBase64URI: true,
          format: 'json',
          filesVariableName: 'files',
          url: FILE_BASE_URL,
          prepareData: (data) => {
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
                const pattern = `\\.${extension}`;
                const matches = file.name.match(new RegExp(pattern, 'g'));

                if (matches.length > 1) {
                  data.set(key, file, file.name.replace(new RegExp(`${pattern}$`), ''));
                }
              }
            }

            return data;
          },
          isSuccess: response => !response.error,
          process: (response) => {
            // TODO: move to methods
            this.uploadedFiles.push(...response[this.editor.options.uploader.filesVariableName]);

            return {
              files: response[this.editor.options.uploader.filesVariableName] || [],
              baseurl: `${FILE_BASE_URL}/`,
              error: response.error,
              msg: response.msg,
            };
          },
          error: (err) => {
            this.editor.events.fire('errorPopap', [err.getMessage(), 'error', 4000]);
          },
          defaultHandlerSuccess: (data, resp) => {
            if (resp.files && resp.files.length) {
              resp.files.forEach((file, index) => {
                const { id, fileName } = file;
                const [tagName, attr] =
                  resp.isImages && resp.isImages[index]
                    ? ['img', 'src']
                    : ['a', 'href'];

                const elm = this.editor.create.inside.element(tagName);

                elm.setAttribute(attr, resp.baseurl + id);

                if (tagName === 'a') {
                  elm.innerText = fileName;
                }

                if (isJoditObject(this.jodit)) {
                  if (tagName === 'img') {
                    this.editor.selection.insertImage(
                      elm,
                      null,
                      this.editor.options.imageDefaultWidth,
                    );
                  } else {
                    this.editor.selection.insertNode(elm);
                  }
                }
              });
            }
          },
          defaultHandlerError: (response) => {
            this.editor.events.fire('errorPopap', [this.options.getMessage(response)]);
          },
        },
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

    console.log(this.uploadedFiles);
    console.log(this.value);
  },
  methods: {
    onChange(value) {
      this.$emit('input', value);
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
