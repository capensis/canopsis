<template lang="pug">
  div
    modal(name="create-ack-event", :dialogProps="{ maxWidth: 700, lazy: true }")
      create-ack-event
    modal(name="create-cancel-event", :dialogProps="{ maxWidth: 700, lazy: true }")
      create-cancel-event
    modal(name="create-change-state-event", :dialogProps="{ maxWidth: 700, lazy: true }")
      create-change-state-event
    modal(name="create-snooze-event", :dialogProps="{ maxWidth: 700, lazy: true }")
      create-snooze-event
    modal(name="create-pbehavior", :dialogProps="{ maxWidth: 700, lazy: true }")
      create-pbehavior
    modal(name="pbehavior-list", :dialogProps="{ maxWidth: 1280, lazy: true }")
      pbehavior-list
</template>

<script>
import Modal from './layouts/modal.vue';
import CreateAckEvent from './create-ack-event.vue';
import CreateCancelEvent from './create-cancel-event.vue';
import CreateChangeStateEvent from './create-change-state-event.vue';
import CreateSnoozeEvent from './create-snooze-event.vue';
import CreatePbehavior from './create-pbehavior.vue';
import PbehaviorList from './pbehavior-list.vue';

export default {
  components: {
    Modal,
    CreateAckEvent,
    CreateCancelEvent,
    CreateChangeStateEvent,
    CreateSnoozeEvent,
    CreatePbehavior,
    PbehaviorList,
  },
};
</script>
