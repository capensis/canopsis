<template lang="pug">
  div {{ map | json }}
</template>

<script>
export default {
  props: {
    map: {
      type: Object,
      required: true,
    },
  },
};
</script>
