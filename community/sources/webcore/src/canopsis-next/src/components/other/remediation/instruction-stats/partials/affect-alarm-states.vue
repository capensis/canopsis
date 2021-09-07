<template lang="pug">
  div.d-flex
    c-state-count-changes-chips(:affect="alarmStates.minor", :color="$config.COLORS.state.minor")
    c-state-count-changes-chips.ml-2(:affect="alarmStates.major", :color="$config.COLORS.state.major")
    c-state-count-changes-chips.ml-2(:affect="alarmStates.critical", :color="$config.COLORS.state.critical")
</template>

<script>
export default {
  props: {
    alarmStates: {
      type: Object,
      required: true,
    },
  },
};
</script>
