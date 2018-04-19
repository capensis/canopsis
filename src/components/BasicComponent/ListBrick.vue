<template>

  <div @click="onClick">
    <transition name="reduce">
      <div v-if="!isExpanded">
        <slot name="reduced" :props="item"></slot>
      </div>
    </transition>
    <transition name="expand">
      <div v-if="!isExpanded">
        <slot v-if="isExpanded" name="expanded" :props="item"></slot>
      </div>
    </transition>
  </div>

</template>

<script>
export default {
  name: 'ListBrick',
  data() {
    return {
      isExpanded: false,
    };
  },
  props: ['item'],
  methods: {
    onClick: () => {
      console.log('lol');
      this.isExpanded = !this.isExpanded;
    },
  },
};
</script>

<style scoped>
  .expand-enter-active {
    animation: bounce-in .5s;
  }
  .expand-leave-active {
    animation: bounce-in .5s reverse;
  }
  @keyframes bounce-in {
    0% {
      transform: scale(0);
    }
    50% {
      transform: scale(1.5);
    }
    100% {
      transform: scale(1);
    }
  }
</style>
