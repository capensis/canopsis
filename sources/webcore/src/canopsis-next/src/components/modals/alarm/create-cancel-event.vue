<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper
      template(slot="title")
        span {{ title }}
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
              :label="$t('modals.createCancelEvent.fields.output')",
              :error-messages="errors.collect('output')",
              name="output"
            )
      template(slot="actions")
        v-btn(depressed, flat, @click="$modals.hide") {{ $t('common.cancel') }}
        v-btn.primary(
          :loading="submitting",
          :disabled="isDisabled",
          type="submit"
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
  name: MODALS.createCancelEvent,
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
  computed: {
    title() {
      return this.config.title || this.$t('modals.createCancelEvent.title');
    },

    eventType() {
      return this.config.eventType || EVENT_ENTITY_TYPES.cancel;
    },
  },
  methods: {
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        const data = { ...this.form };

        if (this.eventType === EVENT_ENTITY_TYPES.cancel) {
          data.cancel = 1;
        }

        await this.createEvent(this.eventType, this.items, data);

        this.$modals.hide();
      }
    },
  },
};
</script>
