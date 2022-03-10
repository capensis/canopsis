<template lang="pug">
  v-layout
    v-flex.pa-2.mr-2
      c-pattern-operator-information {{ $t('common.and') }}
    v-flex(xs11)
      v-layout(column)
        c-pattern-rules-field(
          v-field="group.rules",
          :attributes="attributes",
          :rules-map="rulesMap",
          @remove="$emit('remove')"
        )
</template>

<script>
import { formArrayMixin } from '@/mixins/form';

export default {
  mixins: [formArrayMixin],
  model: {
    prop: 'group',
    event: 'input',
  },
  props: {
    group: {
      type: Object,
      required: true,
    },
    attributes: {
      type: Array,
      required: true,
    },
    rulesMap: {
      type: Object,
      default: () => ({}),
    },
    name: {
      type: String,
      default: 'group',
    },
  },
};
</script>
