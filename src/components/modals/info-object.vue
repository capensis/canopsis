<template lang="pug">
  v-card
    info-object-form(@submit="handleSubmit", :infoObject="config.infoObject", :forbiddenNames="forbiddenNames")
</template>

<script>
import cloneDeep from 'lodash/cloneDeep';
import InfoObjectForm from '@/components/other/info-object/form.vue';
import InnerModalMixin from '@/mixins/modal/modal-inner';
import ContextMixin from '@/mixins/context/list';

export default {
  components: {
    InfoObjectForm,
  },
  mixins: [
    InnerModalMixin,
    ContextMixin,
  ],
  computed: {
    forbiddenNames() {
      return Object.keys(this.config.entity.props.infos);
    },
  },
  methods: {
    async handleSubmit(infoObjectData) {
      const updatedEntity = cloneDeep(this.config.entity.props);
      updatedEntity.infos[infoObjectData.name] = infoObjectData;

      if (this.config.infoObject) {
        delete updatedEntity.infos[this.config.infoObject.name];
      }

      await this.updateContextEntity({
        entity: updatedEntity,
      });

      await this.fetchList();

      this.hideModal();
    },
  },
};
</script>
