<template lang="pug">
  v-layout(column)
    c-description-field(v-field="form.description")
    c-duration-field(v-field="form.duration", long)
</template>

<script>
export default {
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      default: () => ({}),
    },
  },
};
</script>
