<template lang="pug">
  div
    .addContainer
        .label {{ label }}:&ensp;
        .entities.scrollbar
          span(v-for="entity in entities") {{ entity }}
        v-btn(icon @click="showList=!showList")
          v-icon {{ listIcon }}
        v-btn(icon @click="clear")
          v-icon clear
    context-general-list(
      v-show="showList",
      @update:selectedIds="updateEntities($event)",
    )
</template>

<script>
import ContextGeneralList from '@/components/other/context/context-general-list.vue';
import union from 'lodash/union';

/**
 * Component to select entities for impact/dependencies
 *
 * @module context
 *
 * @prop {String} [label] - "Impacts" or "Dependencies"
 *
 * @event selectedIds#update
 */
export default {
  components: { ContextGeneralList },
  props: {
    label: {
      type: String,
      required: true,
    },
    entities: {
      type: Array,
      default() {
        return [];
      },
    },
  },
  data() {
    return {
      showList: false,
    };
  },
  computed: {
    listIcon() {
      return this.showList ? 'remove' : 'add';
    },
  },
  methods: {
    updateEntities(entities) {
      const entitiesIds = entities.map(entity => entity._id);
      const selectedEntities = union(entitiesIds, ...this.entities);
      this.$emit('updateEntities', selectedEntities);
    },
    clear() {
      this.entities = [];
      this.$emit('updateEntities', this.entities);
    },
  },
};
</script>

<style scoped>
  .addContainer {
    display :flex;
    align-items: center;
    align-content: flex-start;
  }
  .entities{
    width: 50%;
    white-space: nowrap;
    overflow: auto;
  }
  .label{
    width: 16%;
  }
  .scrollbar::-webkit-scrollbar-track
  {
    -webkit-box-shadow: inset 0 0 6px rgba(0,0,0,0.3);
    border-radius: 10px;
    background-color: #F5F5F5;
  }

  .scrollbar::-webkit-scrollbar
  {
    height: 3px;
  }

  .scrollbar::-webkit-scrollbar-thumb
  {
    border-radius: 10px;
    -webkit-box-shadow: inset 0 0 6px rgba(0,0,0,.3);
    background-color: darkgray;
    height: 0;
  }
</style>
