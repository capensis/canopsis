<template lang="pug">
  v-btn(@click.stop="deleteInfoObject", icon, small)
    v-icon delete
</template>

<script>
import cloneDeep from 'lodash/cloneDeep';
import InnerModalMixin from '@/mixins/modal/modal-inner';
import ContextMixin from '@/mixins/context/list';
import { MODALS } from '@/constants';

export default {
  mixins: [
    InnerModalMixin,
    ContextMixin,
  ],
  props: {
    infoObjectName: {
      type: String,
      required: true,
    },
    contextEntity: {
      type: Object,
      required: true,
    },
  },
  methods: {
    deleteInfoObject() {
      this.showModal({
        name: MODALS.confirmation,
        config: {
          action: async () => {
            const updatedEntity = cloneDeep(this.contextEntity);
            delete updatedEntity.props.infos[this.infoObjectName];

            await this.updateContextEntity({
              entity: updatedEntity,
            });

            this.fetchList();
          },
        },
      });
    },
  },
};
</script>
