<template lang="pug">
  v-layout(column)
    c-name-field(v-field="form.name", required)
    v-sheet(light)
      cluster-graph-editor(v-field="form.parameters")
</template>

<script>
import ClusterGraphEditor from '@/components/common/chart/cluster-graph/cluster-graph-editor.vue';

export default {
  components: {
    ClusterGraphEditor,
  },
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      required: true,
    },
  },
};
</script>
