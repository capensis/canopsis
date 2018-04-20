<template lang="pug">
  .home
    img(src="../assets/logo.png")
    hello-world(msg="Test props")
    MFilterEditor
</template>

<script>
import HelloWorld from '@/components/hello-world.vue';
import MFilterEditor from '@/components/mfilter-editor.vue';

export default {
  name: 'home',
  components: {
    HelloWorld,
    MFilterEditor,
  },
};
</script>
