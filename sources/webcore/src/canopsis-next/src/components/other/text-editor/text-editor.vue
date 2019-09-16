<template lang="pug">
  div
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
  },
  data() {
    return {
      editor: null,
    };
  },
  computed: {
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
    this.editor = new Jodit(this.$el, this.editorConfig);
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
