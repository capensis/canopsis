<template lang="pug">
  span
    template(v-if="isDropDown")
      v-list-tile(@click.stop="method")
        v-list-tile-title
          v-icon.pr-3(left, small) {{ icon }}
          span.body-1 {{ $t(title) }}
    template(v-else)
      v-btn.mx-1(flat, icon, @click.stop="method", :title="$t(title)")
        v-icon {{ icon }}
</template>

<script>
export default {
  props: {
    icon: {
      type: String,
      required: true,
    },
    method: {
      type: Function,
      required: true,
    },
    title: {
      type: String,
      default: '',
    },
    isDropDown: {
      type: Boolean,
      default: false,
    },
  },
};
</script>
