<template lang="pug">
  v-layout(column)
    c-name-field.mx-4(v-field="form.name")
    example-map(style="height: 600px;")
</template>

<script>
import ExampleMap from '@/components/examples/example-map.vue';

export default {
  inject: ['$validator'],
  components: { ExampleMap },
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
