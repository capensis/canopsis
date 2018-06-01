<template lang="pug">
  div
    modal(v-for="modal in modals", :key="modal.name", :name="modal.name", :dialogProps="modal.dialogProps")
      component(:is="modal.name")
</template>

<script>
import { MODALS } from '@/constants';

import Modal from './layouts/modal.vue';
import CreateAckEvent from './create-ack-event.vue';
import CreateAssociateTicketEvent from './create-associate-ticket-event.vue';
import CreateCancelEvent from './create-cancel-event.vue';
import CreateChangeStateEvent from './create-change-state-event.vue';
import CreateDeclareTicketEvent from './create-declare-ticket-event.vue';
import CreateSnoozeEvent from './create-snooze-event.vue';
import CreatePbehavior from './create-pbehavior.vue';
import PbehaviorList from './pbehavior-list.vue';

/**
 * Wrapper for all modal windows
 */
export default {
  components: {
    Modal,
    CreateAckEvent,
    CreateAssociateTicketEvent,
    CreateCancelEvent,
    CreateChangeStateEvent,
    CreateDeclareTicketEvent,
    CreateSnoozeEvent,
    CreatePbehavior,
    PbehaviorList,
  },
  data() {
    const dialogProps = { maxWidth: 700, lazy: true };

    return {
      modals: [
        { name: MODALS.createAckEvent, dialogProps },
        { name: MODALS.createAssociateTicketEvent, dialogProps },
        { name: MODALS.createCancelEvent, dialogProps },
        { name: MODALS.createChangeStateEvent, dialogProps },
        { name: MODALS.createDeclareTicketEvent, dialogProps },
        { name: MODALS.createSnoozeEvent, dialogProps },
        { name: MODALS.createPbehavior, dialogProps },
        { name: MODALS.pbehaviorList, dialogProps: { maxWidth: 1280, lazy: true } },
      ],
    };
  },
};
</script>
