<template lang="pug">
  div.flowchart.fill-height
    flowchart-sidebar(v-model="shapes")
    div.flowchart__editor
      flowchart-editor(v-model="shapes")
</template>

<script>
import FlowchartEditor from './flowchart-editor.vue';
import FlowchartSidebar from './flowchart-sidebar.vue';

export default {
  components: {
    FlowchartSidebar,
    FlowchartEditor,
  },
  data() {
    return {
      shapes: {},
    };
  },
};
</script>

<style lang="scss">
.flowchart {
  display: flex;

  &__editor {
    flex-grow: 1;
  }
}
</style>
