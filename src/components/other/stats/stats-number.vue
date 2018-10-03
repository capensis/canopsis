<template lang="pug">
  div {{ item.value }}
    sub {{ item.trend }}
</template>

<script>
export default {
  props: {
    item: {
      type: Object,
      required: true,
    },
  },
};
</script>
