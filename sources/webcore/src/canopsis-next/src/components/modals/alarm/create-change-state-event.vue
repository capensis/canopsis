<template lang="pug">
  v-form(data-test="createChangeStateEventModal", @submit.prevent="submit")
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
              name="output",
              data-test="createChangeStateEventNote"
            )
      template(slot="actions")
        v-btn(
          data-test="createChangeStateEventCancelButton",
          depressed,
          flat,
          @click="$modals.hide"
        ) {{ $t('common.cancel') }}
        v-btn.primary(
          :loading="submitting",
          :disabled="isDisabled",
          data-test="createChangeStateEventSubmitButton",
          type="submit"
        ) {{ $t('common.actions.saveChanges') }}
</template>

<script>
import { omit } from 'lodash';

import { MODALS, ENTITIES_STATES, EVENT_ENTITY_TYPES } from '@/constants';

import modalInnerItemsMixin from '@/mixins/modal/inner-items';
import eventActionsAlarmMixin from '@/mixins/event-actions/alarm';
import submittableMixin from '@/mixins/submittable';
import entitiesInfoMixin from '@/mixins/entities/info';

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
  mixins: [entitiesInfoMixin, modalInnerItemsMixin, eventActionsAlarmMixin, submittableMixin()],
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
      return this.allowChangeSeverityToInfo ? ENTITIES_STATES : omit(ENTITIES_STATES, ['ok']);
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

        this.$modals.hide();
      }
    },
  },
};
</script>
