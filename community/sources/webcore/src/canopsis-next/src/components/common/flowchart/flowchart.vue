<template lang="pug">
  div.flowchart.fill-height
    flowchart-sidebar(v-model="shapes", :view-box="viewBox")
    div.flowchart__editor
      flowchart-editor(v-model="shapes", :view-box.sync="viewBox")
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
      viewBox: {
        x: 0,
        y: 0,
        width: 0,
        height: 0,
      },
    };
  },
};
</script>

<style lang="scss">
.flowchart {
  position: relative;
  display: flex;

  &__editor {
    background: red;
    flex-grow: 1;
    height: 100%;

    position: absolute;
    left: 300px;
    top: 0;
    right: 0;
    bottom: 0;
  }
}
</style>
