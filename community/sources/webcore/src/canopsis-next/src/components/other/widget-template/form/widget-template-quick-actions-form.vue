<template>
  <v-layout column>
    <span class="text-body-2 my-2">{{ $tc('common.column', 2) }}</span>
    <v-flex xs12>
      <v-alert
        :value="errors.has(name)"
        color="info"
      >
        {{ $t('widgetTemplate.errors.quickActionsRequired') }}
      </v-alert>
    </v-flex>
    <quick-alarm-actions-form
      v-field="form.quickActions"
      :massive="massive"
      @input="validate"
    />
  </v-layout>
</template>

<script>
import { useValidationAttachRequiredForField } from '@/hooks/validator/validation-attach-required';

import QuickAlarmActionsForm from '@/components/common/actions-panel/quick-alarm-actions-form.vue';

export default {
  inject: ['$validator'],
  components: {
    QuickAlarmActionsForm,
  },
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      required: true,
    },
    massive: {
      type: Boolean,
      default: false,
    },
  },
  setup(props) {
    const name = 'quick-actions';

    const {
      asyncValidateRequiredRule: validate,
    } = useValidationAttachRequiredForField(name, () => !!props.form.quickActions?.length);

    return {
      name,
      validate,
    };
  },
};
</script>
