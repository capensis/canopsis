<template>
  <c-layout column>
    <c-name-field
      :model-value="form.name"
      autofocus
      required
      @update:model-value="updateField('name', $event)"
    />
    <v-sheet>
      <cluster-graph-editor
        :form="form.parameters"
        @update:form="updateField('parameters', $event)"
      />
    </v-sheet>
  </c-layout>
</template>

<script>
import { formMixin } from '@/mixins/form';

import ClusterGraphEditor from '@/components/common/chart/cluster-graph/cluster-graph-editor.vue';

export default {
  components: {
    ClusterGraphEditor,
  },
  mixins: [formMixin],
  props: {
    form: {
      type: Object,
      required: true,
    },
  },
  emits: ['update:form'],
};
</script>
