export const flowchartTextEditorMixin = {
  data() {
    return {
      editing: false,
    };
  },
  methods: {
    enableEditingMode() {
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
