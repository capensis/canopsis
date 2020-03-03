<template lang="pug">
  v-form(data-test="createEventModal", @submit.prevent="submit")
    modal-wrapper
      template(slot="title")
        span {{ config.title }}
      template(slot="text")
        v-container
          v-layout(row)
            v-flex.text-xs-center
              alarm-general-table(:items="items")
          v-layout(row)
            v-divider.my-3
          v-layout(row)
            v-text-field(
              v-model="form.output",
              v-validate="'required'",
              :label="$t('modals.createEvent.fields.output')",
              :error-messages="errors.collect('output')",
              name="output"
            )
      template(slot="actions")
        v-btn(
          data-test="createEventCancelButton",
          depressed,
          flat,
          @click="$modals.hide"
        ) {{ $t('common.cancel') }}
        v-btn.primary(
          :loading="submitting",
          :disabled="isDisabled",
          type="submit",
          data-test="createEventSubmitButton"
        ) {{ $t('common.actions.saveChanges') }}
</template>

<script>
import { MODALS, EVENT_ENTITY_TYPES } from '@/constants';

import modalInnerItemsMixin from '@/mixins/modal/inner-items';
import eventActionsAlarmMixin from '@/mixins/event-actions/alarm';
import submittableMixin from '@/mixins/submittable';

import AlarmGeneralTable from '@/components/other/alarm/alarm-general-list.vue';

import ModalWrapper from '../modal-wrapper.vue';

/**
 * Modal to cancel an alarm
 */
export default {
  name: MODALS.createEvent,
  $_veeValidate: {
    validator: 'new',
  },
  components: { AlarmGeneralTable, ModalWrapper },
  mixins: [modalInnerItemsMixin, eventActionsAlarmMixin, submittableMixin()],
  data() {
    return {
      form: {
        output: '',
      },
    };
  },
  methods: {
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        const data = { ...this.form };

        if (this.config.eventType === EVENT_ENTITY_TYPES.cancel) {
          data.cancel = 1;
        }

        await this.createEvent(this.eventType, this.items, data);

        this.$modals.hide();
      }
    },
  },
};
</script>
