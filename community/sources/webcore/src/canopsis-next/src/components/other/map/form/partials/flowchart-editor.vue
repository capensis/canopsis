<template lang="pug">
  v-layout(column)
    flowchart.flowchart-map-editor.mb-2(
      :shapes="form.shapes",
      :background-color="form.background_color",
      :style="editorStyles",
      :cursor-style="addOnClick ? 'crosshair' : undefined",
      @input="updateShapes",
      @update:backgroundColor="updateBackgroundColor"
    )
      template(#sidebar-prepend="{ data }")
        add-location-btn(v-model="addOnClick")
        v-divider
      template(#layers="{ data }")
        flowchart-points-editor(
          v-field="form.points",
          :shapes="data",
          :add-on-click="addOnClick"
        )
      c-help-icon(
        :text="$t('flowchart.panzoom.helpText')",
        size="32",
        icon-class="flowchart-map-editor__help-icon",
        color="secondary",
        icon="help",
        top
      )
    v-messages(v-if="hasChildrenError", :value="errorMessages", color="error")
</template>

<script>
import { COLORS } from '@/config';

import { formMixin, validationChildrenMixin } from '@/mixins/form';

import Flowchart from '@/components/common/flowchart/flowchart.vue';

import FlowchartPointsEditor from './flowchart-points-editor.vue';
import AddLocationBtn from './add-location-btn.vue';

export default {
  inject: ['$validator'],
  components: { Flowchart, FlowchartPointsEditor, AddLocationBtn },
  mixins: [formMixin, validationChildrenMixin],
  model: {
    prop: 'form',
    event: 'input',
  },
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
  data() {
    return {
      addOnClick: false,
      readonly: false,
    };
  },
  computed: {
    errorMessages() {
      return [this.$t('flowchart.errors.pointsRequired')];
    },

    editorStyles() {
      return {
        borderColor: this.hasChildrenError ? COLORS.error : undefined,
      };
    },
  },
  watch: {
    form: {
      deep: true,
      handler() {
        if (this.hasChildrenError) {
          this.$validator.validate(this.name);
        }
      },
    },
  },
  mounted() {
    this.attachRequiredRule();
  },
  beforeDestroy() {
    this.detachRules();
  },
  methods: {
    attachRequiredRule() {
      this.$validator.attach({
        name: this.name,
        rules: 'required:true',
        getter: () => !!this.form.points.length,
        vm: this,
      });
    },

    detachRules() {
      this.$validator.detach(this.name);
    },

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
$borderColor: #e5e5e5;

.flowchart-map-editor {
  position: relative;
  border: 1px solid $borderColor;
  height: 800px;

  &__help-icon {
    position: absolute;
    right: 10px;
    bottom: 10px;
  }
}
</style>
