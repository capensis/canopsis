<template lang="pug">
  v-list-group
    v-list-tile(slot="activator") Weather entity template
    quill-editor(:config="editorConfig")
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
