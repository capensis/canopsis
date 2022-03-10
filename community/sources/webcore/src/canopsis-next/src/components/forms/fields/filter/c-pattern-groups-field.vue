<template lang="pug">
  v-layout(column)
    v-layout(v-for="(group, index) in groups", :key="group.key", wrap, row)
      v-flex(xs12)
        c-pattern-group-field(
          v-field="groups[index]",
          :attributes="attributes",
          :rules-map="rulesMap",
          @remove="removeItemFromArray(index)"
        )
      v-layout(v-show="index !== groups.length - 1", justify-center)
        c-pattern-operator-chip {{ $t('common.or') }}
    v-layout
      v-btn.mx-0(color="primary", @click="addFilterGroup") {{ $t('patterns.addGroup') }}
</template>

<script>
import { filterGroupToForm } from '@/helpers/forms/filter';

import { formArrayMixin } from '@/mixins/form';

export default {
  mixins: [formArrayMixin],
  model: {
    prop: 'groups',
    event: 'input',
  },
  props: {
    groups: {
      type: Array,
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
      default: 'groups',
    },
  },
  methods: {
    addFilterGroup() {
      const [firstAttribute] = this.attributes;

      this.addItemIntoArray(filterGroupToForm({
        rules: [
          { attribute: firstAttribute?.value },
        ],
      }));
    },
  },
};
</script>
