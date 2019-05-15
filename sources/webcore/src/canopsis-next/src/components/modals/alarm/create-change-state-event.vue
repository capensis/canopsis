<template lang="pug">
  v-form(@submit.prevent="submit")
    v-card
      v-card-title.primary.white--text
        v-layout(justify-space-between, align-center)
          span.headline {{ $t('modals.createChangeStateEvent.title') }}
      v-card-text
        v-container
          v-layout(row)
            state-criticity-field(v-model="form.state", :stateValues="availableStateValues")
          v-layout.mt-4(row)
            v-text-field(
            :label="$t('modals.createChangeStateEvent.fields.output')",
            :error-messages="errors.collect('output')",
            v-model="form.output",
            v-validate="'required'",
            data-vv-name="output"
            )
      v-divider
      v-layout.py-1(justify-end)
        v-btn(@click="hideModal", depressed, flat) {{ $t('common.cancel') }}
        v-btn.primary(type="submit", :disabled="errors.any()") {{ $t('common.actions.saveChanges') }}
</template>

<script>
import { omit } from 'lodash';

import { MODALS, ENTITIES_STATES, EVENT_ENTITY_TYPES } from '@/constants';

import modalInnerItemsMixin from '@/mixins/modal/inner-items';
import eventActionsAlarmMixin from '@/mixins/event-actions/alarm';

import StateCriticityField from '@/components/forms/fields/state-criticity-field.vue';

/**
 * Modal to create a 'change-state' event
 */
export default {
  name: MODALS.createChangeStateEvent,

  $_veeValidate: {
    validator: 'new',
  },
  components: { StateCriticityField },
  mixins: [modalInnerItemsMixin, eventActionsAlarmMixin],
  data() {
    return {
      form: {
        output: '',
        state: ENTITIES_STATES.major,
      },
    };
  },
  computed: {
    availableStateValues() {
      return omit(ENTITIES_STATES, ['ok']);
    },
  },
  mounted() {
    this.form.state = this.firstItem.v.state.val;
  },
  methods: {
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        await this.createEvent(EVENT_ENTITY_TYPES.changeState, this.items, this.form);

        this.hideModal();
      }
    },
  },
};
</script>
