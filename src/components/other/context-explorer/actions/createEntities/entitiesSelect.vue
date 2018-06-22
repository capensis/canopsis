<template lang="pug">
  div
    .addContainer
        span {{ label }}:&ensp;
        .entities
          span(v-for="entity in entities") {{ entity._id }},&emsp;
        v-btn(icon @click="showList=!showList")
          v-icon add
        v-btn(icon @click="clear")
          v-icon clear
    context-list(
      v-show="showList",
      @update:selectedIds="updateEntities($event)",
    )
</template>

<script>
import ContextList from '@/components/other/context-explorer/context-general-list.vue';
import union from 'lodash/union';


export default {
  components: { ContextList },
  props: {
    label: {
      type: String,
      required: true,
    },
  },
  data() {
    return {
      showList: false,
      entities: [],
    };
  },
  computed: {
  },
  methods: {
    updateEntities(entities) {
      this.entities = union(entities, this.entities);
      this.$emit('update:entities', this.entities);
    },
    clear() {
      this.entities = [];
      this.$emit('update:selectedIds', this.entities);
    },
  },
};
</script>

<style scoped>
  .addContainer {
    display :flex;
    align-items: center;
  }
  .entities{
    width: 50%;
    white-space: nowrap;
    overflow: hidden;
  }
</style>
