<template lang="pug">
  span {{ name }}
</template>

<script>

export default {
  props: {
    name: {
      type: String,
      required: true,
    },
  },
};
</script>
