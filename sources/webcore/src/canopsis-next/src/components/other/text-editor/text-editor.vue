<template lang="pug">
  quill-editor(:value="value", @input="$emit('input', $event)", :disabled="disabled")
</template>

<script>
import { quillEditor as QuillEditor } from 'vue-quill-editor';

export default {
  components: {
    QuillEditor,
  },
  props: {
    value: {
      type: String,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
  },
};
</script>

