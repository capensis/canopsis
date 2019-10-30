<template lang="pug">
  v-form(data-test="createSnoozeEventModal", @submit.prevent="submit")
    v-card
      v-card-title.primary.white--text
        v-layout(justify-space-between, align-center)
          span.headline {{ $t('modals.createSnoozeEvent.title') }}
      v-card-text
        v-container
          duration-field(v-model="form")
      v-divider
      v-layout.py-1(justify-end)
        v-btn(
          data-test="createSnoozeEventCancelButton",
          @click="hideModal",
          depressed,
          flat
        ) {{ $t('common.cancel') }}
        v-btn(
          data-test="createSnoozeEventSubmitButton",
          type="submit",
          :disabled="errors.any()",
          color="primary"
        ) {{ $t('common.actions.saveChanges') }}
</template>

<script>
import moment from 'moment';

import { MODALS, EVENT_ENTITY_TYPES, DURATION_UNITS } from '@/constants';

import modalInnerItemsMixin from '@/mixins/modal/inner-items';
import eventActionsAlarmMixin from '@/mixins/event-actions/alarm';

import DurationField from '@/components/forms/fields/duration.vue';

/**
 * Modal to put a snooze on an alarm
 */
export default {
  name: MODALS.createSnoozeEvent,

  $_veeValidate: {
    validator: 'new',
  },
  components: {
    DurationField,
  },
  mixins: [modalInnerItemsMixin, eventActionsAlarmMixin],
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

        this.hideModal();
      }
    },
  },
};
</script>
