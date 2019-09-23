<template lang="pug">
  div
    template(v-if="form.generalParameters.type === $constants.ACTION_TYPES.snooze")
      v-textarea(
      v-model="form.snoozeParameters.message",
      :label="$t('modals.createAction.fields.message')",
      )
      duration-field(:duration="form.snoozeParameters.duration")
    template(v-if="form.generalParameters.type === $constants.ACTION_TYPES.pbehavior")
      pbehavior-form(v-model="form.pbehaviorParameters", :author="$constants.ACTION_AUTHOR", :noFilter="true")
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
    form: {
      type: Object,
      required: true,
    },
  },
};
</script>
