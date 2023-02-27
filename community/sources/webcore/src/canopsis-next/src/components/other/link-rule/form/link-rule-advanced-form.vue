<template lang="pug">
  java-script-code-editor.fill-min-height(v-field="value")
</template>

<script>
import { formValidationHeaderMixin } from '@/mixins/form/validation-header';

import JavaScriptCodeEditor from './partials/javascript-code-editor.vue';

export default {
  components: { JavaScriptCodeEditor },
  mixins: [formValidationHeaderMixin],
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: String,
      default: '',
    },
  },
};
</script>
