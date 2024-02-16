<template>
  <v-layout column>
    <p class="font-italic grey--text my-3">
      {{ $t('linkRule.sourceCodeAlert') }}
    </p>
    <java-script-code-editor
      class="java-script-code-editor"
      v-field="value"
      resettable
    />
  </v-layout>
</template>

<script>
import { formValidationHeaderMixin } from '@/mixins/form/validation-header';

import JavaScriptCodeEditor from './fields/javascript-code-editor.vue';

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

<style lang="scss" scoped>
.java-script-code-editor {
  min-height: 400px;
}
</style>
