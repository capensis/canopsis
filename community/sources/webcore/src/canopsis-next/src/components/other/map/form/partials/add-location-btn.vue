<template lang="pug">
  v-btn-toggle(:value="value", light, @change="inputHandler")
    v-btn.ma-0(:value="true", flat, large, block)
      v-icon(left) add_location
      span {{ $t('map.addPoint') }}
</template>

<script>
import { formBaseMixin } from '@/mixins/form';

export default {
  mixins: [formBaseMixin],
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: Boolean,
      required: true,
    },
  },
  methods: {
    inputHandler(value) {
      this.updateModel(!!value);
    },
  },
};
</script>
