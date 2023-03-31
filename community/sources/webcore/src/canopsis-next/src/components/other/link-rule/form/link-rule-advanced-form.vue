<template lang="pug">
  v-layout(column)
    p.font-italic.grey--text.my-3 {{ $t('linkRule.sourceCodeAlert') }}
    java-script-code-editor.java-script-code-editor(v-field="value", resettable)
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

<style lang="scss" scoped>
.java-script-code-editor {
  min-height: 400px;
}
</style>
