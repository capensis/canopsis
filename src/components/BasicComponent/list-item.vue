<template lang="pug">
  div(id="content", @click="clickOnListItem")
    transition(name="expand", mode="out-in")
      v-card(v-if="!isExpanded", key="reduced")
        slot(name="reduced", :props="item")
      v-card.expanded(v-else, key="expanded", :raised="true")
        slot(name="expanded", :props="item")
</template>

<script>
export default {
  name: 'BrickList',
  data() {
    return {
      isExpanded: false,
    };
  },
  methods: {
    /**
     * A Function to prevent the expansion of a row when you highlight/select something in it
     */
    clickOnListItem() {
      const selection = window.getSelection();
      if (selection.toString().length === 0) {
        this.isExpanded = !this.isExpanded;
      }
    },
  },
  props: ['item'],
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
