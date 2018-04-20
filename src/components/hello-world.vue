<template lang="pug">
  div
    h1.hello {{ $t('common.hello') }}
    h2 {{ msg }}
    v-container(fluid)
      v-layout(row, wrap)
        v-flex(xs6)
          v-subheader Example
        v-flex(xs6)
</template>

<script>
export default {
  name: 'HelloWorld',
  props: {
    msg: {
      type: String,
      required: true,
    },
  },
};
</script>
