<template lang="pug">
  div(id="content", @click="clickOnListItem")
    v-card(key="reduced", :color="color")
      slot(name="checkbox")
      slot(name="reduced", :props="item")
    transition(name="expand", mode="out-in")
      v-card.expanded(v-if="isExpanded", key="expanded", :raised="true", :color="color")
        slot(name="expanded", :props="item")
</template>

<script>
export default {
  props: {
    item: {
      type: Object,
    },
    expanded: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      isExpanded: false,
    };
  },
  computed: {
    color() {
      if (this.isExpanded) {
        return 'grey lighten-5';
      }
      return 'white';
    },
  },
  methods: {
    /**
       * A Function to prevent the expansion of a row when you highlight/select something in it
       */
    clickOnListItem() {
      const selection = window.getSelection();
      if (selection.toString().length === 0) {
        if (this.expanded) {
          this.isExpanded = !this.isExpanded;
        }
      }
    },
  },
};
</script>

<style scoped>
  .expand-enter-active, .expand-leave-active {
    transition: opacity .01s ease;
  }

  .expanded {
    margin: 15px;
  }
</style>
