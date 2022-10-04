export const flowchartTextEditorMixin = {
  props: {
    readonly: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      editing: false,
    };
  },
  methods: {
    enableEditingMode() {
      if (this.readonly) {
        return;
      }

      this.editing = true;

      if (this.$refs.editor) {
        this.$refs.editor.focus();
      }
    },

    disableEditingMode(event) {
      this.$emit('update', {
        text: event.target.innerHTML,
      });

      this.editing = false;
    },
  },
};
