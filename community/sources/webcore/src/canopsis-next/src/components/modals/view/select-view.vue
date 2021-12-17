<template lang="pug">
  modal-wrapper(data-test="selectViewModal", close)
    template(slot="title")
      span {{ $t('modals.view.select.title') }}
    template(slot="text")
      v-list.py-0(dark)
        v-list-group(v-for="group in groups", :key="group._id")
          v-list-tile(slot="activator")
            v-list-tile-title {{ group.title }}
          v-list-tile(
            v-for="view in group.views",
            :key="view._id",
            @click="selectView(view._id)"
          )
            v-list-tile-title.pl-2 {{ view.title }}
</template>

<script>
import { MODALS } from '@/constants';

import { modalInnerMixin } from '@/mixins/modal/inner';
import { entitiesViewGroupMixin } from '@/mixins/entities/view/group';
import { permissionsEntitiesGroupMixin } from '@/mixins/permissions/entities/group';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.selectView,
  components: { ModalWrapper },
  mixins: [modalInnerMixin, entitiesViewGroupMixin, permissionsEntitiesGroupMixin],
  methods: {
    async selectView(viewId) {
      if (this.config.action) {
        await this.config.action(viewId);
      }

      this.$modals.hide();
    },
  },
};
</script>
