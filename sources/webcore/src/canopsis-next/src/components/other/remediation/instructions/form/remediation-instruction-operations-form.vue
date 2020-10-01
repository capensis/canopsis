<template lang="pug">
  v-layout.mt-2(column)
    v-layout(v-for="operation in operations", :key="operation.key") {{ operation }}
</template>

<script>
export default {
  inject: ['$validator'],
  model: {
    prop: 'operations',
    event: 'input',
  },
  props: {
    operations: {
      type: Array,
      default: () => ([]),
    },
    endpoint: {
      type: String,
      default: '',
    },
  },
};
</script>
