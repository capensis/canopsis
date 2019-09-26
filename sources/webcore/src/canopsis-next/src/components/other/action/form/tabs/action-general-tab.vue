<template lang="pug">
  div
    template(v-if="value.generalParameters.type === $constants.ACTION_TYPES.snooze")
      v-textarea(
        :value="value.snoozeParameters.message",
        @input="updateField('snoozeParameters.message', $event)",
        :label="$t('modals.createAction.fields.message')"
      )
      duration-field(
        :value="value.snoozeParameters.duration",
        @input="updateField('snoozeParameters.duration', $event)"
      )
    template(v-if="value.generalParameters.type === $constants.ACTION_TYPES.pbehavior")
      pbehavior-form(
        :form="value.pbehaviorParameters",
        @input="updateField('pbehaviorParameters', $event)",
        :author="$constants.ACTION_AUTHOR",
        :noFilter="true"
      )
</template>

<script>
import formValidationHeaderMixin from '@/mixins/form/validation-header';

import DurationField from '@/components/forms/fields/duration.vue';
import PbehaviorForm from '@/components/other/pbehavior/form/pbehavior-form.vue';

export default {
  inject: ['$validator'],
  components: {
    DurationField,
    PbehaviorForm,
  },
  mixins: [formValidationHeaderMixin],
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    value: {
      type: Object,
      required: true,
    },
  },
};
</script>
