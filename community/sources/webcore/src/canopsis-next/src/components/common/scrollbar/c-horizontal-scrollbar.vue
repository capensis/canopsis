<template>
  <div class="horizontal-scrollbar" @scroll="onScroll">
    <div ref="content" class="horizontal-scrollbar__content" />
  </div>
</template>

<script>
export default {
  beforeDestroy() {
    this.disconnect();
  },
  methods: {
    connect(el) {
      this.disconnect();

      this.$connectedEl = el;
      this.$connectedEl?.addEventListener('scroll', this.onConnectedElScroll);

      this.setContentWidth(el?.scrollWidth ?? 0);
    },

    disconnect() {
      this.$connectedEl?.removeEventListener('scroll', this.onConnectedElScroll);
      delete this.$connectedEl;
    },

    setContentWidth(width) {
      if (this.$refs.content) {
        this.$refs.content.style.width = `${width}px`;
      }
    },

    setOpacity(opacity) {
      this.$el.style.opacity = opacity;
    },

    setTranslateY(translateY) {
      this.$el.style.transform = `translateY(${translateY}px)`;
    },

    onScroll() {
      if (this.$connectedEl) {
        this.$connectedEl.scrollLeft = this.$el.scrollLeft;
      }
    },

    onConnectedElScroll() {
      this.$el.scrollLeft = this.$connectedEl.scrollLeft;
    },
  },
};
</script>

<style lang="scss">
.horizontal-scrollbar {
  width: 100%;
  overflow-y: hidden;
  overflow-x: auto;
  display: block;
  transition: 0.3s cubic-bezier(0.25, 0.8, 0.5, 1);
  transition-property: opacity;

  &__content {
    height: 1px;
  }
}
</style>
