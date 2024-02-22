<template>
  <v-layout column>
    <c-name-field
      v-field="form.name"
      required
    />
    <v-sheet>
      <cluster-graph-editor v-field="form.parameters" />
    </v-sheet>
  </v-layout>
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
