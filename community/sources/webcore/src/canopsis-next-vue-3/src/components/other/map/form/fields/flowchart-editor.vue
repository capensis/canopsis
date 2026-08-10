<template>
  <c-layout column>
    <flowchart
      :shapes="form.shapes"
      :background-color="form.background_color"
      :style="editorStyles"
      :cursor-style="addOnClick ? 'crosshair' : undefined"
      class="flowchart-map-editor mb-2"
      @update:shapes="updateShapes"
      @update:background-color="updateBackgroundColor"
    >
      <template #sidebar-prepend>
        <add-location-btn :value="addOnClick" @input="addOnClick = $event" />
        <v-divider />
      </template>
      <template #layers="{ data }">
        <flowchart-points-editor
          :value="form.points"
          :shapes="data"
          :add-on-click="addOnClick"
          @input="updateField('points', $event)"
        />
      </template>
      <c-help-icon
        :text="$t('flowchart.panzoom.helpText')"
        size="32"
        icon-class="flowchart-map-editor__help-icon"
        color="secondary"
        icon="help"
        top
      />
    </flowchart>
  </c-layout>
</template>

<script>
import Flowchart from '@/components/common/flowchart/flowchart.vue';
import { formMixin } from '@/mixins/form';

import AddLocationBtn from './add-location-btn.vue';
import FlowchartPointsEditor from './flowchart-points-editor.vue';

export default {
  components: { Flowchart, FlowchartPointsEditor, AddLocationBtn },
  mixins: [formMixin],
  props: {
    form: {
      type: Object,
      required: true,
    },
    name: {
      type: String,
      default: 'parameters',
    },
  },
  emits: ['update:form'],
  data() {
    return {
      addOnClick: false,
      readonly: false,
    };
  },
  computed: {
    editorStyles() {
      return {};
    },
  },
  methods: {
    updateShapes(shapes) {
      const points = this.form.points.filter(point => (point.shape ? shapes[point.shape] : true));

      this.updateModel({
        ...this.form,
        shapes,
        points,
      });
    },

    updateBackgroundColor(color) {
      this.updateField('background_color', color);
    },
  },
};
</script>

<style lang="scss">
.flowchart-map-editor {
  --flowchart-map-editor-border-color: #e5e5e5;

  position: relative;
  border: 1px solid var(--flowchart-map-editor-border-color);
  height: 800px;

  &__help-icon.v-icon {
    position: absolute;
    right: 10px;
    bottom: 10px;
  }
}
</style>
