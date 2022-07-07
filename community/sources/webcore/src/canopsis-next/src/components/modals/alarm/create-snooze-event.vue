<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper(close)
      template(#title="")
        span {{ $t('modals.createSnoozeEvent.title') }}
      template(#text="")
        v-container
          snooze-event-form(v-model="form", :is-note-required="isNoteRequired")
      template(#actions="")
        v-btn(
          depressed,
          flat,
          @click="$modals.hide"
        ) {{ $t('common.cancel') }}
        v-btn.primary(
          :loading="submitting",
          :disabled="isDisabled",
          type="submit"
        ) {{ $t('common.saveChanges') }}
</template>

<script>
import { MODALS, EVENT_ENTITY_TYPES } from '@/constants';

import { modalInnerMixin } from '@/mixins/modal/inner';
import { modalInnerItemsMixin } from '@/mixins/modal/inner-items';
import { eventActionsAlarmMixin } from '@/mixins/event-actions/alarm';
import { submittableMixinCreator } from '@/mixins/submittable';
import { confirmableModalMixinCreator } from '@/mixins/confirmable-modal';

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
    modalInnerMixin,
    modalInnerItemsMixin,
    eventActionsAlarmMixin,
    submittableMixinCreator(),
    confirmableModalMixinCreator(),
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
