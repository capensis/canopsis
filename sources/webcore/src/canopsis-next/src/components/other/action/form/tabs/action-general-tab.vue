<template lang="pug">
  div
    template(v-if="form.generalParameters.type === $constants.ACTION_TYPES.snooze")
      v-textarea(
        v-field="form.snoozeParameters.message",
        :label="$t('modals.createAction.fields.message')"
      )
      duration-field(v-field="form.snoozeParameters.duration")
    template(v-if="form.generalParameters.type === $constants.ACTION_TYPES.pbehavior")
      pbehavior-form(
        v-field="form.pbehaviorParameters",
        :author="$constants.ACTION_AUTHOR",
        noFilter
      )
    template(v-if="form.generalParameters.type === $constants.ACTION_TYPES.changeState")
      div.mt-3
        v-layout(row)
          state-criticity-field(
            v-field="form.changeStateParameters.state",
            :stateValues="availableStateValues"
          )
        v-layout.mt-4(row)
          v-text-field(
            v-field="form.changeStateParameters.author",
            v-validate="'required'",
            :error-messages="errors.collect('author')",
            :label="$t('common.author')",
            name="author"
          )
        v-layout(row)
          v-textarea(
            v-field="form.changeStateParameters.message",
            :label="$t('modals.createAction.fields.message')"
          )
</template>

<script>
import { omit } from 'lodash';

import { ENTITIES_STATES } from '@/constants';

import formValidationHeaderMixin from '@/mixins/form/validation-header';

import DurationField from '@/components/forms/fields/duration.vue';
import PbehaviorForm from '@/components/other/pbehavior/form/pbehavior-form.vue';
import StateCriticityField from '@/components/forms/fields/state-criticity-field.vue';

export default {
  inject: ['$validator'],
  components: {
    DurationField,
    PbehaviorForm,
    StateCriticityField,
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
  computed: {
    availableStateValues() {
      return omit(ENTITIES_STATES, ['ok']);
    },
  },
};
</script>
