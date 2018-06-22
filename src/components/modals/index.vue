<template lang="pug">
  div
    modal(
    v-for="modal in modals",
    :key="modal.name",
    :name="modal.name",
    :dialogProps="modal.dialogProps || defaultDialogProps"
    )
      component(:is="modal.name")
</template>

<script>
// COMPONENTS
import Modal from '@/components/modals/layouts/modal.vue';
import CreateAckEvent from '@/components/modals/create-ack-event.vue';
import CreateAssociateTicketEvent from '@/components/modals/create-associate-ticket-event.vue';
import CreateCancelEvent from '@/components/modals/create-cancel-event.vue';
import CreateChangeStateEvent from '@/components/modals/create-change-state-event.vue';
import CreateDeclareTicketEvent from '@/components/modals/create-declare-ticket-event.vue';
import CreateSnoozeEvent from '@/components/modals/create-snooze-event.vue';
import CreatePbehavior from '@/components/modals/create-pbehavior.vue';
import PbehaviorList from '@/components/modals/pbehavior-list.vue';
import EditLiveReporting from '@/components/modals/edit-live-reporting.vue';
import MoreInfos from '@/components/modals/more-infos.vue';
import CreateEntity from '@/components/modals/create-entity.vue';
// OTHERS
import { MODALS } from '@/constants';

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
    EditLiveReporting,
    MoreInfos,
    CreateEntity,
  },
  data() {
    return {
      modals: [
        { name: MODALS.createAckEvent },
        { name: MODALS.createAssociateTicketEvent },
        { name: MODALS.createCancelEvent },
        { name: MODALS.createChangeStateEvent },
        { name: MODALS.createDeclareTicketEvent },
        { name: MODALS.createSnoozeEvent },
        { name: MODALS.createPbehavior },
        { name: MODALS.pbehaviorList, dialogProps: { maxWidth: 1280, lazy: true } },
        { name: MODALS.editLiveReporting },
        { name: MODALS.moreInfos },
        { name: MODALS.createEntity },
      ],
      defaultDialogProps: { maxWidth: 700, lazy: true },
    };
  },
};
</script>
