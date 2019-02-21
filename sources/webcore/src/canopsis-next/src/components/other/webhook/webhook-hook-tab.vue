<template lang="pug" functional>
  v-card
    v-card-text
      v-list
        v-list-tile
          v-list-tile-content
            v-list-tile-title.text-xs-center
              strong Triggers:
              |  {{ props.hook.triggers.join(', ') }}
      v-tabs(fixed-tabs)
        v-tab {{ parent.$t('webhook.tabs.hook.fields.eventPatterns') }}
        v-tab {{ parent.$t('webhook.tabs.hook.fields.alarmPatterns') }}
        v-tab {{ parent.$t('webhook.tabs.hook.fields.entityPatterns') }}
        v-tab-item
          v-textarea(
          v-for="(pattern, index) in props.hook.event_pattern",
          :key="`pattern-event-${index}`",
          :value="JSON.stringify(pattern, null, 4)",
          read-only,
          disabled
          )
        v-tab-item
          v-textarea(
          v-for="(pattern, index) in props.hook.alarm_pattern",
          :key="`pattern-alarm-${index}`",
          :value="JSON.stringify(pattern, null, 4)",
          read-only,
          disabled
          )
        v-tab-item
          v-textarea(
          v-for="(pattern, index) in props.hook.entity_pattern",
          :key="`pattern-eneity-${index}`",
          :value="JSON.stringify(pattern, null, 4)",
          read-only,
          disabled
          )
</template>

<script>
export default {
  props: {
    hook: {
      type: Object,
      required: true,
    },
  },
};
</script>
