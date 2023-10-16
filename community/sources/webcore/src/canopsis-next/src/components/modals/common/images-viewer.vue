<template>
  <div>
    <c-images-viewer
      :images="imagesWithKey"
      :active="config.active"
      item-key="key"
    />
    <v-btn
      class="ma-1 images-viewer__close-btn"
      icon
      @click="$modals.hide"
    >
      <v-icon color="white">
        close
      </v-icon>
    </v-btn>
  </div>
</template>

<script>
import { MODALS } from '@/constants';

import { addKeyInEntities } from '@/helpers/array';

import { modalInnerMixin } from '@/mixins/modal/inner';

export default {
  name: MODALS.imagesViewer,
  mixins: [modalInnerMixin],
  computed: {
    imagesWithKey() {
      return addKeyInEntities(this.config.images);
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
