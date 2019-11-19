<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper
      template(slot="title")
        span {{ $t('modals.createSnoozeEvent.title') }}
      template(slot="text")
        v-container
          duration-field(v-model="form")
      template(slot="actions")
        v-btn(depressed, flat, @click="$modals.hide") {{ $t('common.cancel') }}
        v-btn.primary(
          :disabled="errors.any() || submitting",
          type="submit"
        ) {{ $t('common.actions.saveChanges') }}
</template>

<script>
import moment from 'moment';

import { MODALS, EVENT_ENTITY_TYPES, DURATION_UNITS } from '@/constants';

import modalInnerItemsMixin from '@/mixins/modal/inner-items';
import eventActionsAlarmMixin from '@/mixins/event-actions/alarm';

import DurationField from '@/components/forms/fields/duration.vue';

import ModalWrapper from '../modal-wrapper.vue';

/**
 * Modal to put a snooze on an alarm
 */
export default {
  name: MODALS.createSnoozeEvent,
  $_veeValidate: {
    validator: 'new',
  },
  components: { DurationField, ModalWrapper },
  mixins: [modalInnerItemsMixin, eventActionsAlarmMixin],
  data() {
    return {
      submitting: false,
      form: {
        duration: 1,
        durationType: DURATION_UNITS.minute.value,
      },
    };
  },
  methods: {
    async submit() {
      try {
        this.submitting = true;

        const isFormValid = await this.$validator.validateAll();

        if (isFormValid) {
          const duration = moment.duration(
            parseInt(this.form.duration, 10),
            this.form.durationType,
          ).asSeconds();

          await this.createEvent(EVENT_ENTITY_TYPES.snooze, this.items, { duration });

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
