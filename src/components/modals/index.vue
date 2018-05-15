<template lang="pug">
  div
    modal(name="add-ack-event", :dialogProps="{ maxWidth: 700, lazy: true }")
      add-ack-event
    modal(name="add-cancel-event", :dialogProps="{ maxWidth: 700, lazy: true }")
      add-cancel-event
    modal(name="add-change-state-event", :dialogProps="{ maxWidth: 700, lazy: true }")
      add-change-state-event
    modal(name="add-snooze-event", :dialogProps="{ maxWidth: 700, lazy: true }")
      add-snooze-event
    modal(name="add-pbehavior", :dialogProps="{ maxWidth: 700, lazy: true }")
      add-pbehavior
    modal(name="pbehavior-list", :dialogProps="{ maxWidth: 1280, lazy: true }")
      pbehavior-list
</template>

<script>
import Modal from './layouts/modal.vue';
import AddAckEvent from './add-ack-event.vue';
import AddCancelEvent from './add-cancel-event.vue';
import AddChangeStateEvent from './add-change-state-event.vue';
import AddSnoozeEvent from './add-snooze-event.vue';
import AddPbehavior from './add-pbehavior.vue';
import PbehaviorList from './pbehavior-list.vue';

export default {
  components: {
    Modal,
    AddAckEvent,
    AddCancelEvent,
    AddChangeStateEvent,
    AddSnoozeEvent,
    AddPbehavior,
    PbehaviorList,
  },
};
</script>
