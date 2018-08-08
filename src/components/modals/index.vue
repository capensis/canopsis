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
import { MODALS } from '@/constants';

import Modal from './layouts/modal.vue';
import CreateAckEvent from './alarm/create-ack-event.vue';
import CreateAssociateTicketEvent from './alarm/create-associate-ticket-event.vue';
import CreateCancelEvent from './alarm/create-cancel-event.vue';
import CreateChangeStateEvent from './alarm/create-change-state-event.vue';
import CreateDeclareTicketEvent from './alarm/create-declare-ticket-event.vue';
import CreateSnoozeEvent from './alarm/create-snooze-event.vue';
import EditLiveReporting from './alarm/edit-live-reporting.vue';
import MoreInfos from './alarm/more-infos.vue';
import CreatePbehavior from './pbehavior/create-pbehavior.vue';
import PbehaviorList from './pbehavior/pbehavior-list.vue';
import Confirmation from './common/confirmation.vue';
import CreateWidget from './common/create-widget.vue';
import CreateWatcher from './context/create-watcher.vue';
import CreateEntity from './context/create-entity.vue';
import WeatherWatcher from './weather-watcher/weather-watcher.vue';

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
    Confirmation,
    CreateEntity,
    CreateWatcher,
    CreateWidget,
    WeatherWatcher,
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
        { name: MODALS.confirmation },
        { name: MODALS.createEntity },
        { name: MODALS.createWatcher },
        { name: MODALS.weatherWatcher },
        { name: MODALS.createWidget, dialogProps: { maxWidth: 500, lazy: true } },
      ],
      defaultDialogProps: { maxWidth: 700, lazy: true },
    };
  },
};
</script>
