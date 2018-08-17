<template lang="pug">
  v-list-group
    v-list-tile(slot="activator") Weather modal template
    quill-editor(:value="value", @input="$emit('update:content', $event)")
</template>

<script>
import 'quill/dist/quill.core.css';
import 'quill/dist/quill.snow.css';
import 'quill/dist/quill.bubble.css';

import { quillEditor } from 'vue-quill-editor';

export default {
  components: {
    quillEditor,
  },
  props: {
    content: {
      type: String,
    },
    value: {
      type: String,
    },
  },
  data() {
    return {
      editorConfig: {
        modules: {
          toolbar: ['bold'],
        },
      },
    };
  },
};
</script>
