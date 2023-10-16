<template>
  <v-btn-toggle
    :value="value"
    @change="inputHandler"
  >
    <v-btn
      class="ma-0"
      :value="true"
      light="light"
      text
      large="large"
      block="block"
    >
      <v-icon left="left">
        add_location
      </v-icon><span>{{ $t('map.addPoint') }}</span>
    </v-btn>
  </v-btn-toggle>
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
