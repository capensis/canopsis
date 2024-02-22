<template>
  <div class="compiled-template__wrapper background darken-1 pa-1">
    <c-compiled-template :template="template" />
  </div>
</template>

<script>
export default {
  props: {
    template: {
      type: String,
      default: '',
    },
  },
};
</script>

<style lang="scss" scoped>
.compiled-template__wrapper {
  min-height: 5em;
}
</style>
