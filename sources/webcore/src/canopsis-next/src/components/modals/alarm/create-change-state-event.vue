<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper
      template(slot="title")
        span {{ $t('modals.createChangeStateEvent.title') }}
      template(slot="text")
        v-container
          v-layout(row)
            state-criticity-field(v-model="form.state", :stateValues="availableStateValues")
          v-layout.mt-4(row)
            v-text-field(
              v-model="form.output",
              v-validate="'required'",
              :label="$t('modals.createChangeStateEvent.fields.output')",
              :error-messages="errors.collect('output')",
              name="output"
            )
      template(slot="actions")
        v-btn(depressed, flat, @click="$modals.hide") {{ $t('common.cancel') }}
        v-btn.primary(
          :disabled="errors.any() || submitting",
          type="submit"
        ) {{ $t('common.actions.saveChanges') }}
</template>

<script>
import { omit } from 'lodash';

import { MODALS, ENTITIES_STATES, EVENT_ENTITY_TYPES } from '@/constants';

import modalInnerItemsMixin from '@/mixins/modal/inner-items';
import eventActionsAlarmMixin from '@/mixins/event-actions/alarm';

import StateCriticityField from '@/components/forms/fields/state-criticity-field.vue';

import ModalWrapper from '../modal-wrapper.vue';

/**
 * Modal to create a 'change-state' event
 */
export default {
  name: MODALS.createChangeStateEvent,
  $_veeValidate: {
    validator: 'new',
  },
  components: { StateCriticityField, ModalWrapper },
  mixins: [modalInnerItemsMixin, eventActionsAlarmMixin],
  data() {
    return {
      submitting: false,
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
      try {
        this.submitting = true;

        const isFormValid = await this.$validator.validateAll();

        if (isFormValid) {
          await this.createEvent(EVENT_ENTITY_TYPES.changeState, this.items, this.form);

          this.$modals.hide();
        }
      } catch (err) {
        this.$popups.error({ text: err.description || this.$t('error.default') });
      } finally {
        this.submitting = false;
      }
    },
  },
};
</script>
