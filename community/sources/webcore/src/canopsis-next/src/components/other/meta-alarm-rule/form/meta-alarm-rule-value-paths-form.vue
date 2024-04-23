<template>
  <v-layout class="gap-3" column>
    <c-information-block :title="$tc('metaAlarmRule.valuePath', 2)">
      <meta-alarm-value-paths-field
        v-field="form.value_paths"
        :label="$tc('metaAlarmRule.valuePath', 1)"
        :name="name"
        required
      />
    </c-information-block>
  </v-layout>
</template>

<script>
import MetaAlarmValuePathsField from './fields/meta-alarm-value-paths-field.vue';

export default {
  components: { MetaAlarmValuePathsField },
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      default: () => ({}),
    },
    name: {
      type: String,
      default: 'value_paths',
    },
  },
};
</script>
