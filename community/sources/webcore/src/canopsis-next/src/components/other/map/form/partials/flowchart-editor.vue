<template lang="pug">
  v-layout(column)
    flowchart.flowchart-editor.mb-2(
      v-field="form.shapes",
      :background-color.sync="form.backgroundColor",
      :style="editorStyles"
    )
      template(#layers="")
        flowchart-points-editor(v-field="form.points")
    v-messages
</template>

<script>
import { COLORS } from '@/config';

import { formMixin, validationChildrenMixin } from '@/mixins/form';

import Flowchart from '@/components/common/flowchart/flowchart.vue';

import FlowchartPointsEditor from './flowchart-points-editor.vue';

export default {
  inject: ['$validator'],
  components: { Flowchart, FlowchartPointsEditor },
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
  computed: {
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
        getter: () => true,
        context: () => this,
        vm: this,
      });
    },

    detachRules() {
      this.$validator.detach(this.name);
    },
  },
};
</script>

<style lang="scss">
$borderColor: #e5e5e5;

.flowchart-editor {
  border: 1px solid $borderColor;
  height: 800px;
}
</style>
