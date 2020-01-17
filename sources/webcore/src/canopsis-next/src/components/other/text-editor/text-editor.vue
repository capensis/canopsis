<template lang="pug">
  div.text-editor(:class="{ 'error--text': hasError }")
    div(ref="textEditor")
    div.text-editor__details
      div.v-messages.theme--light.error--text
        div.v-messages__wrapper
          div.v-messages__message(v-for="errorMessage in errorMessages") {{ errorMessage }}
</template>

<script>
import { Jodit } from 'jodit';

import 'jodit/build/jodit.min.css';

export default {
  props: {
    value: {
      type: String,
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
        uploader: {
          insertImageAsBase64URI: true,
        },

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
