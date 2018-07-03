<template lang="pug">
  v-container
    v-expansion-panel
      v-expansion-panel-content.grey.darken-2.white--text
        div.white--text(slot="header") {{ label }}
        v-card
          v-card-text
            v-layout
              v-chip(v-for="entity in entities", :key="entity._id", close, @input="removeEntity(entity)") {{ entity._id }}
            v-btn.red.white--text(v-show="entities.length", @click="clear", small) Clear
            context-general-list(
              @update:selectedIds="updateEntities($event)",
            )
</template>

<script>
import ContextGeneralList from '@/components/other/context/context-general-list.vue';
import union from 'lodash/union';
import remove from 'lodash/remove';


export default {
  components: { ContextGeneralList },
  props: {
    label: {
      type: String,
      required: true,
    },
  },
  data() {
    return {
      entities: [],
    };
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
    removeEntity(entity) {
      const entities = [...this.entities];
      remove(entities, item => item._id === entity._id);
      this.entities = entities;
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
