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
import { computed } from 'vue';

import {
  LINK_RULE_ADVANCED_ALARM_COMPLETIONS,
  LINK_RULE_ADVANCED_ENTITY_COMPLETIONS,
  LINK_RULE_TYPES,
} from '@/constants';

import { useValidationHeader } from '@/hooks/validator/validation-header';
import { useTemplateVars } from '@/hooks/store/modules/template-vars';

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
    type: {
      type: String,
      default: LINK_RULE_TYPES.alarm,
    },
  },
  setup(props) {
    const { hasAnyError } = useValidationHeader();
    const { templateVars } = useTemplateVars();

    const completions = computed(() => ({
      ...(
        props.type === LINK_RULE_TYPES.alarm
          ? LINK_RULE_ADVANCED_ALARM_COMPLETIONS
          : LINK_RULE_ADVANCED_ENTITY_COMPLETIONS
      ),
      env: templateVars?.value,
    }));

    return {
      hasAnyError,
      completions,
    };
  },
};
</script>

<style lang="scss" scoped>
.java-script-code-editor {
  min-height: 400px;
}
</style>
