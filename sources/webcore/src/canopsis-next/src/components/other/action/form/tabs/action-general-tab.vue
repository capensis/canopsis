<template lang="pug">
  div
    template(v-if="form.generalParameters.type === $constants.ACTION_TYPES.snooze")
      v-textarea(
        :value="form.snoozeParameters.message",
        @input="updateField('snoozeParameters.message', $event)",
        :label="$t('modals.createAction.fields.message')"
      )
      duration-field(
        :value="form.snoozeParameters.duration",
        @input="updateField('snoozeParameters.duration', $event)"
      )
    template(v-if="form.generalParameters.type === $constants.ACTION_TYPES.pbehavior")
      pbehavior-form.mt-1(
        :form="form.pbehaviorParameters",
        @input="updateField('pbehaviorParameters', $event)",
        :author="$constants.ACTION_AUTHOR",
        :noFilter="true"
      )
</template>

<script>
import formMixin from '@/mixins/form';
import formValidationHeaderMixin from '@/mixins/form/validation-header';

import DurationField from '@/components/forms/fields/duration.vue';
import PbehaviorForm from '@/components/other/pbehavior/form/pbehavior-form.vue';

export default {
  inject: ['$validator'],
  components: {
    DurationField,
    PbehaviorForm,
  },
  mixins: [formMixin, formValidationHeaderMixin],
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      required: true,
    },
  },
};
</script>
