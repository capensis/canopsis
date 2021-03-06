<template lang="pug">
  v-form(data-test="createSnoozeEventModal", @submit.prevent="submit")
    modal-wrapper(close)
      template(slot="title")
        span {{ $t('modals.createSnoozeEvent.title') }}
      template(slot="text")
        v-container
          snooze-event-form(v-model="form", :is-note-required="isNoteRequired")
      template(slot="actions")
        v-btn(
          data-test="createSnoozeEventCancelButton",
          depressed,
          flat,
          @click="$modals.hide"
        ) {{ $t('common.cancel') }}
        v-btn.primary(
          :loading="submitting",
          :disabled="isDisabled",
          data-test="createSnoozeEventSubmitButton",
          type="submit"
        ) {{ $t('common.actions.saveChanges') }}
</template>

<script>
import { MODALS, EVENT_ENTITY_TYPES } from '@/constants';

import modalInnerItemsMixin from '@/mixins/modal/inner-items';
import eventActionsAlarmMixin from '@/mixins/event-actions/alarm';
import submittableMixin from '@/mixins/submittable';
import confirmableModalMixin from '@/mixins/confirmable-modal';

import { formToSnooze, snoozeToForm } from '@/helpers/forms/snooze-event';

import SnoozeEventForm from '@/components/widgets/alarm/forms/snooze-event-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

/**
 * Modal to put a snooze on an alarm
 */
export default {
  name: MODALS.createSnoozeEvent,
  $_veeValidate: {
    validator: 'new',
  },
  components: {
    ModalWrapper,
    SnoozeEventForm,
  },
  mixins: [
    modalInnerItemsMixin,
    eventActionsAlarmMixin,
    submittableMixin(),
    confirmableModalMixin(),
  ],
  data() {
    return {
      form: snoozeToForm(),
    };
  },
  computed: {
    isNoteRequired() {
      return this.config.isNoteRequired;
    },
  },
  methods: {
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        const form = formToSnooze(this.form);

        await this.createEvent(EVENT_ENTITY_TYPES.snooze, this.items, form);

        this.$modals.hide();
      }
    },
  },
};
</script>
