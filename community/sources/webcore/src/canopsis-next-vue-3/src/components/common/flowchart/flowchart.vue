<template>
  <div class="flowchart fill-height">
    <flowchart-sidebar
      v-if="!readonly"
      v-model:selected="selected"
      :shapes="shapes"
      :view-box="viewBox"
      :background-color="backgroundColor"
      :readonly="readonly"
      class="flowchart__sidebar"
      @update:shapes="updateShapes"
      @update:backgroundColor="$emit('update:backgroundColor', $event)"
    >
      <template #prepend="">
        <slot name="sidebar-prepend" />
      </template>
    </flowchart-sidebar>
    <c-zoom-overlay
      :class="{ 'flowchart__editor--readonly': readonly }"
      class="flowchart__editor"
      skip-alt
      skip-shift
    >
      <flowchart-editor
        v-model:view-box="viewBox"
        v-model:selected="selected"
        :shapes="shapes"
        :background-color="backgroundColor"
        :readonly="readonly"
        :cursor-style="cursorStyle"
        @update:shapes="updateShapes"
      >
        <template #layers="{ data }">
          <slot :data="data" name="layers" />
        </template>
      </flowchart-editor>
    </c-zoom-overlay>
    <div v-if="selected.length" class="flowchart__properties">
      <flowchart-properties
        v-if="!readonly"
        :shapes="shapes"
        :selected="selected"
        @update:shapes="updateShapes"
      />
    </div>
    <slot />
  </div>
</template>

<script>
import FlowchartEditor from './flowchart-editor.vue';
import FlowchartProperties from './flowchart-properties.vue';
import FlowchartSidebar from './flowchart-sidebar.vue';

export default {
  components: {
    FlowchartSidebar,
    FlowchartEditor,
    FlowchartProperties,
  },
  model: {
    prop: 'shapes',
    event: 'input',
  },
  emits: ['input', 'update:shapes', 'update:backgroundColor'],
  props: {
    shapes: {
      type: Object,
      required: true,
    },
    cursorStyle: {
      type: String,
      required: false,
    },
    backgroundColor: {
      type: String,
      required: false,
    },
    readonly: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      selected: [],
      viewBox: {
        x: 0,
        y: 0,
        width: 0,
        height: 0,
      },
    };
  },
  methods: {
    updateShapes(shapes) {
      this.$emit('input', shapes);
      this.$emit('update:shapes', shapes);
    },
  },
};
</script>

<style lang="scss">
$sidebarWidth: 320px;

.flowchart {
  position: relative;
  display: flex;
  background: white;

  &__sidebar {
    width: $sidebarWidth !important;
  }

  &__editor {
    flex-grow: 1;
    width: auto;
    height: 100%;

    position: absolute;
    left: $sidebarWidth;
    top: 0;
    right: 0;
    bottom: 0;

    user-select: none;

    &--readonly {
      left: 0;
    }
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
