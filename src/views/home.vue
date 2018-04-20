<template lang="pug">
  .home
    img(src="../assets/logo.png")
    hello-world(msg="Test props")
    test
</template>

<script>
import HelloWorld from '@/components/hello-world.vue';
import Test from '@/components/Context/Test.vue';

export default {
  name: 'home',
  components: {
    HelloWorld,
    Test,
  },
};
</script>
