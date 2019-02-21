<template lang="pug" functional>
  v-card
    v-card-text
      v-list
        v-list-tile
          v-list-tile-content
            v-list-tile-title.text-xs-center
              strong {{ parent.$t('webhook.tabs.request.fields.method') }}:
              |  {{ props.request.method }}
        v-list-tile
          v-list-tile-content
            v-list-tile-title.text-xs-center
              strong {{ parent.$t('webhook.tabs.request.fields.url') }}:
              |  {{ props.request.url }}
        v-list-tile
          v-list-tile-content
            v-list-tile-title.text-xs-center
              strong {{ parent.$t('webhook.tabs.request.fields.headers') }}
        v-list-tile(v-for="(value, key) in props.request.headers", :key="key")
          v-list-tile-content
            v-list-tile-title.text-xs-center
              strong {{ key }}:
              |  {{ value }}
        v-list-tile
          v-list-tile-content
            v-list-tile-title.text-xs-center
              strong {{ parent.$t('webhook.tabs.request.fields.payload') }}
        v-list-tile
          v-list-tile-content
            v-textarea(:value="props.request.payload", read-only, disabled)
</template>

<script>
export default {
  props: {
    request: {
      type: Object,
      required: true,
    },
  },
};
</script>
