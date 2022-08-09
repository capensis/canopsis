<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper(close)
      template(#title="")
        span {{ $t('modals.createDeclareTicket.title') }}
      template(#text="")
        v-container
          v-layout(row)
            v-flex.text-xs-center
              alarm-general-table(:items="items")
          v-layout(row)
            v-divider.my-3
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
        ) {{ $t('common.reportIncident') }}
</template>

<script>
import { MODALS, EVENT_ENTITY_TYPES, DECLARE_TICKET_OUTPUT } from '@/constants';

import { modalInnerMixin } from '@/mixins/modal/inner';
import { modalInnerItemsMixin } from '@/mixins/modal/inner-items';
import { eventActionsAlarmMixin } from '@/mixins/event-actions/alarm';
import { submittableMixinCreator } from '@/mixins/submittable';

import AlarmGeneralTable from '@/components/widgets/alarm/alarm-general-list.vue';

import ModalWrapper from '../modal-wrapper.vue';

/**
 * Modal to declare a ticket
 */
export default {
  name: MODALS.createDeclareTicketEvent,
  components: { AlarmGeneralTable, ModalWrapper },
  mixins: [
    modalInnerMixin,
    modalInnerItemsMixin,
    eventActionsAlarmMixin,
    submittableMixinCreator(),
  ],
  methods: {
    async submit() {
      await this.createEvent(EVENT_ENTITY_TYPES.declareTicket, this.items, {
        output: DECLARE_TICKET_OUTPUT,
      });

      this.$modals.hide();
    },
  },
};
</script>
