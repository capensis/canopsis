<template>
  <v-layout class="gap-3" column>
    <c-information-block :title="$tc('metaAlarmRule.valuePath', 2)">
      <text-fields
        v-field="form.value_paths"
        :label="$tc('metaAlarmRule.valuePath', 1)"
        :error="hasValuePathsErrors"
        :name="name"
        required
        @input="validateValuePaths"
      />
    </c-information-block>
  </v-layout>
</template>

<script>
import { computed, nextTick, onBeforeUnmount } from 'vue';

import TextFields from '@/components/forms/fields/text-fields.vue';
import CInformationBlock from '@/components/common/block/c-information-block.vue';

import { useInjectValidator } from '@/hooks/form/useValidationChildren';
import { useValidationAttachRequired } from '@/hooks/form/useValidationAttachRequired';

export default {
  inject: ['$validator'],
  components: { CInformationBlock, TextFields },
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      default: () => ({}),
    },
    name: {
      type: String,
      default: 'value_paths',
    },
  },
  setup(props) {
    const validator = useInjectValidator();
    const {
      attachRequiredRule,
      detachRequiredRule,
      validateRequiredRule,
    } = useValidationAttachRequired(props.name);

    const hasValuePathsErrors = computed(() => validator.errors.has(props.name));

    const validateValuePaths = () => {
      nextTick(validateRequiredRule);
    };

    const isValuePathsExist = () => props.form.value_paths && props.form.value_paths.length > 0;

    attachRequiredRule(isValuePathsExist);
    onBeforeUnmount(detachRequiredRule);

    return {
      hasValuePathsErrors,
      validateValuePaths,
    };
  },
};
</script>
