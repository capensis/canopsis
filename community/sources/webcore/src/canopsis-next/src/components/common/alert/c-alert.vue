<template lang="pug">
  v-alert.c-alert(v-bind="$attrs", :value="true")
    slot
</template>

<script>
export default {
  inheritAttrs: false,
};
</script>

<style>
.c-alert {
  width: 100%;
}
</style>
