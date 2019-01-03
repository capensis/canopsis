<template lang="pug">
  quill-editor(:value="value", @input="$emit('input', $event)")
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
  },
};
</script>

