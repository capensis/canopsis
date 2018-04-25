<template lang="pug">
  .home
    img(src="../assets/logo.png")
    hello-world(msg="Test props")
</template>

<script>
import HelloWorld from '@/components/hello-world.vue';

export default {
  name: 'home',
  components: {
    HelloWorld,
  },
};
</script>
