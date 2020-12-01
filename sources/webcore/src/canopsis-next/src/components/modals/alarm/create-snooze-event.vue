<template lang="pug">
  v-form(data-test="createSnoozeEventModal", @submit.prevent="submit")
    modal-wrapper(close)
      template(slot="title")
        span {{ $t('modals.createSnoozeEvent.title') }}
      template(slot="text")
        v-container
          old-duration-field(v-model="form")
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
import moment from 'moment';

import { MODALS, EVENT_ENTITY_TYPES, DURATION_UNITS } from '@/constants';

import modalInnerItemsMixin from '@/mixins/modal/inner-items';
import eventActionsAlarmMixin from '@/mixins/event-actions/alarm';
import submittableMixin from '@/mixins/submittable';

import OldDurationField from '@/components/forms/fields/old-duration.vue';

import ModalWrapper from '../modal-wrapper.vue';

/**
 * Modal to put a snooze on an alarm
 */
export default {
  name: MODALS.createSnoozeEvent,
  $_veeValidate: {
    validator: 'new',
  },
  components: { OldDurationField, ModalWrapper },
  mixins: [modalInnerItemsMixin, eventActionsAlarmMixin, submittableMixin()],
  data() {
    return {
      form: {
        duration: 1,
        durationType: DURATION_UNITS.minute.value,
      },
    };
  },
  methods: {
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        const duration = moment.duration(
          parseInt(this.form.duration, 10),
          this.form.durationType,
        ).asSeconds();

        await this.createEvent(EVENT_ENTITY_TYPES.snooze, this.items, { duration });

        this.$modals.hide();
      }
    },
  },
};
</script>
