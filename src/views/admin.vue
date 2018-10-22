<template lang="pug">
  div
    v-tabs(
    v-model="activeTab",
    centered,
    dark,
    slider-color="green"
    )
      v-tab Rights
      v-tab Users
      v-tab Roles
      v-tab Parameters
</template>

<script>
export default {
  data() {
    return {
      activeTab: 0,
    };
  },
};
</script>

