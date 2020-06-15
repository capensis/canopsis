<template lang="pug">
  ds-calendar-app
    ds-calendar-event-create-popover(
      slot="eventCreatePopover",
      slot-scope="{ placeholder, calendar, close }",
      :calendar-event="placeholder",
      :calendar="calendar",
      :close="close"
    )
</template>

<script>
export default {
  data() {
    return {
      events: [],
    };
  },
  methods: {},
};
</script>
