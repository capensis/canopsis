<template>
  <v-btn-toggle
    :model-value="value ? true : undefined"
    variant="text"
    @update:model-value="inputHandler"
  >
    <v-btn :value="true" class="ma-0" block>
      <v-icon start> add_location </v-icon>
      <span>{{ $t('map.addPoint') }}</span>
    </v-btn>
  </v-btn-toggle>
</template>

<script>
import { formBaseMixin } from '@/mixins/form';

export default {
  mixins: [formBaseMixin],
  props: {
    value: {
      type: Boolean,
      required: true,
    },
  },
  emits: ['input'],
  methods: {
    inputHandler(value) {
      this.updateModel(!!value);
    },
  },
};
</script>
