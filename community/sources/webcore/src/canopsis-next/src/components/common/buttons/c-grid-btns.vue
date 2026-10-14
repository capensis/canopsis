<template>
  <div>
    <c-action-btn
      v-if="editing"
      :tooltip="$t('alarm.tooltips.resetChangeColumns')"
      icon="restart_alt"
      @click="$emit('reset')"
    />
    <c-action-btn
      :icon="editing ? 'lock_open' : 'lock_outline'"
      :tooltip="$t(`alarm.tooltips.${editing ? 'finishChangeColumns' : 'startChangeColumns'}`)"
      @click="$emit('toggle')"
    />
  </div>
</template>

<script>
export default {
  props: {
    editing: {
      type: Boolean,
      default: false,
    },
  },

  emits: ['reset', 'toggle'],
};
</script>
