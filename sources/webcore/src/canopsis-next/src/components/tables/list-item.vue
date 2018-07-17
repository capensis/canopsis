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

/**
* Wrapper component for a row on basic-list
*
* @prop {Object} [item] - Item (alarm, entity) to display on the row
* @prop {Boolean} [expanded] - Boolean to determine if the row is expanded
*/
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
