<template>
  <v-layout column>
    <p class="font-italic grey--text my-3">
      {{ $t('linkRule.sourceCodeAlert') }}
    </p>
    <java-script-code-editor
      v-field="value"
      :completions="completions"
      class="java-script-code-editor"
      resettable
    />
  </v-layout>
</template>

<script>
import { LINK_RULE_ADVANCED_COMPLETIONS } from '@/constants';

import { useValidationHeader } from '@/hooks/validator/validation-header';

import JavaScriptCodeEditor from './fields/javascript-code-editor.vue';

export default {
  components: { JavaScriptCodeEditor },
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
  setup() {
    const { hasAnyError } = useValidationHeader();

    return {
      hasAnyError,

      completions: LINK_RULE_ADVANCED_COMPLETIONS,
    };
  },
};
</script>

<style lang="scss" scoped>
.java-script-code-editor {
  min-height: 400px;
}
</style>
