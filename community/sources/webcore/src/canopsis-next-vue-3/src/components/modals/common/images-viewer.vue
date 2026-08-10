<template>
  <div>
    <c-images-viewer :images="imagesWithKey" :active="config.active" item-key="key" />

    <v-btn class="ma-1 images-viewer__close-btn" icon @click="close">
      <c-icon color="white">
        close
      </c-icon>
    </v-btn>
  </div>
</template>

<script>
import CImagesViewer from '@/components/common/images-viewer/c-images-viewer.vue';
import { addKeyInEntities } from '@/helpers/array';

export default {
  components: { CImagesViewer },

  props: {
    modal: {
      type: Object,
      required: true,
    },
  },

  computed: {
    config() {
      return this.modal.config ?? {};
    },

    imagesWithKey() {
      return addKeyInEntities(this.config.images);
    },
  },

  methods: {
    close() {
      this.$modals.hide({ id: this.modal.id });
    },
  },
};
</script>

<style lang="scss">
.v-dialog__images-viewer {
  width: auto;
  position: relative;
  overflow: hidden;

  .images-viewer__close-btn {
    position: absolute;
    top: 0;
    right: 0;
  }
}
</style>
