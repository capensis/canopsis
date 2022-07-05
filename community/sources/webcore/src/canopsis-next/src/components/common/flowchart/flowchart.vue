<template lang="pug">
  div.flowchart.fill-height
    flowchart-sidebar(
      v-model="shapes",
      :view-box="viewBox",
      :selected.sync="selected"
    )
    div.flowchart__editor
      flowchart-editor(
        v-model="shapes",
        :view-box.sync="viewBox",
        :selected.sync="selected"
      )
    div.flowchart__properties(v-show="selected.length")
      flowchart-properties(v-model="shapes", :selected="selected")
</template>

<script>
import FlowchartEditor from './flowchart-editor.vue';
import FlowchartSidebar from './flowchart-sidebar.vue';
import FlowchartProperties from './flowchart-properties.vue';

export default {
  components: {
    FlowchartSidebar,
    FlowchartEditor,
    FlowchartProperties,
  },
  data() {
    return {
      shapes: {},
      selected: [],
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
    flex-grow: 1;
    height: 100%;

    position: absolute;
    left: 300px;
    top: 0;
    right: 0;
    bottom: 0;
  }

  &__properties {
    flex-grow: 1;
    width: 350px;

    position: absolute;
    top: 0;
    right: 0;
  }
}
</style>
