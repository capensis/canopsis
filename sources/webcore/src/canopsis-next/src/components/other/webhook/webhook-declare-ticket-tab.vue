<template lang="pug" functional>
  v-card
    v-card-text
      v-list
        v-list-tile(v-for="(value, key) in props.declareTicket", :key="key")
          v-list-tile-content
            v-list-tile-title.text-xs-center
              strong {{ key }}:
              |  {{ value }}
</template>

<script>
export default {
  props: {
    declareTicket: {
      type: Object,
      default: () => ({}),
    },
  },
};
</script>
