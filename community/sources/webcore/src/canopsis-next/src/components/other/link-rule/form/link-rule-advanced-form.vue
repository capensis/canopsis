<template>
  <c-layout gap="3" column>
    <p class="font-italic text-medium-emphasis">
      {{ $t('linkRule.sourceCodeAlert') }}
    </p>
    <java-script-code-editor
      v-model:value="sourceCodeField"
      :completions="completions"
      class="java-script-code-editor"
      resettable
      @update:value="$emit('change')"
    />
    <v-messages v-if="!!errorMessages.length" :messages="errorMessages" color="error" />
  </c-layout>
</template>

<script>
import { computed } from 'vue';

import {
  LINK_RULE_ADVANCED_ALARM_COMPLETIONS,
  LINK_RULE_ADVANCED_ENTITY_COMPLETIONS,
  LINK_RULE_TYPES,
} from '@/constants';

import { useFormModel } from '@/hooks/form';
import { useTemplateVars } from '@/hooks/store/modules/template-vars';
import { useValidationHeader } from '@/hooks/validator/validation-header';

import JavaScriptCodeEditor from '@/components/common/code-editor/javascript-code-editor.vue';

export default {
  components: { JavaScriptCodeEditor },
  props: {
    type: {
      type: String,
      default: LINK_RULE_TYPES.alarm,
    },
    errorMessages: {
      type: Array,
      default: () => [],
    },
  },

  emits: ['change'],

  setup(props) {
    const sourceCodeField = useFormModel('source_code');
    const { hasAnyError } = useValidationHeader();
    const { templateVars } = useTemplateVars();

    const completions = computed(() => ({
      ...(props.type === LINK_RULE_TYPES.alarm
        ? LINK_RULE_ADVANCED_ALARM_COMPLETIONS
        : LINK_RULE_ADVANCED_ENTITY_COMPLETIONS),
      env: templateVars?.value,
    }));

    return {
      sourceCodeField,
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
