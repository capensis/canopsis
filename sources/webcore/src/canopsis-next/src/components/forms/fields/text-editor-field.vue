<template lang="pug">
  div
    text-editor(
      v-if="!focusable",
      v-field="value",
      :label="label",
      :buttons="buttons",
      :extraButtons="extraButtons",
      :config="config",
      :errorMessages="errorMessages"
    )
    div(v-else)
      text-editor(
        ref="textEditor",
        v-if="focused",
        v-field="value",
        v-click-outside.same="{ handler: toggleFocusedOff, closeConditional }",
        :label="label",
        :buttons="buttons",
        :extraButtons="extraButtons",
        :config="config",
        :errorMessages="errorMessages"
      )
      text-editor-blurred(
        v-else,
        :value="value",
        :label="label",
        @click="toggleFocusedOn"
      )
</template>

<script>
import TextEditor from '@/components/other/text-editor/text-editor.vue';
import TextEditorBlurred from '@/components/other/text-editor/text-editor-blurred.vue';

export default {
  components: {
    TextEditor,
    TextEditorBlurred,
  },
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: String,
      default: '',
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
    focusable: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      focused: !this.focusable,
    };
  },
  methods: {
    closeConditional(e) {
      return this.focused && (this.$refs.textEditor && !this.$refs.textEditor.$el.contains(e.target));
    },

    toggleFocusedOn() {
      this.focused = true;

      this.$nextTick(() => {
        if (this.$refs.textEditor) {
          this.$refs.textEditor.editor.selection.focus();
        }
      });
    },

    toggleFocusedOff() {
      this.focused = false;
    },
  },
};
</script>
