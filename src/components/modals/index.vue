<template lang="pug">
  div
    modal(:name="createAckEvent", :dialogProps="{ maxWidth: 700, lazy: true }")
      create-ack-event
    modal(:name="createAssociateTicketEvent", :dialogProps="{ maxWidth: 700, lazy: true }")
      create-associate-ticket-event
    modal(:name="createCancelEvent", :dialogProps="{ maxWidth: 700, lazy: true }")
      create-cancel-event
    modal(:name="createChangeStateEvent", :dialogProps="{ maxWidth: 700, lazy: true }")
      create-change-state-event
    modal(:name="createDeclareTicketEvent", :dialogProps="{ maxWidth: 700, lazy: true }")
      create-declare-ticket-event
    modal(:name="createSnoozeEvent", :dialogProps="{ maxWidth: 700, lazy: true }")
      create-snooze-event
    modal(:name="createPbehavior", :dialogProps="{ maxWidth: 700, lazy: true }")
      create-pbehavior
    modal(:name="pbehaviorList", :dialogProps="{ maxWidth: 1280, lazy: true }")
      pbehavior-list
    modal(name="more-infos", :dialogProps="{ maxWidth: 1280, lazy: true }")
      more-infos
</template>

<script>
import { MODALS } from '@/config';

import Modal from './layouts/modal.vue';
import CreateAckEvent from './create-ack-event.vue';
import CreateAssociateTicketEvent from './create-associate-ticket-event.vue';
import CreateCancelEvent from './create-cancel-event.vue';
import CreateChangeStateEvent from './create-change-state-event.vue';
import CreateDeclareTicketEvent from './create-declare-ticket-event.vue';
import CreateSnoozeEvent from './create-snooze-event.vue';
import CreatePbehavior from './create-pbehavior.vue';
import PbehaviorList from './pbehavior-list.vue';
import MoreInfos from './more-infos.vue';

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
    MoreInfos,
  },
  data() {
    return {
      ...MODALS,
    };
  },
};
</script>
