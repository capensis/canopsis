<template lang="pug">
  div
    add-ack-event
    add-cancel-event
    add-change-state-event
    add-snooze-event
    add-pbehavior
</template>

<script>
import AddAckEvent from './add-ack-event.vue';
import AddCancelEvent from './add-cancel-event.vue';
import AddChangeStateEvent from './add-change-state-event.vue';
import AddSnoozeEvent from './add-snooze-event.vue';
import AddPbehavior from './add-pbehavior.vue';

export default {
  components: {
    AddAckEvent,
    AddCancelEvent,
    AddChangeStateEvent,
    AddSnoozeEvent,
    AddPbehavior,
  },
};
</script>
