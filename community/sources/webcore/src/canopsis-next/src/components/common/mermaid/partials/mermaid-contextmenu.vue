<template lang="pug">
  v-list(dense)
    v-list-tile(v-for="item in items", :key="item.event", @click="applyEvent(item.event)")
      v-list-tile-content
        v-list-tile-title {{ item.text }}
</template>

<script>
export default {
  props: {
    editing: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    items() {
      if (this.editing) {
        return [
          {
            text: this.$t('mermaid.editPoint'),
            event: 'edit:point',
          },
          {
            text: this.$t('mermaid.removePoint'),
            event: 'remove:point',
          },
        ];
      }

      return [
        {
          text: this.$t('mermaid.addPoint'),
          event: 'add:point',
        },
      ];
    },
  },
  methods: {
    applyEvent(event) {
      this.$emit(event);
    },
  },
};
</script>
